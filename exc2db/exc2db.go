package exc2db

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	mics "github.com/bkzy/micscript"
)

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

func (ef *ExcelFile) GetValues() ([]*DbValues, error) {
	switch ef.SyncType {
	case "cell":
		return ef.getcellsvalues()
	case "row":
		if err := ef.get_colnames_index(); err != nil {
			return nil, fmt.Errorf("ColNames in Row error:[%s]", err.Error())
		}
		return ef.getrowvalues()
	case "column":
		var err error
		ef.Column.firstcolindex, err = excelize.ColumnNameToNumber(ef.Column.FirstCol)
		if err != nil {
			return nil, fmt.Errorf("FirstCol in Column error:[%s]", err.Error())
		}
		return ef.getcolvalues()
	default:
		return nil, nil
	}
}

func (ef *ExcelFile) get_colnames_index() error {
	if ef.SyncType == "row" {
		for i, colname := range ef.Row.ColNames {
			n, err := excelize.ColumnNameToNumber(colname)
			if err != nil {
				return err
			}
			ef.Row.colindexs[i] = n - 1
		}
	}
	return nil
}

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

func (ef *ExcelFile) getcellvalues(sheetname string, i int) (*DbValues, error) {
	dbv := new(DbValues)
	dbv.TableName = ef.Cells[i].DbDest.TableName
	var values []string
	for j, axis := range ef.Cells[i].Axis {
		v, err := ef.exfile.GetCellValue(sheetname, axis)
		if err != nil {
			return dbv, err
		}
		dbv.ColNames = append(dbv.ColNames, ef.Cells[i].DbDest.ColNames[j])
		values = append(values, v)
	}

	dbv.Values = append(dbv.Values, values)

	return dbv, nil
}

func (ef *ExcelFile) getrowvalues() ([]*DbValues, error) {
	var dbvs []*DbValues
	dbv := new(DbValues)
	dbv.TableName = ef.Row.DbDest.TableName
	dbv.ColNames = ef.Row.DbDest.ColNames
	for _, sheet := range ef.Sheets {
		rows, err := ef.exfile.Rows(sheet)
		if err != nil {
			return dbvs, err
		}
		rowi := 0
		for rows.Next() {
			rowi += 1
			if rowi < ef.Row.FirstRow {
				continue
			}
			row, err := rows.Columns()
			if err != nil {
				return dbvs, fmt.Errorf("get rows.Columns fail:[%s] ", err.Error())
			}
			rl := len(row)
			var values []string
			for _, colindex := range ef.Row.colindexs {
				if colindex < rl {
					values = append(values, row[colindex])
				} else {
					values = append(values, "")
				}
			}
			dbv.Values = append(dbv.Values, values)
		}
	}
	dbvs = append(dbvs, dbv)
	return dbvs, nil
}

func (ef *ExcelFile) getcolvalues() ([]*DbValues, error) {
	var dbvs []*DbValues
	dbv := new(DbValues)
	dbv.TableName = ef.Column.DbDest.TableName
	dbv.ColNames = ef.Column.DbDest.ColNames
	for _, sheet := range ef.Sheets {
		cols, err := ef.exfile.Cols(sheet)
		if err != nil {
			return dbvs, err
		}
		coli := 0
		for cols.Next() {
			coli += 1
			if coli < ef.Column.firstcolindex {
				continue
			}
			col, err := cols.Rows()
			if err != nil {
				return dbvs, fmt.Errorf("get cols.Rows fail:[%s]", err.Error())
			}
			cl := len(col)
			var values []string
			for _, rowindex := range ef.Column.Rows {
				if rowindex <= cl {
					values = append(values, col[rowindex-1])
				} else {
					values = append(values, "")
				}
			}
			dbv.Values = append(dbv.Values, values)
		}
	}
	dbvs = append(dbvs, dbv)
	return dbvs, nil
}
