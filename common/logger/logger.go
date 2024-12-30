package logger

import "log"

const (
	Orange = "\033[38;5;214m"
	Red    = "\033[31m"
	Reset  = "\033[0m"
)

func Logger(err string, funcName string) {
	log.Println(Orange+funcName+Reset," ",Red+err+Reset)
}