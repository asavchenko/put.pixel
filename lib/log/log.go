package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

func Info(msgs ...interface{}) {
	baseDir := getBaseDir()
	formattedMsg := time.Now().UTC().Format("15:04:05.999 02-01-2006 ")
	hasStackTrace := false
	for i := 1; ; i++ {
		_, f, line, _ := runtime.Caller(i)
		if f == "" {
			break
		}
		if !strings.Contains(f, baseDir) {
			break
		}
		if strings.Contains(f, "main.go") {
			break
		}

		hasStackTrace = true
		relativePath := strings.Replace(f, baseDir+"/", "", -1)
		if i == 1 {
			formattedMsg += fmt.Sprintln(append([]interface{}{fmt.Sprintf(" %s:%d", relativePath, line)}, msgs...)...)
		} else {
			formattedMsg += fmt.Sprintf("%s:%d\n", relativePath, line)
		}
	}
	if !hasStackTrace {
		formattedMsg += " " + fmt.Sprintln(msgs...)
	}
	fmt.Println(formattedMsg)
}

func Error(msgs ...interface{}) {
	baseDir := getBaseDir()
	formattedMsg := time.Now().UTC().Format("15:04:05.999 02-01-2006 ")
	hasStackTrace := false
	for i := 1; ; i++ {
		_, f, line, _ := runtime.Caller(i)
		if f == "" {
			break
		}
		if !strings.Contains(f, baseDir) {
			break
		}
		if strings.Contains(f, "main.go") {
			break
		}

		hasStackTrace = true
		relativePath := strings.Replace(f, baseDir+"/", "", -1)
		if i == 1 {
			formattedMsg += fmt.Sprintln(append([]interface{}{fmt.Sprintf(" %s:%d", relativePath, line)}, msgs...)...)
		} else {
			formattedMsg += fmt.Sprintf("%s:%d\n", relativePath, line)
		}
	}
	if !hasStackTrace {
		formattedMsg += " " + fmt.Sprintln(msgs...)
	}
	fmt.Println("*********ERROR", formattedMsg)
}

func getBaseDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	return pwd
}
