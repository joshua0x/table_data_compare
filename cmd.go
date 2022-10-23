package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/joshua0x/table_data_compare/config"
	"github.com/joshua0x/table_data_compare/db"
	"github.com/joshua0x/table_data_compare/generate_model"
	"github.com/joshua0x/table_data_compare/tab_models"
)

//yaml cfgs , hostA,hostB ,指定table-list ,scan-all-tables and compared,
/*
	1.getTables schema , converted to golang ?? sql to go ?, input is ddl-str , output is DDL-go ,


	2.对比流程参考当前的不变，处理下写死ID , 通过
		2.1 config-info-parsed in-yaml ,hosta,hostb ,table-list ,
		2.2 Db-conns-Init ,sqls,
*/

var cmd = flag.String("cmd", "genmodel/diff", "生成table models")



func main() {
	//
	flag.Parse()
	cfg := config.InitCfg()
	db.InitDB(*cfg)
	jointRes := []*db.TableLevelDiff{}
	if *cmd == "genmodel" {
		generate_model.Gen(cfg)
	}else if *cmd == "diff"{
		//compare-rows,
		//init-Caches,
		for _,tname := range cfg.TableList {
			tabler := tab_models.GetTable(tname)
			db.InitCache([]db.Tabler{tabler})
			difft := db.CompareTable(context.TODO(),tabler,cfg)
			jointRes = append(jointRes,difft)

		}
		//
		res,_ := json.MarshalIndent(jointRes," "," ")
		fmt.Println(string(res))
	}
}
