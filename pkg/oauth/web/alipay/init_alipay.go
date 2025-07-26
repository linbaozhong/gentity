package alipay

import (
	"github.com/linbaozhong/alipay/v3"
	"github.com/linbaozhong/gentity/pkg/oauth/web"
)

type ali struct {
	appid           string
	appSecret       string
	appPrivateKey   string
	alipayPublicKey string
	returnUrl       string
	notifyUrl       string
}

var (
	aliClient *alipay.Client
)

type option func(a *ali)

// WithAppId 支付宝 appid
func WithAppId(appid string) option {
	return func(a *ali) {
		a.appid = appid
	}
}

// WithAppSecret 支付宝内容加密秘钥
func WithAppSecret(appSecret string) option {
	return func(a *ali) {
		a.appSecret = appSecret
	}
}

// WithAppPrivateKey 支付宝应用私钥
func WithAppPrivateKey(privateKey string) option {
	return func(a *ali) {
		a.appPrivateKey = privateKey
	}
}

// WithAlipayPublicKey 支付宝公钥
func WithAlipayPublicKey(publicKey string) option {
	return func(a *ali) {
		a.alipayPublicKey = publicKey
	}
}

// WithReturnUrl 支付宝授权回调地址
func WithReturnUrl(returnUrl string) option {
	return func(a *ali) {
		a.returnUrl = returnUrl
	}
}

// WithNotifyUrl 支付宝异步通知地址
func WithNotifyUrl(notifyUrl string) option {
	return func(a *ali) {
		a.notifyUrl = notifyUrl
	}
}

func New(opts ...option) web.Platformer {
	a := &ali{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *ali) client() *alipay.Client {
	if aliClient == nil {
		var e error
		aliClient, e = alipay.New(a.appid, a.appPrivateKey, true)
		if e != nil {
			panic(e)
		}
		aliClient.LoadAliPayPublicKey(a.alipayPublicKey)
		aliClient.SetEncryptKey(a.appSecret)
	}

	return aliClient
}
