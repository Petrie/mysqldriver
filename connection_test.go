package mysqldriver

import "testing"

func TestQuery(T *testing.T){
	connect()
	defer close()
	ConnHandler(mc.con, mc.cfg)
	mc.query("select * from userinfo")
}