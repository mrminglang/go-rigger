package news_main_repository

import (
	"fmt"
	"github.com/mrminglang/go-rigger/connect/gmysql"
)

// QueryRecord 测试
func QueryRecord(sSql string, vMysqlData *[]map[string]string) error {

	rows, err := gmysql.DbMySQLConn.Query(sSql)
	if err != nil {
		return err
	}

	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	// 数据列
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// 列的个数
	count := len(columns)

	// 一条数据的各列的值（需要指定长度为列的个数，以便获取地址）
	values := make([]interface{}, count)
	// 一条数据的各列的值的地址
	valPointers := make([]interface{}, count)
	for rows.Next() {
		// 获取各列的值的地址
		for i := 0; i < count; i++ {
			valPointers[i] = &values[i]
		}
		// 获取各列的值，放到对应的地址中
		rows.Scan(valPointers...)
		// 一条数据的Map (列名和值的键值对)
		entry := make(map[string]string)

		// Map 赋值
		for i, col := range columns {
			// 值复制给val(所以Scan时指定的地址可重复使用)

			val := values[i]
			fmt.Println("111111:::", i, col, val)
			b, ok := val.([]byte)
			if ok {
				// 字符切片转为字符串
				entry[col] = string(b)
				fmt.Println("222222:::", i, col, string(b))
			} else {
				entry[col] = fmt.Sprintf("%v", values[i])
				if entry[col] == "<nil>" {
					entry[col] = ""
				}
				fmt.Println("333333:::", i, col, entry[col], values[i], fmt.Sprintf("%v", values[i]))
			}
		}

		*vMysqlData = append(*vMysqlData, entry)
	}

	fmt.Println(fmt.Sprintf("sql:%s, result size:%d", sSql, len(*vMysqlData)), vMysqlData)

	return nil
}
