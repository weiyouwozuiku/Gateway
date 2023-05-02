package middleware

type Error struct {
	errno  ResponseCode
	errmsg string
}
type ResponseCode int

const (
	SuccessCode  ResponseCode = iota
	innerErrCode              = 1000 + iota
	undefErrCode
	validErrCode
	paramErrCode
	dBErrCode
	getGormPoolFailed
	gROUPALL_SAVE_FLOWERROR
	adminLoginFailed
	sessionParseFailed
	gormQueryFailed
	gormSaveFailed
	sessionOptFailed
)

var (
	InnerErr        = Error{errno: innerErrCode, errmsg: "内部错误"}
	ParamErr        = Error{errno: paramErrCode, errmsg: "参数错误"}
	DBErr           = Error{errno: dBErrCode, errmsg: "数据库错误"}
	ValidErr        = Error{errno: validErrCode, errmsg: "参数不合法"}
	GetGormPoolErr  = Error{errno: getGormPoolFailed, errmsg: "获取gorm失败"}
	QueryGormErr    = Error{errno: gormQueryFailed, errmsg: "查询gorm失败"}
	SaveGormErr     = Error{errno: gormSaveFailed, errmsg: "更新数据库失败"}
	AdminLoginErr   = Error{errno: adminLoginFailed, errmsg: "管理员登录失败"}
	SessionParseErr = Error{errno: sessionParseFailed, errmsg: "session解析失败"}
	SessionOptErr   = Error{errno: sessionOptFailed, errmsg: "session操作失败"}
)
