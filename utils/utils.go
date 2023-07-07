package utils

import (
	"cPing/conf"
	"errors"
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
		cmd = exec.Command("ping", "-w", "5000", "-n", "4", ip)
	} else {
		cmd = exec.Command("ping", "-w", "5", "-c", "4", ip)
	}

	// 执行ping命令
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func MTR(address string) (string, error) {
	// 检查 mtr 是否已安装
	_, err := exec.LookPath("mtr")
	if err != nil {
		return "", errors.New("mtr not found, please install it")
	}

	// 根据操作系统构建 mtr 命令
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("mtr.exe", "--report", "--report-cycles", "4", address)
	} else {
		cmd = exec.Command("mtr", "--report", "--report-cycles", "4", address)
	}

	// 执行 mtr 命令
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
