package utils

import (
	"testing"
)

func init() {
	InitLogger()
}

func TestLogger(t *testing.T) {
	Log("this is log")
	ErrorMsg("this is error")
	Warn("this is warn")
}
