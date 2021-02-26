package util

import (
	"testing"
)

func TestAesEncryptStr(t *testing.T) {

	key := []byte("fe023f_9fd&fwfl0")

	str, err := AesEncryptStr([]byte("1908273366332802"), key)
	if err != nil {

		t.Error("测试失败")
	}

	t.Log(str)
}
