package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/utils"
	"github.com/google/uuid"
	"io"
	"loan/internal/config"
	"loan/internal/tool"
	"math"
	"os"
	"path/filepath"
	"strings"

	"loan/internal/cache"
	"loan/internal/dao"
	"loan/internal/database"
	"loan/internal/ecode"
	"loan/internal/model"
	"loan/internal/types"
)

var _ LoanRepaymentTransactionsHandler = (*loanRepaymentTransactionsHandler)(nil)

// LoanRepaymentTransactionsHandler defining the handler interface
type LoanRepaymentTransactionsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)

	DeleteByIDs(c *gin.Context)
	GetByCondition(c *gin.Context)
	ListByIDs(c *gin.Context)
	ListByLastID(c *gin.Context)
	DetailByScheduleID(c *gin.Context)
	History(c *gin.Context)
	UploadVoucher(c *gin.Context)
	GetVoucherBase64(c *gin.Context)
}

type loanRepaymentTransactionsHandler struct {
	iDao dao.LoanRepaymentTransactionsDao
}

// NewLoanRepaymentTransactionsHandler creating the handler interface
func NewLoanRepaymentTransactionsHandler() LoanRepaymentTransactionsHandler {
	return &loanRepaymentTransactionsHandler{
		iDao: dao.NewLoanRepaymentTransactionsDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanRepaymentTransactionsCache(database.GetCacheType()),
		),
	}
}

func (h *loanRepaymentTransactionsHandler) GetVoucherBase64(c *gin.Context) {
	// 1. 获取请求参数（file_name）
	fileName := c.Param("file_name") // 从URL查询参数获取，也可改用PostForm
	if fileName == "" {
		response.Error(c, ecode.InvalidParams)
		return
	}

	// 2. 定义存储目录并清理格式（和上传接口保持一致）
	storageDir := config.Get().Storage.Voucher
	cleanStorageDir := filepath.Clean(storageDir)

	// 3. 拼接文件完整路径并校验合法性（防止路径穿越）
	filePath := filepath.Join(storageDir, fileName)
	cleanFilePath := filepath.Clean(filePath)
	if !strings.HasPrefix(cleanFilePath, cleanStorageDir) {
		response.Error(c, ecode.ErrInvalidFilePath)
		return
	}

	// 4. 检查文件是否存在
	if _, err := os.Stat(cleanFilePath); os.IsNotExist(err) {
		response.Error(c, ecode.FileNotFound) // 需新增文件不存在的错误码
		return
	} else if err != nil {
		response.Error(c, ecode.ErrReadFile) // 需新增文件读取失败的错误码
		return
	}

	// 5. 读取文件内容
	fileContent, err := os.ReadFile(cleanFilePath)
	if err != nil {
		response.Error(c, ecode.ErrReadFile)
		return
	}

	// 6. 将文件内容编码为Base64
	base64Str := base64.StdEncoding.EncodeToString(fileContent)

	// 7. （可选）生成带前缀的Base64（方便前端直接渲染图片）
	// 获取文件后缀，拼接Data URI前缀
	fileExt := strings.ToLower(filepath.Ext(fileName))
	mimeType := ""
	switch fileExt {
	case ".png":
		mimeType = "image/png"
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	default:
		mimeType = "application/octet-stream"
	}
	base64WithPrefix := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str)

	// 8. 返回结果
	response.Success(c, gin.H{
		//"file_name":          fileName,
		//"base64":             base64Str,        // 纯Base64编码
		"base64_with_prefix": base64WithPrefix, // 带Data URI前缀的Base64（前端可直接用）
		"file_size":          len(fileContent), // 文件大小
	})
}

func (h *loanRepaymentTransactionsHandler) DetailByScheduleID(c *gin.Context) {
	form := &types.DetailByScheduleIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)
	detail, err := h.iDao.DetailByScheduleID(ctx, form.ScheduleID)
	if err != nil {
		response.Error(c, ecode.ErrGetByConditionLoanRepaymentTransactions)
		return
	}
	response.Success(c, detail)
}

// UploadVoucher 上传交易凭证图片
func (h *loanRepaymentTransactionsHandler) UploadVoucher(c *gin.Context) {
	// 1. 从表单获取上传的文件（表单字段名：voucher）
	file, fileHeader, err := c.Request.FormFile("voucher")
	if err != nil {
		response.Error(c, ecode.InvalidParams) // 替换为你实际的参数错误码
		return
	}
	// 优化：使用 defer 关闭上传文件流（原代码已有，但需确保执行）
	defer func() {
		_ = file.Close() // 忽略关闭错误，或根据业务记录日志
	}()

	// 2. 校验文件类型（仅允许 png/jpg/jpeg）
	allowedExts := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
	}
	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExts[fileExt] {
		response.Error(c, ecode.UnsupportedFileType)
		return
	}

	// 3. 确保存储目录存在（不存在则创建）
	storageDir := config.Get().Storage.Voucher
	cleanStorageDir := filepath.Clean(storageDir) // 此时变为 "storage/transaction-vouchers"

	if err := os.MkdirAll(cleanStorageDir, 0755); err != nil {
		response.Error(c, ecode.ErrCreateFileFolder)
		return
	}

	// 4. 生成唯一文件名（避免覆盖）
	uniqueID := uuid.New().String()
	filename := fmt.Sprintf("%s%s", uniqueID, fileExt)
	filePath := filepath.Join(storageDir, filename)

	// 5. 创建本地文件（优化：先判断路径合法性，避免路径穿越）
	if !strings.HasPrefix(filepath.Clean(filePath), cleanStorageDir) {
		response.Error(c, ecode.ErrInvalidFilePath)
		return
	}
	dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		response.Error(c, ecode.ErrSaveFile)
		return
	}
	// 优化：安全关闭文件（先判断非 nil，避免 panic）
	defer func() {
		if dstFile != nil {
			_ = dstFile.Close()
		}
	}()

	// 6. 拷贝文件内容（优化：用 io.Copy 替代 ReadFrom，更通用）
	_, err = io.Copy(dstFile, file)
	if err != nil {
		response.Error(c, ecode.ErrSaveFile)
		return
	}

	// 7. 返回成功结果（相对路径供后续 POST 使用）
	//relativePath := fmt.Sprintf("/storage/transaction-vouchers/%s", filename)
	response.Success(c, gin.H{
		"file_name": filename,
		//"file_path": relativePath,
		"size": fileHeader.Size,
	})
}

func (h *loanRepaymentTransactionsHandler) History(c *gin.Context) {
	form := &types.DetailByScheduleIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)

	transactions, err := h.iDao.GetByScheduleID(ctx, form.ScheduleID)
	if err != nil {
		response.Error(c, ecode.ErrGetByConditionLoanRepaymentTransactions)
		return
	}
	response.Success(c, transactions)
}

// Create a new loanRepaymentTransactions
// @Summary Create a new loanRepaymentTransactions
// @Description Creates a new loanRepaymentTransactions entity using the provided data in the request body.
// @Tags loanRepaymentTransactions
// @Accept json
// @Produce json
// @Param data body types.CreateLoanRepaymentTransactionsRequest true "loanRepaymentTransactions information"
// @Success 200 {object} types.CreateLoanRepaymentTransactionsReply{}
// @Router /api/v1/loanRepaymentTransactions [post]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) Create(c *gin.Context) {
	ctx := middleware.WrapCtx(c)

	uid, ok := getUIDFromClaims(c)
	if !ok || uid == 0 {
		response.Out(c, ecode.Unauthorized)
		return
	}

	form := &types.CreateLoanRepaymentTransactionsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	otpCode := strings.TrimSpace(form.MfaCode)
	ok, err = tool.ValidateMFA(c, uid, otpCode)
	if err != nil || !ok {
		logger.Warn("ValidateMFA error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return
	}

	loanRepaymentTransactions := &model.LoanRepaymentTransactions{}
	loanRepaymentTransactions.CollectOrderNo = generateOrderNo("PI")
	loanRepaymentTransactions.CreatedBy = uid
	loanRepaymentTransactions.PayMethod = "IMPORT"
	
	err = copier.Copy(loanRepaymentTransactions, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanRepaymentTransactions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	err = h.iDao.Create(ctx, loanRepaymentTransactions)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanRepaymentTransactions.ID})
}

// DeleteByID delete a loanRepaymentTransactions by id
// @Summary Delete a loanRepaymentTransactions by id
// @Description Deletes a existing loanRepaymentTransactions identified by the given id in the path.
// @Tags loanRepaymentTransactions
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanRepaymentTransactionsByIDReply{}
// @Router /api/v1/loanRepaymentTransactions/{id} [delete]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanRepaymentTransactionsIDFromPath(c)
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

// UpdateByID update a loanRepaymentTransactions by id
// @Summary Update a loanRepaymentTransactions by id
// @Description Updates the specified loanRepaymentTransactions by given id in the path, support partial update.
// @Tags loanRepaymentTransactions
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanRepaymentTransactionsByIDRequest true "loanRepaymentTransactions information"
// @Success 200 {object} types.UpdateLoanRepaymentTransactionsByIDReply{}
// @Router /api/v1/loanRepaymentTransactions/{id} [put]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanRepaymentTransactionsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanRepaymentTransactionsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanRepaymentTransactions := &model.LoanRepaymentTransactions{}
	err = copier.Copy(loanRepaymentTransactions, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanRepaymentTransactions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanRepaymentTransactions)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanRepaymentTransactions by id
// @Summary Get a loanRepaymentTransactions by id
// @Description Gets detailed information of a loanRepaymentTransactions specified by the given id in the path.
// @Tags loanRepaymentTransactions
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanRepaymentTransactionsByIDReply{}
// @Router /api/v1/loanRepaymentTransactions/{id} [get]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanRepaymentTransactionsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRepaymentTransactions, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanRepaymentTransactionsObjDetail{}
	err = copier.Copy(data, loanRepaymentTransactions)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanRepaymentTransactions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanRepaymentTransactions": data})
}

// List get a paginated list of loanRepaymentTransactionss by custom conditions
// @Summary Get a paginated list of loanRepaymentTransactionss by custom conditions
// @Description Returns a paginated list of loanRepaymentTransactions based on query filters, including page number and size.
// @Tags loanRepaymentTransactions
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanRepaymentTransactionssReply{}
// @Router /api/v1/loanRepaymentTransactions/list [post]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) List(c *gin.Context) {
	form := &types.ListLoanRepaymentTransactionssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRepaymentTransactionss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanRepaymentTransactionss(loanRepaymentTransactionss)
	if err != nil {
		response.Error(c, ecode.ErrListLoanRepaymentTransactions)
		return
	}

	response.Success(c, gin.H{
		"loanRepaymentTransactionss": data,
		"total":                      total,
	})
}

// DeleteByIDs batch delete loanRepaymentTransactions by ids
// @Summary Batch delete loanRepaymentTransactions by ids
// @Description Deletes multiple loanRepaymentTransactions by a list of id
// @Tags loanRepaymentTransactions
// @Param data body types.DeleteLoanRepaymentTransactionssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanRepaymentTransactionssByIDsReply{}
// @Router /api/v1/loanRepaymentTransactions/delete/ids [post]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanRepaymentTransactionssByIDsRequest{}
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

// GetByCondition get a loanRepaymentTransactions by custom condition
// @Summary Get a loanRepaymentTransactions by custom condition
// @Description Returns a single loanRepaymentTransactions that matches the specified filter conditions.
// @Tags loanRepaymentTransactions
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanRepaymentTransactionsByConditionReply{}
// @Router /api/v1/loanRepaymentTransactions/condition [post]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanRepaymentTransactionsByConditionRequest{}
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
	loanRepaymentTransactions, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanRepaymentTransactionsObjDetail{}
	err = copier.Copy(data, loanRepaymentTransactions)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanRepaymentTransactions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanRepaymentTransactions": data})
}

// ListByIDs batch get loanRepaymentTransactions by ids
// @Summary Batch get loanRepaymentTransactions by ids
// @Description Returns a list of loanRepaymentTransactions that match the list of id.
// @Tags loanRepaymentTransactions
// @Param data body types.ListLoanRepaymentTransactionssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanRepaymentTransactionssByIDsReply{}
// @Router /api/v1/loanRepaymentTransactions/list/ids [post]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanRepaymentTransactionssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRepaymentTransactionsMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanRepaymentTransactionss := []*types.LoanRepaymentTransactionsObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanRepaymentTransactionsMap[id]; ok {
			record, err := convertLoanRepaymentTransactions(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanRepaymentTransactions)
				return
			}
			loanRepaymentTransactionss = append(loanRepaymentTransactionss, record)
		}
	}

	response.Success(c, gin.H{
		"loanRepaymentTransactionss": loanRepaymentTransactionss,
	})
}

// ListByLastID get a paginated list of loanRepaymentTransactionss by last id
// @Summary Get a paginated list of loanRepaymentTransactionss by last id
// @Description Returns a paginated list of loanRepaymentTransactionss starting after a given last id, useful for cursor-based pagination.
// @Tags loanRepaymentTransactions
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanRepaymentTransactionssReply{}
// @Router /api/v1/loanRepaymentTransactions/list [get]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) ListByLastID(c *gin.Context) {
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
	loanRepaymentTransactionss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanRepaymentTransactionss(loanRepaymentTransactionss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanRepaymentTransactions)
		return
	}

	response.Success(c, gin.H{
		"loanRepaymentTransactionss": data,
	})
}

func getLoanRepaymentTransactionsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanRepaymentTransactions(loanRepaymentTransactions *model.LoanRepaymentTransactions) (*types.LoanRepaymentTransactionsObjDetail, error) {
	data := &types.LoanRepaymentTransactionsObjDetail{}
	err := copier.Copy(data, loanRepaymentTransactions)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanRepaymentTransactionss(fromValues []*model.LoanRepaymentTransactions) ([]*types.LoanRepaymentTransactionsObjDetail, error) {
	toValues := []*types.LoanRepaymentTransactionsObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanRepaymentTransactions(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
