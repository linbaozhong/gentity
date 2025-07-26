package wechat

import (
	"context"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/oauth/web"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	wxoption "github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

type wx struct {
	appid                      string
	appSecret                  string
	lang                       string
	redirectURI                string // 回调地址
	mchID                      string
	mchCertificateSerialNumber string
	mchPrivateKey              string
	mchAPIv3Key                string
}
type option func(w *wx)

// 微信 API 地址
const (
	wechatAuthURL     = "https://open.weixin.qq.com/connect/qrconnect"
	wechatTokenURL    = "https://api.weixin.qq.com/sns/oauth2/access_token"
	wechatUserInfoURL = "https://api.weixin.qq.com/sns/userinfo"
)

func WithAppId(appid string) option {
	return func(w *wx) {
		w.appid = appid
	}
}

func WithLang(lang string) option {
	return func(w *wx) {
		w.lang = lang
	}
}

func WithAppSecret(appSecret string) option {
	return func(w *wx) {
		w.appSecret = appSecret
	}
}

func WithRedirectURI(redirectURI string) option {
	return func(w *wx) {
		w.redirectURI = redirectURI
	}
}

func New(opts ...option) web.Platformer {
	w := &wx{}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (w *wx) client() *h5.H5ApiService {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKey(w.mchPrivateKey)
	if err != nil {
		log.Fatal("load merchant private key error")
	}
	_opts := []core.ClientOption{
		wxoption.WithWechatPayAutoAuthCipher(w.mchID, w.mchCertificateSerialNumber, mchPrivateKey, w.mchAPIv3Key),
	}
	_cli, err := core.NewClient(context.Background(), _opts...)
	if err != nil {
		log.Fatal("new wechat client error")
	}
	return &h5.H5ApiService{
		Client: _cli,
	}
}
