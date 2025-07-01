package util

import "fmt"

func GetFileNameWithLine(fileName string, lineNumber int) string {
	return fmt.Sprintf("%v: %d", fileName, lineNumber)
}