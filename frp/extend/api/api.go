package api

import (
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	// "net/http"
	// "net/url"
	// "strconv"
	// "github.com/fatedier/frp/pkg/util/util"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/config/types"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/fatedier/frp/pkg/msg"
	"strconv"
)

// Service sakurafrp api servie
// type Service struct {
// 	Host string
// 	// Host url.URL
// }

// // NewService crate sakurafrp api servie
// func NewService(host string) (s *Service, err error) {
// 	u, err := url.Parse(host)
// 	if err != nil {
// 		return
// 	}
// 	return &Service{*u}, nil
// }
// 封装 GET 请求函数
func HttpGet(url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second, // 设置超时时间
	}

	// 创建 GET 请求
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// 封装 POST 请求函数（发送 JSON 数据）
func HttpPost(url string, data map[string]interface{}) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second, // 设置超时时间
	}

	// 将请求数据编码为 JSON 格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// 创建 POST 请求
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// 定义将要解析到的结构体
type CheckTokenResp struct {
    Status int  	 `json:"status"`
    Success   bool    `json:"success"`
}

type Response struct {
    Status  int    `json:"status"`
    Success bool   `json:"success"`
    Message string `json:"message"`
}

type Message struct {
    Inbound  int `json:"inbound"`
    Outbound int `json:"outbound"`
    Type     int `json:"type"`
}

// CheckToken 校验客户端 token   user=client中的token
func CheckToken(apiurl string,user string, John_San_Token string) (ok bool, err error) {
	
	var CheckToken_url string = apiurl + "?action=checktoken&user="+user + "&apitoken=" + John_San_Token
	fmt.Println("John_San_检查客户端登录 ",CheckToken_url)
	resp,err := HttpGet(CheckToken_url)
	if err != nil {
		return false, err
	}
	// 将要解析到的变量
    var r CheckTokenResp
 
    // 解析JSON字符串到结构体
    err = json.Unmarshal([]byte(resp), &r)
    if (err != nil) {
        fmt.Println("解析checktoken失败",resp , err)
        return false, err
    }
    fmt.Println("api调用结束",resp)	
    return r.Success == true,nil
}

// CheckProxy 校验客户端代理
func CheckProxy(apiurl string, user string, pxyConf v1.ProxyConfigurer, pMsg *msg.NewProxy,John_San_Token string) (ok bool, err error) {
	// fmt.Println("下面是 pxyConf")
	// util.PrintHexAndDecFields(pxyConf)	
	// fmt.Println("下面是 pxyConf.GetBaseConfig()")
	// util.PrintHexAndDecFields(pxyConf.GetBaseConfig())	
	// fmt.Println("下面是 pMsg")
	// util.PrintHexAndDecFields(pMsg)	


	// fmt.Println("结束",len(pMsg.CustomDomains),len(pMsg.SubDomain))
	

	var geturl string = "?action=checkproxy&user=" + user + 
	    "&apitoken=" + John_San_Token +
	    "&proxy_name=" + pxyConf.GetBaseConfig().Name +
	    "&proxy_type=" + pxyConf.GetBaseConfig().Type

	//自定义域名
	if (len(pMsg.CustomDomains) > 0) {
		geturl = geturl + "&customdomains=" + string(pMsg.CustomDomains[0])
	}
	//自定义二级域名
	if (len(pMsg.SubDomain) > 0) {
		geturl = geturl + "&customdomains=" + string(pMsg.SubDomain[0])
	}
	//远程端口 可能为空？
	remotp := strconv.Itoa(pMsg.RemotePort)
	if (remotp != "" && remotp != "0" ) {
		geturl = geturl + "&remote_port=" + remotp
	}
	fmt.Println("apiurl=", apiurl + geturl)

	//http://127.0.0.1:7899/api?action=checkproxy&user=46414e6574f2969e&apitoken=wobushitoken9527|1&proxy_name=46414e6574f2969e.MyProxy&proxy_type=tcp&remote_port=46008
	//http://127.0.0.1:7899/api?action=checkproxy&user=46414e6574f2969e&apitoken=wobushitoken9527|1&proxy_name=46414e6574f2969e.test_http&proxy_type=http&customdomains=aaaa.aaa.comi&remote_port=0
	resp, err := HttpGet(apiurl + geturl)
	fmt.Println("CheckProxy 代理检查",resp, err)
	if err != nil {
        fmt.Println("CheckProxy Error:", resp,err)
        return false,err
    }
	var response Response
    err = json.Unmarshal([]byte(resp), &response)
    if err != nil {
        fmt.Println("CheckProxy decode Error:", resp,err)
        return false,err
    }
    if (response.Success != true){
    	fmt.Println("CheckProxy Error:", resp,err)
        return false,err
    }

    // 解析 message 字段
    var msg Message
    err = json.Unmarshal([]byte(response.Message), &msg)
    if err != nil {
        fmt.Println("CheckProxy decode_message Error:", resp,err)
        return false,err
    }
    // fmt.Printf("Status: %d\n", response.Status)
    // fmt.Printf("Success: %t\n", response.Success)
    // fmt.Printf("Inbound: %d\n", msg.Inbound)
    // fmt.Printf("Outbound: %d\n", msg.Outbound)
    // fmt.Printf("Type: %d\n", msg.Type)
    if (msg.Outbound <= 0) {
		return false,nil
    }
    outboundWithKB := strconv.Itoa(msg.Outbound) + "KB"
    // fmt.Println("本代理限速",outboundWithKB,msg.Outbound)
	pxyConf.GetBaseConfig().Transport.BandwidthLimit, _ = types.NewBandwidthQuantity(outboundWithKB)
	fmt.Println("John_San_检查代理这里可以最终限速:", pxyConf.GetBaseConfig().Transport.BandwidthLimit)
	
	
	return true, nil
}