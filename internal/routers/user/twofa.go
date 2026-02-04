package user

import (
	"bytes"
	"encoding/base64"
	"image/png"

	"github.com/gin-gonic/gin"
	"github.com/gocronx-team/gocron/internal/models"
	"github.com/gocronx-team/gocron/internal/modules/i18n"
	"github.com/gocronx-team/gocron/internal/modules/logger"
	"github.com/gocronx-team/gocron/internal/routers/base"
	"github.com/pquerna/otp/totp"
)

// Setup2FA 设置2FA
func Setup2FA(c *gin.Context) {
	uid := Uid(c)
	username := Username(c)

	userModel := new(models.User)
	err := userModel.Find(uid)
	if err != nil || userModel.Id == 0 {
		base.RespondError(c, i18n.T(c, "user_not_found"))
		return
	}

	// 生成TOTP密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Gocron",
		AccountName: username,
	})
	if err != nil {
		logger.Error("生成2FA密钥失败", err)
		base.RespondError(c, i18n.T(c, "generate_2fa_key_failed"))
		return
	}

	// 生成二维码
	img, err := key.Image(200, 200)
	if err != nil {
		logger.Error("生成二维码失败", err)
		base.RespondError(c, i18n.T(c, "generate_qrcode_failed"))
		return
	}

	// 将图片转为base64
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		logger.Error("编码二维码失败", err)
		base.RespondError(c, i18n.T(c, "generate_qrcode_failed"))
		return
	}
	qrCode := base64.StdEncoding.EncodeToString(buf.Bytes())

	base.RespondSuccess(c, i18n.T(c, "get_success"), map[string]interface{}{
		"secret":  key.Secret(),
		"qr_code": "data:image/png;base64," + qrCode,
	})
}

// Enable2FAForm 启用2FA表单
type Enable2FAForm struct {
	Secret string `form:"secret" json:"secret" binding:"required"`
	Code   string `form:"code" json:"code" binding:"required,len=6"`
}

// Enable2FA 启用2FA
func Enable2FA(c *gin.Context) {
	var form Enable2FAForm
	if err := c.ShouldBind(&form); err != nil {
		base.RespondValidationError(c, err)
		return
	}

	uid := Uid(c)

	// 验证TOTP码
	valid := totp.Validate(form.Code, form.Secret)
	if !valid {
		base.RespondError(c, i18n.T(c, "verification_code_error"))
		return
	}

	// 保存密钥并启用2FA
	userModel := new(models.User)
	_, err := userModel.Update(uid, models.CommonMap{
		"two_factor_key": form.Secret,
		"two_factor_on":  1,
	})
	if err != nil {
		base.RespondError(c, i18n.T(c, "enable_failed"), err)
		return
	}

	base.RespondSuccess(c, i18n.T(c, "2fa_enabled"), nil)
}

// Disable2FAForm 禁用2FA表单
type Disable2FAForm struct {
	Code string `form:"code" json:"code" binding:"required,len=6"`
}

// Disable2FA 禁用2FA
func Disable2FA(c *gin.Context) {
	var form Disable2FAForm
	if err := c.ShouldBind(&form); err != nil {
		base.RespondValidationError(c, err)
		return
	}

	uid := Uid(c)

	userModel := new(models.User)
	err := userModel.Find(uid)
	if err != nil || userModel.Id == 0 {
		base.RespondError(c, i18n.T(c, "user_not_found"))
		return
	}

	if userModel.TwoFactorOn == 0 {
		base.RespondError(c, i18n.T(c, "2fa_not_enabled"))
		return
	}

	// 验证TOTP码
	valid := totp.Validate(form.Code, userModel.TwoFactorKey)
	if !valid {
		base.RespondError(c, i18n.T(c, "verification_code_error"))
		return
	}

	// 禁用2FA
	_, err = userModel.Update(uid, models.CommonMap{
		"two_factor_key": "",
		"two_factor_on":  0,
	})
	if err != nil {
		base.RespondError(c, i18n.T(c, "disable_failed"), err)
		return
	}

	base.RespondSuccess(c, i18n.T(c, "2fa_disabled"), nil)
}

// Get2FAStatus 获取2FA状态
func Get2FAStatus(c *gin.Context) {
	uid := Uid(c)

	userModel := new(models.User)
	err := userModel.Find(uid)
	if err != nil || userModel.Id == 0 {
		base.RespondError(c, i18n.T(c, "user_not_found"))
		return
	}

	base.RespondSuccess(c, i18n.T(c, "get_success"), map[string]interface{}{
		"enabled": userModel.TwoFactorOn == 1,
	})
}
