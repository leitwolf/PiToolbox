package lib

//
// 一个cookie容器
//
import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// CookieContainer cookies管理
type CookieContainer struct {
	// 当前cookie所在路径
	filepath string
	cookies  []http.Cookie
}

// Load 加载
func (cc *CookieContainer) Load(filepath string) (err error) {
	cc.filepath = filepath
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &cc.cookies)
	return
}

// Update 更新写入cookies文件
func (cc *CookieContainer) Update(cookies []*http.Cookie) {
	if cookies == nil || len(cookies) == 0 {
		return
	}
	for i := 0; i < len(cookies); i++ {
		c := cookies[i]
		// 对应的
		has := false
		var c2 http.Cookie
		for j := 0; j < len(cc.cookies); j++ {
			c1 := cc.cookies[j]
			if c.Name == c1.Name {
				c2 = c1
				has = true
				break
			}
		}
		if has {
			c2.Name = c.Name
			c2.Value = c.Value
			c2.Expires = c.Expires
			// c2.Path = c.Path
			// c2.Domain = c.Domain
			// c2.MaxAge = c.MaxAge
			// c2.Raw = c.Raw
			// c2.RawExpires = c.RawExpires
			// c2.Secure = c.Secure
			// c2.Unparsed = c.Unparsed
		} else {
			cc.cookies = append(cc.cookies, *c)
		}
	}
	// 写入
	b, err := json.Marshal(cc.cookies)
	if err == nil {
		// println("write cookies")
		ioutil.WriteFile(cc.filepath, b, 777)
	}
}

// GetValueByName 查找指定Name的cookie的值
func (cc *CookieContainer) GetValueByName(name string) (value string) {
	for _, cookie := range cc.cookies {
		if cookie.Name == name {
			value = cookie.Value
			return
		}
	}
	return
}

// AddToReqeust 添加cookie给request
func (cc *CookieContainer) AddToReqeust(request *http.Request) {
	for _, cookie := range cc.cookies {
		request.AddCookie(&cookie)
	}
}

// GetHeaderStr 生成请求头里的Cookie
// 如: Cookie: xxx=xxx; xxx=xxx
func (cc *CookieContainer) GetHeaderStr() (str string) {
	str = "Cookie:"
	for i := 0; i < len(cc.cookies); i++ {
		cookie := cc.cookies[i]
		str += " " + cookie.Name + "=" + cookie.Value
		if i < len(cc.cookies)-1 {
			str += ";"
		}
	}
	return
}

// NewCookieContainer 新建一个CookieJar
func NewCookieContainer(path string) (cc *CookieContainer, err error) {
	cc = &CookieContainer{cookies: make([]http.Cookie, 0)}
	err = cc.Load(path)
	return
}
