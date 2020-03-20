package request

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	netUrl "net/url"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"openapiai/utils/encrypt"
)

const (
	Boundary            = "PP**ASR**LIB"
	DecodeServerTimeout = 30
)

const (
	proxyUrl = "http://127.0.0.1:8888/"
	useProxy = false
)

type Request struct {
	Boundary            string
	DecodeServerTimeout int
}

func NewRequest() *Request {
	return &Request{
		DecodeServerTimeout: 10,
	}
}

func (r *Request) init() {
	if r.DecodeServerTimeout == 0 {
		r.DecodeServerTimeout = 10
	}
}

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Duration(2*time.Second))
}

func GetBaseHeader() map[string]string {
	header := make(map[string]string)
	return header
}

func GetGwSignHeader() map[string]string {
	header := make(map[string]string)
	appId := beego.AppConfig.String("api_gw_app_id")
	appKey := beego.AppConfig.String("api_gw_app_key")
	encrypt.GetQwSignHeader(header, appId, appKey)
	return header
}

func (r *Request) PostDecoder(url string, header map[string]string, req_bytes []byte) (b []byte, err error) {
	send_req, err := http.NewRequest("POST", url, bytes.NewBuffer(req_bytes))
	if err != nil {
		return b, err
	}
	send_req.Close = true

	// header
	send_req.Header.Add("Content-Type", "multipart/form-data; boundary="+r.Boundary)
	send_req.Header.Add("Content-Length", fmt.Sprintf("%d", len(req_bytes)))
	for k, v := range header {
		send_req.Header.Add(k, v)
	}

	// http setting
	tr := &http.Transport{
		Dial:                  dialTimeout,
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: time.Second * time.Duration(r.DecodeServerTimeout),
	}
	httpClient := &http.Client{Transport: tr}

	// do http
	//logs.Debug(send_req)
	resp, err := httpClient.Do(send_req)
	if err != nil {
		return b, err
	} else {
		defer resp.Body.Close()
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logs.Error(err)
			return b, err
		}
		logs.Notice("URL:", url, "Req:", string(req_bytes), "Resp:", string(b))
	}
	if resp.StatusCode >= 400 {
		err = errors.New(resp.Status)
	}
	return b, err
}

func (r *Request) DoHttpPost(url string, req_bytes []byte) (b []byte, err error) {
	logs.Info(string(req_bytes))
	send_req, err := http.NewRequest("POST", url, bytes.NewBuffer(req_bytes))
	if err != nil {
		return b, err
	}
	send_req.Close = true
	send_req.Header.Add("Content-Type", "application/json")
	send_req.Header.Add("Content-Length", fmt.Sprintf("%d", len(req_bytes)))
	tr := &http.Transport{
		Dial:                  dialTimeout,
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: time.Second * time.Duration(r.DecodeServerTimeout),
	}
	httpClient := &http.Client{Transport: tr}
	//logs.Debug(send_req)
	resp, err := httpClient.Do(send_req)
	if err != nil {
		return b, err
	} else {
		defer resp.Body.Close()
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logs.Error(err)
			return b, err
		}
		logs.Notice("URL:", url, "Req:", string(req_bytes), "Resp:", string(b))
	}
	if resp.StatusCode >= 400 {
		err = errors.New(resp.Status)
	}
	return b, err
}

func (r *Request) DoHttpPostWithHeaders(url string, header map[string]string, req_bytes []byte) (b []byte, err error) {
	logs.Info(string(req_bytes))
	send_req, err := http.NewRequest("POST", url, bytes.NewBuffer(req_bytes))
	if err != nil {
		return b, err
	}
	send_req.Close = true

	// header
	send_req.Header.Add("Content-Type", "application/json")
	send_req.Header.Add("Content-Length", fmt.Sprintf("%d", len(req_bytes)))

	for k, v := range header {
		send_req.Header.Add(k, v)
	}

	tr := &http.Transport{
		Dial:                  dialTimeout,
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: time.Second * time.Duration(r.DecodeServerTimeout),
	}

	if useProxy {
		proxy, err := netUrl.Parse(proxyUrl)
		if err != nil {
			logs.Error(err)
		}
		tr.Proxy = http.ProxyURL(proxy)
	}

	httpClient := &http.Client{Transport: tr}
	logs.Debug(send_req)
	resp, err := httpClient.Do(send_req)
	if err != nil {
		logs.Notice("URL:", url, "Req:", string(req_bytes), "Resp:", string(b), "Error:", err)
		return b, err
	} else {
		defer resp.Body.Close()
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logs.Error(err)
			return b, err
		}
		logs.Notice("URL:", url, "Req:", string(req_bytes), "Resp:", string(b))
	}

	if resp.StatusCode >= 400 {
		err = errors.New(resp.Status)
	}
	return b, err
}

func (r *Request) DoHttpPostWithHeadersAndFormUrlEncode(url string, header map[string]string, param map[string]string) (b []byte, err error) {
	logs.Info(param)

	// var params []string
	params := netUrl.Values{}

	for k, v := range param {
		params.Add(k, v)
	}

	paramStr := params.Encode()

	send_req, err := http.NewRequest("POST", url, strings.NewReader(paramStr))
	if err != nil {
		return b, err
	}
	send_req.Close = true

	send_req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range header {
		send_req.Header.Add(k, v)
	}

	tr := &http.Transport{
		Dial:                  dialTimeout,
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: time.Second * time.Duration(r.DecodeServerTimeout),
	}

	if useProxy {
		proxy, err := netUrl.Parse(proxyUrl)
		if err != nil {
			logs.Error(err)
		}
		tr.Proxy = http.ProxyURL(proxy)
	}

	httpClient := &http.Client{Transport: tr}
	logs.Debug(send_req)
	resp, err := httpClient.Do(send_req)
	if err != nil {
		return b, err
	} else {
		defer resp.Body.Close()
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logs.Error(err)
			return b, err
		}
		logs.Notice("URL:", url, "Req:", paramStr, "Resp:", string(b))
	}
	if resp.StatusCode >= 400 {
		err = errors.New(resp.Status)
	}
	return b, err
}

func (r *Request) DoHttpPostByMultipart(url string, header map[string]string, req io.Reader, contentLength int64) (b []byte, err error) {
	send_req, err := http.NewRequest("POST", url, req)
	if err != nil {
		return b, err
	}
	send_req.ContentLength = contentLength
	send_req.Close = true
	for k, v := range header {
		send_req.Header.Add(k, v)
	}

	tr := &http.Transport{
		Dial:                  dialTimeout,
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: time.Second * time.Duration(r.DecodeServerTimeout),
	}

	if useProxy {
		proxy, err := netUrl.Parse(proxyUrl)
		if err != nil {
			logs.Error(err)
		}
		tr.Proxy = http.ProxyURL(proxy)
	}

	httpClient := &http.Client{Transport: tr}
	logs.Debug(send_req)
	resp, err := httpClient.Do(send_req)
	if err != nil {
		return b, err
	} else {
		defer resp.Body.Close()
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logs.Error(err)
			return b, err
		}
		logs.Debug(string(b))
		logs.Notice("URL:", url, "Req:", req, "Resp:", string(b))
	}
	return b, err
}

func (r *Request) DoHttpGet(url string, params map[string]string, header map[string]string) (b []byte, err error) {
	for k, v := range params {
		url += k + "=" + v + "&"
	}
	send_req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return b, err
	}
	send_req.Close = true
	for k, v := range header {
		send_req.Header.Add(k, v)
	}
	tr := &http.Transport{
		Dial:                  dialTimeout,
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: time.Second * time.Duration(r.DecodeServerTimeout),
	}

	if useProxy {
		proxy, err := netUrl.Parse(proxyUrl)
		if err != nil {
			logs.Error(err)
		}
		tr.Proxy = http.ProxyURL(proxy)
	}

	httpClient := &http.Client{Transport: tr}
	logs.Debug(send_req)
	resp, err := httpClient.Do(send_req)
	if err != nil {
		return b, err
	} else {
		defer resp.Body.Close()
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logs.Error(err)
			return b, err
		}
		logs.Notice("URL:", url, "Params:", params, "Resp:", string(b))
	}
	if resp.StatusCode >= 400 {
		err = errors.New(resp.Status)
	}
	return b, err
}

func (r *Request) DoHttpPut(url string, header map[string]string, req_bytes []byte) (b []byte, err error) {
	send_req, err := http.NewRequest("PUT", url, bytes.NewBuffer(req_bytes))
	if err != nil {
		return b, err
	}
	send_req.Close = true
	for k, v := range header {
		send_req.Header.Add(k, v)
	}
	tr := &http.Transport{
		Dial:                  dialTimeout,
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: time.Second * time.Duration(r.DecodeServerTimeout),
	}

	if useProxy {
		proxy, err := netUrl.Parse(proxyUrl)
		if err != nil {
			logs.Error(err)
		}
		tr.Proxy = http.ProxyURL(proxy)
	}

	httpClient := &http.Client{Transport: tr}
	//logs.Debug(send_req)
	resp, err := httpClient.Do(send_req)
	if err != nil {
		return b, err
	} else {
		defer resp.Body.Close()
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logs.Error(err)
			return b, err
		}
		logs.Notice("URL:", url, "Header:", header, "Resp:", string(b))
	}
	if resp.StatusCode >= 400 {
		err = errors.New(resp.Status)
	}
	return b, err
}

func GetReqBytes(boundary string, jsonPartBytes []byte, contentBytes []byte) []byte {
	BOUNDARY := boundary
	a_boundary := "\r\n--" + BOUNDARY + "_a" + "\r\n"
	b_boundary := "\r\n--" + BOUNDARY + "_b" + "\r\n"
	c_boundary := "\r\n--" + BOUNDARY + "_c" + "\r\n"

	buffer := bytes.Buffer{}
	buffer.WriteString(a_boundary)
	buffer.Write(jsonPartBytes)
	buffer.WriteString(b_boundary)
	buffer.Write(contentBytes)
	buffer.WriteString(c_boundary)

	return buffer.Bytes()
}
