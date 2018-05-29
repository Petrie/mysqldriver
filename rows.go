package mysqldriver


type value interface {}

type mysqlField struct {
	name      string
	flags     fieldFlag
	fieldType fieldType
}


type rows struct {
	columns []mysqlField
	mc *mConnect
	val []value
}

