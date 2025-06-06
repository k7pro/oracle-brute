package tools

import (
	"bufio"
	"fmt"
	"os"
)

func FileWrite(content string) {
	// 以追加模式打开文件，权限设置为0666（所有人可读写）
	file, err := os.OpenFile("result.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	file.WriteString(content + "\n")
}

func FileReadLines(filepath string) []string {
	// 定义一个字符串切片用于存储读取的内容
	var fileSlice []string

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("文件打开错误：", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileSlice = append(fileSlice, scanner.Text())
	}

	return fileSlice
}
