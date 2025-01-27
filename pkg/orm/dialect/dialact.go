package dialect

import "reflect"

// Dialect 這邊聲明結口是為了方便擴充,未來只要滿足街口需求都可以調用底下函數
type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExists(tableName string) (string, []interface{})
}

// key 存 drive的名稱
var dialectsMap = map[string]Dialect{}

func RegisterDialect(d Dialect, name string) {
	dialectsMap[name] = d
}

func GetDialect(name string) (d Dialect, ok bool) {
	d, ok = dialectsMap[name]
	return
}
