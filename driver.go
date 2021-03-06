package mysqldriver

import (
	"fmt"
	"net"
	"time"
)

type mConnect struct {
	con net.Conn
	cfg *MysqlConfig
	readData []byte
	readDataIndex int
	readDataLength int
	writeData []byte
	sequence int
}




var mc *mConnect

func connect() {
	mc = &mConnect{
		readData:make([]byte, 4096),
		writeData:make([]byte, 4096),
	}

	mc.cfg = &MysqlConfig{
		User:   "root",
		Passwd: "123123",
		DBname: "test",
		Port:   "3306",
		Addr:   "127.0.0.1",
	}
	conn, err := net.Dial("tcp", mc.cfg.Addr+":"+mc.cfg.Port)
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
	}
	mc.con = conn
}

func Open () []byte {
	connect()
	defer close()
	data := ConnHandler(mc.con, mc.cfg)
	time.Sleep(5*time.Second)
	return data

}

func GetSystemVar(){
	connect()
	defer close()
	ConnHandler(mc.con, mc.cfg)
	mc.getSystemVar("SELECT @@version")
	mc.getSystemVar("SELECT @@max_allowed_packet")

}

func NewGetSystemVar(){
	connect()
	defer close()
	ConnHandler(mc.con, mc.cfg)
	mc.newGetSystemVar("SELECT @@version")
	mc.newGetSystemVar("SELECT @@max_allowed_packet")
}

func close(){
	mc.con.Close()
}