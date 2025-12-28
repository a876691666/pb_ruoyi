package tools

import (
	"image/color"

	"github.com/mojocn/base64Captcha"
	"github.com/pocketbase/pocketbase/core"
)

var white = &color.RGBA{255, 255, 255, 255}

// GenerateBase64Captcha 使用 base64Captcha 生成图片验证码，返回验证码ID、base64图片与答案。
//
// 参数：
// - length: 验证码位数，<=0 则默认 6 位
// - width, height: 图片宽高，<=0 时使用默认 120x40
//
// 返回值：
// - id: 验证码ID（用于后续验证时与用户输入关联）
// - b64img: 形如 "data:image/png;base64,xxx" 的图片数据
// - answer: 验证码正确答案（服务端可用于调试或直接校验；如不需要可忽略）
// - err: 错误信息
func GenerateBase64Captcha(length, width, height int) (id string, b64img string, answer string, err error) {
	if length <= 0 {
		length = 4
	}
	if width <= 0 {
		width = 150
	}
	if height <= 0 {
		height = 60
	}

	// 使用数字型验证码驱动（还可根据需要换成 DriverString/DriverMath 等）
	driver := base64Captcha.NewDriverMath(height, width, 20, 2, white, nil, []string{})
	c := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)

	id, b64img, answer, err = c.Generate()
	if err != nil {
		return "", "", "", err
	}
	return
}

// VerifyCaptcha 从请求中提取 uuid 与 code 并校验图片验证码。
// 支持以下来源（依次优先）：
// - 请求体（application/json、x-www-form-urlencoded、multipart/form-data）的字段：uuid、code
// - URL 查询参数：uuid、code
// - 请求头：X-Captcha-Id、X-Captcha-Code
// 返回：true 表示校验通过，false 表示失败或参数缺失。
func VerifyCaptcha(r *core.RequestEvent) bool {
	if r == nil {
		return false
	}

	// 1) 优先尝试从请求体解析
	var payload struct {
		UUID string `json:"uuid" form:"uuid"`
		Code string `json:"code" form:"code"`
	}
	if v, _, err := ParseBody[struct {
		UUID string `json:"uuid" form:"uuid"`
		Code string `json:"code" form:"code"`
	}](r.Request); err == nil {
		payload = v
	}

	uuid := payload.UUID
	code := payload.Code

	// 2) 其次从查询参数获取
	if uuid == "" || code == "" {
		q := r.Request.URL.Query()
		if uuid == "" {
			uuid = q.Get("uuid")
		}
		if code == "" {
			code = q.Get("code")
		}
	}

	// 3) 再次从请求头获取
	if uuid == "" {
		uuid = r.Request.Header.Get("X-Captcha-Id")
	}
	if code == "" {
		code = r.Request.Header.Get("X-Captcha-Code")
	}

	if uuid == "" || code == "" {
		return false
	}

	// 使用 DefaultMemStore 进行校验（第三个参数为 true 表示校验后清除）
	return base64Captcha.DefaultMemStore.Verify(uuid, code, true)
}
