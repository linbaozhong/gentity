package alipay

import (
	"context"
	"github.com/linbaozhong/alipay/v3"
	"github.com/linbaozhong/gentity/pkg/oauth/web"
	"github.com/linbaozhong/gentity/pkg/types"
	"net/http"
	"net/url"
)

func (a *ali) PagePay(ctx context.Context, req *web.PagePayReq) (types.Smap, error) {
	_trade := alipay.TradePagePay{}
	_trade.OutTradeNo = req.Bill.String()
	_trade.TotalAmount = req.Amount.Yuan().String()
	_trade.ProductCode = "FAST_INSTANT_TRADE_PAY"
	_trade.Subject = req.Name.String()
	_trade.Body = req.Desc.String()
	_trade.NotifyURL = a.notifyUrl
	_trade.ReturnURL = a.returnUrl
	_trade.TimeoutExpress = "30m"
	// _trade.QRPayMode = "4" // 前置模式
	// _trade.QRCodeWidth = "120"
	_trade.QRPayMode = "2" // 跳转模式
	_trade.PassbackParams = url.QueryEscape(
		req.Sku.String() + web.Passbackchar +
			req.Sharer.String() + web.Passbackchar +
			req.Buyer.String() + web.Passbackchar +
			req.Seller.String())
	_url, err := a.client().TradePagePay(_trade)
	if err != nil {
		return nil, err
	}
	return types.NewSmap().Set("url", _url.String()), nil
}

func (a *ali) Notify(ctx context.Context, req *http.Request) (*web.NotifyResp, error) {
	return nil, nil
}
