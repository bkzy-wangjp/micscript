package exc2sql

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	mics "github.com/bkzy/micscript"
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
		return nil, fmt.Errorf("FirstCol in Column error:[%s]", err.Error())
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
	if len(ef.Sheets) == 0 {
		ef.Sheets = ef.exfile.GetSheetList()
	}
	return err
}

//根据配置获取excel中的数据
func (ef *ExcelFile) GetValues() ([]*DbValues, error) {
	switch ef.SyncType {
	case "cell":
		return ef.getcellsvalues()
	case "row":
		return ef.getrowvalues()
	case "column":
		return ef.getcolvalues()
	default:
		return nil, nil
	}
}

//解析按单元格同步时配置的单元格map
func (ef *ExcelFile) unmarshal_cell_axismaps() {
	if ef.SyncType == "cell" {
		for i, cell := range ef.Cells {
			for excel_axis, dbtb_col := range cell.Axismaps {
				ef.Cells[i].axis = append(ef.Cells[i].axis, excel_axis)
				ef.Cells[i].DbDest.colnames = append(ef.Cells[i].DbDest.colnames, dbtb_col)
			}
		}
	}
}

//解析按行同步时配置的列名map
func (ef *ExcelFile) unmarshal_row_colmaps() error {
	if ef.SyncType == "row" {
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
	if ef.SyncType == "column" {
		var err error
		//获取首列索引
		ef.Column.firstcolindex, err = excelize.ColumnNameToNumber(ef.Column.FirstCol)
		if err != nil {
			return err
		}
		//忽略列转换为索引
		for _, colname := range ef.Column.IgnoreCols {
			colid, err := excelize.ColumnNameToNumber(colname)
			if err != nil {
				return err
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
	dbvs = append(dbvs, dbv)
	return dbvs, nil
}

//获取列类型的数据
func (ef *ExcelFile) getcolvalues() ([]*DbValues, error) {
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
			col, err := cols.Rows()
			if err != nil {
				return dbvs, fmt.Errorf("get cols.Rows fail:[%s]", err.Error())
			}
			cl := len(col)
			var values []string
			for _, rowindex := range ef.Column.rows {
				if rowindex <= cl {
					values = append(values, col[rowindex-1])
				} else {
					values = append(values, "")
				}
			}
			values = append(values, constv...)
			dbv.Values = append(dbv.Values, values)
		}
	}
	dbvs = append(dbvs, dbv)
	return dbvs, nil
}
