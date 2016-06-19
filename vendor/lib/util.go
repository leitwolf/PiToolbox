package lib

//
// 一些工具类
//
import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

// WriteFile 保存文件，可自动建文件夹
func WriteFile(filename string, data []byte) (err error) {
	dirStr := path.Dir(filename)
	_, err = os.Stat(dirStr)
	b := os.IsNotExist(err)
	if b {
		err = os.MkdirAll(dirStr, 0777)
		if err != nil {
			return
		}
	}
	err = ioutil.WriteFile(filename, data, 0777)
	return
}

const gnumber float32 = 1024 * 1024 * 1024
const mnumber float32 = 1024 * 1024
const knumber float32 = 1024

// GetReadableSize 把容量转成可读性比较高的方式 以K M G形式
func GetReadableSize(origin string) (readable string) {
	nn, err := strconv.Atoi(origin)
	if err != nil {
		return
	}
	n := float32(nn)
	if n >= gnumber {
		r := n / gnumber
		readable = strconv.FormatFloat(float64(r), 'f', 2, 32) + "G"
	} else if n >= mnumber {
		r := n / mnumber
		readable = strconv.FormatFloat(float64(r), 'f', 2, 32) + "M"
	} else if n >= knumber {
		r := n / knumber
		readable = strconv.FormatFloat(float64(r), 'f', 2, 32) + "K"
	} else {
		readable = origin
	}
	return
}

// Account 一个账户
type Account struct {
	// 名称，默认为default
	Name string
	// cookie内容
	CookieContainer *CookieContainer
}

// LoadAccountList 加载账户列表
// url: config/cookies_accountType_xxx.json
// @param accountType [xunlei,baidu]
func LoadAccountList(accountType string) (accountList []Account, err error) {
	accountList = make([]Account, 0)
	dir, err := ioutil.ReadDir("config")
	if err != nil {
		return
	}
	prefix := "cookies_" + accountType
	defaultAccount := prefix + ".json"
	// println("default", defaultAccount)
	for i := 0; i < len(dir); i++ {
		f := dir[i]
		name := f.Name()
		if !f.IsDir() && strings.HasPrefix(name, prefix) && strings.HasSuffix(name, ".json") {
			var accountName string
			// 只一个，default
			if name == defaultAccount {
				accountName = "default"
			} else {
				accountName = name[len(prefix)+1 : len(name)-5]
			}
			if accountName == "" {
				accountName = "empty"
			}
			println("accountName", accountName)
			path := "config/" + name
			cc, err1 := NewCookieContainer(path)
			if err1 != nil {
				log.Println("Load cookie " + path + " fail!")
				err = err1
				return
			}
			account := Account{Name: accountName, CookieContainer: cc}
			accountList = append(accountList, account)
		}
	}
	return
}
