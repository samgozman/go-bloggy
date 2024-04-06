package captcha

import "github.com/kataras/hcaptcha"

// Client is a type alias for the HCaptcha client.
type Client = hcaptcha.Client

type ClientInterface interface {
	VerifyToken(tkn string) (response hcaptcha.Response)
}
