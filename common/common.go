package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

func Sha(encode string) string {
	// 定义需要加密的字符串
	str := encode
	// 计算 SHA256 哈希值
	hash := sha256.Sum256([]byte(str))
	// 将哈希值转换为十六进制字符串
	hashString := hex.EncodeToString(hash[:])
	// 输出加密结果
	return hashString
}

func Post(url string, data []byte) string {
	body := bytes.NewBuffer(data)

	// 发送 POST 请求
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		fmt.Println("Failed to send POST request:", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应数据
	var respData []byte
	respData, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response data:", err)
		return ""
	}

	return string(respData)
}
