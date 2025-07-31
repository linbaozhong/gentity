package wechat

import (
	"context"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/oauth/web"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	wxoption "github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

type PaymentType string

const (
	PaymentTypeJsapi  PaymentType = "jsapi"
	PaymentTypeH5     PaymentType = "h5"
	PaymentTypeNative PaymentType = "native"
)

type wx struct {
	appid       string
	appSecret   string
	lang        string
	redirectURI string // 授权回调地址
	// 支付相关
	mchID                      string      // 商户号
	mchCertificateSerialNumber string      // 商户证书序列号
	mchPrivateKey              string      // 商户私钥
	mchAPIv3Key                string      // 商户APIv3密钥
	mchPaymentType             PaymentType // 支付类型
	// 通知相关
	notifyURL string // 通知地址

}

var (
	wxClient *core.Client
	_        web.Platformer = &wx{}
)

type option func(w *wx)

// 微信 API 地址
const (
	wechatAuthURL     = "https://open.weixin.qq.com/connect/oauth2/authorize"
	wechatTokenURL    = "https://api.weixin.qq.com/sns/oauth2/access_token"
	wechatUserInfoURL = "https://api.weixin.qq.com/sns/userinfo"
)

// WithAppId 设置appid
func WithAppId(appid string) option {
	return func(w *wx) {
		w.appid = appid
	}
}

// WithLang 设置语言
func WithLang(lang string) option {
	return func(w *wx) {
		w.lang = lang
	}
}

// WithAppSecret 设置appSecret
func WithAppSecret(appSecret string) option {
	return func(w *wx) {
		w.appSecret = appSecret
	}
}

// WithRedirectURI 设置授权回调地址
func WithRedirectURI(redirectURI string) option {
	return func(w *wx) {
		w.redirectURI = redirectURI
	}
}

// WithMchId 设置商户号
func WithMchId(mchID string) option {
	return func(w *wx) {
		w.mchID = mchID
	}
}

// WithMchCertificateSerialNumber 设置商户证书序列号
func WithMchCertificateSerialNumber(mchCertificateSerialNumber string) option {
	return func(w *wx) {
		w.mchCertificateSerialNumber = mchCertificateSerialNumber
	}
}

// WithMchPrivateKey 设置商户私钥
func WithMchPrivateKey(mchPrivateKey string) option {
	return func(w *wx) {
		w.mchPrivateKey = mchPrivateKey
	}
}

// WithMchAPIv3Key 设置商户APIv3密钥
func WithMchApiv3Key(mchAPIv3Key string) option {
	return func(w *wx) {
		w.mchAPIv3Key = mchAPIv3Key
	}
}

// WithMchPaymentType 设置商户支付类型
func WithMchPaymentType(mchPaymentType PaymentType) option {
	return func(w *wx) {
		w.mchPaymentType = mchPaymentType
	}
}

// WithMchNotifyUrl 设置通知地址
func WithMchNotifyUrl(notifyURL string) option {
	return func(w *wx) {
		w.notifyURL = notifyURL
	}
}

func New(opts ...option) web.Platformer {
	w := &wx{}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (w *wx) client() *core.Client {
	if wxClient == nil {
		var err error
		// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
		mchPrivateKey, err := utils.LoadPrivateKeyWithPath(w.mchPrivateKey)
		if err != nil {
			panic(fmt.Sprintf("load merchant private key error: %v", err))
			return nil
		}
		_opts := []core.ClientOption{
			wxoption.WithWechatPayAutoAuthCipher(w.mchID, w.mchCertificateSerialNumber, mchPrivateKey, w.mchAPIv3Key),
		}
		wxClient, err = core.NewClient(context.Background(), _opts...)
		if err != nil {
			panic(fmt.Sprintf("new wechat client error: %v", err))
			return nil
		}
	}
	return wxClient
}
