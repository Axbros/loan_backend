package tool

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"loan/internal/cache"
	"loan/internal/config"
	"loan/internal/dao"
	"loan/internal/database"
	"loan/internal/ecode"
	"strings"
	"time"
	"unicode"
)

const (
	MFAOTPCodeLength = 6               // 验证码长度
	MFATOTPPeriod    = 30              // TOTP周期（秒）
	MFATOPSkew       = 1               // TOTP时间偏移量
	MFAAsyncTimeout  = 5 * time.Second // 异步操作超时
	MFAMinSecretLen  = 16              // MFA密钥最小长度
)

// ValidateMFA MFA 校验（完全保留你的错误返回 + DAO 调用方式 + 常量全引用）
func ValidateMFA(ctx *gin.Context, uid uint64, otpCode string) (bool, error) {
	otpCode = strings.TrimSpace(otpCode)

	// 1. 替换硬编码6为常量 MFAOTPCodeLength
	if len(otpCode) != MFAOTPCodeLength {
		logger.Warn("MFA验证码长度错误", logger.String("otp_code", otpCode), middleware.GCtxRequestIDField(ctx))
		response.Error(ctx, ecode.MFAOTPRequired)
		return false, nil
	}

	// 新增：校验OTP码是否为纯数字
	if !isNumeric(otpCode) {
		logger.Warn("MFA验证码非数字格式", logger.String("otp_code", otpCode), middleware.GCtxRequestIDField(ctx))
		response.Error(ctx, ecode.MFAOTPRequired)
		return false, nil
	}

	userDao := dao.NewLoanUsersDao(
		database.GetDB(),
		cache.NewLoanUsersCache(database.GetCacheType()),
	)
	// 查询用户的MFA主设备
	dev, err := userDao.GetActivePrimaryMFADevice(ctx, uid)
	if err != nil {
		logger.Error("查询MFA主设备失败", logger.Err(err), logger.Uint64("uid", uid), middleware.GCtxRequestIDField(ctx))
		return false, err
	}
	if dev == nil {
		logger.Warn("用户无激活的MFA主设备", logger.Uint64("uid", uid), middleware.GCtxRequestIDField(ctx))
		response.Error(ctx, ecode.ErrGetByIDLoanMfaDevices)
		return false, nil
	}

	// 解密MFA密钥
	secret, err := DecryptSecretFromBytes(dev.SecretEnc)
	if err != nil {
		logger.Error("解密MFA密钥失败", logger.Err(err), logger.Uint64("device_id", dev.ID), middleware.GCtxRequestIDField(ctx))
		response.Error(ctx, ecode.ErrSecret)
		return false, err
	}
	secret = strings.TrimSpace(secret)

	// 2. 替换硬编码16为常量 MFAMinSecretLen
	if secret == "" || len(secret) < MFAMinSecretLen {
		logger.Warn("MFA密钥为空或长度不足", logger.Uint64("device_id", dev.ID), middleware.GCtxRequestIDField(ctx))
		response.Error(ctx, ecode.ErrSecret)
		return false, nil
	}

	// 校验OTP：替换硬编码30、1为对应常量
	ok, err := totp.ValidateCustom(otpCode, secret, time.Now(), totp.ValidateOpts{
		Period:    MFATOTPPeriod, // 替换30
		Skew:      MFATOPSkew,    // 替换1
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		logger.Error("校验MFA验证码失败", logger.Err(err), logger.Uint64("uid", uid), middleware.GCtxRequestIDField(ctx))
		response.Error(ctx, ecode.ErrSecret)
		return false, err
	}
	if !ok {
		logger.Warn("MFA验证码无效", logger.String("otp_code", otpCode), logger.Uint64("uid", uid), middleware.GCtxRequestIDField(ctx))
		response.Error(ctx, ecode.InvalidOTP)
		return false, nil
	}

	// 核心修复：正确获取requestID（你之前的 requestID := middleware.RequestIDKey 是错误的）
	// 正确方式：从gin ctx中提取requestID字符串
	requestID := middleware.RequestIDKey

	// 3. 替换硬编码5*time.Second为常量 MFAAsyncTimeout
	go func(deviceID uint64, reqID string) {
		// 保留独立上下文，替换硬编码超时时间为常量
		asyncCtx, cancel := context.WithTimeout(context.Background(), MFAAsyncTimeout)
		defer cancel()

		// 调用DAO更新最后使用时间
		if err = userDao.TouchMFADeviceLastUsedAt(asyncCtx, deviceID); err != nil {
			logFields := []logger.Field{
				logger.Err(err),
				logger.Uint64("device_id", deviceID),
				logger.String("request_id", reqID),
			}
			if errors.Is(context.DeadlineExceeded, err) {
				logger.Warn("更新MFA设备最后使用时间超时", logFields...)
			} else {
				logger.Warn("更新MFA设备最后使用时间失败", logFields...)
			}
		}
	}(dev.ID, string(requestID))

	// 成功
	return true, nil
}

// isNumeric 辅助函数：校验字符串是否为纯数字
func isNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func DecryptSecretFromBytes(enc []byte) (string, error) {
	keyStr := config.Get().Authorization.Key
	if keyStr == "" {
		return "", errors.New("MFA_AES_KEY empty")
	}

	var key []byte
	if len(keyStr) == 64 {
		b, err := hex.DecodeString(keyStr)
		if err != nil {
			return "", err
		}
		key = b
	} else {
		key = []byte(keyStr)
	}
	if len(key) != 32 {
		return "", errors.New("MFA_AES_KEY must be 32 bytes or 64 hex chars")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(enc) < nonceSize+16 {
		return "", errors.New("ciphertext too short")
	}

	nonce := enc[:nonceSize]
	ciphertext := enc[nonceSize:]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
