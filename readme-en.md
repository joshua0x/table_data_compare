# mysql data compare 👏

## Compare the data of different mysql database tables based on the primary key
- Precondition: The table name and table structure to be compared are the same.

- Configure the dsn connection information to be compared. As shown below, it supports comparing the user_info and user_mail_token data of host_a and host_b.


````json
{

  "host_a": "user:password@(localhost:3306)/src_db?parseTime=true",
  "host_b": "user:password@(localhost:3306)/dst_db?parseTime=true",
  "table_list": ["user_info","user_mail_token"],
  "scantable_batch_size": 1000 // Batch length for scanning and reading data,
  "scan_sleep_period": 1 // Interval for scanning and reading data, ms level, to prevent reading qps from being too high,
}

````


- The output comparison results are as follows:
````json
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
````


## How to use
- step1 : git clone && go mod tidy && go build
  - Generate schema ddl definitions, the generated .go files are located in ./tab_models/
    - ./table_data_compare -cmd=genmodel
- step2 : comparison
    - go build
    - ./table_data_compare -cmd=diff
    
## Comparative performance data
- table definition
```sql
CREATE TABLE `user_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'User ID',
  `username` varchar(32) NOT NULL,
  `nickname` varchar(32) NOT NULL,
  `profile_photo_url` varchar(255) DEFAULT NULL,
  `email` varchar(255) NOT NULL,
  `password` char(32) DEFAULT NULL,
  `user_status` tinyint(4) DEFAULT NULL,
  `follower_cnt` int(11) NOT NULL DEFAULT '0',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Last modified time',
  `test_max` varchar(1024) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uname` (`username`),
  KEY `mtime_key` (`mtime`)
) ENGINE=InnoDB
````
- Number of data rows: 2 million
- Number of lines where diff exists: 2
- Comparison completion time : 5min

## Implementation principle
- First generate the golang gorm definition of table through show create table and lexical analysis. (based on https://github.com/cascax/sql2gorm, modified some internal implementations to support generating PrimaryKey Column Name
)
- Comparison of row-level fields based on reflection and reflect.Deepeuqal.