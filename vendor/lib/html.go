package lib

//
// 获取html内容
//
import "net/http"

// GetHTML 获取连接响应
func GetHTML(url string) (res *http.Response, err error) {
	return GetHTMLWithCookies(url, nil)
}

// GetHTMLWithCookies 获取带cookies的连接响应
func GetHTMLWithCookies(url string, cc *CookieContainer) (res *http.Response, err error) {
	client := http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	if cc != nil {
		cc.AddToReqeust(request)
	}
	res, err = client.Do(request)
	return
}
