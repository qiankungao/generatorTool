package tools

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"unicode"
)
func CreateFileIfHasDel(filename string) {
	if IsDirOrFileExist(filename) == true {
		DelFile(filename)
	}

	CreateFile(filename)
}

func IsDirOrFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func DelFile(filename string) {
	_ = os.Remove(filename)
}

func CreateFile(path string) bool {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		return false
	}
	return true
}

func ReadFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(data)
}

//追加写文件
func WriteAppendFile(path, data string) (err error) {
	if _, err := WriteFileAppend(path, data); err == nil {
		fmt.Printf("Generate success:%s\n", path)
		return nil
	} else {
		return err
	}
}


//追加写文件
func  WriteFileAppend(filename string, data string) (count int, err error) {
	var f *os.File
	if IsDirOrFileExist(filename) == false {
		f, err = os.Create(filename)
		if err != nil {
			return
		}
	} else {
		f, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666)
	}
	defer f.Close()
	count, err = io.WriteString(f, data)
	if err != nil {
		return
	}
	return
}

// 添加``符号
func AddQuote(str string) string {
	return "`" + str + "`"
}
// 去掉 `符号
func CleanQuote(str string) string {
	return strings.Replace(str, "`", "", -1)
}

// 去掉 "符号
func CleanDoubleQuotes(str string) string {
	return strings.Replace(str, "\"", "", -1)
}

//替换字符
func ReplaceStr(str, old, new string) string {
	return strings.Replace(str, old, new, -1)
}

//首字母小写
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

//FMT 格式代码
func Gofmt(path string) bool {
	if IsDirOrFileExist(path) {
		if !ExecCommand("goimports", "-l", "-w", path) {
			if !ExecCommand("gofmt", "-l", "-w", path) {
				return ExecCommand("go", "fmt", path)
			}
		}
		return true
	}
	return false
}
func ExecCommand(name string, args ...string) bool {
	cmd := exec.Command(name, args...)
	_, err := cmd.Output()
	if err != nil {
		return false
	}
	return true
}