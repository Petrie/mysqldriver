package mysqldriver

import (
	"fmt"
	"encoding/binary"
)

func (mc *mConnect) readPacket() ([]byte, error) {
	if mc.readDataLength == 0 {
		mc.readDataIndex = 0
		n, err := mc.con.Read(mc.readData)
		if err != nil {
			return nil, err
		}
		mc.readDataLength = n

	}
	//package header
	pctLen := int(uint32(mc.readData[mc.readDataIndex]) |
		uint32(mc.readData[mc.readDataIndex+1])<<8 |
		uint32(mc.readData[mc.readDataIndex+2])<<16)
	//log.Println("sequence:",mc.readData[3])

	if pctLen > mc.readDataLength && mc.readDataIndex != 0 {
		copy(mc.readData[:mc.readDataLength], mc.readData[mc.readDataIndex:mc.readDataLength])
		mc.readDataIndex = 0
		n, err := mc.con.Read(mc.readData[mc.readDataLength:])
		if err != nil {
			return nil, err
		}
		mc.readDataLength += n

	}
	mc.sequence++

	data := mc.readData[mc.readDataIndex+4 : mc.readDataIndex+4+pctLen]
	//这里暂时存在一个maxpacketsize相关的bug,涉及到扩充 readData的长度，这里后续再说。
	mc.readDataIndex = mc.readDataIndex + 4 + pctLen
	mc.readDataLength = mc.readDataLength - 4 - pctLen

	return data, nil
}

func (mc *mConnect) readColumn() ( []mysqlField, error) {
	data, _ := mc.readPacket()
	num, _, _ := readLengthEncodedInteger(data)
	columns := make([]mysqlField,num)
	for i:=0 ; ;i++ {
		data, err := mc.readPacket()
		//fmt.Println(data)
		if err != nil {
			return nil, err
		}

		if data[0] == iEOF && (len(data) == 5 || len(data) == 1) {
			if i == int(num) {
				return columns, nil
			}
			return nil, fmt.Errorf("column count mismatch n:%d len:%d", num, len(columns))
		}

		pos, err := skipLengthEncodedString(data)
		n, err := skipLengthEncodedString(data[pos:])
		pos += n
		n, err = skipLengthEncodedString(data[pos:])
		pos += n

		n, err = skipLengthEncodedString(data[pos:])
		pos += n
		name, _, n, err := readLengthEncodedString(data[pos:])
		columns[i].name = string(name)

		// Original name [len coded string]
		n, err = skipLengthEncodedString(data[pos:])
		if err != nil {
			return nil, err
		}
		pos += n

		// Filler [uint8]
		pos++

		// Charset [charset, collation uint8]
		pos += 2

		// Length [uint32]
		pos += 4

		// Field type [uint8]
		columns[i].fieldType = fieldType(data[pos])
		pos++

		// Flags [uint16]
		columns[i].flags = fieldFlag(binary.LittleEndian.Uint16(data[pos : pos+2]))
		pos += 2

	}
	return columns, nil
}

func (mc *mConnect) readRow(val []string) {
	data, _ := mc.readPacket()
	// NULL-bitmap,  [(column-count + 7 + 2) / 8 bytes]
	pos := 0
	for i := range val {
		// Read bytes and convert to string
		d, _, n, _ := readLengthEncodedString(data[pos:])
		val[i] = string(d)
		//fmt.Println(string(d))
		pos += n
	}
}