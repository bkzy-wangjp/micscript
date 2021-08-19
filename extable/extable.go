// extable project extable.go
package extable

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

/********************************************
功能:打开Excel文件
输入:[filename] 文件名（含路径）
	[calcformula] 打开时是否计算Excel公式
	[calcdeep int] 计算公式时候的循环次数,0时不限制
	[sheet] 工作表名(可选,不填写是默认Sheet1)
输出:错误信息
说明:
编辑:wang_jp
时间:2020年5月29日
*******************************************/
func (t *Table) OpenFile(filename string, calcformula bool, calcdeep int, sheet ...string) error {
	sheetname := "Sheet1" //默认打开Sheet1
	if len(sheet) > 0 {
		sheetname = sheet[0]
	}
	f, err := excelize.OpenFile(filename) //打开excel
	if err != nil {
		return fmt.Errorf("打开文件失败:[%s]", err.Error())
	}
	//获取合并的单元格
	mcs, err := f.GetMergeCells(sheetname)
	if err != nil {
		return fmt.Errorf("读取工作表的合并单元格失败:[%s]", err.Error())
	}
	//获取单元格所有行列的值
	rows, err := f.GetRows(sheetname)
	if err != nil {
		return fmt.Errorf("读取工作表内容失败:[%s]", err.Error())
	}
	//var totalwidth float64 = 0
	for ir, row := range rows { //遍历所有行
		tr := new(TableRow)                       //新建行
		h, err := f.GetRowHeight(sheetname, ir+1) //读取行高
		if err != nil {                           //读取行高失败
			h = 20 //默认行高
		}
		v, err := f.GetRowVisible(sheetname, ir+1) //读取行可见性
		if err != nil {                            //读取失败
			v = true //默认可见
		}
		tr.Init(h, ir+1, v)          //初始化行
		t.Rows = append(t.Rows, *tr) //添加行到表格中
		var cellsInRow []TableCell
		for ic, colCell := range row { //遍历行中的每一列
			caxis, _ := excelize.ColumnNumberToName(ic + 1)
			w, err := f.GetColWidth(sheetname, caxis) //读取列宽
			if err != nil {
				w = 100 //设置默认列宽
			}
			v, err := f.GetColVisible(sheetname, caxis) //读取列可见性
			if err != nil {                             //读取失败
				v = true //默认可见
			}
			//第一行的时候读取列属性
			if ir == 0 { //是第一行
				tc := new(TableCol)
				tc.Init(caxis, w, v)
				t.Cols = append(t.Cols, *tc) //添加列到表格中
			}
			axis, _ := excelize.CoordinatesToCellName(ic+1, ir+1) //坐标号
			styleindex, err := f.GetCellStyle(sheetname, axis)    //获取样式索引
			if err != nil {
				styleindex = 0
			}
			formula, err := f.GetCellFormula(sheetname, axis) //获取公式
			if err != nil {
				formula = ""
			}

			//解析设置单元格
			cell := new(TableCell)
			cell.Axis = axis

			cell.Width = w               //列宽
			cell.Visible = v             //所在列的可见性
			cell.StyleIndex = styleindex //样式索引
			cell.IsCellInMerged(mcs)

			if calcformula && len(formula) > 0 {
				if !cell.IsInMerged || (cell.IsInMerged && cell.IsMergeStart) {
					cell.NeedCalcFormula = true
				}
			}
			if stringContains(colCell, "#DIV/0!", "#NAME?", "#N/A", "#NUM!", "#VALUE!", "#REF!", "#NULL", "#SPILL!", "#CALC!", "#GETTING_DATA") {
				var str string
				colCell = str
			}
			cell.Value = colCell   //值
			cell.Formula = formula //公式

			cellsInRow = append(cellsInRow, *cell)
		}
		t.Cells = append(t.Cells, cellsInRow)
	}
	t.RowCnt = len(t.Rows)
	t.ColCnt = len(t.Cells)
	t.FindHead() //设置表头行标记
	if calcformula {
		t.CalcFormula(f, sheetname, calcdeep) //计算公式
	}

	return f.Save()
}

/********************************************
功能:计算单元格的公式
输入:[f *excelize.File] 文件指针
	[sheetname string] Sheet名称
	[calcdeep int] 计算公式时候的循环次数,0时不限制
输出:计算结果保存到结构中
说明:
编辑:wang_jp
时间:2020年5月29日
*******************************************/
func (t *Table) CalcFormula(f *excelize.File, sheetname string, calcdeep int) {
	var haveuncalc bool               //还有未被计算的单元格
	for row, tbrow := range t.Cells { //遍历行
		for col, cell := range tbrow { //遍历列
			if !cell.NeedCalcFormula || len(cell.Formula) == 0 { //不含需要计算的公式
				continue
			}
			cellstring, err := f.CalcCellValue(sheetname, cell.Axis) //计算单元格的公式
			if err != nil {                                          //计算结果有错误
				t.Cells[row][col].Err = err.Error()
				if !stringContains(err.Error(), "not support") {
					t.Cells[row][col].CalcFormulaScanCnt++                               //扫描次数自加
					if calcdeep > 0 && t.Cells[row][col].CalcFormulaScanCnt > calcdeep { //扫描次数大于设定值,不再尝试计算
						t.Cells[row][col].NeedCalcFormula = false
					} else {
						haveuncalc = true //设置有尚未计算的单元格标记
					}
				}
			} else { //计算结果没有错误
				t.Cells[row][col].Err = ""
				t.Cells[row][col].NeedCalcFormula = false
			}

			if stringContains(cellstring, "#DIV/0!", "#NAME?", "#N/A", "#NUM!", "#VALUE!", "#REF!", "#NULL", "#SPILL!", "#CALC!", "#GETTING_DATA") {
				var str string
				cellstring = str
			}
			t.Cells[row][col].Value = cellstring
			f.SetCellValue(sheetname, cell.Axis, cellstring)
		}
	}
	if haveuncalc {
		haveuncalc = false
		t.CalcFormula(f, sheetname, calcdeep)
	}
}

// func (t *Table) CalcFormula() {
// 	reg_func := regexp.MustCompile(_FUNC)
// 	reg_cell := regexp.MustCompile(_cell)
// 	reg_func_para := regexp.MustCompile(fmt.Sprintf("%s|%s", _cell, _number)) //Excel函数参数

// 	var haveuncalc bool               //还有未被计算的单元格
// 	for row, tbrow := range t.Cells { //遍历行
// 		for col, cell := range tbrow { //遍历列
// 			if cell.NeedCalcFormula == false || len(cell.Formula) == 0 { //不含需要计算的公式
// 				continue
// 			}

// 			var tokens string //需要送给计算器的元组字符串
// 			var uncalc bool   //依赖单元格未被计算
// 			exfm := new(Formula)
// 			ok := exfm.GetTokens(cell.Formula) //提取公式中的元组
// 			if ok == false {                   //不可计算
// 				t.Cells[row][col].Value = fmt.Sprint("#VALUE!")
// 				t.Cells[row][col].NeedCalcFormula = false
// 				continue
// 			}
// 			//fmt.Println(exfm.TokenStr)          //==========================
// 			for _, tok := range exfm.TokenStr { //遍历提取到的元组
// 				var tokvalue string

// 				if reg_func.MatchString(tok) { //匹配Excel函数
// 					var functype string
// 					//判断函数类型
// 					if strings.Contains(tok, "SUM") {
// 						functype = "SUM"
// 					} else {
// 						if strings.Contains(tok, "AVERAGE") {
// 							functype = "AVERAGE"
// 						} else {
// 							if strings.Contains(tok, "MIN") {
// 								functype = "MIN"
// 							} else {
// 								if strings.Contains(tok, "MAX") {
// 									functype = "MAX"
// 								} else {
// 									if strings.Contains(tok, "MEDIAN") {
// 										functype = "MEDIAN"
// 									} else {
// 										if strings.Contains(tok, "PRODUCT") {
// 											functype = "PRODUCT"
// 										}
// 									}
// 								}
// 							}
// 						}
// 					}
// 					paras := reg_func_para.FindAllString(tok, -1) //提取函数参数
// 					var vals numgo.Array
// 					for _, p := range paras {
// 						vstr, uc, canuse := t.GetCellValue(p)
// 						if uc { //依赖单元格尚未计算
// 							uncalc = uc
// 							break //跳出循环
// 						}
// 						if canuse == false {
// 							continue
// 						}
// 						v, e := strconv.ParseFloat(vstr, 64) //转换为浮点数
// 						if e == nil {
// 							vals = append(vals, v)
// 						}
// 					}
// 					if uncalc { //依赖单元格尚未计算
// 						break //跳出循环
// 					}
// 					//fmt.Println("数组", vals) //==========================
// 					switch functype {
// 					case "SUM":
// 						tokvalue = fmt.Sprint(vals.Sum())
// 					case "AVERAGE":
// 						tokvalue = fmt.Sprint(vals.Mean())
// 					case "MIN":
// 						tokvalue = fmt.Sprint(vals.Min())
// 					case "MAX":
// 						tokvalue = fmt.Sprint(vals.Max())
// 					case "MEDIAN":
// 						tokvalue = fmt.Sprint(vals.Median())
// 					case "PRODUCT":
// 						tokvalue = fmt.Sprint(vals.Product())
// 					default:
// 					}
// 				} else {
// 					if reg_cell.MatchString(tok) { //匹配单元格
// 						tokvalue, uncalc, _ = t.GetCellValue(tok)
// 						if uncalc { //依赖单元格尚未计算
// 							break //跳出循环
// 						}
// 					} else {
// 						tokvalue = tok
// 					}
// 				}
// 				tokens += tokvalue
// 			}
// 			if uncalc { //依赖单元格尚未计算
// 				haveuncalc = true
// 				break //跳出循环
// 			}
// 			//fmt.Println("Tokens", tokens) //=================================================
// 			cellv, err := calc.Calc(tokens)
// 			//fmt.Println("计算结果和错误", cellv, err) //=================================================
// 			if err == nil {
// 				t.Cells[row][col].Value = fmt.Sprint(cellv)
// 			}
// 			t.Cells[row][col].NeedCalcFormula = false
// 		}
// 	}
// 	if haveuncalc {
// 		haveuncalc = false
// 		t.CalcFormula()
// 	}
// }

/********************************************
功能:获取单元格的值
输入:[axis] 单元格坐标
输出:[string] 单元格的值
	[bool] 单元格中有公式,但还没有计算
	[bool] 如果是合并单元格，但不是首格,false；其他情况为true
说明:
编辑:wang_jp
时间:2020年5月12日
*******************************************/
func (t *Table) GetCellValue(axis string) (string, bool, bool) {
	var cellok bool = true
	col, row, err := excelize.CellNameToCoordinates(axis) ///换为数字坐标
	if err != nil {                                       //转换错误
		return axis, false, cellok
	}
	if col > t.ColCnt || row > t.RowCnt { //坐标超范围
		return "", false, false
	}

	if t.Cells[row-1][col-1].IsInMerged && !t.Cells[row-1][col-1].IsMergeStart {
		cellok = false
	}

	return t.Cells[row-1][col-1].Value,
		t.Cells[row-1][col-1].NeedCalcFormula,
		cellok
}

/********************************************
功能:初始化表格的行
输入:[rowheight] 行高
    [row] 行号(从1开始)
    [visible] 可见性
输出:
说明:
编辑:wang_jp
时间:2020年5月12日
*******************************************/
func (tr *TableRow) Init(rowheight float64, row int, visible bool) {
	tr.Height = rowheight
	tr.Index = row
	tr.Visible = visible //行可见性
}

/********************************************
功能:初始化表格的列
输入:[axis] 按Excel列的命名方式:A、B、C等
    [width] 列宽
    [visible] 可见性
输出:
说明:
编辑:wang_jp
时间:2020年5月12日
*******************************************/
func (tc *TableCol) Init(axis string, width float64, visible bool) {
	tc.Axis = axis
	tc.Width = width
	tc.Visible = visible
}

/********************************************
功能:检验单元格是否是合并单元格的一部分
输入:[merges] SHEET中合并单元格的列表
输出:[bool] 如果是合并单元格的一部分,为true,否则为false
	[bool] 如果是合并单元格的左上角第一格,为true,否则为false
说明:同时将合并的行数、列数、单元格的值、起始单元格保存到结构体中
编辑:wang_jp
时间:2020年5月12日
*******************************************/
func (t *TableCell) IsCellInMerged(merges []excelize.MergeCell) (bool, bool) {
	col, row, _ := excelize.CellNameToCoordinates(t.Axis)
	t.IndexC = col
	t.IndexR = row
	var inMerg, isMgStart bool
	var cellValue string
	for _, mg := range merges { //遍历合并的单元格
		msc, msr, _ := excelize.CellNameToCoordinates(mg.GetStartAxis())
		mec, mer, _ := excelize.CellNameToCoordinates(mg.GetEndAxis())
		inMerg = ((col >= msc && col <= mec) && (row >= msr && row <= mer)) //单元格坐标是否在给定单元格内
		if inMerg {
			inMerg = true                             //如果单元格坐标在给定单元格内
			isMgStart = (t.Axis == mg.GetStartAxis()) //是否起始单元格
			cellValue = mg.GetCellValue()             //输出单元格的值
			t.IsInMerged = inMerg
			t.IsMergeStart = isMgStart
			t.ColSpan = mec - msc + 1
			t.RowSpan = mer - msr + 1
			t.MergeStart = mg.GetStartAxis()
			t.MergeEnd = mg.GetEndAxis()
			t.Value = cellValue
			break //跳出循环
		}
	}
	return inMerg, isMgStart
}

/********************************************
功能:寻找表头行
输入:无
输出:HTML
说明:
编辑:wang_jp
时间:2020年5月12日
*******************************************/
func (t *Table) FindHead() {
	headspan := 0
	if t.RowCnt > 0 {
		headspan = 1
	}
	for i, row := range t.Cells {
		if i == 0 { //取第一行
			for _, cell := range row { //遍历第一行每一个格
				if cell.RowSpan > headspan { //合并的行数是否大于表头行数
					headspan = cell.RowSpan //如果大于,替换表头行数
				}
			}
			break
		}
	}
	for i := 0; i < headspan; i++ { //遍历表头行
		for j := range t.Cells[i] {
			t.Cells[i][j].IsHead = true //将表头行的每个单元格设置表头标志
		}
	}
}

/********************************************
功能:单元格格式化为HTML描述
输入:无
输出:HTML
说明:
编辑:wang_jp
时间:2020年5月12日
*******************************************/
func (t *TableCell) FomatToHtml() string {
	cellTag := "td"
	if t.IsHead { //属于表头
		cellTag = "th"
	}
	cellValue := t.Value
	fv, e := strconv.ParseFloat(t.Value, 64)
	if e == nil {
		strs := strings.Split(cellValue, ".") //以小数点分割
		if len(strs) > 1 {                    //能分割成2段以上
			if len(strs[1]) > 3 { //第二段的长度大于3
				cellValue = fmt.Sprintf("%.3f", fv) //精确到小数点后3位
			}
		}
	}
	if t.Visible {
		if !t.IsInMerged { //不是合并单元格
			t.Html = fmt.Sprintf("\t<%s id='%s'>%s</%s>\n",
				cellTag, t.Axis, cellValue, cellTag)
		} else { //是合并的单元格
			if t.IsMergeStart { //合并单元格的起始格
				t.Html = fmt.Sprintf("\t<%s id='%s' rowspan='%d' colspan='%d' >%s</%s>\n",
					cellTag, t.Axis, t.RowSpan, t.ColSpan, cellValue, cellTag)
			}
		}
	}
	return t.Html
}

/********************************************
功能:单元格格式化为HTML描述
输入:无
输出:HTML
说明:
编辑:wang_jp
时间:2020年5月12日
*******************************************/
func (t *TableCell) FomatToHtmlWithBtn() string {
	cellTag := "td"
	if t.IsHead { //属于表头
		cellTag = "th"
	}
	cellValue := t.Value
	fv, e := strconv.ParseFloat(t.Value, 64)
	if e == nil {
		strs := strings.Split(cellValue, ".") //以小数点分割
		if len(strs) > 1 {                    //能分割成2段以上
			if len(strs[1]) > 3 { //第二段的长度大于3
				cellValue = fmt.Sprintf("%.3f", fv) //精确到小数点后3位
			}
		}
	}
	if t.Visible {
		if !t.IsInMerged { //不是合并单元格
			t.Html = fmt.Sprintf("\t<%s id='%s' onclick='onClickCell(\"%s\")'>%s</%s>\n",
				cellTag, t.Axis, t.Axis, cellValue, cellTag)
		} else { //是合并的单元格
			if t.IsMergeStart { //合并单元格的起始格
				t.Html = fmt.Sprintf("\t<%s rowspan='%d' colspan='%d' id='%s' onclick='onClickCell(\"%s\")'>%s</%s>\n",
					cellTag, t.RowSpan, t.ColSpan, t.Axis, t.Axis, cellValue, cellTag)
			}
		}
	}
	return t.Html
}

/********************************************
功能:获取表格中所有单元格的Map
输入:无
输出:map[string]TableCell,以单元格坐标字符串为Key
说明:
编辑:wang_jp
时间:2020年5月24日
*******************************************/
func (t *Table) GetCellMap() map[string]TableCell {
	cellmap := make(map[string]TableCell)
	for _, tbrow := range t.Cells {
		for _, cell := range tbrow {
			cellmap[cell.Axis] = cell
		}
	}
	return cellmap
}

/********************************************
功能:表格格式化为HTML描述
输入:无
输出:HTML
说明:
编辑:wang_jp
时间:2020年5月12日
*******************************************/
func (t *Table) FomatToHtml() string {
	t.Html = `<div class="excel excel-table">`
	t.Html += `<table><thead>`
	//模拟Excel表头
	t.Html += `<tr class="thead" style="height: 25px;">`
	t.Html += `<td class="drug-ele-td drug-ele-td-vertical" style="width: 60px; text-align: center;" id="coordinate"></td>`
	for _, tbcol := range t.Cols {
		if tbcol.Visible { //如果列可见
			w := 100.0 //tbcol.Width * _WidthFactor
			t.Html += fmt.Sprintf("\t<td class='drug-ele-td drug-ele-td-horizontal' style='text-align: center;width:%.0fpx;'>%s</td>\n", w, tbcol.Axis)
		}
	}
	t.Html += `</tr>`

	for r, row := range t.Cells { //遍历每行
		if t.Rows[r].Visible {
			t.Html += fmt.Sprintf("<tr id='%d' style='height:%.0fpx;'>\n", r+1, math.Ceil(t.Rows[r].Height))                            //行开始
			t.Html += fmt.Sprintf(`<td class="drug-ele-td drug-ele-td-vertical" style="width: 60px; text-align: center;">%d</td>`, r+1) //行开始
			for _, cell := range row {                                                                                                  //遍历行中的每个单元格
				t.Html += cell.FomatToHtmlWithBtn() //格式化单元格
			}
			t.Html += fmt.Sprintln(`</tr>`) //行结束
		}
	}
	t.Html += fmt.Sprintln(`</tbody></table></div>`)

	cellmap, err := json.Marshal(t.GetCellMap())
	if err == nil {
		t.Html += `<SCRIPT type="text/javascript">`
		t.Html += fmt.Sprintf("CELLSMAP = %s;", cellmap)
		t.Html += `</SCRIPT>`
	}
	return t.Html
}

/********************************************
功能:设置单元格的值
输入:[filename] Excel文件名(含路径)
	[sheetname] 工作表Sheet名称
输出:[error] 错误信息
说明:
编辑:wang_jp
时间:2020年5月24日
*******************************************/
func (t *TableCell) SetCellValue(filename string, sheet ...string) error {
	sheetname := "Sheet1" //默认打开Sheet1
	if len(sheet) > 0 {
		sheetname = sheet[0]
	}
	f, err := excelize.OpenFile(filename) //打开excel
	if err != nil {
		return fmt.Errorf("打开文件失败:[%s]", err.Error())
	}
	fvalue, err := strconv.ParseFloat(t.Value, 64) //转换成浮点数
	if err == nil {                                //如果能转换
		err = f.SetCellValue(sheetname, t.Axis, fvalue) //设置单元格的值
		if err != nil {
			return fmt.Errorf("设置单元格[%s]的值失败:[%s]", t.Axis, err.Error())
		}
	} else { //如果不能转换
		err = f.SetCellValue(sheetname, t.Axis, t.Value) //设置单元格的值
		if err != nil {
			return fmt.Errorf("设置单元格[%s]的值失败:[%s]", t.Axis, err.Error())
		}
	}
	err = f.Save()
	if err != nil {
		return fmt.Errorf("保存文件失败:[%s]", err.Error())
	}
	return nil
}

/********************************************
功能:设置单元格的公式
输入:[filename] Excel文件名(含路径)
	[sheetname] 工作表Sheet名称
输出:[error] 错误信息
说明:
编辑:wang_jp
时间:2020年5月24日
*******************************************/
func (t *TableCell) SetCellFormula(filename string, sheet ...string) error {
	sheetname := "Sheet1" //默认打开Sheet1
	if len(sheet) > 0 {
		sheetname = sheet[0]
	}
	f, err := excelize.OpenFile(filename) //打开excel
	if err != nil {
		return fmt.Errorf("打开文件失败:[%s]", err.Error())
	}
	err = f.SetCellFormula(sheetname, t.Axis, t.Formula) //设置单元格的公式
	if err != nil {
		return fmt.Errorf("设置单元格[%s]的公式失败:[%s]", t.Axis, err.Error())
	}
	err = f.Save()
	if err != nil {
		return fmt.Errorf("保存文件失败:[%s]", err.Error())
	}
	return nil
}

/*******************************************************************************
- 功能:字符串是否包含子字符串
- 参数:[s string] 待检测字符串
	[substrs ...string] 子字符串切片
- 输出:如果字符串中含有任意一个子字符串，则返回true,否则返回false
- 备注:
- 时间: 2020年6月27日
*******************************************************************************/
func stringContains(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
