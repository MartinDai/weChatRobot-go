package logger

import (
	"fmt"
	"testing"
)

func TestLogger(t *testing.T) {
	SetLevel("DEBUG")
	Debug("this is a debug msg", "key1", "value1", "key2", "value2")
	Info("this is a info msg", "key1", "value1", "key2", "value2")
	Warn("this is a warn msg", "key1", "value1", "key2", "value2")
	Error(fmt.Errorf("error reason"), "this is a error msg", "key1", "value1", "key2", "value2")
	FatalWithErr(fmt.Errorf("fatal reason"), "this is a fatal msg", "key1", "value1", "key2", "value2")
}
