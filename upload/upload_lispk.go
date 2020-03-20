package upload

import (
	"openapiai/utils/encrypt"
	"openapiai/utils/request"

	"github.com/astaxie/beego/logs"
)

func DoUploadFileLispk(fileBytes []byte, ext string) (fileUrl string, err error) {
	fileUrl = ""
	fileName := encrypt.Md5s(string(fileBytes)) + "." + ext
	tr, err := GetUploadToken(fileName)
	if err != nil {
		logs.Error(err.Error())
		return fileUrl, err
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
				return fileUrl, err
			}
		}
	}
	return fileUrl, err
}
