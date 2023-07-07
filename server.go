package main

import (
	"cPing/assets"
	"cPing/common"
	"cPing/conf"
	"encoding/json"
	"fmt"
	"github.com/oklog/ulid"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var WebPort = "7789"

const ServerToken = "csz.net"

func main() {
	web()
}

func web() {
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"add": add,
	}).ParseFS(assets.Templates, "templates/*"))
	http.HandleFunc("/reg", func(w http.ResponseWriter, r *http.Request) {
		var req conf.Site
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
			test := do(req.Client+"/ping", "baidu.com", req.Token)
			if test.Status {
				fmt.Fprint(w, "Reg success")
				fmt.Println("Reg success")
				reg(req)
			} else {
				fmt.Fprint(w, "Reg fail:\n"+test.Result)
				fmt.Println("Reg fail")
			}
		} else {
			fmt.Fprint(w, "error")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	http.HandleFunc("/site/", func(w http.ResponseWriter, r *http.Request) {
		_, data := getClient()
		// 解析动态路径参数
		id := r.URL.Path[len("/site/"):]
		if common.Sha(r.FormValue("s")+ServerToken+r.FormValue("to")) == r.FormValue("t") && common.StampPass(r.FormValue("s"), 20) {
			// 在数据中查找对应的信息
			for _, site := range data.List {
				if site.Id == id {
					// 返回对应的信息
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(do(site.Client+"/"+r.FormValue("do"), r.FormValue("to"), site.Token))
					//json.NewEncoder(w).Encode(site)
					return
				}
			}
			// 如果未找到对应信息，返回404 Not Found
			http.NotFound(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})

	//Ping首页
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// 显示页面
			if r.FormValue("to") == "" {
				tmpl.ExecuteTemplate(w, "ping.tmpl", nil)
			} else {
				_, client := getClient()
				s := strconv.FormatInt(time.Now().Unix(), 10)
				t := common.Sha(s + ServerToken + r.FormValue("to"))
				data := struct {
					To     string
					Client conf.Client
					S      string
					T      string
				}{
					To:     r.FormValue("to"),
					Client: client,
					T:      t,
					S:      s,
				}
				tmpl.ExecuteTemplate(w, "ping_res.tmpl", data)
			}

		} else if r.Method == "POST" {
			// 处理表单提交
			url := "/ping?to=" + r.FormValue("url")
			http.Redirect(w, r, url, http.StatusSeeOther)
		}
	})
	fmt.Println("Listening on :" + WebPort + "...")
	http.ListenAndServe(":"+WebPort, nil)
}
func add(a, b int) int {
	return a + b
}

func do(url, to, token string) conf.ExRes {
	req := conf.ExReq{To: to, Stamp: strconv.FormatInt(time.Now().Unix(), 10)}
	req.Token = common.Sha(req.Stamp + token + req.To)
	var res conf.ExRes
	res.Status = false
	// 构造 POST 请求的数据
	data, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Failed to marshal request data:", err)
		return res
	}
	err = json.Unmarshal([]byte(common.Post(url, data)), &res)
	if err != nil {
		fmt.Println("Failed to unmarshal JSON data:", err)
		return res
	}
	// 输出响应数据
	return res
}

func reg(site conf.Site) {
	var list conf.Client
	// 判断是否存在 client.json 文件
	if _, err := os.Stat("client.json"); os.IsNotExist(err) {
		// 文件不存在，创建并写入数据
		data := list
		// 自动分配 ULID 作为 ID
		id := ulid.MustNew(ulid.Now(), nil).String()
		site.Id = id // 将新的 ID 存入 ID 字段中
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

		// 判断是否已存在相同的数据
		for _, oldSite := range oldList.List {
			if oldSite.Id != site.Id &&
				oldSite.Client == site.Client &&
				oldSite.Server == site.Server &&
				oldSite.Token == site.Token {
				fmt.Printf("Site already exists: %+v\n", oldSite.Id)
				return
			}
		}

		// 自动分配 ULID 作为 ID
		id := ulid.MustNew(ulid.Now(), nil).String()
		site.Id = id // 将新的 ID 存入 ID 字段中
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
func getClient() (error, conf.Client) {
	file, err := os.Open("client.json")
	var oldList conf.Client
	if err != nil {
		fmt.Println("Failed to open config file:", err)
		return err, oldList
	}
	defer file.Close()

	// 获取文件大小
	stat, err := file.Stat()
	if err != nil {
		fmt.Println("Failed to get file size:", err)
		return err, oldList
	}

	// 创建足够大的缓冲区
	bytes := make([]byte, stat.Size())

	// 读取文件内容到缓冲区
	_, err = io.ReadFull(file, bytes)
	if err != nil {
		fmt.Println("Failed to read file:", err)
		return err, oldList
	}

	// 解析 JSON 数据为 Client 类型的变量 oldList

	err = json.Unmarshal(bytes, &oldList)
	if err != nil {
		fmt.Println("Failed to unmarshal JSON data:", err)
		return err, oldList
	}
	return nil, oldList
}
