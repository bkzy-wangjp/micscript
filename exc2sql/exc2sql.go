package exc2sql

import (
	"encoding/json"
	"fmt"
	"strconv"

	mics "github.com/bkzy/micscript"
	"github.com/xuri/excelize/v2"
)

func ExcelToSql(cfg string) ([]string, error) {
	ef, err := NewExcelFile(cfg)
	if err != nil {
		return nil, err
	}
	err = ef.OpenFile()
	if err != nil {
		return nil, err
	}
	dbvs, err := ef.GetValues()
	if err != nil {
		return nil, err
	}
	var sqls []string
	for _, dbv := range dbvs {
		sqls = append(sqls, dbv.FormatInsertSql())
	}
	return sqls, nil
}

//通过配置文件创建对象
func NewExcelFile(cfg string) (*ExcelFile, error) {
	ef := new(ExcelFile)
	err := json.Unmarshal([]byte(cfg), ef)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config json fail:[%s]", err.Error())
	}
	ef.unmarshal_cell_axismaps()
	if err := ef.unmarshal_row_colmaps(); err != nil {
		return nil, fmt.Errorf("ColNames in Row error:[%s]", err.Error())
	}
	if err := ef.unmarshal_column_rowmaps(); err != nil {
		return nil, err
	}
	return ef, nil
}

//打开excel文件
func (ef *ExcelFile) OpenFile() error {
	var err error
	if len(ef.Password) > 0 {
		ef.exfile, err = excelize.OpenFile(ef.FileName, excelize.Options{Password: ef.Password})
	} else {
		ef.exfile, err = excelize.OpenFile(ef.FileName)
	}
	if err != nil {
		err = fmt.Errorf("open file fail:[%s]", err.Error())
		return err
	}
	if len(ef.Sheets) == 0 {
		ef.Sheets = ef.exfile.GetSheetList()
	}
	return nil
}

//根据配置获取excel中的数据
func (ef *ExcelFile) GetValues() ([]*DbValues, error) {
	var dbvs []*DbValues
	cells, err := ef.getcellsvalues()
	if err != nil {
		return nil, err
	}
	for _, v := range cells {
		if v != nil {
			dbvs = append(dbvs, v)
		}
	}

	rows, err := ef.getrowvalues()
	if err != nil {
		return nil, err
	}
	for _, v := range rows {
		if v != nil {
			dbvs = append(dbvs, v)
		}
	}
	cols, err := ef.getcolvalues()
	if err != nil {
		return nil, err
	}
	for _, v := range cols {
		if v != nil {
			dbvs = append(dbvs, v)
		}
	}
	//防止一次性插入过多数据
	maxlen := 65000
	for i, v := range dbvs {
		if len(v.ColNames)*len(v.Values) > maxlen {
			mrow := maxlen / len(v.ColNames)
			var rslts []*DbValues
			for j := 0; j < len(v.Values); j += mrow {
				rslt := new(DbValues)
				rslt.TableName = v.TableName
				rslt.ColNames = v.ColNames
				ed := j + mrow
				if ed > len(v.Values) {
					ed = len(v.Values)
				}
				rslt.Values = v.Values[j:ed]
				rslts = append(rslts, rslt)
			}
			dbvs = append(dbvs[:i], dbvs[i+1:]...)
			dbvs = append(dbvs, rslts...)
		}
	}
	return dbvs, nil
}

//解析按单元格同步时配置的单元格map
func (ef *ExcelFile) unmarshal_cell_axismaps() {
	for i, cell := range ef.Cells {
		if cell != nil {
			for excel_axis, dbtb_col := range cell.Axismaps {
				ef.Cells[i].axis = append(ef.Cells[i].axis, excel_axis)
				ef.Cells[i].DbDest.colnames = append(ef.Cells[i].DbDest.colnames, dbtb_col)
			}
		}
	}
}

//解析按行同步时配置的列名map
func (ef *ExcelFile) unmarshal_row_colmaps() error {
	if ef.Row != nil {
		for excel_col, dbtb_col := range ef.Row.Colmaps {
			n, err := excelize.ColumnNameToNumber(excel_col)
			if err != nil {
				return err
			}
			ef.Row.colnames = append(ef.Row.colnames, excel_col)
			ef.Row.DbDest.colnames = append(ef.Row.DbDest.colnames, dbtb_col)
			ef.Row.colindexs = append(ef.Row.colindexs, n-1)
		}
	}
	return nil
}

//解析按行同步时配置的列名map
func (ef *ExcelFile) unmarshal_column_rowmaps() error {
	if ef.Column != nil {
		var err error
		//获取首列索引
		if ef.Column.FirstCol != "" {
			ef.Column.firstcolindex, err = excelize.ColumnNameToNumber(ef.Column.FirstCol)
			if err != nil {
				return fmt.Errorf("FirstCol in Column error:[%s]", err.Error())
			}
		} else {
			ef.Column.firstcolindex = 0
		}
		//获取末尾列索引
		if ef.Column.LastCol != "" {
			ef.Column.lastcolindex, err = excelize.ColumnNameToNumber(ef.Column.LastCol)
			if err != nil {
				return fmt.Errorf("LastCol in Column error:[%s]", err.Error())
			}
		} else {
			ef.Column.lastcolindex = 0
		}
		//指定转换列转换为索引
		for _, colname := range ef.Column.LoadCols {
			colid, err := excelize.ColumnNameToNumber(colname)
			if err != nil {
				return fmt.Errorf("LoadCols in Column error:[%s]", err.Error())
			}
			ef.Column.loadcolindexs = append(ef.Column.loadcolindexs, colid)
		}
		//忽略列转换为索引
		for _, colname := range ef.Column.IgnoreCols {
			colid, err := excelize.ColumnNameToNumber(colname)
			if err != nil {
				return fmt.Errorf("IgnoreCols in Column error:[%s]", err.Error())
			}
			ef.Column.ignorecolindexs = append(ef.Column.ignorecolindexs, colid)
		}
		//获取行map
		for excel_row_id, dbtb_col := range ef.Column.Rowmaps {
			rowid, err := strconv.ParseInt(excel_row_id, 10, 64)
			if err != nil {
				return fmt.Errorf("Column.Rowmaps key must be number:[%s]", err.Error())
			}
			ef.Column.rows = append(ef.Column.rows, int(rowid))
			ef.Column.DbDest.colnames = append(ef.Column.DbDest.colnames, dbtb_col)
		}
	}
	return nil
}

//获取单元格类型的数据
func (ef *ExcelFile) getcellsvalues() ([]*DbValues, error) {
	var dbvs []*DbValues
	for _, sheet := range ef.Sheets {
		var dbv *DbValues
		for i := range ef.Cells {
			vs, err := ef.getcellvalues(sheet, i)
			if err != nil {
				return dbvs, err
			}
			if i == 0 {
				dbv = vs
			} else {
				same := false
				if dbv.TableName == vs.TableName { //数据表相同
					if mics.StringsCompiler(dbv.ColNames, vs.ColNames) { //字段相同
						same = true
					}
				}
				if same {
					dbv.Values = append(dbv.Values, vs.Values...)
				} else {
					dbvs = append(dbvs, dbv)
					dbv = vs
				}
			}
		}
		dbvs = append(dbvs, dbv)
	}
	return dbvs, nil
}

//获取单个单元格的数据
func (ef *ExcelFile) getcellvalues(sheetname string, i int) (*DbValues, error) {
	dbv := new(DbValues)
	dbv.TableName = ef.Cells[i].DbDest.TableName
	var values []string
	for j, axis := range ef.Cells[i].axis {
		v, err := ef.exfile.GetCellValue(sheetname, axis)
		if err != nil {
			return dbv, err
		}
		colname := ef.Cells[i].DbDest.colnames[j]
		dbv.ColNames = append(dbv.ColNames, colname)
		values = append(values, v)
	}
	for k, v := range ef.Cells[i].DbDest.Consts {
		dbv.ColNames = append(dbv.ColNames, k)
		values = append(values, v)
	}
	dbv.Values = append(dbv.Values, values)

	return dbv, nil
}

//获取行类型的数据
func (ef *ExcelFile) getrowvalues() ([]*DbValues, error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	if ef.Row != nil {
		var dbvs []*DbValues
		dbv := new(DbValues)
		dbv.TableName = ef.Row.DbDest.TableName
		dbv.ColNames = ef.Row.DbDest.colnames
		var constv []string
		//取常数列
		for k, v := range ef.Row.DbDest.Consts {
			dbv.ColNames = append(dbv.ColNames, k)
			constv = append(constv, v)
		}
		for _, sheet := range ef.Sheets { //遍历sheet
			//设定了读取行的列表
			if len(ef.Row.LoadRows) > 0 {
				//读取每一行
				rows, err := ef.exfile.GetRows(sheet)
				if err != nil {
					return nil, err
				}
				rcnt := len(rows) //总行数
				//按照指定的行号遍历
				for _, rid := range ef.Row.LoadRows {
					var values []string
					if rid > rcnt { //如果行号大于总行数
						v := make([]string, len(ef.Row.colindexs))
						values = append(values, v...) //填写为空值
					} else { //行号不大于总行数
						row := rows[rid-1] //获取行号指定的行数据
						rl := len(row)     //计算行的列数
						//按照指定的列索引，在当前行读取指定列
						for _, colindex := range ef.Row.colindexs {
							if colindex < rl { //列索引不大于当前行的列数
								values = append(values, row[colindex])
							} else { //列索引大于当前行的列数
								values = append(values, "") //填写为空
							}
						}
					}
					//在值列表中添加常数项
					values = append(values, constv...)
					dbv.Values = append(dbv.Values, values)
				}
			} else { //没有指定具体的行
				rows, err := ef.exfile.Rows(sheet)
				if err != nil {
					return nil, err
				}
				rowi := 0
				//遍历行
				for rows.Next() {
					rowi += 1 //行号
					row, err := rows.Columns()
					if err != nil {
						return nil, fmt.Errorf("get rows.Columns fail:[%s] ", err.Error())
					}
					//行号小于首行
					if rowi < ef.Row.FirstRow {
						continue
					}
					//如果指定了最后一行
					if ef.Row.LastRow > 0 && rowi > ef.Row.LastRow {
						break
					}
					//行号是忽略行
					if mics.IsExistItem(rowi, ef.Row.IgnoreRows) {
						continue
					}
					rl := len(row)
					var values []string
					//根据设定的列索引，逐列取值
					for _, colindex := range ef.Row.colindexs {
						if colindex < rl {
							values = append(values, row[colindex])
						} else {
							values = append(values, "")
						}
					}
					//在值列表中添加常数项
					values = append(values, constv...)
					dbv.Values = append(dbv.Values, values)
				}
			}
		}
		dbvs = append(dbvs, dbv)
		return dbvs, nil
	} else {
		return nil, nil
	}

}

//获取列类型的数据
func (ef *ExcelFile) getcolvalues() ([]*DbValues, error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	if ef.Column != nil {
		var dbvs []*DbValues
		dbv := new(DbValues)
		dbv.TableName = ef.Column.DbDest.TableName
		dbv.ColNames = ef.Column.DbDest.colnames
		var constv []string
		//取常数列
		for k, v := range ef.Column.DbDest.Consts {
			dbv.ColNames = append(dbv.ColNames, k)
			constv = append(constv, v)
		}
		//遍历工作表
		for _, sheet := range ef.Sheets {
			//设定了读取列的列表
			if len(ef.Column.LoadCols) > 0 {
				cols, err := ef.exfile.GetCols(sheet)
				if err != nil {
					return dbvs, err
				}
				ccnt := len(cols) //总列数
				//按照指定的列序号遍历
				for _, colid := range ef.Column.loadcolindexs {
					var values []string
					if colid > ccnt { //如果列号大于总列数
						v := make([]string, len(ef.Column.rows))
						values = append(values, v...) //填写为空值
					} else { //列号不大于总列数
						col := cols[colid-1] //获取列号指定的列数据
						rl := len(col)       //计算列的行数
						//按照指定的行id，在当前列读取指定行
						for _, rid := range ef.Column.rows {
							if rid <= rl { //列索引不大于当前列的行数
								values = append(values, col[rid-1])
							} else { //列索引大于当前列的行数
								values = append(values, "") //填写为空
							}
						}
					}
					//在值列表中添加常数项
					values = append(values, constv...)
					dbv.Values = append(dbv.Values, values)
				}
			} else {
				cols, err := ef.exfile.Cols(sheet)
				if err != nil {
					return dbvs, err
				}
				coli := 0
				//遍历列
				for cols.Next() {
					coli += 1
					//未到起始列
					if coli < ef.Column.firstcolindex {
						continue
					}
					//是否忽略列
					if mics.IsExistItem(coli, ef.Column.ignorecolindexs) {
						continue
					}
					//是否到了设定的最后一列
					if ef.Column.lastcolindex > 0 && coli > ef.Column.lastcolindex {
						break
					}
					col, err := cols.Rows()
					if err != nil {
						return dbvs, fmt.Errorf("get cols.Rows fail:[%s]", err.Error())
					}
					cl := len(col)
					var values []string
					for _, rid := range ef.Column.rows {
						if rid <= cl {
							values = append(values, col[rid-1])
						} else {
							values = append(values, "")
						}
					}
					values = append(values, constv...)
					dbv.Values = append(dbv.Values, values)
				}
			}
		}
		dbvs = append(dbvs, dbv)
		return dbvs, nil
	} else {
		return nil, nil
	}
}
