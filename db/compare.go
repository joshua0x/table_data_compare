package db

import (
	"context"
	"fmt"
	"github.com/joshua0x/table_data_compare/config"
	"reflect"
	"strings"
	"time"
)
//??
type Tabler interface {
	TableName() string
	RecordList() interface{}
	Record() interface{}
	PkName() string
}

//reported-results ??,pk-,parser-checks,
type Splited struct {
	//?done,Id maybe string
	OnlyInSrcIdList []interface{}
	OnlyInDstIdList []interface{}
	NeedUpdate []interface{}
}

func CompareTable(ctx context.Context,table Tabler,cfg *config.DbConfig) *TableLevelDiff{
	report := new(TableLevelDiff)
	report.TableName = table.TableName()

	srcpkList := table.RecordList()
	dstpkList := table.RecordList()

	srcdb,dstdb := GetDb()

	srcdb.Table(table.TableName()).Select(table.PkName()).Find(srcpkList)
	dstdb.Table(table.TableName()).Select(table.PkName()).Find(dstpkList)
	res := divideByGroup(table.TableName(),srcpkList,dstpkList)

	diff := checkdiff(ctx,table,res.NeedUpdate,cfg.ScantablebatchSize,cfg.ScanSleepPeriod)
	report.RowGotDiff = diff
	setUpdiff(report,&res)
	return report
}

func getPkList(src interface{},pkOffset int) map[interface{}]bool {
	srcMap := map[interface{}]bool{}

	//must-be-slice ?, []IdObj, return first-field ,Ptr *[]-> [],
	convsrc := reflect.ValueOf(src).Elem()
	for k:=0 ; k <  convsrc.Len() ;k++{
		structObj := convsrc.Index(k).Elem() // struct,
		idVal := structObj.Field(pkOffset)
		//convert to uint64 ??,
		idInt := idVal.Interface()
		srcMap[idInt] = true
	}
	return srcMap
}

func divideByGroup(tname string,src, dst interface{}) Splited {
	ret := Splited{OnlyInSrcIdList: []interface{}{}, OnlyInDstIdList: []interface{}{}, NeedUpdate: []interface{}{}}

	srcMap := map[interface{}]bool{}
	dstMap := map[interface{}]bool{}
	//must-be-slice ?, []IdObj, return first-field ,Ptr *[]-> [],
	pkOffset := getPkOffsetFromCache(tname)
	srcMap = getPkList(src,pkOffset)
	dstMap = getPkList(dst,pkOffset)
	//Pk-is-string or int ??,
	for k, _ := range srcMap {
		if _, exist := dstMap[k]; !exist {
			ret.OnlyInSrcIdList = append(ret.OnlyInSrcIdList, k)
		} else {
			//exist,
			ret.NeedUpdate = append(ret.NeedUpdate, k)
		}
	}

	for k, _ := range dstMap {
		if _, exist := srcMap[k]; !exist {
			ret.OnlyInDstIdList = append(ret.OnlyInDstIdList, k)
		}
	}
	return ret
}

func convertListToMap(srcList reflect.Value) map[interface{}]interface{} {
	ret := map[interface{}]interface{}{}

	for k := 0; k < srcList.Len(); k++ {
		//range-over-fields,
		pk := getPk(srcList.Index(k).Interface())
		ret[pk] = srcList.Index(k).Interface()
	}
	return ret
}

func getColName(tag string) string {
	if tag == "" {
		return ""
	}
	if tagArr := strings.Split(tag, ";"); len(tagArr) >= 1 {
		spl := strings.Split(tagArr[0],":")
		if spl[0] == "column" {
			return spl[1]
		}
	}
	return ""
}

func checkPk(tag string) bool {
	return strings.Index(tag,"primary_key") >= 0
}

func getPkOffset(src interface{}) int {
	srcv := reflect.ValueOf(src).Elem()
	t := reflect.TypeOf(src).Elem()

	for i := 0; i < srcv.NumField(); i++ {
		fieldName := (t.Field(i).Tag.Get("gorm"))
		if checkPk(fieldName) {
			return i
		}
	}

	return 0
}

func getPk(src interface{}) interface{} {
	srcv := reflect.ValueOf(src).Elem()
	t := reflect.TypeOf(src).Elem()

	if t.Kind() == reflect.Struct {
		for i := 0; i < srcv.NumField(); i++ {
			fieldName := (t.Field(i).Tag.Get("gorm"))
			if checkPk(fieldName) {
				return srcv.Field(i).Interface()
			}
		}
	}
	return nil
}

func checkdiff(ctx context.Context,table Tabler,srcPkIdList []interface{},batchsize,sleepPeriod int) []*RowLevelDiff{
	diff := []*RowLevelDiff{}

	start := 0
	srcdb,dstdb := GetDb()
	//sort.Slice(srcPkIdList, func(i, j int) bool {
	//	return srcPkIdList[i] < srcPkIdList[j]
	//})
	colList := getColListFromCache(table.TableName())
	round := 0
	totalRound := len(srcPkIdList) / batchsize

	for {
		round++
		if round % 10 == 0 {
			fmt.Printf("table:%v ts:%v progess:%v/%v=%.2f%% \n",table.TableName(),time.Now().Format(time.Stamp),round,totalRound,100*float32(round) / float32(totalRound ))
		}
		if start >= len(srcPkIdList) {
			break
		}
		//select compare,
		srcList := table.RecordList()
		dstList := table.RecordList()

		groupList := make([]interface{}, 0)
		//[start,start+pageSize),
		for index := start; index < start+batchsize; index++ {
			if index >= len(srcPkIdList) {
				break
			}
			groupList = append(groupList, srcPkIdList[index])
		}

		if err := srcdb.Table(table.TableName()).Where(table.PkName()+" in (?"+strings.Repeat(",?", len(groupList)-1)+")", groupList...).Find(srcList).Error; err != nil {
			panic(err)
		}

		if err := dstdb.Table(table.TableName()).Where( table.PkName() + " in (?"+strings.Repeat(",?", len(groupList)-1)+")", groupList...).Find(dstList).Error; err != nil {
			panic(err)
		}
		//compared-slice-obj,长度不一致 异常情况，需要考虑，
		srcListConverted := reflect.ValueOf(srcList).Elem() // []record ,
		dstListConverted := reflect.ValueOf(dstList).Elem() // []record,

		srcMap := convertListToMap(srcListConverted)
		dstMap := convertListToMap(dstListConverted)
		//debug-modes

		for pk, record := range srcMap {

			if dstRecord, exist := dstMap[pk]; !exist {
				continue
			} else {

				updmap := compareRow(ctx, record, dstRecord,colList)
				updmap.PkId = pk
				if len(updmap.ColDiffs) > 0 {
					diff = append(diff,updmap)
				}
				//
			}
		}
		start += batchsize
		time.Sleep(time.Millisecond*time.Duration(sleepPeriod))
	}
	return diff

}

func getRealType(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}
func getRealValue(val reflect.Value) reflect.Value {
	for val.Kind() == reflect.Ptr  {
		val = val.Elem()
	}
	return val
}

func compareRow(ctx context.Context, src, dst interface{},colList []string) *RowLevelDiff {
	m := new(RowLevelDiff)
	m.ColDiffs = []*ColDiff{}
	srcv := getRealValue(reflect.ValueOf(src))
	dstv := getRealValue(reflect.ValueOf(dst))

	for i := 0; i < srcv.NumField(); i++ {
		//todo,get-column-name-cached-in-mem ??,
		fieldName := colList[i]
		//
		srcFieldV := srcv.Field(i)
		dstFieldV := dstv.Field(i)
		if !srcFieldV.CanInterface() || !dstFieldV.CanInterface() {
			continue
		}
		if reflect.DeepEqual(srcFieldV.Interface(), dstFieldV.Interface()) {
		} else {
			//m[fieldName] = srcFieldV.Interface()
			m.ColDiffs = append(m.ColDiffs,&ColDiff{
				ColName: fieldName,
				SrcVal:  srcFieldV.Interface(),
				DstVal:  dstFieldV.Interface(),
			})
		}
	}
	return m

}
