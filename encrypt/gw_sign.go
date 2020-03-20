package encrypt

import (
	"strconv"
	"time"
)

func GetQwSignHeader(header map[string]string, appId string, appKey string) {
	msTimestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	header["X-Auth-Appid"] = appId
	header["X-Auth-TimeStamp"] = msTimestamp
	header["X-Auth-Sign"] = Md5s(appId + "&" + msTimestamp + appKey)
}
