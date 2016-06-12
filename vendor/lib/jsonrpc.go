package lib

//
// 处理jsonrpc的请求与返回
//
import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// JSONRPCRequest jsonrpc请求
type JSONRPCRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// JSONRPCResponse jsonrpc响应返回
type JSONRPCResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Error   interface{} `json:"error"`
	Result  interface{} `json:"result"`
}

// CallJSONRPC 请求
func CallJSONRPC(url string, req *JSONRPCRequest) (res *JSONRPCResponse, err error) {
	b, err := json.Marshal(*req)
	if err != nil {
		return
	}
	// fmt.Printf("%v\n", *req)
	// println("req", string(b))
	body := bytes.NewBuffer(b)
	response, err := http.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		return
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	// ioutil.WriteFile("test/aria2.json", content, 777)
	// println("res", string(content))
	res = &JSONRPCResponse{}
	err = json.Unmarshal(content, res)
	// if err != nil {
	// 	println(err.Error())
	// }
	return
}
