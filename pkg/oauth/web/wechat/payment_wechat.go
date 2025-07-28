package wechat

import (
	"context"
	"errors"
	"github.com/linbaozhong/gentity/pkg/oauth/web"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"time"
)

// PagePay 统一下单支付
func (w *wx) PagePay(ctx context.Context, req *web.PagePayReq) (string, error) {
	switch w.mchPaymentType {
	case PaymentTypeH5:
		return w.h5(ctx, req)
	case PaymentTypeNative:
		return w.native(ctx, req)
	}
	return "", errors.New("invalid payment type")
}

// h5 浏览器打开支付
func (w *wx) h5(ctx context.Context, req *web.PagePayReq) (string, error) {
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

	_cli := &h5.H5ApiService{Client: w.client()}
	resp, result, err := _cli.Prepay(ctx, _trade)
	if err != nil {
		return "", err
	}
	if result.Response.StatusCode != 200 {
		return "", errors.New(result.Response.Status)
	}
	return *resp.H5Url, nil
}

// native 扫码支付
func (w *wx) native(ctx context.Context, req *web.PagePayReq) (string, error) {
	_trade := native.PrepayRequest{}
	_trade.OutTradeNo = core.String(req.Bill.String())
	_trade.Amount = &native.Amount{
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

	_cli := &native.NativeApiService{Client: w.client()}
	resp, result, err := _cli.Prepay(ctx, _trade)
	if err != nil {
		return "", err
	}
	if result.Response.StatusCode != 200 {
		return "", errors.New(result.Response.Status)
	}
	return *resp.CodeUrl, nil
}
