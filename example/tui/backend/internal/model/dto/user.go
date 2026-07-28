package dto

import "github.com/linbaozhong/gentity/pkg/types"

// BillReq
// @request
type BillReq struct {
	Page *int `json:"page" valid:"min(1)"`
	Size *int `json:"size" valid:"min(1)"`
}

// BillResp
// @response
type BillResp struct {
	Total TotalResp  `json:"total"`
	List  []FlowResp `json:"list"`
}

type ProfileReq struct {
}

// ProfileResp
// @response
type ProfileResp struct {
	ID     *types.BigInt `json:"id"`
	Nick   *types.String `json:"nick"`
	Gender *types.String `json:"gender"`
	Avatar *types.String `json:"avatar"`
}

// TotalReq
// @request
type TotalReq struct {
}

// TotalResp
// @response
type TotalResp struct {
	ID           *types.BigInt `json:"id"`
	User         *types.BigInt `json:"user"`
	Currency     *types.String `json:"currency"`
	Income       *types.Money  `json:"income"`
	FreezeIncome *types.Money  `json:"freeze_income"`
	TotalIncome  *types.Money  `json:"total_income"`
	Amount       *types.Money  `json:"amount"`
	Freeze       *types.Money  `json:"freeze"`
	TotalAmount  *types.Money  `json:"total_amount"`
}

// FlowReq
// @request
type FlowReq struct {
}

// FlowResp
// @response
type FlowResp struct {
	ID              *types.BigInt `json:"id"`
	Type            *types.String `json:"type"`
	TradeType       *types.Int8   `json:"trade_type"`
	TradeTypeString *types.String `json:"trade_type_string"`
	Currency        *types.String `json:"currency"`
	IncomeBefore    *types.Money  `json:"income_before"`
	Income          *types.Money  `json:"income"`
	IncomeAfter     *types.Money  `json:"income_after"`
	Fee             *types.Money  `json:"fee"`
	AmountBefore    *types.Money  `json:"amount_before"`
	Amount          *types.Money  `json:"amount"`
	AmountAfter     *types.Money  `json:"amount_after"`
	Descr           *types.String `json:"descr"`
	Ctime           *types.Time   `json:"ctime"`
}

// WithdrawReq
// @request
type WithdrawReq struct {
	BankCode *string      `json:"bank_code" valid:"required,stringlength(1|10)"` // 银行代码
	BankName *string      `json:"bank_name" valid:"required,stringlength(1|20)"` // 银行名称
	Name     *string      `json:"name" valid:"required,stringlength(1|45)"`      // 户名
	Number   *string      `json:"number" valid:"required,stringlength(1|20)"`    // 账号
	Amount   *types.Money `json:"amount" valid:"required,min(0.01)"`             // 提现金额
}

// WithdrawResp
// @response
type WithdrawResp struct {
}

// BankResp
// @response
type BankResp struct {
	BankCode *types.String `json:"bank_code"` // 银行代码
	BankName *types.String `json:"bank_name"` // 银行名称
	Name     *types.String `json:"name"`      // 户名
	Number   *types.String `json:"number"`    // 账号
}
