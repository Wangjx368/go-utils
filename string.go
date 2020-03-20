package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func interface2String(inter interface{}) (res string) {
	res = ""
	switch inter.(type) {

	case string:
		res = inter.(string)
		break
	case int:
		res = string(inter.(int))
		break
	case float64:
		fmt.Println("float64", inter.(float64))
		res = ""
		break
	}
	return res
}

func B2S(bs []uint8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		b[i] = byte(v)
	}
	return string(b)
}

func InArray(str string, strs []string) int {
	for k, v := range strs {
		if v == str {
			return k
		}
	}
	return -1
}

func InArrayInt(str int, strs []int) int {
	for k, v := range strs {
		if v == str {
			return k
		}
	}
	return -1
}

func Between(str string, strTmp string, starting string, ending string) (string, int) {
	s := strings.Index(strTmp, starting)
	if s < 0 {
		return "", len(str)
	}
	s += len(starting)
	e := strings.Index(strTmp[s:], ending)
	if e < 0 {
		return "", len(str)
	}
	return strTmp[s : s+e], s + e + len(str) - len(strTmp) + 1
}

func RandNum(start int, end int) int {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(end-start+1) + start
	return num
}

func MatchReplace(str string, s1 string) string {
	reg, err := regexp.Compile(s1)
	if err != nil {
		return ""
	}
	temps := reg.FindAllString(str, -1)
	stemTmp := str
	for _, v := range temps {
		stemTmp = strings.Replace(stemTmp, v, " ", -1)
	}
	return stemTmp
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
