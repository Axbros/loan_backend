package tool

import (
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
	"loan/internal/model"
	"strings"
	"time"
)

// MFA 校验
func ValidateMFA(ctx *gin.Context, uid uint64, otpCode string) (bool, *model.LoanMfaDevices, error) {

	otpCode = strings.TrimSpace(otpCode)
	if len(otpCode) != 6 {
		logger.Warn("MFA验证码长度错误", logger.String("otp_code", otpCode), middleware.GCtxRequestIDField(ctx))
		response.Error(ctx, ecode.MFAOTPRequired)
		return false, nil, nil
	}

	userDao := dao.NewLoanUsersDao(
		database.GetDB(),
		cache.NewLoanUsersCache(database.GetCacheType()),
	)
	// 查询用户的MFA主设备
	dev, err := userDao.GetActivePrimaryMFADevice(ctx, uid)
	if err != nil {
		logger.Error("查询MFA主设备失败", logger.Err(err), logger.Uint64("uid", uid), middleware.GCtxRequestIDField(ctx))
		return false, nil, err
	}
	if dev == nil {
		logger.Warn("用户无激活的MFA主设备", logger.Uint64("uid", uid), middleware.GCtxRequestIDField(ctx))
		response.Error(ctx, ecode.ErrGetByIDLoanMfaDevices)
		return false, nil, nil
	}

	// 解密MFA密钥
	secret, err := DecryptSecretFromBytes(dev.SecretEnc)
	if err != nil {
		logger.Error("解密MFA密钥失败", logger.Err(err), logger.Uint64("device_id", dev.ID), middleware.GCtxRequestIDField(ctx))
		response.Error(ctx, ecode.ErrSecret)
		return false, nil, err
	}
	secret = strings.TrimSpace(secret)
	if secret == "" {
		logger.Warn("MFA密钥为空", logger.Uint64("device_id", dev.ID), middleware.GCtxRequestIDField(ctx))
		response.Error(ctx, ecode.ErrSecret)
		return false, nil, nil
	}

	// 校验OTP
	ok, err := totp.ValidateCustom(otpCode, secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		logger.Error("校验MFA验证码失败", logger.Err(err), logger.Uint64("uid", uid), middleware.GCtxRequestIDField(ctx))
		response.Error(ctx, ecode.ErrSecret)
		return false, nil, err
	}
	if !ok {
		logger.Warn("MFA验证码无效", logger.String("otp_code", otpCode), logger.Uint64("uid", uid), middleware.GCtxRequestIDField(ctx))
		response.Error(ctx, ecode.InvalidOTP)
		return false, nil, nil
	}

	// 成功
	return true, dev, nil
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
