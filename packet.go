package mysqldriver

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

	if pctLen > mc.readDataLength && mc.readDataIndex != 0{
		copy(mc.readData[:mc.readDataLength], mc.readData[mc.readDataIndex:mc.readDataLength])
		mc.readDataIndex = 0
		n, err := mc.con.Read(mc.readData[mc.readDataLength:])
		if err != nil {
			return nil, err
		}
		mc.readDataLength += n

	}
	mc.sequence++

	data := mc.readData[mc.readDataIndex + 4: mc.readDataIndex + 4+ pctLen]
    //这里暂时存在一个maxpacketsize相关的bug,涉及到扩充 readData的长度，这里后续再说。
	mc.readDataIndex = mc.readDataIndex + 4 + pctLen
	mc.readDataLength = mc.readDataLength - 4 - pctLen

	return  data, nil

}