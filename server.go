package main

import (
	"bytes"
	"cPing/common"
	"cPing/conf"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var SiteConf conf.Conf
var WebPort = "7789"

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
	web()
}
func ttt() {
	url := "http://127.0.0.1:7788/ping"
	req := conf.ExReq{To: "127.0.0.1", Stamp: strconv.FormatInt(time.Now().Unix(), 10)}
	req.Token = common.Sha(req.Stamp + SiteConf.Token)
	// 构造 POST 请求的数据
	data, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Failed to marshal request data:", err)
		return
	}
	body := bytes.NewBuffer(data)
	// 发送 POST 请求
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		fmt.Println("Failed to send POST request:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应数据
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response data:", err)
		return
	}

	// 输出响应数据
	fmt.Println(string(respData))
}
func web() {
	http.HandleFunc("/reg", func(w http.ResponseWriter, r *http.Request) {
		var req conf.Conf
		if r.Method == "POST" {
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
			time.Sleep(3 * time.Second)
			if test(req.Client, req.Token) {
				fmt.Fprint(w, "Reg success")
				fmt.Println("Reg success")
				reg(req)
			} else {
				fmt.Fprint(w, "Reg fail")
				fmt.Println("Reg fail")
			}
		} else {
			fmt.Fprint(w, "error")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	fmt.Println("Listening on :" + WebPort + "...")
	http.ListenAndServe(":"+WebPort, nil)
}

func test(url, token string) bool {
	url = url + "/ping"
	req := conf.ExReq{To: "127.0.0.1", Stamp: strconv.FormatInt(time.Now().Unix(), 10)}
	req.Token = common.Sha(req.Stamp + token)
	// 构造 POST 请求的数据
	data, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Failed to marshal request data:", err)
		return false
	}
	var res conf.ExRes
	err = json.Unmarshal([]byte(common.Post(url, data)), &res)
	if err != nil {
		fmt.Println("Failed to unmarshal JSON data:", err)
		return false
	}
	// 输出响应数据
	return res.Status
}

func reg(site conf.Conf) {
	var list conf.Client
	// 判断是否存在 client.json 文件
	if _, err := os.Stat("client.json"); os.IsNotExist(err) {
		// 文件不存在，创建并写入数据
		data := list
		data.List = append(data.List, site)
		file, err := os.Create("client.json")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		jsonBytes, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			panic(err)
		}

		_, err = file.Write(jsonBytes)
		if err != nil {
			panic(err)
		}

		fmt.Println("client.json created")
	} else { // 打开 client.json 文件
		file, err := os.Open("client.json")
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

		// 创建足够大的缓冲区
		bytes := make([]byte, stat.Size())

		// 读取文件内容到缓冲区
		_, err = io.ReadFull(file, bytes)
		if err != nil {
			fmt.Println("Failed to read file:", err)
			return
		}

		// 解析 JSON 数据为 Client 类型的变量 oldList
		var oldList conf.Client
		err = json.Unmarshal(bytes, &oldList)
		if err != nil {
			fmt.Println("Failed to unmarshal JSON data:", err)
			return
		}
		oldList.List = append(oldList.List, site)
		// 将更新后的 oldList 变量写入 client.json 文件中
		file, err = os.Create("client.json")
		if err != nil {
			fmt.Println("Failed to create file:", err)
			return
		}
		defer file.Close()

		newBytes, err := json.MarshalIndent(oldList, "", "    ")
		if err != nil {
			fmt.Println("Failed to marshal JSON data:", err)
			return
		}

		_, err = file.Write(newBytes)
		if err != nil {
			fmt.Println("Failed to write file:", err)
			return
		}

		fmt.Println("Site added to client.json")

	}
}
