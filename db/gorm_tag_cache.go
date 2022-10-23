package db

import (
	"reflect"
)

//cache for columnname
type GormCache struct {
	PkOffset int
	ColumnName []string //offsets to columnNames ,
}

var gormCache = map[string]*GormCache{}

/*

func (l *UserInfo) Record() interface{} {
	return &UserInfo{}
}

	m := make(map[string]interface{})
	t := getRealType(reflect.TypeOf(src))
	srcv := getRealValue(reflect.ValueOf(src))
	dstv := getRealValue(reflect.ValueOf(dst))

	if t.Kind() == reflect.Struct {
		for i := 0; i < srcv.NumField(); i++ {
			//todo,get-column-name-cached-in-mem ??,
			fieldName := parseGormTag(t.Field(i).Tag.Get("gorm"))
			if fieldName == "" {
				continue
			}

*/
func InitCache(tableList []Tabler) {
	for _,tabler := range tableList {

		obj := reflect.ValueOf(tabler.Record()).Elem()
		tinfo := reflect.TypeOf(tabler.Record()).Elem()
		cache := new(GormCache)
		colNameList := make([]string,obj.NumField())
		for k := 0 ; k < obj.NumField(); k++ {
			colName := getColName(tinfo.Field(k).Tag.Get("gorm"))
			colNameList[k] = colName
		}
		cache.ColumnName = colNameList
		cache.PkOffset = getPkOffset(tabler.Record())
		gormCache[tabler.TableName()] = cache
	}
}

func getPkOffsetFromCache(tname string) int{
	return gormCache[tname].PkOffset
}

func getColListFromCache(tname string) []string{
	return gormCache[tname].ColumnName
}