package main

import (
	"cPing/common"
	"cPing/conf"
	"cPing/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

var SiteConf conf.Site
var WebPort = "7788"

func init() {
	// 读取配置文件
	file, err := os.Open("conf.json")
	if err != nil {
		fmt.Println("Failed to open config file:", err)
		return
	}
	defer file.Close()
	// 获取文件大小
	stat, err := file.Stat()
	if err != nil {
		fmt.Println("Failed to get file size:", err)
		return
	}
	// 解析配置文件
	// 分配足够的空间来存储文件内容
	bytes := make([]byte, stat.Size())
	_, err = io.ReadFull(file, bytes)
	if err != nil {
		fmt.Println("Failed to read file:", err)
		return
	}
	err = json.Unmarshal(bytes, &SiteConf)
	if err != nil {
		fmt.Println("Failed to parse config file:", err)
		return
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go web(&wg)
	reg()
	wg.Wait()
}

func web(wg *sync.WaitGroup) {
	http.HandleFunc("/mtr", func(w http.ResponseWriter, r *http.Request) {
		mtr, _ := utils.MTR(r.FormValue("to"))
		fmt.Fprint(w, mtr)
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		var req conf.ExReq
		if r.Method == "GET" {
			req.To = r.FormValue("to")
			req.Stamp = r.FormValue("stamp")
			req.Token = r.FormValue("token")
		} else {
			// 读取响应数据
			// 读取请求体
			var body []byte
			var err error
			if r.ContentLength > 0 {
				body = make([]byte, r.ContentLength)
				_, err = io.ReadFull(r.Body, body)
				if err != nil {
					fmt.Println("Failed to read request body:", err)
					return
				}
			}
			// 解析响应数据
			err = json.Unmarshal(body, &req)
			if err != nil {
				fmt.Println("Failed to unmarshal response data:", err)
				return
			}
		}
		if common.Sha(req.Stamp+SiteConf.Token) == req.Token && common.StampPass(req.Stamp, 5) {
			// 将 Person 实例转换为 JSON 格式的字符串
			jsonBytes, err := json.Marshal(utils.ExPing(req.To))
			if err != nil {
				fmt.Println("Failed to encode person:", err)
				return
			}
			response := string(jsonBytes)
			// 发送响应消息
			fmt.Fprint(w, response)
			return
		} else {
			fmt.Fprint(w, "error")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	fmt.Println("Listening on :" + WebPort + "...")
	http.ListenAndServe(":"+WebPort, nil)
}

func reg() {
	url := SiteConf.Server + "/reg"
	req := SiteConf
	// 构造 POST 请求的数据
	data, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Failed to marshal request data:", err)
		return
	}
	// 输出响应数据
	fmt.Println(common.Post(url, data))
}
