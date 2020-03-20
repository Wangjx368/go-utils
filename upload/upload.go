package upload

import (
	"encoding/base64"
	"encoding/json"
	"openapiai/utils/encrypt"
	"openapiai/utils/request"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type tokenResp struct {
	Stat   int           `json:"stat"`
	Errmsg string        `json:"errmsg"`
	Data   []interface{} `json:"data"`
}

type tokenRespDataOssObj struct {
	Type          string            `json:"type"`
	Method        string            `json:"method"`
	Host          string            `json:"host"`
	VpcHost       string            `json:"vpc_host"`
	FileUrl       string            `json:"file_url"`
	RequestHeader map[string]string `json:"request_header"`
}

const (
	UploadApi = ""
	RootDir   = ""
)

var (
	appId  = beego.AppConfig.String("upload_app_id")
	appKey = beego.AppConfig.String("upload_app_key")
)

func GenerateUploadSign() string {
	key := appKey                                                          //密钥串
	time := strconv.FormatInt(time.Now().Unix(), 10)                       //当前UNIX时间戳
	cipher := encrypt.Md5s(key + "&" + time)                               //生成加密串
	sign := base64.StdEncoding.EncodeToString([]byte(cipher + "&" + time)) //生成sign
	return sign
}

func GetUploadToken(fileName string) (tr tokenResp, err error) {
	tr = tokenResp{}
	params := make(map[string]string)
	headers := make(map[string]string)
	params["dst_path"] = RootDir + "/" + time.Now().Format(("2006-01-02")) + "/" + fileName
	headers["APPID"] = appId
	headers["SIGN"] = GenerateUploadSign()
	req := request.NewRequest()
	respByte, err := req.DoHttpGet(UploadApi, params, headers)
	if err != nil {
		logs.Error(err.Error())
		return tr, err
	}

	json.Unmarshal(respByte, &tr)
	logs.Debug(tr)

	return tr, nil
}

func DoUploadFile(fileBytes []byte, ext string) (fileUrl string) {
	fileUrl = ""
	fileName := encrypt.Md5s(string(fileBytes)) + "." + ext
	tr, err := GetUploadToken(fileName)
	if err != nil {
		logs.Error(err.Error())
		return
	}

	if tr.Stat == 1 {
		ossObj := tr.Data[0].(map[string]interface{})
		_type := ossObj["type"].(string)
		if _type == "OSS" {
			host := ossObj["host"].(string)
			fileUrl = ossObj["file_url"].(string)
			headers := ossObj["request_header"].(map[string]interface{})
			header := make(map[string]string)
			for k, v := range headers {
				header[k] = v.(string)
			}
			req := request.NewRequest()
			_, err := req.DoHttpPut(host, header, fileBytes)
			if err != nil {
				logs.Error(err.Error())
				fileUrl = ""
				return
			}
		}
	}
	return
}
