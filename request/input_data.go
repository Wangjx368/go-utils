package request

import (
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	_ "github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net"
	"openapiai/models"
	"openapiai/utils/meta"
	"openapiai/utils/wslock"
	"strings"
	"time"
)

func GetClientIP(ctx *context.Context) string {
	ipStr := ctx.Request.Header.Get("X-Forwarded-For")
	ips := strings.Split(ipStr, ",")
	ip := ips[0]
	if strings.Contains(ip, "127.0.0.1") || ip == "" {
		ip = ctx.Request.Header.Get("X-Real-IP")
	}

	if ip == "" {
		return "127.0.0.1"
	}

	return ip
}

func GetIntranetIp() (ip string) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {

	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return
}

func GetIPLastSeq(ip string) (lastIpSeq string) {
	var lastIndex = 1
	var seqs []string = strings.Split(ip, ".")
	_len := len(seqs)
	if _len > 0 {
		_len = _len - 1
	}
	lastIpSeq = seqs[lastIndex]
	return
}

// app_key
func SetInputDataAppKey(ctx *context.Context, appKey string) {
	ctx.Input.SetData("appKey", appKey)
}

func GetInputDataAppKey(ctx *context.Context) string {
	return ctx.Input.GetData("appKey").(string)
}

// user
func SetInputDataUser(ctx *context.Context, user models.AiUser) {
	ctx.Input.SetData("user", user)
}

func GetInputDataUser(ctx *context.Context) models.AiUser {
	return ctx.Input.GetData("user").(models.AiUser)
}

// validSign
func SetInputDataValidSign(ctx *context.Context, validSign bool) {
	ctx.Input.SetData("validSign", validSign)
}

func GetInputDataValidSign(ctx *context.Context) bool {
	return ctx.Input.GetData("validSign").(bool)
}

// signParams
func SetInputDataSignParams(ctx *context.Context, signParams map[string]string) {
	ctx.Input.SetData("signParams", signParams)
}

func GetInputDataSignParams(ctx *context.Context) map[string]string {
	return ctx.Input.GetData("signParams").(map[string]string)
}

// signExtraParams
func SetInputDataSignExtraParams(ctx *context.Context, signExtraParams map[string]string) {
	ctx.Input.SetData("signExtraParams", signExtraParams)
}

func GetInputDataSignExtraParams(ctx *context.Context) map[string]string {
	return ctx.Input.GetData("signExtraParams").(map[string]string)
}

// sign
func SetInputDataSign(ctx *context.Context, sign string) {
	ctx.Input.SetData("sign", sign)
}

func GetInputDataSign(ctx *context.Context) string {
	return ctx.Input.GetData("sign").(string)
}

// startTime
func SetInputDataStartTime(ctx *context.Context) {
	ctx.Input.SetData("startTime", time.Now().UnixNano()/int64(time.Millisecond))
}

func GetInputDataStartTime(ctx *context.Context) int64 {
	return ctx.Input.GetData("startTime").(int64)
}

// respCxt
func SetInputDataRespContext(ctx *context.Context, rc meta.ResponseContext) {
	ctx.Input.SetData("respCxt", rc)
}

func GetInputDataRespContext(ctx *context.Context) meta.ResponseContext {
	return ctx.Input.GetData("respCxt").(meta.ResponseContext)
}

// traceId
func SetInputDataTraceId(ctx *context.Context) {
	traceId, err := uuid.NewV4()
	if err != nil {
		logs.Error(err)
	}
	ctx.Input.SetData("traceId", traceId.String())
}

func GetInputDataTraceId(ctx *context.Context) string {
	return ctx.Input.GetData("traceId").(string)
}

// rpcId
func SetInputDataRpcId(ctx *context.Context) {
	rpcId, err := uuid.NewV4()
	if err != nil {
		logs.Error(err)
	}
	ctx.Input.SetData("rpcId", rpcId.String())
}

func GetInputDataRpcId(ctx *context.Context) string {
	return ctx.Input.GetData("rpcId").(string)
}

// ws conn
func SetInputDataWsConn(ctx *context.Context, wsConn *wslock.WSer) {
	ctx.Input.SetData("wsConn", wsConn)
}

func GetInputDataWsConn(ctx *context.Context) *wslock.WSer {
	return ctx.Input.GetData("wsConn").(*wslock.WSer)
}
