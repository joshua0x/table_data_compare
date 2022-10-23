
package db

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)


type UserInfo struct {
	ID              int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"` // 用户ID
	Username        string    `gorm:"column:username;NOT NULL"`
	Nickname        string    `gorm:"column:nickname;NOT NULL"`
	Extra []byte `gorm:"column:extra"`
	Float1 float32 `gorm:"column:float1"`
}
//updated-bytes-fields-check,utf8,runes,Update-float-vals ,

func TestCmp(t *testing.T) {
	src := UserInfo{
		ID:              0,
		Username:        "x",
		Nickname:        "",

		Extra:           []byte(`123`),
		Float1: 1.32,
	}

	dst := UserInfo{
		ID:              2,
		Username:        "",
		Nickname:        "",

		Extra:           []byte(`xyz`),
		Float1: 2.34,
	}
	res := (compareRow(context.TODO(),src,dst,[]string{"id", "username", "nickname", "extra","float1"}))
	m,_ := json.Marshal(res)
	fmt.Println(string(m))
}