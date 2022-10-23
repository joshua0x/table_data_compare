package config

import (
	"encoding/json"
	"os"
)

type DbConfig struct {
	HostA     string   `json:"host_a"`
	HostB     string   `json:"host_b"`
	TableList []string `json:"table_list"`
	ScantablebatchSize int `json:"scantable_batch_size"`
	ScanSleepPeriod int `json:"scan_sleep_period"`
}


func InitCfg() *DbConfig{
	//json.Unmarshal(),
	cfg,err := os.ReadFile("./config/conf.json")
	if err != nil {
		panic(err)
	}

	ret := new(DbConfig)
	json.Unmarshal(cfg,ret)
	return ret
}