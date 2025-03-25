package db

import (
	"github.com/linbaozhong/gentity/pkg/ace/pool"
	"github.com/linbaozhong/gentity/pkg/types"
)

// tablename company
type Company struct {
	pool.Model
	Id               types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`                       //
	LongName         types.String `json:"long_name,omitempty" db:"'long_name' size:100"`                // 全称
	ShortName        types.String `json:"short_name,omitempty" db:"'short_name' size:100"`              // 简称
	Address          types.String `json:"address,omitempty" db:"'address' size:200"`                    // 注册住所
	Email            types.String `json:"email,omitempty" db:"'email' size:10"`                         // 邮政编码
	ContactName      types.String `json:"contact_name,omitempty" db:"'contact_name' size:20"`           // 联系人姓名
	ContactTelephone types.String `json:"contact_telephone,omitempty" db:"'contact_telephone' size:50"` // 座机号码
	ContactMobile    types.String `json:"contact_mobile,omitempty" db:"'contact_mobile' size:50"`       // 联系人手机号码 或者是座机号码
	ContactEmail     types.String `json:"contact_email,omitempty" db:"'contact_email' size:45"`         // 联系人邮箱
	LegalName        types.String `json:"legal_name,omitempty" db:"'legal_name' size:20"`               // 法人姓名
	Creator          types.BigInt `json:"creator,omitempty" db:"'creator' size:20"`                     //
	State            types.Int8   `json:"state,omitempty" db:"'state' size:3"`                          // 管理状态：-1=已删除，0=禁用，1=启用
	Status           types.Int8   `json:"status,omitempty" db:"'status' size:3"`                        // 状态
	Ctime            types.Time   `json:"ctime,omitempty" db:"'ctime'"`                                 //
	Utime            types.Time   `json:"utime,omitempty" db:"'utime'"`                                 //
}

// tablename company_fadada
type CompanyFadada struct {
	pool.Model
	Id                types.BigInt `json:"id,omitempty" db:"'id' pk size:20"`                                 // 企业id
	CompanyName       types.String `json:"company_name,omitempty" db:"'company_name' size:100"`               // 企业名称
	CustomerId        types.String `json:"customer_id,omitempty" db:"'customer_id' size:100"`                 // 客户编号
	TransactionNo     types.String `json:"transaction_no,omitempty" db:"'transaction_no' size:50"`            // 交易号（需要保存，用于证书申请和实名认证查询）
	Url               types.String `json:"url,omitempty" db:"'url' size:65535"`                               // 地址（需要保存，遇到中途退出认证或页面过期等情况可重新访问
	CertInfo          types.String `json:"cert_info,omitempty" db:"'cert_info' size:4294967295"`              // 实名认证信息
	Status            types.Int8   `json:"status,omitempty" db:"'status' size:3"`                             // 0未认证   1管理员资料已提交  2企业基本资料(没有申请表)已提交  3已提交待审核   4审核通过（认证完成）  5审核不通过   6人工初审通过（认证未完成，还需按提示完成接下来的操作)
	HasCertificate    types.Int8   `json:"has_certificate,omitempty" db:"'has_certificate' size:3"`           // 是否申请了实名证书   1已申请
	AuthSign          types.Int8   `json:"auth_sign,omitempty" db:"'auth_sign' size:3"`                       // 自动签状态  1自动签
	AuthTransactionId types.String `json:"auth_transaction_id,omitempty" db:"'auth_transaction_id' size:100"` // 自动签交易号
	AuthContractId    types.String `json:"auth_contract_id,omitempty" db:"'auth_contract_id' size:100"`       // 自动签合同编号
	AuthResult        types.String `json:"auth_result,omitempty" db:"'auth_result' size:20"`                  // 自动签签章结果
	State             types.Int8   `json:"state,omitempty" db:"'state' size:3"`                               // 状态
	Ctime             types.Time   `json:"ctime,omitempty" db:"'ctime'"`                                      // 注册时间
}

// tablename company_role
type CompanyRole struct {
	pool.Model
	Id        types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`         //
	CompanyId types.BigInt `json:"company_id,omitempty" db:"'company_id' size:20"` // 公司id
	Name      types.String `json:"name,omitempty" db:"'name' size:10"`             // 角色名称
	Descr     types.String `json:"descr,omitempty" db:"'descr' size:45"`           // 角色描述
	Rules     types.String `json:"rules,omitempty" db:"'rules' size:1000"`         // 逗号分割的权限id字符串
	Type      types.Int8   `json:"type,omitempty" db:"'type' size:3"`              // 角色类型：1超级管理员，不可修改，拥有所有权限   0普通角色  2内置角色，不可修改
	State     types.Int8   `json:"state,omitempty" db:"'state' size:3"`            // 状态:1=启用
}

// tablename company_stamp
type CompanyStamp struct {
	pool.Model
	Id          types.Money  `json:"id,omitempty" db:"'id' pk size:19"`                   // 签章图片ID，法大大生成
	CompanyId   types.Money  `json:"company_id,omitempty" db:"'company_id' size:19"`      // 公司id
	Url         types.String `json:"url,omitempty" db:"'url' size:100"`                   // 印章路径
	Genre       types.Int8   `json:"genre,omitempty" db:"'genre' size:3"`                 // 类型  1公章      2合同章        3财务章
	IsDefault   types.Int8   `json:"is_default,omitempty" db:"'is_default' size:3"`       // 是否是默认章  1默认章
	Creator     types.BigInt `json:"creator,omitempty" db:"'creator' size:20"`            // 创建人
	CreatorName types.String `json:"creator_name,omitempty" db:"'creator_name' size:100"` //
	Department  types.String `json:"department,omitempty" db:"'department' size:50"`      // 部门
	Position    types.String `json:"position,omitempty" db:"'position' size:50"`          // 职务
	State       types.Int8   `json:"state,omitempty" db:"'state' size:3"`                 // 管理状态：-1=已删除，0=禁用，1=启用
	Status      types.Int8   `json:"status,omitempty" db:"'status' size:3"`               // 状态
	Ctime       types.Time   `json:"ctime,omitempty" db:"'ctime'"`                        //
	Utime       types.Time   `json:"utime,omitempty" db:"'utime'"`                        //
}

// tablename dispatch_company
type DispatchCompany struct {
	pool.Model
	Id        types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`         //
	CompanyId types.BigInt `json:"company_id,omitempty" db:"'company_id' size:20"` // 企业表id
	Name      types.String `json:"name,omitempty" db:"'name' size:100"`            // 名称
	Address   types.String `json:"address,omitempty" db:"'address' size:200"`      // 注册住所
	Creator   types.BigInt `json:"creator,omitempty" db:"'creator' size:20"`       // 创建人
	State     types.Int8   `json:"state,omitempty" db:"'state' size:3"`            // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime     types.Time   `json:"ctime,omitempty" db:"'ctime'"`                   //
	Utime     types.Time   `json:"utime,omitempty" db:"'utime'"`                   //
}

// tablename document_template
type DocumentTemplate struct {
	pool.Model
	Id       types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`           //
	Genre    types.Int8   `json:"genre,omitempty" db:"'genre' size:3"`              // 类型   1合同属性  2薪酬确认单
	Company  types.Money  `json:"company,omitempty" db:"'company' size:19"`         // 企业id   0为系统文件  默认分配文件使用的字段，该模板可手动分配给其他企业
	Way      types.Int8   `json:"way,omitempty" db:"'way' size:3"`                  // 乙方处理方式   1需要签署   2只需查看
	AuthSign types.Int8   `json:"auth_sign,omitempty" db:"'auth_sign' size:3"`      // 企业是否需要自动盖章   1需要
	Title    types.String `json:"title,omitempty" db:"'title' size:200"`            // 文件名称
	Content  types.String `json:"content,omitempty" db:"'content' size:4294967295"` // 模板内容（需要操作的文件）（法大大不需要这个字段）
	Url      types.String `json:"url,omitempty" db:"'url' size:200"`                // 文件储存路径
	HasForm  types.Int8   `json:"has_form,omitempty" db:"'has_form' size:3"`        // 是否将模板变量配置成了表单（dacument_variable）  1配置了，可动态生成变量表单
	State    types.Int8   `json:"state,omitempty" db:"'state' size:3"`              // 管理状态：-1=已删除，0=禁用，1=启用
	Utime    types.Time   `json:"utime,omitempty" db:"'utime'"`                     // 最后修改时间
}

// tablename account
type Account struct {
	pool.Model
	Id        types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`          //
	LoginName types.String `json:"login_name,omitempty" db:"'login_name' size:200"` // 登录名（手机号）
	Password  types.String `json:"password,omitempty" db:"'password' size:50"`      // 密码
	State     types.Int8   `json:"state,omitempty" db:"'state' size:3"`             // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime     types.Time   `json:"ctime,omitempty" db:"'ctime'"`                    //
	Utime     types.Time   `json:"utime,omitempty" db:"'utime'"`                    //
}

// tablename company_document
type CompanyDocument struct {
	pool.Model
	Id           types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`               //
	CompanyId    types.BigInt `json:"company_id,omitempty" db:"'company_id' size:20"`       // 企业id   0为系统文件
	TemplateId   types.Money  `json:"template_id,omitempty" db:"'template_id' size:19"`     // document_template表id
	Title        types.String `json:"title,omitempty" db:"'title' size:200"`                // 文件名称（利于人工查看）
	Classify     types.Int8   `json:"classify,omitempty" db:"'classify' size:3"`            // 文件分类  10入职材料  20 离职材料...  详情在字典
	Genre        types.Int8   `json:"genre,omitempty" db:"'genre' size:3"`                  // 员工人事关系类型 0未分类  1通用 10劳动  20劳务   30返聘...  详情在字典
	HandleGenre  types.Int8   `json:"handle_genre,omitempty" db:"'handle_genre' size:3"`    // 操作类型 根据classify详情见各自的genre类型字典   例如离职操作有 1试用期终止  2合同到期  3协商解除……
	Job          types.String `json:"job,omitempty" db:"'job' size:50"`                     // 岗位
	VariableMode types.Int32  `json:"variable_mode,omitempty" db:"'variable_mode' size:10"` // 文件变量获取处理方式
	IsDefault    types.Int8   `json:"is_default,omitempty" db:"'is_default' size:3"`        // 是否默认分配  1默认分配  2不分配 （随入离职和在职事件，直接默认配置的待签署查看文件）
	CanDefault   types.Int8   `json:"can_default,omitempty" db:"'can_default' size:3"`      // 是否可以修改默认  1可修改   2不可修改
	Modifier     types.Money  `json:"modifier,omitempty" db:"'modifier' size:19"`           // 修改人id
	State        types.Int8   `json:"state,omitempty" db:"'state' size:3"`                  // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime        types.Time   `json:"ctime,omitempty" db:"'ctime'"`                         // 创建时间
}

// tablename company_man
type CompanyMan struct {
	pool.Model
	Id        types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`         // 操作员id
	AccountId types.BigInt `json:"account_id,omitempty" db:"'account_id' size:20"` // 账号表id
	CompanyId types.BigInt `json:"company_id,omitempty" db:"'company_id' size:20"` // 企业表id
	RealName  types.String `json:"real_name,omitempty" db:"'real_name' size:50"`   // 姓名
	Email     types.String `json:"email,omitempty" db:"'email' size:100"`          // 邮箱
	Roles     types.String `json:"roles,omitempty" db:"'roles' size:100"`          // 逗号分割的角色id字符串
	State     types.Int8   `json:"state,omitempty" db:"'state' size:3"`            // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime     types.Time   `json:"ctime,omitempty" db:"'ctime'"`                   //
	Utime     types.Time   `json:"utime,omitempty" db:"'utime'"`                   //
}

// tablename company_rule
type CompanyRule struct {
	pool.Model
	Id        types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`        //
	Pid       types.BigInt `json:"pid,omitempty" db:"'pid' size:20"`              // 父级id
	Path      types.String `json:"path,omitempty" db:"'path' size:45"`            // 标识
	Title     types.String `json:"title,omitempty" db:"'title' size:20"`          // 名称
	Type      types.Uint8  `json:"type,omitempty" db:"'type' size:3"`             // 前端ui类型:0=菜单,1=其它
	IsPrivate types.Uint8  `json:"is_private,omitempty" db:"'is_private' size:3"` // 权限是否公开    0公开  1私有定制
	State     types.Int8   `json:"state,omitempty" db:"'state' size:3"`           // 状态:1=启用
	Descr     types.String `json:"descr,omitempty" db:"'descr' size:45"`          // 描述
	Belong    types.Int8   `json:"belong,omitempty" db:"'belong' size:3"`         // 权限所属端   (目前没有使用)  0通用  3松鼠用工企业用户   4松鼠极聘企业用户
}
