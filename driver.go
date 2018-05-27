package mysqldriver

import (
	"fmt"
	"net"
)




func Open () []byte {
	cfg := &MysqlConfig{
		User:"root",
		Passwd:"123123",
		DBname:"test",
		Port:"3306",
		Addr:"127.0.0.1",
	}
	conn, err := net.Dial("tcp", cfg.Addr+":"+cfg.Port)
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
	}
	defer conn.Close()

	return ConnHandler(conn, cfg)
}
