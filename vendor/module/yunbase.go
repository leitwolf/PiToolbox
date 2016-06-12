package module

//
// 云盘基类
//
import "lib"

// YunBase 各种云盘基类
type YunBase struct {
	// 账户类型 [xunlei,yun360]
	accountType string
	// 账户列表
	accountList []lib.Account
}

// GetAccountList 获取账户列表
// 返回 [Account.Name] 列表
func (base *YunBase) GetAccountList(sender *Sender) {
	println("account", base.accountType)
	list, err := lib.LoadAccountList(base.accountType)
	if err != nil {
		sender.Err = err.Error()
		return
	}
	base.accountList = list
	var names []string
	for i := 0; i < len(list); i++ {
		names = append(names, list[i].Name)
	}
	sender.Data = names
}

// getCookieContainer 获取指定的cookie
func (base *YunBase) getCookieContainer(accountName string) (cc *lib.CookieContainer) {
	for i := 0; i < len(base.accountList); i++ {
		if base.accountList[i].Name == accountName {
			cc = base.accountList[i].CookieContainer
			return
		}
	}
	return
}
