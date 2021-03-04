package errcode

var (
	// 成功
	Success = New("00000", "Success")

	// A 组
	// 客户端错误
	ClientError        = New("A0001", "用户端错误")
	InvalidParamsError = New("A0002", "请求参数错误")

	// 鉴权错误
	UserPasswordError = New("A0100", "用户密码不正确")
	// token
	AuthNotExistError      = New("A0101", "鉴权失败, 找不到对应的 AppKey 和 AppSecret")
	AuthTokenGenerateError = New("A0102", "Token 生成失败")
	AuthTokenTimeout       = New("A0103", "Token 超时")
	AuthTokenError         = New("A0104", "Token 错误")
	AuthTokenParseError    = New("A0105", "Token 解析失败")
	AuthTokenNotObtained   = New("A0106", "未获取到 token")
	// cookie session
	CookieSessionError = New("A0107", "CookieSession 错误")
	CreateSessionError = New("A0108", "创建 Session 错误")
	ClearSessionError  = New("A0109", "删除 Session 错误")

	// B 组
	// 服务端错误
	ServerError = New("B0001", "系统执行出错")

	// 标签错误
	CreateTagError = New("B0100", "创建标签失败")
	EditTagError   = New("B0101", "编辑标签失败")
	DeleteTagError = New("B0102", "删除标签失败")
	GetTagError    = New("B0103", "获取标签失败")

	// 用户错误
	CreateUserError = New("B0104", "创建用户失败")
	EditUserError   = New("B0105", "编辑用户失败")
	DeleteUserError = New("B0106", "删除用户失败")
	GetUserError    = New("B0107", "获取用户失败")

	// 上传文件错误
	UploadFileError = New("B0200", "上传文件失败")

	// C 组
	// 第三方调用错误
	ThirdPartyCallError = New("C0001", "第三方调用错误")

	//NotFound                  = New(10000002, "找不到")

	//TooManyRequests           = New(10000007, "请求过多")
)
