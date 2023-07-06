package utils

import (
	"cPing/conf"
	"fmt"
	"net"
	"os/exec"
	"runtime"
)

func ExPing(to string) conf.ExRes {
	//ping
	res := conf.ExRes{}
	result, err := Ping(to)
	if err != nil {
		res.Status = false
		res.Result = fmt.Sprintf("Error:", err)
	} else {
		res.Status = true
		res.Result = "Ping result:\n" + result
	}
	return res
}

func Ping(address string) (string, error) {
	// 解析域名，获取IP
	ips, err := net.LookupIP(address)
	if err != nil {
		return "", err
	}

	ip := ips[0].String()

	// 根据操作系统构建ping命令
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "4", ip)
	} else {
		cmd = exec.Command("ping", "-c", "4", ip)
	}

	// 执行ping命令
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
