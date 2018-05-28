package mysqldriver_test

import (
	"testing"
	"github.com/petrie/mysqldriver"
)

func TestOpen(t *testing.T) {
	testByte := []byte {7,0,0,2,0,0,0,2,0,0,0} //mysql handshake `Ok packet`
	res := mysqldriver.Open()

	if (string(testByte) == string(res)){
		t.Log("连接成功")
	}else{
		t.Error("连接失败")
	}

}

func TestGetSystemVar(t *testing.T){
	mysqldriver.GetSystemVar()
}

func TestNewGetSystemVar(t *testing.T){
	mysqldriver.NewGetSystemVar()
}