package util

import (
	"os/exec"
	"os"
	"strings"
)

func GetCurrentPath()string{
	s,err := exec.LookPath(os.Args[0])
	CheckErr(err)
	i := strings.LastIndex(s,"\\")
	path:=string(s[0:i+1])
	return path
}


