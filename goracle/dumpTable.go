package goracle

import (
	// "fmt"
	"log"
	"os"
	"strconv"
	"strings"

	// "common/goracle"
	"common/goracle/connect"
)

func DumpTable(uid, sql, cols, keys, deli string) (map[string]map[string]string, error) {
	// map[string]map[string]string
	//     每行的键    列      值

	/*
		uid 	scott/oracle@goracle
		sql 	select * from user_tables
		cols 	col1,col2,col3  列的别名
		keys 	0,1 	键值序号
		deli	拼接key时的分隔符
	*/

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	rowMap := make(map[string]map[string]string)

	alias := strings.Split(strings.Replace(cols, ` `, ``, -1), ",")

	conn, err := connect.GetRawConnection(uid)
	if err != nil {
		logger.Printf("连接数据库发生错误.\n连接信息为: %s\n错误信息为: %s", uid, strings.Split(err.Error(), "\n")[1])
		return nil, err
	}
	defer conn.Close()

	cur := conn.NewCursor()
	defer cur.Close()

	err = cur.Execute(sql, nil, nil)
	if err != nil {
		logger.Printf("执行sql 语句发生错误.\nsql 语句为: %s\n错误信息为: %s", sql, strings.Split(err.Error(), "\n")[1])
		return nil, err
	}

	// 获取sql 的列别名
	columns, err := GetColumns(cur)
	if err != nil {
		logger.Printf("获取列信息发生错误.\n错误信息为: %s", strings.Split(err.Error(), "\n")[1])
		return nil, err
	}

	// records 为全部记录 records[i][j]=v
	records, err := cur.FetchMany(1000)
	for err == nil && len(records) > 0 {
		for _, record := range records {
			// 组织本行的关键字

			k := ""
			for _, p := range strings.Split(strings.Replace(keys, ` `, ``, -1), ",") {
				idx, _ := strconv.Atoi(p)
				k = k + deli + columns[idx].String(record[idx])
			}
			k = k[1:]

			// k = "" + strconv.Itoa(i)

			// 将第i行 记录映射为一个 map[string]string
			m := make(map[string]string)

			for j, v := range record {
				m[alias[j]] = columns[j].String(v)
			}

			rowMap[k] = m

		}
		records, err = cur.FetchMany(1000)
	}
	if err != nil {
		logger.Printf("获取结果集失败.\n错误信息为: %s", strings.Split(err.Error(), "\n")[1])
		return nil, err
	}

	return rowMap, nil

}
