// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package tblcompanyfadada

import (
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
)

var (
	Id                = dialect.Field{Name: "id", Json: "id", Table: "company_fadada", Type: "types.BigInt"}
	CompanyName       = dialect.Field{Name: "company_name", Json: "company_name", Table: "company_fadada", Type: "types.String"}
	CustomerId        = dialect.Field{Name: "customer_id", Json: "customer_id", Table: "company_fadada", Type: "types.String"}
	TransactionNo     = dialect.Field{Name: "transaction_no", Json: "transaction_no", Table: "company_fadada", Type: "types.String"}
	Url               = dialect.Field{Name: "url", Json: "url", Table: "company_fadada", Type: "types.String"}
	CertInfo          = dialect.Field{Name: "cert_info", Json: "cert_info", Table: "company_fadada", Type: "types.String"}
	Status            = dialect.Field{Name: "status", Json: "status", Table: "company_fadada", Type: "types.Int8"}
	HasCertificate    = dialect.Field{Name: "has_certificate", Json: "has_certificate", Table: "company_fadada", Type: "types.Int8"}
	AuthSign          = dialect.Field{Name: "auth_sign", Json: "auth_sign", Table: "company_fadada", Type: "types.Int8"}
	AuthTransactionId = dialect.Field{Name: "auth_transaction_id", Json: "auth_transaction_id", Table: "company_fadada", Type: "types.String"}
	AuthContractId    = dialect.Field{Name: "auth_contract_id", Json: "auth_contract_id", Table: "company_fadada", Type: "types.String"}
	AuthResult        = dialect.Field{Name: "auth_result", Json: "auth_result", Table: "company_fadada", Type: "types.String"}
	State             = dialect.Field{Name: "state", Json: "state", Table: "company_fadada", Type: "types.Int8"}
	Ctime             = dialect.Field{Name: "ctime", Json: "ctime", Table: "company_fadada", Type: "types.Time"}
	// 主键
	PrimaryKey = Id

	// 可写列
	WritableFields = []dialect.Field{
		Id,
		CompanyName,
		CustomerId,
		TransactionNo,
		Url,
		CertInfo,
		Status,
		HasCertificate,
		AuthSign,
		AuthTransactionId,
		AuthContractId,
		AuthResult,
		State,
		Ctime,
	}
	// 可读列
	ReadableFields = []dialect.Field{
		Id,
		CompanyName,
		CustomerId,
		TransactionNo,
		Url,
		CertInfo,
		Status,
		HasCertificate,
		AuthSign,
		AuthTransactionId,
		AuthContractId,
		AuthResult,
		State,
		Ctime,
	}
)
