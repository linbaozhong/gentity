package wechat

import (
	"context"
	"errors"
	"github.com/linbaozhong/gentity/pkg/oauth/web"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"time"
)

func (w *wx) PagePay(ctx context.Context, req *web.PagePayReq) (string, error) {
	_trade := h5.PrepayRequest{}
	_trade.OutTradeNo = core.String(req.Bill.String())
	_trade.Amount = &h5.Amount{
		Total:    core.Int64(req.Amount.Int64()),
		Currency: core.String(req.Currency),
	}
	_trade.Description = core.String(req.Desc)
	_trade.NotifyUrl = core.String(req.NotifyUrl)
	_trade.TimeExpire = core.Time(time.Now().Add(time.Minute * 30))
	_trade.Attach = core.String(
		req.Sku.String() + web.Passbackchar +
			req.Sharer.String() + web.Passbackchar +
			req.Buyer.String() + web.Passbackchar +
			req.Seller.String())

	resp, result, err := w.client().Prepay(ctx, _trade)
	if err != nil {
		return "", err
	}
	if result.Response.StatusCode != 200 {
		return "", errors.New(result.Response.Status)
	}
	return *resp.H5Url, nil
}
