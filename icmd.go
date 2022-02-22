package icmd

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"strings"
)

//执行常规命令
func Exec(args string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", args)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	if err := cmd.Run(); err != nil {
		return "", errors.New(stderr.String())
	}
	return strings.Trim(stdout.String(), "\n"), nil
}

//管道响应
func Pipe(args string, ch *chan []byte) {
	cmd := exec.Command("/bin/bash", "-c", args)
	// 命令的错误输出和标准输出都连接到同一个管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("错误：", err)
		return
	}
	cmd.Stderr = cmd.Stdout

	// 从管道中实时获取输出并打印到终端
	for {
		tmp := make([]byte, 1024*1024)
		_, err = stdout.Read(tmp)
		if err != nil {
			log.Println("错误：", err)
			_ = stdout.Close()
			break
		}

		tmp = bytes.Trim(tmp, "\x00")
		arr := bytes.Split(tmp, []byte("\n"))
		for _, item := range arr {
			if len(item) > 0 {
				*ch <- item
			}
		}
	}
	if err = cmd.Wait(); err != nil {
		log.Println(err)
	}
}
