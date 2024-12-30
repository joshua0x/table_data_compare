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

var cmd = flag.String("cmd", "genmodel/diff", "生成table models")

func main() {
	flag.Parse()
	cfg := config.InitCfg()
	db.InitDB(*cfg)
	jointRes := []*db.TableLevelDiff{}
	if *cmd == "genmodel" {
		generate_model.Gen(cfg)
	} else if *cmd == "diff" {
		for _, tname := range cfg.TableList {
			tabler := tab_models.GetTable(tname)
			db.InitCache([]db.Tabler{tabler})
			difft := db.CompareTable(context.TODO(), tabler, cfg)
			jointRes = append(jointRes, difft)

		}
		res, _ := json.MarshalIndent(jointRes, " ", " ")
		fmt.Println(string(res))
	}
}
