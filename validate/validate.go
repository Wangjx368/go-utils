package validate

import (
	"github.com/astaxie/beego/logs"
)

const (
	SidMinLen = 16
	SidMaxLen = 64

	UrlMaxLen = 255
)

func ValidateSid(sid string) bool {
	if sid == "" {
		return false
	}

	runes := []rune(sid)
	if len(runes) < SidMinLen {
		logs.Error("sid:%s is too short\n", sid)
		return false
	}

	if len(runes) > SidMaxLen {
		logs.Error("sid:%s is too long\n", sid)
		return false
	}

	return true
}

func ValidateURL(url string) bool {
	if url == "" {
		return false
	}

	runes := []rune(url)
	if len(runes) > UrlMaxLen {
		logs.Error("post url:%s is too long\n", url)
		return false
	}
	return true
}
