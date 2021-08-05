package exc2sql

import "fmt"

func (dbv *DbValues) Print() {
	if dbv != nil && len(dbv.ColNames) > 0 {
		maxlen := make([]int, len(dbv.ColNames))
		for i, col := range dbv.ColNames {
			maxlen[i] = len(col) + 2
		}
		for _, row := range dbv.Values {
			for i, v := range row {
				if (len(v) + 2) > maxlen[i] {
					maxlen[i] = len(v) + 2
				}
			}
		}
		colformat := make([]string, len(dbv.ColNames))
		for i := range dbv.ColNames {
			colformat[i] = fmt.Sprintf("|%%%ds", maxlen[i])
		}
		//开始打印输出
		fmt.Println()
		fmt.Print(dbv.TableName)
		printline(maxlen)
		for i, key := range dbv.ColNames {
			fmt.Printf(colformat[i], key)
		}
		fmt.Print("|")
		printline(maxlen)

		for r, row := range dbv.Values {
			if r > 0 {
				fmt.Println()
			}
			for i, v := range row {
				fmt.Printf(colformat[i], v)
			}
			fmt.Print("|")
		}
		printline(maxlen)
	}
}

func printline(linelen []int) {
	fmt.Println()
	line := ""
	for _, sl := range linelen {
		for i := 0; i <= sl; i++ {
			c := "-"
			if i == 0 {
				c = "+"
			}
			line += c
		}
	}
	line += "+"
	fmt.Println(line)
}

func (dbv *DbValues) FormatInsertSql() string {
	if dbv != nil {
		fields := ""
		for i, field := range dbv.ColNames {
			if i == len(dbv.ColNames)-1 {
				fields += field
			} else {
				fields += fmt.Sprintf("%s,", field)
			}
		}
		values := ""
		for i, rval := range dbv.Values {
			vstr := "("
			for j, val := range rval {
				if j == len(rval)-1 {
					vstr += fmt.Sprintf("'%s'", val)
				} else {
					vstr += fmt.Sprintf("'%s',", val)
				}
			}
			if i == len(dbv.Values)-1 {
				vstr += ")"
			} else {
				vstr += "),"
			}
			values += vstr
		}
		sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;", dbv.TableName, fields, values)
		return sql
	} else {
		return ""
	}
}
