package wechat

import (
	"context"
	"errors"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/oauth/web"
	"github.com/linbaozhong/gentity/pkg/types"
	"github.com/linbaozhong/gentity/pkg/util"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"io"
	"net/http"
	"time"
)

// PagePay 统一下单支付
func (w *wx) PagePay(ctx context.Context, req *web.PagePayReq) (types.Smap, error) {
	switch w.mchPaymentType {
	case PaymentTypeH5:
		return w.h5(ctx, req)
	case PaymentTypeNative:
		return w.native(ctx, req)
	case PaymentTypeJsapi:
		return w.jsapi(ctx, req)
	}
	return nil, errors.New("invalid payment type")
}

// h5 浏览器打开支付
func (w *wx) h5(ctx context.Context, req *web.PagePayReq) (types.Smap, error) {
	_trade := h5.PrepayRequest{}
	_trade.Appid = core.String(w.appid)
	_trade.Mchid = core.String(w.mchID)
	_trade.NotifyUrl = core.String(w.notifyURL)
	_trade.OutTradeNo = core.String(req.Bill.String())
	_trade.Amount = &h5.Amount{
		Total:    core.Int64(req.Amount.Int64()),
		Currency: core.String(req.Currency.String()),
	}
	_trade.Description = core.String(req.Desc.String())
	_trade.NotifyUrl = core.String(w.notifyURL)
	_trade.TimeExpire = core.Time(time.Now().Add(time.Minute * 30))
	_trade.Attach = core.String(
		req.Sku.String() + web.Passbackchar +
			req.Sharer.String() + web.Passbackchar +
			req.Buyer.String() + web.Passbackchar +
			req.Seller.String())

	_cli := &h5.H5ApiService{Client: w.client()}
	resp, result, err := _cli.Prepay(ctx, _trade)
	if err != nil {
		return nil, err
	}
	if result.Response.StatusCode != 200 {
		return nil, errors.New(result.Response.Status)
	}
	return types.NewSmap().Set("url", *resp.H5Url), nil
}

// native 扫码支付
func (w *wx) native(ctx context.Context, req *web.PagePayReq) (types.Smap, error) {
	_trade := native.PrepayRequest{}
	_trade.Appid = core.String(w.appid)
	_trade.Mchid = core.String(w.mchID)
	_trade.NotifyUrl = core.String(w.notifyURL)
	_trade.OutTradeNo = core.String(req.Bill.String())
	_trade.Amount = &native.Amount{
		Total:    core.Int64(req.Amount.Int64()),
		Currency: core.String(req.Currency.String()),
	}
	_trade.Description = core.String(req.Desc.String())
	_trade.NotifyUrl = core.String(w.notifyURL)
	_trade.TimeExpire = core.Time(time.Now().Add(time.Minute * 30))
	_trade.Attach = core.String(
		req.Sku.String() + web.Passbackchar +
			req.Sharer.String() + web.Passbackchar +
			req.Buyer.String() + web.Passbackchar +
			req.Seller.String())

	_cli := &native.NativeApiService{Client: w.client()}
	resp, result, err := _cli.Prepay(ctx, _trade)
	if err != nil {
		return nil, err
	}
	if result.Response.StatusCode != 200 {
		return nil, errors.New(result.Response.Status)
	}
	return types.NewSmap().Set("url", *resp.CodeUrl), nil
}

// jsapi 浏览器打开支付
func (w *wx) jsapi(ctx context.Context, req *web.PagePayReq) (types.Smap, error) {
	_trade := jsapi.PrepayRequest{}
	_trade.Appid = core.String(w.appid)
	_trade.Mchid = core.String(w.mchID)
	_trade.NotifyUrl = core.String(w.notifyURL)
	_trade.OutTradeNo = core.String(req.Bill.String())
	_trade.Amount = &jsapi.Amount{
		Total:    core.Int64(req.Amount.Int64()),
		Currency: core.String(req.Currency.String()),
	}
	_trade.Description = core.String(req.Desc.String())
	_trade.NotifyUrl = core.String(w.notifyURL)
	_trade.TimeExpire = core.Time(time.Now().Add(time.Minute * 30))
	_trade.Attach = core.String(
		req.Sku.String() + web.Passbackchar +
			req.Sharer.String() + web.Passbackchar +
			req.Buyer.String() + web.Passbackchar +
			req.Seller.String())
	_cli := &jsapi.JsapiApiService{Client: w.client()}
	resp, result, err := _cli.Prepay(ctx, _trade)
	if err != nil {
		return nil, err
	}
	if result.Response.StatusCode != 200 {
		return nil, errors.New(result.Response.Status)
	}

	// 生成 JSAPI 支付参数
	nonceStr := util.GetRandLowerString(32)
	timeStamp := fmt.Sprintf("%d", time.Now().Unix())
	packageStr := fmt.Sprintf("prepay_id=%s", *resp.PrepayId)
	signType := "RSA"

	signStr := fmt.Sprintf("%s\n%s\n%s\n%s\n", w.appid, timeStamp, nonceStr, packageStr)
	sign, err := _cli.Client.Sign(ctx, signStr)
	if err != nil {
		return nil, err
	}

	return types.NewSmap(6).
		Set("appId", w.appid).
		Set("timeStamp", timeStamp).
		Set("nonceStr", nonceStr).
		Set("package", packageStr).
		Set("signType", signType).
		Set("paySign", sign.Signature), nil
}

// Notify 支付成功的异步通知
func (w *wx) Notify(ctx context.Context, req *http.Request, resp *web.NotifyResp) error {
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(w.mchID)
	_cli, e := notify.NewRSANotifyHandler(w.mchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	if e != nil {
		resp.Code = http.StatusInternalServerError
		return fmt.Errorf("create notify handler failed: %w", e)
	}
	// 解密通知数据
	transaction := new(payments.Transaction)
	notifyReq, err := _cli.ParseNotifyRequest(ctx, req, transaction)
	// 如果验签未通过，或者解密失败
	if err != nil {
		resp.Code = http.StatusBadRequest
		return fmt.Errorf("parse notify request failed: %w", err)
	}
	if notifyReq.EventType != "TRANSACTION.SUCCESS" {
		resp.Code = http.StatusBadRequest
		return fmt.Errorf("The event type is not TRANSACTION.SUCCESS: %s", notifyReq.EventType)
	}
	// 处理通知数据
	resp.Code = http.StatusOK
	resp.Id = notifyReq.ID
	resp.CreateTime = notifyReq.CreateTime.Format(time.RFC3339)
	resp.EventType = notifyReq.EventType
	resp.Summary = notifyReq.Summary
	resp.Transaction = web.Transaction{
		Amount: web.Amount{
			Currency:      *transaction.Amount.Currency,
			PayerCurrency: *transaction.Amount.PayerCurrency,
			PayerTotal:    *transaction.Amount.PayerTotal,
			Total:         *transaction.Amount.Total,
		},
		Appid:          *transaction.Appid,
		Attach:         *transaction.Attach,
		BankType:       *transaction.BankType,
		Mchid:          *transaction.Mchid,
		OutTradeNo:     *transaction.OutTradeNo,
		Payer:          web.Payer{Openid: *transaction.Payer.Openid},
		SuccessTime:    *transaction.SuccessTime,
		TradeState:     *transaction.TradeState,
		TradeStateDesc: *transaction.TradeStateDesc,
		TradeType:      *transaction.TradeType,
		TransactionId:  *transaction.TransactionId,
	}
	_body, _ := io.ReadAll(req.Body)
	resp.Message = string(_body)
	return nil
}

func (w *wx) Refund(ctx context.Context, req *web.RefundReq, resp *web.RefundResp) error {
	_cli := refunddomestic.RefundsApiService{Client: w.client()}
	rsp, result, e := _cli.Create(ctx, refunddomestic.CreateRequest{
		OutTradeNo:  core.String(req.OutTradeNo),
		OutRefundNo: core.String(req.OutRefundNo),
		Reason:      core.String(req.Reason),
		Amount: &refunddomestic.AmountReq{
			Refund: core.Int64(req.Amount.Int64()),
			Total:  core.Int64(req.Amount.Int64()),
		},
		NotifyUrl: core.String(req.NotifyUrl),
	})
	if e != nil {
		return e
	}
	if result.Response.StatusCode != 200 {
		return errors.New(result.Response.Status)
	}
	resp.RefundId = *rsp.RefundId
	resp.OutRefundNo = *rsp.OutRefundNo
	resp.TransactionId = *rsp.TransactionId
	resp.OutTradeNo = *rsp.OutTradeNo
	resp.Channel = string(*rsp.Channel)
	resp.UserReceivedAccount = *rsp.UserReceivedAccount
	resp.SuccessTime = *rsp.SuccessTime
	resp.CreateTime = *rsp.CreateTime
	resp.Status = string(*rsp.Status)
	resp.Amount = types.Money(*rsp.Amount.Refund)
	resp.FundsAccount = string(*rsp.FundsAccount)
	_body, _ := io.ReadAll(result.Response.Body)
	resp.Message = string(_body)
	return nil
}
