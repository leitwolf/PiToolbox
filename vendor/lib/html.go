package lib

//
// 获取html内容
//
import (
	"bytes"
	"io"
	"net/http"
)

// GetHTML 获取连接响应
func GetHTML(url string, cc *CookieContainer) (res *http.Response, err error) {
	req, err := MakeRequest("GET", url, nil, cc)
	if err != nil {
		return
	}
	res, err = FetchHTML(req, cc)
	return
}

// PostHTML post连接
func PostHTML(url string, body []byte, cc *CookieContainer) (res *http.Response, err error) {
	req, err := MakeRequest("POST", url, body, cc)
	if err != nil {
		return
	}
	res, err = FetchHTML(req, cc)
	return
}

// MakeRequest 构造一个http.Reqeust
// 可以自己加上一些header
func MakeRequest(method string, urlStr string, body []byte, cc *CookieContainer) (req *http.Request, err error) {
	var body1 io.Reader
	if body != nil {
		body1 = bytes.NewBuffer(body)
	}
	req, err = http.NewRequest(method, urlStr, body1)
	if err != nil {
		return
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.84 Safari/537.36")
	if cc != nil {
		cc.AddToReqeust(req)
	}
	return
}

// FetchHTML 根据http.Reqeust获取html响应
func FetchHTML(req *http.Request, cc *CookieContainer) (res *http.Response, err error) {
	client := http.Client{}
	res, err = client.Do(req)
	if cc != nil && err == nil && res.StatusCode == 200 {
		// 成功的响应，写入cookies
		cc.Update(res.Cookies())
	}
	return
}
