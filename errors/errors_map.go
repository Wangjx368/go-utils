package errors

var (
	// ErrorsMap 错误码表
	ErrorsMap = map[int]string{
		-14000: "添加数据库失败",
		-14001: "查询数据失败",

		0: "请求成功",

		14000: "参数非法",
		14001: "缺少app_key参数",
	}
)
