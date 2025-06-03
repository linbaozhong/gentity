package do

import (
	"github.com/linbaozhong/gentity/pkg/ace/pool"
	"github.com/linbaozhong/gentity/pkg/types"
)

// Account 账号表
// @tablename account
type Account struct {
	pool.Model
	Id        types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`          //
	LoginName types.String `json:"login_name,omitempty" db:"'login_name' size:200"` // 登录名（手机号）
	Password  types.String `json:"password,omitempty" db:"'password' size:50"`      // 密码
	State     types.Int8   `json:"state,omitempty" db:"'state' size:3"`             // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime     types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                 //
	Utime     types.Time   `json:"utime,omitempty" db:"'utime' <-"`                 //
}

// Communication 通讯表
// @tablename communication
type Communication struct {
	pool.Model
	Id       types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`       // id
	Company  types.BigInt `json:"company,omitempty" db:"'company' size:20"`     // 企业id
	Phone    types.String `json:"phone,omitempty" db:"'phone' size:20"`         // 手机号
	ManName  types.String `json:"man_name,omitempty" db:"'man_name' size:50"`   // 发送人姓名
	UserName types.String `json:"user_name,omitempty" db:"'user_name' size:50"` // 接收人信息
	Platform types.String `json:"platform,omitempty" db:"'platform' size:100"`  // 平台   短信
	Genre    types.Int8   `json:"genre,omitempty" db:"'genre' size:3"`          // 类型
	Result   types.String `json:"result,omitempty" db:"'result' size:200"`      // 结果
	MsgId    types.String `json:"msg_id,omitempty" db:"'msg_id' size:100"`      // 消息结果回填id
	State    types.Int8   `json:"state,omitempty" db:"'state' size:3"`          // 管理状态：-1=已删除，0=禁用，1=启用，2=隐藏，3=历史
	Ctime    types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`              //
	Utime    types.Time   `json:"utime,omitempty" db:"'utime' <-"`              //
}

// Company 企业表
// @tablename company
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
	Ctime            types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                              //
	Utime            types.Time   `json:"utime,omitempty" db:"'utime' <-"`                              //
}

// CompanyDocument 企业文件
// @tablename company_document
type CompanyDocument struct {
	pool.Model
	Id           types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`               //
	Company      types.BigInt `json:"company,omitempty" db:"'company' size:20"`             // 企业id   0为系统文件
	TemplateId   types.BigInt `json:"template_id,omitempty" db:"'template_id' size:20"`     // document_template表id
	Title        types.String `json:"title,omitempty" db:"'title' size:200"`                // 文件名称（利于人工查看）
	Classify     types.Int8   `json:"classify,omitempty" db:"'classify' size:3"`            // 文件分类
	Genre        types.Int8   `json:"genre,omitempty" db:"'genre' size:3"`                  // 员工人事关系类型 0未分类  1通用 10劳动  20劳务   30返聘...  详情在字典
	HandleGenre  types.Int8   `json:"handle_genre,omitempty" db:"'handle_genre' size:3"`    // 操作类型 根据classify详情见各自的genre类型字典   例如离职操作有 1试用期终止  2合同到期  3协商解除……
	Job          types.String `json:"job,omitempty" db:"'job' size:50"`                     // 岗位
	VariableMode types.Int32  `json:"variable_mode,omitempty" db:"'variable_mode' size:10"` // 文件变量获取处理方式
	IsDefault    types.Int8   `json:"is_default,omitempty" db:"'is_default' size:3"`        // 是否默认分配  1默认分配  2不分配 （随入离职和在职事件，直接默认配置的待签署查看文件）
	CanDefault   types.Int8   `json:"can_default,omitempty" db:"'can_default' size:3"`      // 是否可以修改默认  1可修改   2不可修改
	Modifier     types.BigInt `json:"modifier,omitempty" db:"'modifier' size:20"`           // 修改人id
	State        types.Int8   `json:"state,omitempty" db:"'state' size:3"`                  // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime        types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                      // 创建时间
}

// CompanyFadada 法大大企业表
// @tablename company_fadada
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
	Ctime             types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                                   // 注册时间
}

// CompanyMan 管理员表
// @tablename company_man
type CompanyMan struct {
	pool.Model
	Id         types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`          // 操作员id
	AccountId  types.BigInt `json:"account_id,omitempty" db:"'account_id' size:20"`  // 账号表id
	Company    types.BigInt `json:"company,omitempty" db:"'company' size:20"`        // 企业表id
	RealName   types.String `json:"real_name,omitempty" db:"'real_name' size:50"`    // 姓名
	Email      types.String `json:"email,omitempty" db:"'email' size:100"`           // 邮箱
	Roles      types.String `json:"roles,omitempty" db:"'roles' size:100"`           // 逗号分割的角色id字符串
	Gender     types.String `json:"gender,omitempty" db:"'gender' size:20"`          // 性别  男 女
	Genre      types.Int8   `json:"genre,omitempty" db:"'genre' size:3"`             // 账号类型  1主账号，拥有管理员权限，角色不能修改  2子账号，可分配权限
	IsActivate types.Int8   `json:"is_activate,omitempty" db:"'is_activate' size:3"` // 是否激活（未激活不能访问）  1已激活  2未激活
	LoginTime  types.Time   `json:"login_time,omitempty" db:"'login_time'"`          // 最后登录时间
	State      types.Int8   `json:"state,omitempty" db:"'state' size:3"`             // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime      types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                 //
	Utime      types.Time   `json:"utime,omitempty" db:"'utime' <-"`                 //
}

// CompanyMaterial 企业资料表
// @tablename company_material
type CompanyMaterial struct {
	pool.Model
	Id         types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`          // id
	Company    types.BigInt `json:"company,omitempty" db:"'company' size:20"`        // 派遣企业id
	Name       types.String `json:"name,omitempty" db:"'name' size:100"`             // 名称
	Genre      types.Int8   `json:"genre,omitempty" db:"'genre' size:3"`             // 类型  1系统内置  2企业添加
	IsChecked  types.Int8   `json:"is_checked,omitempty" db:"'is_checked' size:3"`   // 是否选中  1选中   2未选中
	IsRequired types.Int8   `json:"is_required,omitempty" db:"'is_required' size:3"` // 是否必填  1必填   2非必填
	Remark     types.String `json:"remark,omitempty" db:"'remark' size:100"`         // 备注
	State      types.Int8   `json:"state,omitempty" db:"'state' size:3"`             // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime      types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                 // 创建时间
}

// CompanyRole 公司管理员角色表
// @tablename company_role
type CompanyRole struct {
	pool.Model
	Id      types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`   //
	Company types.BigInt `json:"company,omitempty" db:"'company' size:20"` // 公司id
	Name    types.String `json:"name,omitempty" db:"'name' size:10"`       // 角色名称
	Descr   types.String `json:"descr,omitempty" db:"'descr' size:45"`     // 角色描述
	Rules   types.String `json:"rules,omitempty" db:"'rules' size:1000"`   // 逗号分割的权限id字符串
	Type    types.Int8   `json:"type,omitempty" db:"'type' size:3"`        // 角色类型：1管理员角色    2普通角色
	State   types.Int8   `json:"state,omitempty" db:"'state' size:3"`      // 状态:1=启用
}

// CompanyRule 公司管理员权限表
// @tablename company_rule
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

// CompanyStamp 企业印章表
// @tablename company_stamp
type CompanyStamp struct {
	pool.Model
	Id          types.BigInt `json:"id,omitempty" db:"'id' pk size:20"`                   // 签章图片ID，法大大生成
	Company     types.BigInt `json:"company,omitempty" db:"'company' size:20"`            // 公司id
	Url         types.String `json:"url,omitempty" db:"'url' size:100"`                   // 印章路径
	Genre       types.Int8   `json:"genre,omitempty" db:"'genre' size:3"`                 // 类型  1公章      2合同章        3财务章
	IsDefault   types.Int8   `json:"is_default,omitempty" db:"'is_default' size:3"`       // 是否是默认章  1默认章
	Creator     types.BigInt `json:"creator,omitempty" db:"'creator' size:20"`            // 创建人
	CreatorName types.String `json:"creator_name,omitempty" db:"'creator_name' size:100"` //
	State       types.Int8   `json:"state,omitempty" db:"'state' size:3"`                 // 管理状态：-1=已删除，0=禁用，1=启用
	Status      types.Int8   `json:"status,omitempty" db:"'status' size:3"`               // 状态
	Ctime       types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                     //
	Utime       types.Time   `json:"utime,omitempty" db:"'utime' <-"`                     //
}

// Department 部门
// @tablename department
type Department struct {
	pool.Model
	Id      types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`   //
	Name    types.String `json:"name,omitempty" db:"'name' size:15"`       // 部门名称
	Descr   types.String `json:"descr,omitempty" db:"'descr' size:100"`    // 部门简介
	State   types.Int8   `json:"state,omitempty" db:"'state' size:3"`      // 状态：-1=删除,0=禁用,1=可用
	Creator types.BigInt `json:"creator,omitempty" db:"'creator' size:20"` // 创建人
	Ctime   types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`          // 创建时间
}

// DispatchedEmployee 被派遣员工表
// @tablename dispatched_employee
type DispatchedEmployee struct {
	pool.Model
	Id                types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`                         //
	EmploymentCompany types.BigInt `json:"employment_company,omitempty" db:"'employment_company' size:20"` // 用工企业id
	Employee          types.BigInt `json:"employee,omitempty" db:"'employee' size:20"`                     // 派遣企业员工id
	Company           types.BigInt `json:"company,omitempty" db:"'company' size:20"`                       // 派遣企业id
	User              types.BigInt `json:"user,omitempty" db:"'user' size:20"`                             // 劳动者id
	HireDate          types.Time   `json:"hire_date,omitempty" db:"'hire_date'"`                           // 入职日期
	Position          types.String `json:"position,omitempty" db:"'position' size:100"`                    // 职务
	Job               types.String `json:"job,omitempty" db:"'job' size:100"`                              // 岗位
	WorkAddress       types.String `json:"work_address,omitempty" db:"'work_address' size:100"`            // 工作地点
	SalaryDay         types.String `json:"salary_day,omitempty" db:"'salary_day' size:20"`                 // 薪酬支付每月(日)
	Salary            types.Money  `json:"salary,omitempty" db:"'salary' size:19"`                         // 薪酬工资(月)
	SalaryTrial       types.Money  `json:"salary_trial,omitempty" db:"'salary_trial' size:19"`             // 试用期薪资
	SalaryScale       types.String `json:"salary_scale,omitempty" db:"'salary_scale' size:50"`             // 薪资调整标准
	WorkSchedules     types.String `json:"work_schedules,omitempty" db:"'work_schedules' size:50"`         // 工时制度
	Creator           types.BigInt `json:"creator,omitempty" db:"'creator' size:20"`                       // 创建人
	Status            types.Int8   `json:"status,omitempty" db:"'status' size:3"`                          // 状态 1派遣中   2未派遣
	SignStatus        types.Int8   `json:"sign_status,omitempty" db:"'sign_status' size:3"`                // 签署状态 1已签署  2未签署  3空 （签署结果是通用接口，此状态签署完没有更新，所以每回使用前需要先更新）
	EndDate           types.Time   `json:"end_date,omitempty" db:"'end_date'"`                             // 结束时间
	State             types.Int8   `json:"state,omitempty" db:"'state' size:3"`                            // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime             types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                                //
	Utime             types.Time   `json:"utime,omitempty" db:"'utime' <-"`                                //
}

// DispatchedEmployeeContract 被派遣员工派遣合同表
// @tablename dispatched_employee_contract
type DispatchedEmployeeContract struct {
	pool.Model
	Id               types.BigInt `json:"id,omitempty" db:"'id' pk size:20"`                            //
	ContractSubject  types.Money  `json:"contract_subject,omitempty" db:"'contract_subject' size:19"`   // 合同主体，企业id
	ContractType     types.Int8   `json:"contract_type,omitempty" db:"'contract_type' size:3"`          // 合同类型
	ContractBegin    types.Time   `json:"contract_begin,omitempty" db:"'contract_begin'"`               // 合同开始日期
	ContractEnd      types.Time   `json:"contract_end,omitempty" db:"'contract_end'"`                   // 合同终止日期
	ContractDuration types.Uint32 `json:"contract_duration,omitempty" db:"'contract_duration' size:10"` // 合同期限（月）
	TrialPeriod      types.Uint8  `json:"trial_period,omitempty" db:"'trial_period' size:3"`            // 试用期（月）
	TrialPeriodDay   types.Int8   `json:"trial_period_day,omitempty" db:"'trial_period_day' size:3"`    // 试用期额外天数（天）
	PeriodEndDay     types.Time   `json:"period_end_day,omitempty" db:"'period_end_day'"`               // 试用期结束日期 ，如没有通过试用期计算
	Creator          types.BigInt `json:"creator,omitempty" db:"'creator' size:20"`                     // 创建人
	State            types.Int8   `json:"state,omitempty" db:"'state' size:3"`                          // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime            types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                              //
}

// DocumentOffline 线下文件表
// @tablename document_offline
type DocumentOffline struct {
	pool.Model
	Id                 types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`                           // id
	Employee           types.BigInt `json:"employee,omitempty" db:"'employee' size:20"`                       // 派遣企业员工id
	DispatchedEmployee types.BigInt `json:"dispatched_employee,omitempty" db:"'dispatched_employee' size:20"` //  被派遣员工id
	Name               types.String `json:"name,omitempty" db:"'name' size:200"`                              // 文件名称
	Url                types.String `json:"url,omitempty" db:"'url' size:65535"`                              // 文件保存地址，多个逗号分割
	State              types.Int8   `json:"state,omitempty" db:"'state' size:3"`                              // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime              types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                                  // 创建时间
}

// DocumentTemplate 文件协议模板
// @tablename document_template
type DocumentTemplate struct {
	pool.Model
	Id       types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`           //
	Genre    types.Int8   `json:"genre,omitempty" db:"'genre' size:3"`              // 类型   1合同属性  2薪酬确认单
	Company  types.BigInt `json:"company,omitempty" db:"'company' size:20"`         // 企业id   0为系统文件  默认分配文件使用的字段，该模板可手动分配给其他企业
	Way      types.Int8   `json:"way,omitempty" db:"'way' size:3"`                  // 乙方处理方式   1需要签署   2只需查看
	AuthSign types.Int8   `json:"auth_sign,omitempty" db:"'auth_sign' size:3"`      // 企业是否需要自动盖章   1需要
	Title    types.String `json:"title,omitempty" db:"'title' size:200"`            // 文件名称
	Content  types.String `json:"content,omitempty" db:"'content' size:4294967295"` // 模板内容（需要操作的文件）（法大大不需要这个字段）
	Url      types.String `json:"url,omitempty" db:"'url' size:200"`                // 文件储存路径
	HasForm  types.Int8   `json:"has_form,omitempty" db:"'has_form' size:3"`        // 是否将模板变量配置成了表单（dacument_variable）  1配置了，可动态生成变量表单
	State    types.Int8   `json:"state,omitempty" db:"'state' size:3"`              // 管理状态：-1=已删除，0=禁用，1=启用
	Utime    types.Time   `json:"utime,omitempty" db:"'utime' <-"`                  // 最后修改时间
}

// EmailRecord email发送记录
// @tablename email_record
type EmailRecord struct {
	pool.Model
	Id         types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`          // id
	To         types.String `json:"to,omitempty" db:"'to' size:100"`                 // 发送地址
	Cc         types.String `json:"cc,omitempty" db:"'cc' size:300"`                 // 抄送地址，多个逗号分割
	Bcc        types.String `json:"bcc,omitempty" db:"'bcc' size:300"`               // 密送地址，多个逗号分割
	Attachment types.String `json:"attachment,omitempty" db:"'attachment' size:20"`  // 附件，多个逗号分割
	HtmlBody   types.String `json:"html_body,omitempty" db:"'html_body' size:65535"` // 内容
	Subject    types.String `json:"subject,omitempty" db:"'subject' size:255"`       // 标题
	IsSms      types.Int8   `json:"is_sms,omitempty" db:"'is_sms' size:3"`           // 是否短信通知   1通知   2不通知
	Mobile     types.String `json:"mobile,omitempty" db:"'mobile' size:30"`          // 手机号
	Creator    types.BigInt `json:"creator,omitempty" db:"'creator' size:20"`        // 创建人id
	Ctime      types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                 // 创建时间
}

// Employee 派遣员工表
// @tablename employee
type Employee struct {
	pool.Model
	Id             types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`                  //
	Company        types.BigInt `json:"company,omitempty" db:"'company' size:20"`                // 派遣企业id
	User           types.BigInt `json:"user,omitempty" db:"'user' size:20"`                      // 劳动者id
	RealName       types.String `json:"real_name,omitempty" db:"'real_name' size:20"`            // 姓名
	Gender         types.String `json:"gender,omitempty" db:"'gender' size:5"`                   // 性别
	IdCard         types.String `json:"id_card,omitempty" db:"'id_card' size:50"`                // 身份证号
	Birthday       types.String `json:"birthday,omitempty" db:"'birthday' size:50"`              // 出生日期
	Mobile         types.String `json:"mobile,omitempty" db:"'mobile' size:20"`                  // 电话号码
	Email          types.String `json:"email,omitempty" db:"'email' size:64"`                    // 邮箱
	Address        types.String `json:"address,omitempty" db:"'address' size:45"`                // 现居住地
	HukouAddress   types.String `json:"hukou_address,omitempty" db:"'hukou_address' size:250"`   // 户口所在地
	HukouType      types.String `json:"hukou_type,omitempty" db:"'hukou_type' size:50"`          // 户籍类型   农业  非农业
	Creator        types.BigInt `json:"creator,omitempty" db:"'creator' size:20"`                // 创建人
	Status         types.Int8   `json:"status,omitempty" db:"'status' size:3"`                   // 状态 10在职   20离职中  30已离职
	DispatchStatus types.Int8   `json:"dispatch_status,omitempty" db:"'dispatch_status' size:3"` // 状态  1派遣中   2未派遣
	State          types.Int8   `json:"state,omitempty" db:"'state' size:3"`                     // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime          types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                         //
	Utime          types.Time   `json:"utime,omitempty" db:"'utime' <-"`                         //
}

// EmployeeMaterial 派遣员工资料表
// @tablename employee_material
type EmployeeMaterial struct {
	pool.Model
	Id              types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`                     // id
	Employee        types.BigInt `json:"employee,omitempty" db:"'employee' size:20"`                 // 员工表id
	CompanyMaterial types.BigInt `json:"company_material,omitempty" db:"'company_material' size:20"` // material_company表id
	Name            types.String `json:"name,omitempty" db:"'name' size:100"`                        // 名称
	Url             types.String `json:"url,omitempty" db:"'url' size:65535"`                        // 文件保存路径，逗号分割的数组
	UrlTwo          types.String `json:"url_two,omitempty" db:"'url_two' size:65535"`                // 特殊类型的文件保存路径，逗号分割的数组 （例如身份证反面）
	Remark          types.String `json:"remark,omitempty" db:"'remark' size:100"`                    // 备注
	IsRequired      types.Int8   `json:"is_required,omitempty" db:"'is_required' size:3"`            // 是否必填  1必填   2非必填
	Ctime           types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                            // 创建时间
}

// EmployeeUpload 员工上传图片表
// @tablename employee_upload
type EmployeeUpload struct {
	pool.Model
	Id       types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`     // id
	Employee types.BigInt `json:"employee,omitempty" db:"'employee' size:20"` // companys_mans表id
	Name     types.String `json:"name,omitempty" db:"'name' size:100"`        // 名称
	Url      types.String `json:"url,omitempty" db:"'url' size:65535"`        // 文件保存路径，逗号分割的数组
	Remark   types.String `json:"remark,omitempty" db:"'remark' size:100"`    // 备注
	State    types.Int8   `json:"state,omitempty" db:"'state' size:3"`        // 管理状态：-1=已删除，0=禁用，1=启用
	Creator  types.BigInt `json:"creator,omitempty" db:"'creator' size:20"`   // 创建人
	Ctime    types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`            // 创建时间
}

// EmploymentCompany 用工企业表
// @tablename employment_company
type EmploymentCompany struct {
	pool.Model
	Id            types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`            //
	Company       types.BigInt `json:"company,omitempty" db:"'company' size:20"`          // 企业表id
	Name          types.String `json:"name,omitempty" db:"'name' size:100"`               // 名称
	ContractNo    types.String `json:"contract_no,omitempty" db:"'contract_no' size:200"` // 合同编号
	ContractStart types.Time   `json:"contract_start,omitempty" db:"'contract_start'"`    // 服务周期起始日期
	ContractEnd   types.Time   `json:"contract_end,omitempty" db:"'contract_end'"`        // 服务周期结束日期
	Creator       types.BigInt `json:"creator,omitempty" db:"'creator' size:20"`          // 创建人
	Status        types.Int8   `json:"status,omitempty" db:"'status' size:3"`             // 状态  1正常  2关闭  3异常
	State         types.Int8   `json:"state,omitempty" db:"'state' size:3"`               // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime         types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                   //
	Utime         types.Time   `json:"utime,omitempty" db:"'utime' <-"`                   //
}

// ExtSignBatch 手动签交易
// @tablename ext_sign_batch
type ExtSignBatch struct {
	pool.Model
	Id         types.BigInt `json:"id,omitempty" db:"'id' pk size:20"`                 // 批量手动签署交易号，长度<=32
	BatchTitle types.String `json:"batch_title,omitempty" db:"'batch_title' size:100"` // 批量请求标题
	CustomerId types.String `json:"customer_id,omitempty" db:"'customer_id' size:100"` // 客户编号
	Url        types.String `json:"url,omitempty" db:"'url' size:65535"`               // 签署url
	Ctime      types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                   // 创建时间
}

// ExtSignTransaction 手动签交易
// @tablename ext_sign_transaction
type ExtSignTransaction struct {
	pool.Model
	Id           types.BigInt `json:"id,omitempty" db:"'id' pk size:20"`                    // 手动签署交易号，长度<=32
	SignDocument types.BigInt `json:"sign_document,omitempty" db:"'sign_document' size:20"` // sign_document表id
	BatchId      types.BigInt `json:"batch_id,omitempty" db:"'batch_id' size:20"`           // ext_sign_batch表id
	Url          types.String `json:"url,omitempty" db:"'url' size:65535"`                  // 签署url
	Ctime        types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                      // 创建时间
}

// ExtTransactionFail 手动签失败交易记录
// @tablename ext_transaction_fail
type ExtTransactionFail struct {
	pool.Model
	Id           types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`               // id
	Transaction  types.BigInt `json:"transaction,omitempty" db:"'transaction' size:20"`     // 手动签署交易号，长度<=32
	SignDocument types.BigInt `json:"sign_document,omitempty" db:"'sign_document' size:20"` // sign_document表id
	Status       types.String `json:"status,omitempty" db:"'status' size:50"`               // 状态
	Result       types.String `json:"result,omitempty" db:"'result' size:50"`               // 结果
	Url          types.String `json:"url,omitempty" db:"'url' size:65535"`                  // 签署url
	Ctime        types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                      // 创建时间
}

// ExtTransactionSucceed 手动签成功交易记录
// @tablename ext_transaction_succeed
type ExtTransactionSucceed struct {
	pool.Model
	Id           types.BigInt `json:"id,omitempty" db:"'id' pk size:20"`                    // 手动签署交易号，长度<=32
	SignDocument types.BigInt `json:"sign_document,omitempty" db:"'sign_document' size:20"` // sign_document表id
	Status       types.String `json:"status,omitempty" db:"'status' size:50"`               // 状态
	Result       types.String `json:"result,omitempty" db:"'result' size:50"`               // 结果
	Url          types.String `json:"url,omitempty" db:"'url' size:65535"`                  // 签署url
	Ctime        types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                      // 创建时间
}

// FadadaLog 法大大回调函数日志
// @tablename fadada_log
type FadadaLog struct {
	pool.Model
	Id            types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`                  //
	CustomerId    types.String `json:"customer_id,omitempty" db:"'customer_id' size:100"`       // 客户编号
	TransactionId types.String `json:"transaction_id,omitempty" db:"'transaction_id' size:100"` // 合同编号
	ContractId    types.String `json:"contract_id,omitempty" db:"'contract_id' size:100"`       // 交易号
	From          types.String `json:"from,omitempty" db:"'from' size:100"`                     //
	Type          types.String `json:"type,omitempty" db:"'type' size:30"`                      // 类型
	Values        types.String `json:"values,omitempty" db:"'values' size:65535"`               //
	Ip            types.String `json:"ip,omitempty" db:"'ip' size:30"`                          // IP地址
	Status        types.String `json:"status,omitempty" db:"'status' size:10"`                  // 成功=succeed；失败=failed
	StatusTime    types.Time   `json:"status_time,omitempty" db:"'status_time'"`                //
	StatusReason  types.String `json:"status_reason,omitempty" db:"'status_reason' size:500"`   //
	Ctime         types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                         // 创建时间
}

// LaborContract 员工合同
// @tablename labor_contract
type LaborContract struct {
	pool.Model
	Id                 types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`                           // id
	Company            types.BigInt `json:"company,omitempty" db:"'company' size:20"`                         // 企业id
	Employee           types.BigInt `json:"employee,omitempty" db:"'employee' size:20"`                       // 派遣员工id
	DispatchedEmployee types.BigInt `json:"dispatched_employee,omitempty" db:"'dispatched_employee' size:20"` //  被派遣员工id
	ContractSubject    types.BigInt `json:"contract_subject,omitempty" db:"'contract_subject' size:20"`       // 合同主体，企业id
	ContractType       types.Int8   `json:"contract_type,omitempty" db:"'contract_type' size:3"`              // 合同类型
	ContractBegin      types.Time   `json:"contract_begin,omitempty" db:"'contract_begin'"`                   // 合同开始日期
	ContractEnd        types.Time   `json:"contract_end,omitempty" db:"'contract_end'"`                       // 合同终止日期
	ContractDuration   types.Uint32 `json:"contract_duration,omitempty" db:"'contract_duration' size:10"`     // 合同期限（月）
	TrialPeriod        types.Int8   `json:"trial_period,omitempty" db:"'trial_period' size:3"`                // 试用期
	TrialPeriodDay     types.Int8   `json:"trial_period_day,omitempty" db:"'trial_period_day' size:3"`        // 试用期额外天数（天）
	Source             types.String `json:"source,omitempty" db:"'source' size:50"`                           // 来源 入职签署 续签签署 在职签署
	DetailTable        types.String `json:"detail_table,omitempty" db:"'detail_table' size:50"`               // 详情涉及的表名
	DetailId           types.BigInt `json:"detail_id,omitempty" db:"'detail_id' size:20"`                     // 详情id
	SignDocument       types.BigInt `json:"sign_document,omitempty" db:"'sign_document' size:20"`             // sign_document表id
	State              types.Uint8  `json:"state,omitempty" db:"'state' size:3"`                              // 管理状态：-1=已删除，0=禁用，1=启用
	Status             types.Uint8  `json:"status,omitempty" db:"'status' size:3"`                            // 状态  1最新使用的合同     2已过期合同    3作废合同
	Ctime              types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                                  // 创建时间
	Utime              types.Time   `json:"utime,omitempty" db:"'utime' <-"`                                  //
}

// Man 平台操作员表
// @tablename man
type Man struct {
	pool.Model
	Id         types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`         //
	Mobile     types.String `json:"mobile,omitempty" db:"'mobile' size:100"`        // 手机号码
	Code       types.String `json:"code,omitempty" db:"'code' size:45"`             // 短信验证码,只用于绑定支付宝或微信,每次绑定成功后,用新的验证码替换,并且open_bind置为disable
	Name       types.String `json:"name,omitempty" db:"'name' size:45"`             // 真实姓名
	Gender     types.Int8   `json:"gender,omitempty" db:"'gender' size:3"`          // 性别
	Email      types.String `json:"email,omitempty" db:"'email' size:45"`           // 邮箱
	Department types.BigInt `json:"department,omitempty" db:"'department' size:20"` // 部门id
	State      types.Int8   `json:"state,omitempty" db:"'state' size:3"`            // 状态:1=启用
	AlipayId   types.String `json:"alipay_id,omitempty" db:"'alipay_id' size:45"`   // 绑定支付宝id
	WeixinId   types.String `json:"weixin_id,omitempty" db:"'weixin_id' size:45"`   // 绑定微信id
	Roles      types.String `json:"roles,omitempty" db:"'roles' size:100"`          // 逗号分割的角色id字符串
	OpenBind   types.Uint8  `json:"open_bind,omitempty" db:"'open_bind' size:3"`    // 是否开启绑定 1=开启
	Creator    types.BigInt `json:"creator,omitempty" db:"'creator' size:20"`       // 创建人
	Ctime      types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                // 创建时间
}

// MaterialTemplate 资料模板表
// @tablename material_template
type MaterialTemplate struct {
	pool.Model
	Id         types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`          // id
	Name       types.String `json:"name,omitempty" db:"'name' size:100"`             // 名称
	IsChecked  types.Int8   `json:"is_checked,omitempty" db:"'is_checked' size:3"`   // 默认是否选中  1选中   2未选中
	IsRequired types.Int8   `json:"is_required,omitempty" db:"'is_required' size:3"` // 默认是否必填  1必填   2非必填
	State      types.Int8   `json:"state,omitempty" db:"'state' size:3"`             // 管理状态：-1=已删除，0=禁用，1=启用
}

// OperLog 操作日志
// @tablename oper_log
type OperLog struct {
	pool.Model
	Id       types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`          //
	Endpoint types.String `json:"endpoint,omitempty" db:"'endpoint' size:50"`      // 终端名
	User     types.BigInt `json:"user,omitempty" db:"'user' size:20"`              // 用户id
	UserName types.String `json:"user_name,omitempty" db:"'user_name' size:45"`    // 用户名称
	Table    types.String `json:"table,omitempty" db:"'table' size:25"`            // 表名
	TableId  types.BigInt `json:"table_id,omitempty" db:"'table_id' size:20"`      // 表id
	OperType types.Int8   `json:"oper_type,omitempty" db:"'oper_type' size:3"`     // 操作方式：-1=删除 0=修改 1=新增
	OperInfo types.String `json:"oper_info,omitempty" db:"'oper_info' size:65535"` // 操作描述
	OperData types.String `json:"oper_data,omitempty" db:"'oper_data'"`            //
	State    types.Int32  `json:"state,omitempty" db:"'state' size:10"`            // 状态
	Ctime    types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                 // 操作时间
	Ip       types.String `json:"ip,omitempty" db:"'ip' size:15"`                  // ip地址
}

// PubDict
// @tablename pub_dict
type PubDict struct {
	pool.Model
	Id    types.Uint32 `json:"id,omitempty" db:"'id' pk auto size:10"` //
	Type  types.Uint8  `json:"type,omitempty" db:"'type' size:3"`      // 类型对应pub_dict -type=0
	Code  types.Int16  `json:"code,omitempty" db:"'code' size:5"`      // 对应字典 id
	Name  types.String `json:"name,omitempty" db:"'name' size:50"`     // 名称
	Memo  types.String `json:"memo,omitempty" db:"'memo' size:300"`    // 说明
	State types.Int8   `json:"state,omitempty" db:"'state' size:3"`    // 状态\n对应pub_dict -type=10
}

// Resign 离职表
// @tablename resign
type Resign struct {
	pool.Model
	Id               types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`                         // 离职id
	Company          types.BigInt `json:"company,omitempty" db:"'company' size:20"`                       // 企业id
	Employee         types.BigInt `json:"employee,omitempty" db:"'employee' size:20"`                     // 员工id
	Platform         types.String `json:"platform,omitempty" db:"'platform' size:50"`                     // 请求来自哪一端   b（公司提出）     c（个人提出）
	Genre            types.Int8   `json:"genre,omitempty" db:"'genre' size:3"`                            // 离职类型   详情看字典
	GenreName        types.String `json:"genre_name,omitempty" db:"'genre_name' size:20"`                 //
	Reason           types.String `json:"reason,omitempty" db:"'reason' size:200"`                        // 离职原因
	ProposeDay       types.Time   `json:"propose_day,omitempty" db:"'propose_day'"`                       // 提出日期
	LastDay          types.Time   `json:"last_day,omitempty" db:"'last_day'"`                             // 最后工作日
	WhetherHandover  types.Int8   `json:"whether_handover,omitempty" db:"'whether_handover' size:3"`      // 是否需要工作交接  1需要
	WorkHandover     types.BigInt `json:"work_handover,omitempty" db:"'work_handover' size:20"`           // 工作交接人
	ToDay            types.Time   `json:"to_day,omitempty" db:"'to_day'"`                                 // 完成工作交接日期
	RelieveDay       types.Time   `json:"relieve_day,omitempty" db:"'relieve_day'"`                       // 离职日期
	Compensation     types.Money  `json:"compensation,omitempty" db:"'compensation' size:19"`             // 经济补偿金/赔偿金
	CompensationDay  types.Time   `json:"compensation_day,omitempty" db:"'compensation_day'"`             // 补偿金发放日期
	BalanceDay       types.Time   `json:"balance_day,omitempty" db:"'balance_day'"`                       // 结算日期
	Creator          types.BigInt `json:"creator,omitempty" db:"'creator' size:20"`                       // 创建人
	ProveCreatorName types.String `json:"prove_creator_name,omitempty" db:"'prove_creator_name' size:50"` // 离职证明开具人
	SendEmail        types.String `json:"send_email,omitempty" db:"'send_email' size:100"`                // 发送离职通知的邮件
	SendMobile       types.String `json:"send_mobile,omitempty" db:"'send_mobile' size:50"`               // 发送离职通知的手机号
	Status           types.Uint8  `json:"status,omitempty" db:"'status' size:3"`                          // 状态
	StatusLog        types.String `json:"status_log,omitempty" db:"'status_log' size:65535"`              // 状态日志
	StatusTime       types.Time   `json:"status_time,omitempty" db:"'status_time'"`                       // 状态时间
	StatusReason     types.String `json:"status_reason,omitempty" db:"'status_reason' size:200"`          // 状态原因
	State            types.Int8   `json:"state,omitempty" db:"'state' size:3"`                            // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime            types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                                // 创建时间
	Utime            types.Time   `json:"utime,omitempty" db:"'utime' <-"`                                //
}

// ResignFlow 离职审核表
// @tablename resign_flow
type ResignFlow struct {
	pool.Model
	Id           types.BigInt `json:"id,omitempty" db:"'id' pk size:20"`                      // resign表id
	Status       types.Int8   `json:"status,omitempty" db:"'status' size:3"`                  // 审核状态
	StatusLog    types.String `json:"status_log,omitempty" db:"'status_log' size:65535"`      // 状态日志
	StatusTime   types.Time   `json:"status_time,omitempty" db:"'status_time'"`               // 状态时间
	Reason       types.String `json:"reason,omitempty" db:"'reason' size:200"`                // 审核不通过原因
	AuditMan     types.Money  `json:"audit_man,omitempty" db:"'audit_man' size:19"`           // 审核人id
	AuditManName types.String `json:"audit_man_name,omitempty" db:"'audit_man_name' size:20"` // 审核人姓名
	Ctime        types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                        // 创建时间
	Utime        types.Time   `json:"utime,omitempty" db:"'utime' <-"`                        //
}

// Role 角色表
// @tablename role
type Role struct {
	pool.Model
	Id    types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"` //
	Name  types.String `json:"name,omitempty" db:"'name' size:10"`     // 角色名称
	State types.Int8   `json:"state,omitempty" db:"'state' size:3"`    // 状态:1=启用
	Descr types.String `json:"descr,omitempty" db:"'descr' size:45"`   // 角色描述
	Rules types.String `json:"rules,omitempty" db:"'rules' size:500"`  // 逗号分割的权限id字符串
}

// Rule 后台权限表
// @tablename rule
type Rule struct {
	pool.Model
	Id    types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"` //
	Pid   types.BigInt `json:"pid,omitempty" db:"'pid' size:20"`       // 父级id
	Path  types.String `json:"path,omitempty" db:"'path' size:45"`     // 标识
	Title types.String `json:"title,omitempty" db:"'title' size:20"`   // 名称
	Type  types.Uint8  `json:"type,omitempty" db:"'type' size:3"`      // 前端ui类型:0=菜单,1=其它
	State types.Int8   `json:"state,omitempty" db:"'state' size:3"`    // 状态:1=启用
	Descr types.String `json:"descr,omitempty" db:"'descr' size:45"`   // 描述
}

// ShareUrlScheme 分享链接
// @tablename share_url_scheme
type ShareUrlScheme struct {
	pool.Model
	Id    types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"` //
	Peer  types.String `json:"peer,omitempty" db:"'peer' size:100"`    // 来自A或B或C
	Type  types.Uint8  `json:"type,omitempty" db:"'type' size:3"`      // 分享类型：1=微信小程序
	User  types.BigInt `json:"user,omitempty" db:"'user' size:20"`     // 用户id
	Path  types.String `json:"path,omitempty" db:"'path' size:200"`    // 路径
	Url   types.String `json:"url,omitempty" db:"'url' size:200"`      // 生成的链接
	Ctime types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`        //
}

// SignDocument 签署文件表
// @tablename sign_document
type SignDocument struct {
	pool.Model
	Id                  types.BigInt `json:"id,omitempty" db:"'id' pk size:20"`                                    // id
	DispatchedEmployee  types.BigInt `json:"dispatched_employee,omitempty" db:"'dispatched_employee' size:20"`     //  被派遣员工id
	Employee            types.BigInt `json:"employee,omitempty" db:"'employee' size:20"`                           // 派遣员工id
	User                types.BigInt `json:"user,omitempty" db:"'user' size:20"`                                   // 用户id
	Company             types.BigInt `json:"company,omitempty" db:"'company' size:20"`                             // 企业id
	Way                 types.Int8   `json:"way,omitempty" db:"'way' size:3"`                                      // 处理方式   1需要签署   2只需查看
	ExtsigTransactionId types.String `json:"extsig_transaction_id,omitempty" db:"'extsig_transaction_id' size:32"` // 手动签署交易号，长度<=32
	ExtsigName          types.String `json:"extsig_name,omitempty" db:"'extsig_name' size:20"`                     // 手动签姓名
	AuthSign            types.Int8   `json:"auth_sign,omitempty" db:"'auth_sign' size:3"`                          // 企业是否需要自动盖章   1需要
	AutoTransactionId   types.String `json:"auto_transaction_id,omitempty" db:"'auto_transaction_id' size:32"`     // 自动签署交易号，长度<=32
	AutoStampId         types.BigInt `json:"auto_stamp_id,omitempty" db:"'auto_stamp_id' size:20"`                 // 自动签署印章id
	AutoStampGenre      types.Int8   `json:"auto_stamp_genre,omitempty" db:"'auto_stamp_genre' size:3"`            // 自动签署印章类型
	AutoSignTime        types.Time   `json:"auto_sign_time,omitempty" db:"'auto_sign_time'"`                       // 自动签署时间
	DetailId            types.BigInt `json:"detail_id,omitempty" db:"'detail_id' size:20"`                         // 详情id
	DetailTable         types.String `json:"detail_table,omitempty" db:"'detail_table' size:50"`                   // 详情涉及的表名
	Name                types.String `json:"name,omitempty" db:"'name' size:200"`                                  // 文件名称
	Url                 types.String `json:"url,omitempty" db:"'url' size:65535"`                                  // 最新的文件保存地址，随合同的签约不断更新
	ViewpdfUrl          types.String `json:"viewpdf_url,omitempty" db:"'viewpdf_url' size:65535"`                  // 最新的文件查看地址，随合同的签约不断更新
	FadadaUpload        types.Int8   `json:"fadada_upload,omitempty" db:"'fadada_upload' size:3"`                  // 合同上传是否上传至法大大   1已上传
	FadadaDownloadUrl   types.String `json:"fadada_download_url,omitempty" db:"'fadada_download_url' size:65535"`  // 法大大合同下载地址
	FadadaViewpdfUrl    types.String `json:"fadada_viewpdf_url,omitempty" db:"'fadada_viewpdf_url' size:65535"`    // 法大大合同查看地址
	PdfUrl              types.String `json:"pdf_url,omitempty" db:"'pdf_url' size:200"`                            // 生成的初始合同地址
	SignatureSha        types.String `json:"signature_sha,omitempty" db:"'signature_sha' size:500"`                // sha1签名
	Document            types.BigInt `json:"document,omitempty" db:"'document' size:20"`                           // document_company表id
	TemplateId          types.BigInt `json:"template_id,omitempty" db:"'template_id' size:20"`                     // document_template表id
	Classify            types.Int8   `json:"classify,omitempty" db:"'classify' size:3"`                            // 分类
	LineType            types.Int8   `json:"line_type,omitempty" db:"'line_type' size:3"`                          // 线上线下  1线上文件   2线下文件老合同  3线下文件绩效
	Status              types.Int8   `json:"status,omitempty" db:"'status' size:3"`                                // 签署查看状态  1已签署或查看  2未签署或查看
	Result              types.String `json:"result,omitempty" db:"'result' size:20"`                               // 手动签结果
	SignTime            types.Time   `json:"sign_time,omitempty" db:"'sign_time'"`                                 // 签署时间
	AuditStatus         types.Int8   `json:"audit_status,omitempty" db:"'audit_status' size:3"`                    // 审核状态  1未审核 2待审核 3已作废  4已通过
	AuditStatusName     types.String `json:"audit_status_name,omitempty" db:"'audit_status_name' size:50"`         // 审核状态名称
	AuditMan            types.BigInt `json:"audit_man,omitempty" db:"'audit_man' size:20"`                         // 审核人id
	AuditTime           types.Time   `json:"audit_time,omitempty" db:"'audit_time'"`                               // 审核时间
	CancelReason        types.String `json:"cancel_reason,omitempty" db:"'cancel_reason' size:200"`                // 作废原因
	Filing              types.Int8   `json:"filing,omitempty" db:"'filing' size:3"`                                // 法大大合同归档  1已归档
	State               types.Int8   `json:"state,omitempty" db:"'state' size:3"`                                  // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime               types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                                      // 创建时间
}

// SignDocumentVariable 文件协议模板变量
// @tablename sign_document_variable
type SignDocumentVariable struct {
	pool.Model
	Id       types.BigInt `json:"id,omitempty" db:"'id' pk size:20"`          // id
	Template types.BigInt `json:"template,omitempty" db:"'template' size:20"` // document_template表id
	Value    types.String `json:"value,omitempty" db:"'value' size:65535"`    // 变量json串
	State    types.Int8   `json:"state,omitempty" db:"'state' size:3"`        // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime    types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`            // 创建时间
	Utime    types.Time   `json:"utime,omitempty" db:"'utime' <-"`            // 最后修改时间
}

// Sms
// @tablename sms
type Sms struct {
	pool.Model
	Id        types.Money  `json:"id,omitempty" db:"'id' pk auto size:19"`         //
	Phone     types.BigInt `json:"phone,omitempty" db:"'phone' size:20"`           //
	Para      types.String `json:"para,omitempty" db:"'para' size:65535"`          // 模板变量
	SmsCode   types.String `json:"sms_code,omitempty" db:"'sms_code' size:20"`     // 短信模板
	RequestId types.String `json:"request_id,omitempty" db:"'request_id' size:50"` //
	Code      types.String `json:"code,omitempty" db:"'code' size:50"`             //
	Message   types.String `json:"message,omitempty" db:"'message' size:100"`      //
	Response  types.String `json:"response,omitempty" db:"'response' size:300"`    //
	Ctime     types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                //
}

// User 用户表
// @tablename user
type User struct {
	pool.Model
	Id                  types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`                                //
	Name                types.String `json:"name,omitempty" db:"'name' size:45"`                                    // 姓名
	Email               types.String `json:"email,omitempty" db:"'email' size:64"`                                  // 邮箱
	Mobile              types.String `json:"mobile,omitempty" db:"'mobile' size:20"`                                // 手机号
	Gender              types.String `json:"gender,omitempty" db:"'gender' size:10"`                                // 性别
	Birthday            types.String `json:"birthday,omitempty" db:"'birthday' size:25"`                            // 出生日期
	Height              types.String `json:"height,omitempty" db:"'height' size:10"`                                // 身高
	Weight              types.String `json:"weight,omitempty" db:"'weight' size:10"`                                // 体重
	Race                types.String `json:"race,omitempty" db:"'race' size:5"`                                     // 民族
	HometownAddressNorm types.String `json:"hometown_address_norm,omitempty" db:"'hometown_address_norm' size:100"` // 籍贯地址（规范化）
	PolitStatus         types.String `json:"polit_status,omitempty" db:"'polit_status' size:25"`                    // 政治面貌
	MaritalStatus       types.String `json:"marital_status,omitempty" db:"'marital_status' size:10"`                // 婚姻状态
	IdCard              types.String `json:"id_card,omitempty" db:"'id_card' size:50"`                              // 身份证号
	HukouAddressNorm    types.String `json:"hukou_address_norm,omitempty" db:"'hukou_address_norm' size:100"`       // 户口地址（规范化）
	HukouType           types.String `json:"hukou_type,omitempty" db:"'hukou_type' size:20"`                        // 户口性质
	LivingAddress       types.String `json:"living_address,omitempty" db:"'living_address' size:200"`               // 当前所在地
	Degree              types.String `json:"degree,omitempty" db:"'degree' size:10"`                                // 学历
	College             types.String `json:"college,omitempty" db:"'college' size:50"`                              // 毕业学校
	CollegeMobile       types.String `json:"college_mobile,omitempty" db:"'college_mobile' size:30"`                // 学校联系电话
	Major               types.String `json:"major,omitempty" db:"'major' size:25"`                                  // 所学专业
	GradTime            types.String `json:"grad_time,omitempty" db:"'grad_time' size:50"`                          // 毕业时间
	ProfessionalTitle   types.String `json:"professional_title,omitempty" db:"'professional_title' size:100"`       // 职称
	EmergencyContact    types.String `json:"emergency_contact,omitempty" db:"'emergency_contact' size:20"`          // 紧急联系人姓名
	EmergencyRelation   types.String `json:"emergency_relation,omitempty" db:"'emergency_relation' size:20"`        // 紧急联系人关系
	EmergencyMobile     types.String `json:"emergency_mobile,omitempty" db:"'emergency_mobile' size:20"`            // 紧急联系人电话
	Creator             types.BigInt `json:"creator,omitempty" db:"'creator' size:20"`                              // 创建人
	Status              types.Int8   `json:"status,omitempty" db:"'status' size:3"`                                 // 状态
	State               types.Int8   `json:"state,omitempty" db:"'state' size:3"`                                   // 管理状态：-1=已删除，0=禁用，1=启用
	Ctime               types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                                       //
	Utime               types.Time   `json:"utime,omitempty" db:"'utime' <-"`                                       //
}

// UserApp 用户应用表
// @tablename user_app
type UserApp struct {
	pool.Model
	Id            types.BigInt `json:"id,omitempty" db:"'id' pk auto size:20"`           //
	User          types.BigInt `json:"user,omitempty" db:"'user' size:20"`               //
	App           types.Int8   `json:"app,omitempty" db:"'app' size:3"`                  // 10=松果派
	OpenId        types.String `json:"open_id,omitempty" db:"'open_id' size:45"`         // openid
	UnionId       types.String `json:"union_id,omitempty" db:"'union_id' size:45"`       //
	SessionKey    types.String `json:"session_key,omitempty" db:"'session_key' size:45"` //
	Agreement     types.Uint8  `json:"agreement,omitempty" db:"'agreement' size:3"`      // 是否同意用户协议;1=同意
	AgreementDate types.Time   `json:"agreement_date,omitempty" db:"'agreement_date'"`   // 同意协议时间
	LoginTime     types.Time   `json:"login_time,omitempty" db:"'login_time' <-"`        // 最后登录时间
	State         types.Int8   `json:"state,omitempty" db:"'state' size:3"`              // 状态
	Ctime         types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                  // 注册时间
}

// UserFadada 法大大用户表
// @tablename user_fadada
type UserFadada struct {
	pool.Model
	Id             types.BigInt `json:"id,omitempty" db:"'id' pk size:20"`                           // user表id
	Name           types.String `json:"name,omitempty" db:"'name' size:100"`                         // 员工姓名
	CustomerId     types.String `json:"customer_id,omitempty" db:"'customer_id' size:100"`           // 客户编号
	TransactionNo  types.String `json:"transaction_no,omitempty" db:"'transaction_no' size:50"`      // 交易号（需要保存，用于证书申请和实名认证查询）
	Url            types.String `json:"url,omitempty" db:"'url' size:65535"`                         // 地址（需要保存，遇到中途退出认证或页面过期等情况可重新访问
	CertInfo       types.String `json:"cert_info,omitempty" db:"'cert_info' size:4294967295"`        // 实名认证信息
	Status         types.Int8   `json:"status,omitempty" db:"'status' size:3"`                       // 0未激活   1未认证  2审核通过  3已提交待审核   4审核不通过
	HasCertificate types.Int8   `json:"has_certificate,omitempty" db:"'has_certificate' size:3"`     // 是否申请了实名证书   1已申请
	PersonName     types.String `json:"person_name,omitempty" db:"'person_name' size:20"`            // 个人姓名
	IdCard         types.String `json:"id_card,omitempty" db:"'id_card' size:50"`                    // 身份证号
	Mobile         types.String `json:"mobile,omitempty" db:"'mobile' size:20"`                      // 手机号
	HeadPhotoPath  types.String `json:"head_photo_path,omitempty" db:"'head_photo_path' size:65535"` // 身份证正面图片uuid
	BackPhotoPath  types.String `json:"back_photo_path,omitempty" db:"'back_photo_path' size:65535"` // 身份证反面图片uuid
	Photo          types.String `json:"photo,omitempty" db:"'photo' size:65535"`                     // 腾讯云返回的照片uuid
	State          types.Int8   `json:"state,omitempty" db:"'state' size:3"`                         // 状态
	Ctime          types.Time   `json:"ctime,omitempty" db:"'ctime' <-"`                             // 注册时间
}
