package http

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"net/http"
	"time"
)

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

// func main() {
// 	// 示例：GET 请求
// 	getUrl := "https://jsonplaceholder.typicode.com/posts/1"
// 	getResponse, err := HttpGet(getUrl)
// 	if err != nil {
// 		fmt.Println("GET 请求失败:", err)
// 	} else {
// 		fmt.Println("GET 响应:", getResponse)
// 	}

// 	// 示例：POST 请求
// 	postUrl := "https://jsonplaceholder.typicode.com/posts"
// 	postData := map[string]interface{}{
// 		"title":  "foo",
// 		"body":   "bar",
// 		"userId": 1,
// 	}
// 	postResponse, err := HttpPost(postUrl, postData)
// 	if err != nil {
// 		fmt.Println("POST 请求失败:", err)
// 	} else {
// 		fmt.Println("POST 响应:", postResponse)
// 	}
// }
