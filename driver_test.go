package mysqldriver_test

import (
	"testing"
	"github.com/petrie/mysqldriver"
)

func TestOpen(t *testing.T) {
	testByte := []byte {7,0,0,2,0,0,0,2,0,0,0}
	res := mysqldriver.Open()

	if (string(testByte) == string(res)){
		t.Log("初始化成功")
	}else{
		t.Error("连接失败")
	}

}