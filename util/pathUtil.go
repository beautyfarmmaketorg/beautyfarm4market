package util

import (
	"os"
	"path/filepath"
)

func GetCurrentPath()string{
	execpath, _ := os.Executable() // 获得程序路径
	return filepath.Dir(execpath)
}


