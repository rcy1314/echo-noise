package models

// 成功相关
const (
	SuccessMessage               = "请求成功"
    GetAllMessagesSuccess        = "获取信息成功"
    GetMessageByIDSuccess        = "获取信息成功"
    PostMessageSuccess           = "发布信息成功"
	CreateSuccessMessage         = "创建成功"
	QuerySuccessMessage          = "查询成功"
	DeleteSuccessMessage         = "删除成功"
	ServerLaunchSuccessMessage   = "服务器启动成功"
    GetMessagesByPageSuccess     = "分页获取信息成功"
	RegisterSuccessMessage       = "注册成功"
	LoginSuccessMessage          = "登录成功"
	UpdateUserSuccessMessage     = "更新用户成功"
	ChangePasswordSuccessMessage = "修改密码成功"
	UpdateSettingSuccessMessage = "系统设置更新成功"
)

// 失败相关
const (
	ErrorMessage                           = "请求失败"
	CreateFailMessage                      = "创建失败"
    GetAllMessagesFailMessage              = "获取信息失败"
    GetMessageByIDFailMessage              = "获取信息失败"
    MessageNotFoundMessage                 = "信息未找到"
    PostMessageFailMessage                 = "发布信息失败"
	DeleteFailMessage                      = "删除失败"
	QueryFailMessage                       = "查询失败"
	InvalidIDMessage                       = "无效的ID"
	NotFoundMessage                        = "资源未找到"
	ValidationErrorMessage                 = "数据验证失败"
	InternalServerErrorMessage             = "内部服务器错误"
	ServerLaunchErrorMessage               = "服务器启动失败"
	LoadConfigErrorMessage                 = "加载配置失败"
	ParseConfigErrorMessage                = "解析配置失败"
	InvalidRequestMessage                  = "无效的请求"
	InvalidRequestBodyMessage              = "无效的请求体"
	ImageUploadErrorMessage                = "图片上传失败"
	NotUploadImageErrorMessage             = "未上传图片"
	NotSupportedImageTypeMessage           = "不支持的图片类型"
	ImageSizeLimitErrorMessage             = "图片大小超过限制"
	InvalidPageParametersMessage           = "无效的分页参数"
	InvalidPageSizeMessage                 = "无效的页码或页大小"
    GetMessagesByPageFailMessage           = "分页获取信息失败"
	TokenNotFoundMessage                   = "未找到令牌,请点击右上角登录"
	TokenInvalidMessage                    = "令牌无效，请点击右上角登录"
	TokenInvalidFormatMessage              = "令牌格式无效" 
	UsernameOrPasswordCannotBeEmptyMessage = "用户名或密码不能为空"
	UsernameCannotBeEmptyMessage           = "用户名不能为空"
	PasswordCannotBeEmptyMessage           = "密码不能为空"
	UsernameAlreadyExistsMessage           = "用户名已存在"
	CreateUserFailMessage                  = "创建用户失败"
	UserNotFoundMessage                    = "用户未找到"
	PasswordIncorrectMessage               = "密码错误"
	GenerateTokenFailMessage               = "生成令牌失败"
	GetAllUsersFailMessage                 = "获取所有用户失败"
	GetStatusFailMessage                   = "获取状态失败"
	GetStatusSuccessMessage                = "获取状态成功"
	CannotBeEmptyMessage                   = "内容不能为空"
	NoPermissionMessage                    = "没有权限,请联系管理员"
	GenerateRSSFailMessage                 = "生成 RSS 失败"
	PasswordCannotBeSameAsBeforeMessage    = "新密码不能与旧密码相同"
)

// 数据库相关
const (
	DatabaseTypeMessage       = "数据库类型错误"
	DatabaseConnectionError   = "数据库连接失败"
	DatabaseMigrationError    = "数据库迁移失败"
	DatabaseInitErrorMessage  = "数据库初始化失败"
	DatabaseErrorMessage      = "数据库操作失败"
	DatabaseConnectionSuccess = "数据库连接成功"
	DuplicateEntryMessage     = "重复的条目"
	RecordNotFoundMessage     = "记录未找到"
)

// 验证相关
const (
	InvalidDataMessage  = "无效数据"
	MissingFieldMessage = "缺少必填字段"
)
