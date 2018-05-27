package mysqldriver

import (
	"fmt"
	"net"
	"bytes"
	"encoding/binary"
)


func ConnHandler(c net.Conn, cfg *MysqlConfig) []byte{

	buf := make([]byte, 1024)

	n , err := c.Read(buf)
	if err != nil || n == 0 {
		fmt.Printf("connect error%s\n", err)
	}
	//fmt.Println("protocol version:",buf[4])

	serverVersionLen :=  5 + bytes.IndexByte(buf[4:],0x00)
	//fmt.Printf("len:%d\n", serverVersionLen)
	//fmt.Printf("server version:%s\n", string(buf[5:serverVersionLen - 1 ]))

	index := serverVersionLen;
	//fmt.Println("connection id:", buf[index:index + 4])

	index += 4
	var auth_plugin_data_part_1  []byte;
	auth_plugin_data_part_1 = make([]byte , 8)
	copy(auth_plugin_data_part_1 , buf[index: index+8])
	cipher := buf[index: index+8]

	index += 8
	//fill
	index +=1

	//capability_flag_1
	capability_flag_1_len := 2
	flags := clientFlag(binary.LittleEndian.Uint16(buf[index: index + 2]))
	index += capability_flag_1_len

	//unused data now
	unusedDatalen := 1 + 2 + 2 + 1 + 10

	index+= unusedDatalen;
	cipher = append(cipher, buf[index: index + 12]...)

	//构造 HandshakeResponse

	responseFlags := clientProtocol41 |
			clientSecureConn |
			clientLongPassword |
			clientTransactions |
			clientLocalFiles |
			clientPluginAuth |
			clientMultiResults |
			clientConnectWithDB |
			flags & clientLongFlag

	/**
	flag :4
	MaxPacketSize:4
	**/

	scrambleBuff := scramblePassword(cipher, []byte(cfg.Passwd))

	responseLen := 4 + 4 + 1 + 23 + len(cfg.User) + 1 + 1 + len(scrambleBuff) +1 + len(cfg.DBname)+ 21 + 1

	responseData := make([]byte, responseLen + 4)

	//0-2 存储包长度
	responseData[0] = byte(responseLen)
	responseData[1] = byte(responseLen >> 8)
	responseData[2] = byte(responseLen >> 16)
	//seq mysql包的序列号
	responseData[3] = 1
	responseData[4] = byte(responseFlags)
	responseData[5] = byte(responseFlags >> 8)
	responseData[6] = byte(responseFlags >> 16)
	responseData[7] = byte(responseFlags >> 24)

	//8-11 存储flag
	responseData[8] = 0x00
	responseData[9] = 0x00
	responseData[10] = 0x00
	responseData[11] = 0x00
	//存储字符集
	responseData[12] = 33
	resIndex := 13

	//used hacked
	for ; resIndex < 13+23; resIndex++ {
		responseData[resIndex] = 0
	}
	//user
	resIndex += copy(responseData[resIndex:], cfg.User)
	responseData[resIndex] = 0x00
	resIndex ++

	//passwd
	responseData[resIndex] = byte(len(scrambleBuff))
	resIndex += 1 + copy(responseData[resIndex+1:], scrambleBuff)


	//dbname
	resIndex += copy(responseData[resIndex:], cfg.DBname)
	responseData[resIndex] = 0x00
	resIndex++

	//固定字符串
	resIndex += copy(responseData[resIndex:], "mysql_native_password")
	responseData[resIndex] = 0x00


	n, err = c.Write(responseData[:])

	if err != nil{
		fmt.Println(err)
	}
	buf2 := make([]byte, 1024)
	n, err = c.Read(buf2);
	if err != nil {
		fmt.Println("read error")
	}
	return buf2[0:n]
}
func (mc *mConnect) getSystemVar(cmd string) {
	//构造包结构
	packlen := len(cmd) + 1
	buf := make([]byte, packlen+4)

	buf[4] = comQuery
	copy(buf[5:], cmd)

	buf[0] = byte(packlen)
	buf[1] = byte(packlen >> 8)
	buf[2] = byte(packlen >> 16)
	buf[3] = 0

	n, err := mc.con.Write(buf)
	if err == nil && n == 4+packlen {
		fmt.Println("send success")
	}

	buf = make([]byte, 4096)
	n, err = mc.con.Read(buf)

	//package column count
	//https://dev.mysql.com/doc/internals/en/packet-OK_Packet.html 解析ok包
	columncount, _, n := readLengthEncodedInteger(buf[4:])

	fmt.Println("columncount:",columncount)
	pos := 4 + 1
	//忽略 column def package
	pos += 4 + int(buf[pos])
	//string<EOF> https://dev.mysql.com/doc/internals/en/string.html#packet-Protocol::RestOfPacketString
	//忽略 column def EOF package
	pos += 4 + int(buf[pos])
	//query packages
	// https://dev.mysql.com/doc/internals/en/com-query-response.html#packet-ProtocolText::Resultset

	fmt.Println(string(buf[pos+5:pos+5+int(buf[pos+4])]))


}

func (mc *mConnect) readPackage() ([]byte, error) {
	return []byte{}, nil
}