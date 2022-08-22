package errmsg

const (
	SUCCSE = 200
	ERROR  = 500

	// captcha= 1000... 用户模块的错误
	ERROR_USERNAME_USED  = 1001
	ERROR_PASSWORD_WRONG = 1002
	ERROR_USER_NOT_EXIST = 1003
	ERROR_TOKEN_EXIST    = 1004
	ERROR_TOKEN_RUNTIME  = 1005
	ERROR_TOKEN_WRONG    = 1006

	// token的格式错误，第一个字符不是Bearer
	ERROR_TOKEN_TYPE_WRONG = 1007
	ERROR_USER_NO_RIGHT    = 1008

	// captcha= 2000... 文章模块的错误
	ERROR_ART_NOT_EXIST = 2001
	// captcha= 3000... 分类模块的错误

	ERROR_CATENAME_USED  = 3001
	ERROR_CATE_NOT_EXIST = 3002
	// captcha=4000... 标签模块的错误
	ERROR_TAG_USED = 4001
	// capthca=5000 验证码模块的错误
	ERROR_CAPTCHA_RUNTIME   = 5001
	ERROR_CAPTCHA_INCORRECT = 5002
	// 前端用户模块的错误
	ERROR_SESSION_RUNTIME    = 6001
	ERROR_SESSION_KEY        = 6002
	ERROR_SESSION_REMOVE     = 6003
	ERROR_USERNAME_NOT_EXIST = 6004
)

var codeMsg = map[int]string{
	SUCCSE:                 "OK",
	ERROR:                  "FAIL",
	ERROR_USERNAME_USED:    "用户名已存在！",
	ERROR_PASSWORD_WRONG:   "密码错误",
	ERROR_USER_NOT_EXIST:   "用户不存在",
	ERROR_TOKEN_EXIST:      "TOKEN不存在,请重新登陆",
	ERROR_TOKEN_RUNTIME:    "TOKEN已过期,请重新登陆",
	ERROR_TOKEN_WRONG:      "TOKEN不正确,请重新登陆",
	ERROR_TOKEN_TYPE_WRONG: "TOKEN格式错误,请重新登陆",
	ERROR_USER_NO_RIGHT:    "该用户无权限",

	ERROR_ART_NOT_EXIST: "文章不存在",

	ERROR_CATENAME_USED:  "该分类已存在",
	ERROR_CATE_NOT_EXIST: "该分类不存在",

	ERROR_TAG_USED: "该标签已经被使用",

	ERROR_CAPTCHA_RUNTIME:   "验证码已过期，请重新获取",
	ERROR_CAPTCHA_INCORRECT: "验证码不正确，请重新输入",

	ERROR_SESSION_RUNTIME:    "会话信息失效，请重新登录",
	ERROR_SESSION_KEY:        "用户ID与会话信息ID不匹配，请检查后重新提交",
	ERROR_SESSION_REMOVE:     "删除会话信息出错",
	ERROR_USERNAME_NOT_EXIST: "修改密码用户不存在，请重新注册",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}
