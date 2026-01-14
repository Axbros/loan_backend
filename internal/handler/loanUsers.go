package handler

import (
	"encoding/base64"
	"errors"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/jwt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"loan/internal/cache"
	"loan/internal/dao"
	"loan/internal/database"
	"loan/internal/ecode"
	"loan/internal/model"
	"loan/internal/types"
)

var _ LoanUsersHandler = (*loanUsersHandler)(nil)

// LoanUsersHandler defining the handler interface
type LoanUsersHandler interface {
	Register(*gin.Context)
	Login(*gin.Context)
	Me(*gin.Context)
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)

	DeleteByIDs(c *gin.Context)
	GetByCondition(c *gin.Context)
	ListByIDs(c *gin.Context)
	ListByLastID(c *gin.Context)
}

type loanUsersHandler struct {
	iDao dao.LoanUsersDao
}

// NewLoanUsersHandler creating the handler interface
func NewLoanUsersHandler() LoanUsersHandler {
	return &loanUsersHandler{
		iDao: dao.NewLoanUsersDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanUsersCache(database.GetCacheType()),
		),
	}
}

func (h *loanUsersHandler) Register(c *gin.Context) {
	form := &types.RegisterRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	username := strings.TrimSpace(form.Username)
	password := form.Password

	if username == "" || len(username) < 3 || len(username) > 64 {
		response.Error(c, ecode.InvalidParams)
		return
	}
	// bcrypt 只处理前72字节，超过会被截断 -> 直接拒绝
	if len(password) < 8 || len(password) > 72 {
		response.Error(c, ecode.InvalidParams)
		return
	}

	username = strings.ToLower(username)

	// hash 密码
	pwHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("GenerateFromPassword error", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrHashPassword)
		return
	}

	loanUsers := &model.LoanUsers{}
	// 如果你希望 RegisterRequest 也能 copier.Copy（比如还有别的字段），可以 copy
	// 但 password 绝对不要 copy 进 model
	_ = copier.Copy(loanUsers, form) // 即使失败也不影响核心字段，我们下面手动赋值
	loanUsers.Username = username
	loanUsers.PasswordHash = string(pwHash)

	// 默认值（你可以按业务调整）
	loanUsers.DepartmentID = 1 // 没有选择部门时给默认部门；以后可改成从 form 传
	loanUsers.MfaEnabled = 0
	loanUsers.MfaRequired = 0
	loanUsers.Status = 1

	now := time.Now()
	loanUsers.CreatedAt = now
	loanUsers.UpdatedAt = now

	ctx := middleware.WrapCtx(c)

	// share_code 生成 + 处理唯一冲突重试
	for i := 0; i < 5; i++ {
		loanUsers.ShareCode = genShareCode(12)

		err = h.iDao.Create(ctx, loanUsers)
		if err == nil {
			response.Success(c, gin.H{
				"id":         loanUsers.ID,
				"username":   loanUsers.Username,
				"share_code": loanUsers.ShareCode,
			})
			return
		}

		// 唯一冲突（MySQL 1062）
		if isDuplicateKeyErr(err) {
			// 如果是用户名冲突：直接返回“用户名已存在”
			if isDuplicateUsernameErr(err) {
				response.Error(c, ecode.UsernameAlreadyExists) // 你需要补这个 ecode
				return
			}
			// 否则大概率是 share_code 冲突：重试生成
			continue
		}

		logger.Error("Register Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	// share_code 连续冲突基本不可能，按内部错误处理
	response.Output(c, ecode.InternalServerError.ToHTTPCode())
}

func (h *loanUsersHandler) Login(c *gin.Context) {
	form := &types.LoginRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	username := strings.ToLower(strings.TrimSpace(form.Username))
	password := form.Password

	if username == "" || password == "" {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)

	// 1) 查用户
	user, err := h.iDao.GetByUsername(ctx, username)
	if err != nil {
		logger.Error("GetByUsername error", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	if user == nil || user.ID == 0 {
		// 不要提示“用户名不存在”，统一返回账号密码错误
		response.Error(c, ecode.UsernameOrPasswordIncorrect)
		return
	}

	// 2) 状态检查
	if user.Status != 1 {
		response.Error(c, ecode.UserDisabled) // 你可以自己定义这个 ecode
		return
	}

	// 3) 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		response.Error(c, ecode.UsernameOrPasswordIncorrect)
		return
	}

	// 4) 生成 token（你封装的函数）
	token, err := generateAccessToken(strconv.FormatInt(int64(user.ID), 10))
	if err != nil {
		logger.Error("GenerateAccessToken error", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	// 5) 返回
	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":            user.ID,
			"username":      user.Username,
			"department_id": user.DepartmentID,
			"mfa_enabled":   user.MfaEnabled,
		},
	})
}

func (h *loanUsersHandler) Me(c *gin.Context) {
	uid, ok := getUIDFromClaims(c)
	if !ok || uid == 0 {
		response.Out(c, ecode.Unauthorized)
		return
	}

	ctx := middleware.WrapCtx(c)

	user, err := h.iDao.GetByID(ctx, uid) // 你实现：SELECT * FROM loan_users WHERE id=? AND deleted_at IS NULL
	if err != nil {
		logger.Error("GetByID error", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	if user == nil || user.ID == 0 || user.Status != 1 {
		response.Out(c, ecode.Unauthorized)
		return
	}

	roles, err := h.iDao.GetRoleCodesByUserID(ctx, uid) // 返回 []string
	if err != nil {
		logger.Error("GetRoleCodesByUserID error", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	perms, err := h.iDao.GetPermCodesByUserID(ctx, uid) // 返回 []string
	if err != nil {
		logger.Error("GetPermCodesByUserID error", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{
		"user": gin.H{
			"id":            user.ID,
			"username":      user.Username,
			"department_id": user.DepartmentID,
			"mfa_enabled":   user.MfaEnabled,
		},
		"roles": roles,
		"perms": perms,
	})
}

// Create a new loanUsers
// @Summary Create a new loanUsers
// @Description Creates a new loanUsers entity using the provided data in the request body.
// @Tags loanUsers
// @Accept json
// @Produce json
// @Param data body types.CreateLoanUsersRequest true "loanUsers information"
// @Success 200 {object} types.CreateLoanUsersReply{}
// @Router /api/v1/loanUsers [post]
// @Security BearerAuth
func (h *loanUsersHandler) Create(c *gin.Context) {
	form := &types.CreateLoanUsersRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanUsers := &model.LoanUsers{}
	err = copier.Copy(loanUsers, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanUsers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanUsers)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanUsers.ID})
}

// DeleteByID delete a loanUsers by id
// @Summary Delete a loanUsers by id
// @Description Deletes a existing loanUsers identified by the given id in the path.
// @Tags loanUsers
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanUsersByIDReply{}
// @Router /api/v1/loanUsers/{id} [delete]
// @Security BearerAuth
func (h *loanUsersHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanUsersIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByID(ctx, id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update a loanUsers by id
// @Summary Update a loanUsers by id
// @Description Updates the specified loanUsers by given id in the path, support partial update.
// @Tags loanUsers
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanUsersByIDRequest true "loanUsers information"
// @Success 200 {object} types.UpdateLoanUsersByIDReply{}
// @Router /api/v1/loanUsers/{id} [put]
// @Security BearerAuth
func (h *loanUsersHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanUsersIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanUsersByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanUsers := &model.LoanUsers{}
	err = copier.Copy(loanUsers, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanUsers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanUsers)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanUsers by id
// @Summary Get a loanUsers by id
// @Description Gets detailed information of a loanUsers specified by the given id in the path.
// @Tags loanUsers
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanUsersByIDReply{}
// @Router /api/v1/loanUsers/{id} [get]
// @Security BearerAuth
func (h *loanUsersHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanUsersIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUsers, err := h.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.LoanUsersObjDetail{}
	err = copier.Copy(data, loanUsers)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanUsers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanUsers": data})
}

// List get a paginated list of loanUserss by custom conditions
// @Summary Get a paginated list of loanUserss by custom conditions
// @Description Returns a paginated list of loanUsers based on query filters, including page number and size.
// @Tags loanUsers
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanUserssReply{}
// @Router /api/v1/loanUsers/list [post]
// @Security BearerAuth
func (h *loanUsersHandler) List(c *gin.Context) {
	form := &types.ListLoanUserssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUserss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanUserss(loanUserss)
	if err != nil {
		response.Error(c, ecode.ErrListLoanUsers)
		return
	}

	response.Success(c, gin.H{
		"loanUserss": data,
		"total":      total,
	})
}

// DeleteByIDs batch delete loanUsers by ids
// @Summary Batch delete loanUsers by ids
// @Description Deletes multiple loanUsers by a list of id
// @Tags loanUsers
// @Param data body types.DeleteLoanUserssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanUserssByIDsReply{}
// @Router /api/v1/loanUsers/delete/ids [post]
// @Security BearerAuth
func (h *loanUsersHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanUserssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err = h.iDao.DeleteByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByCondition get a loanUsers by custom condition
// @Summary Get a loanUsers by custom condition
// @Description Returns a single loanUsers that matches the specified filter conditions.
// @Tags loanUsers
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanUsersByConditionReply{}
// @Router /api/v1/loanUsers/condition [post]
// @Security BearerAuth
func (h *loanUsersHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanUsersByConditionRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	err = form.Conditions.CheckValid()
	if err != nil {
		logger.Warn("Parameters error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUsers, err := h.iDao.GetByCondition(ctx, &form.Conditions)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByCondition not found", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByCondition error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.LoanUsersObjDetail{}
	err = copier.Copy(data, loanUsers)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanUsers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanUsers": data})
}

// ListByIDs batch get loanUsers by ids
// @Summary Batch get loanUsers by ids
// @Description Returns a list of loanUsers that match the list of id.
// @Tags loanUsers
// @Param data body types.ListLoanUserssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanUserssByIDsReply{}
// @Router /api/v1/loanUsers/list/ids [post]
// @Security BearerAuth
func (h *loanUsersHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanUserssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUsersMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanUserss := []*types.LoanUsersObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanUsersMap[id]; ok {
			record, err := convertLoanUsers(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanUsers)
				return
			}
			loanUserss = append(loanUserss, record)
		}
	}

	response.Success(c, gin.H{
		"loanUserss": loanUserss,
	})
}

// ListByLastID get a paginated list of loanUserss by last id
// @Summary Get a paginated list of loanUserss by last id
// @Description Returns a paginated list of loanUserss starting after a given last id, useful for cursor-based pagination.
// @Tags loanUsers
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanUserssReply{}
// @Router /api/v1/loanUsers/list [get]
// @Security BearerAuth
func (h *loanUsersHandler) ListByLastID(c *gin.Context) {
	lastID := utils.StrToUint64(c.Query("lastID"))
	if lastID == 0 {
		lastID = math.MaxInt32
	}
	limit := utils.StrToInt(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}
	sort := c.Query("sort")

	ctx := middleware.WrapCtx(c)
	loanUserss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanUserss(loanUserss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanUsers)
		return
	}

	response.Success(c, gin.H{
		"loanUserss": data,
	})
}

func getLoanUsersIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanUsers(loanUsers *model.LoanUsers) (*types.LoanUsersObjDetail, error) {
	data := &types.LoanUsersObjDetail{}
	err := copier.Copy(data, loanUsers)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanUserss(fromValues []*model.LoanUsers) ([]*types.LoanUsersObjDetail, error) {
	toValues := []*types.LoanUsersObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanUsers(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}

func genShareCode(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	s := base64.RawURLEncoding.EncodeToString(b)
	if len(s) >= n {
		return s[:n]
	}
	return s
}

func isDuplicateKeyErr(err error) bool {
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		return me.Number == 1062
	}
	return false
}

func isDuplicateUsernameErr(err error) bool {
	var me *mysql.MySQLError
	if errors.As(err, &me) && me.Number == 1062 {
		// 你的唯一索引名是 `username`
		return strings.Contains(me.Message, "for key 'username'") ||
			strings.Contains(me.Message, "for key 'loan_users.username'")
	}
	return false
}
func generateAccessToken(uid string) (string, error) {
	_, token, err := jwt.GenerateToken(uid)
	return token, err
}

func getUIDFromClaims(c *gin.Context) (uint64, bool) {
	v, ok := c.Get("claims")
	if !ok || v == nil {
		return 0, false
	}
	claims, ok := v.(*jwt.Claims)
	if !ok || claims == nil {
		return 0, false
	}

	uidStr := claims.UID
	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		return 0, false
	}
	return uid, true
}
