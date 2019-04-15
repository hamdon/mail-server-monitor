package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// 崩溃时需要传递的上下文信息
type panicContext struct {
	function interface{} // 所在函数
}

// 保护方式允许一个函数
func ProtectRun(key string, entry func(s string) interface{}) {
	// 延迟处理的函数
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		if err != nil {
			switch err.(type) {
			case runtime.Error: // 运行时错误
				fmt.Printf("runtime error:%+v\n", err)
			default: // 非运行时错误
				fmt.Printf("error:%+v\n", err)
			}
		}
	}()
	entry(key)
}

func startMailServer(s string) interface{} {
	command := "C:\\Windows\\System32\\taskkill.exe"
	params := []string{"/im", "MailServer.exe", "/f"}
	paramTwos := []string{"/im", "MailCtrl.exe", "/f"}
	cmd := exec.Command(command, params...)
	//显示运行的命令
	fmt.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	} else {
		// 保证关闭输出流
		defer stdout.Close()
		// 运行命令
		if errRun := cmd.Start(); errRun != nil {
			log.Fatal(errRun)
		}
	}
	cmdTwo := exec.Command(command, paramTwos...)
	//显示运行的命令
	fmt.Println(cmdTwo.Args)
	stdoutTwo, errTwo := cmdTwo.StdoutPipe()
	if errTwo != nil {
		log.Fatal(errTwo)
	} else {
		// 保证关闭输出流
		defer stdoutTwo.Close()
		// 运行命令
		if errRunTwo := cmdTwo.Start(); errRunTwo != nil {
			log.Fatal(errRunTwo)
		}
	}

	command = "C:\\Windows\\System32\\net.exe"
	params = []string{"start", "MagicWinmailServer"}
	paramTwos = []string{"start", "WinmailMailServerHTTPService"}

	cmd = exec.Command(command, params...)
	//显示运行的命令
	fmt.Println(cmd.Args)
	stdoutNew, errNew := cmd.StdoutPipe()
	if errNew != nil {
		panic(&panicContext{errNew})
	}
	// 保证关闭输出流
	defer stdoutNew.Close()
	// 运行命令
	if errThird := cmd.Start(); errThird != nil {
		panic(&panicContext{errThird})
	}
	// 读取输出结果
	opBytes, errFour := ioutil.ReadAll(stdoutNew)
	if errFour != nil {
		panic(&panicContext{errFour})
	}
	log.Println(string(opBytes))

	cmdTwo = exec.Command(command, paramTwos...)
	//显示运行的命令
	fmt.Println(cmdTwo.Args)
	stdoutNewTwo, errNewTwo := cmdTwo.StdoutPipe()
	if errNewTwo != nil {
		panic(&panicContext{errNewTwo})
	}
	// 保证关闭输出流
	defer stdoutNewTwo.Close()
	// 运行命令
	if errThirdTwo := cmd.Start(); errThirdTwo != nil {
		panic(&panicContext{errThirdTwo})
	}
	// 读取输出结果
	opBytesTwo, errFourTwo := ioutil.ReadAll(stdoutNewTwo)
	if errFourTwo != nil {
		panic(&panicContext{errFourTwo})
	}
	log.Println(string(opBytesTwo))
	return s
}

func main() {
	ports := []string{"25", "110", "143", "389", "465", "993", "995", "6000", "6020","6989", "6990"}
	chain := []func(string) interface{}{
		startMailServer,
	}
	fileName := "mail_server_fail.log"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
	defer file.Close()
	// timeLayout := "2006-01-02 15:04:05"
	errorTime := 0
	var fd_time, fd_content string
	for {
		for _, key := range ports {
			fmt.Print("telnet port:" + key + ":")
			_, err := net.Dial("tcp", "127.0.0.1:"+key)
			if err != nil {
				fmt.Println("fail")
				errorTime++
				if errorTime == 1 {
					// file.Seek(0, 2)    // 最后增加
					// file.WriteString(time.Now().Format(timeLayout)+"\n")
					fd_time = time.Now().Format("2006-01-02 15:04:05")
					fd_content = strings.Join([]string{"======", fd_time, "=====", "\n"}, "")
					buf := []byte(fd_content)
					file.Write(buf)
				}
				for _, proc := range chain {
					ProtectRun(key, proc)
				}
			} else {
				errorTime = 0
				fmt.Println("pass")
				fmt.Println()
			}
			time.Sleep(time.Second * 5)
		}
		time.Sleep(time.Second * 60)
	}
}
