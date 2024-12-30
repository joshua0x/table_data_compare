# mysql data compare 👏[English-Version](readme-en.md)

## 基于主键 对比不同mysql database  table 的数据(适应于数据迁移场景中使用)
- 前提条件：对比的表名 和 表结构一致。

- 配置需要对比的dsn 连接信息。如下所示, 则支持对比 host_a 和 host_b 的 user_info ,user_mail_token 数据。


```json
{

  "host_a": "user:password@(localhost:3306)/src_db?parseTime=true",
  "host_b": "user:password@(localhost:3306)/dst_db?parseTime=true",
  "table_list": ["user_info","user_mail_token"],
  "scantable_batch_size": 1000 // 扫描读取数据的 批量长度,
  "scan_sleep_period": 1 // 扫描读取数据的间隔，ms 级别，防止读取qps 过高，
}

```


- 输出的对比结果如下：
```json
[
  {
    "TableName": "user_info",
    "IdOnlyInSrc": [
      3999995
    ],
    "IdOnlyInDst": [
      5000000
    ],
    "RowGotDiff": [
      {
        "PkId": 3,
        "ColDiffs": [
          {
            "ColName": "nickname",
            "SrcVal": "aha",
            "DstVal": "1021-test"
          },
          {
            "ColName": "mtime",
            "SrcVal": "2021-08-12T09:53:02Z",
            "DstVal": "2022-10-21T22:53:46Z"
          }
        ]
      }
    ]
  }
]
```


## 使用方式
- step1 : git clone && go mod tidy && go build 
  - 生成schema ddl 定义 ，生成的.go 文件位于 ./tab_models/ 
    - ./table_data_compare -cmd=genmodel 
- step2 : 对比
    - go build 
    - ./table_data_compare -cmd=diff 
    
## 对比的性能数据 
- table 定义
```sql
CREATE TABLE `user_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(32) NOT NULL,
  `nickname` varchar(32) NOT NULL,
  `profile_photo_url` varchar(255) DEFAULT NULL,
  `email` varchar(255) NOT NULL,
  `password` char(32) DEFAULT NULL,
  `user_status` tinyint(4) DEFAULT NULL,
  `follower_cnt` int(11) NOT NULL DEFAULT '0',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `test_max` varchar(1024) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uname` (`username`),
  KEY `mtime_key` (`mtime`)
) ENGINE=InnoDB
```
- 数据行数： 2 million 
- 存在diff 的行数：2 
- 对比完成时间 : 5min

## 实现原理
- 先 通过show create table 和 词法分析来生成 table 的golang gorm 定义。(基于 https://github.com/cascax/sql2gorm 来实现的, 修改了 一些内部实现以支持 生成PrimaryKey Column Name
)
- 基于反射 和 reflect.Deepeuqal 来做行级别 字段的对比。




