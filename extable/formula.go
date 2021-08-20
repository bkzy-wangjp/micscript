// extable project extable.go
package extable

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/xuri/excelize/v2"
)

//Excel公式
type Formula struct {
	TokenStr []string
}

const (
	_cell       = `[a-zA-Z]+[0-9]+`                 //单元格坐标
	_cellZone   = _cell + `:` + _cell               //坐标区域
	_mathSymbol = `(\+|\-|\*|\/|\%|\(\s*|\s*\)){1}` //数学符号：+ - * / % ( )
	_number     = `\d+\.*\d*`                       //数字
	_FUNC       = `(SUM|PRODUCT|AVERAGE|MEDIAN|MIN|MAX){1}\(((` + _cellZone + `\,?)*|(` + _cell + `\,?)*|(` + _number + `\,?)*)*\)`
)

/*******************************************************************************
* 功能:合成正则字符串
* 输入:无
* 输出:正则字符串
* 编辑:wang_jp
* 时间:2020年5月29日
*******************************************************************************/
func (f *Formula) getRegexpString() string {
	//合成正则表达式
	return fmt.Sprintf("%s|%s|%s|%s",
		_FUNC, _cell, _number, _mathSymbol,
	)
}

/*******************************************************************************
* 功能:分解单元格区域
* 输入:单元格区域字符串
* 输出:单元格字符串数组
* 编辑:wang_jp
* 时间:2020年5月29日
*******************************************************************************/
func (f *Formula) DecomposeCellZone(cellzone string) []string {
	var cells []string
	zone := strings.Split(cellzone, ":") //拆分区域

	if len(zone) < 2 { //长度小于2
		return zone
	}
	if zone[0] == zone[1] { //区域的起始和结束相等
		cells = append(cells, zone[0])
		return cells //只返回一个单元格
	}
	//转换问数字坐标
	bg_col, bg_row, err := excelize.CellNameToCoordinates(zone[0])
	if err != nil {
		return cells
	}
	ed_col, ed_row, err := excelize.CellNameToCoordinates(zone[1])
	if err != nil {
		return cells
	}

	//校验，防止前大后小
	if bg_col > ed_col {
		i := bg_col
		bg_col = ed_col
		ed_col = i
	}
	if bg_row > ed_row {
		i := bg_row
		bg_row = ed_row
		ed_row = i
	}
	for r := bg_row; r <= ed_row; r++ {
		for c := bg_col; c <= ed_col; c++ {
			cell, e := excelize.CoordinatesToCellName(c, r)
			if e == nil {
				cells = append(cells, cell)
			}
		}
	}
	return cells
}

/*******************************************************************************
* 功能:分解函数参数,将含有单元格区域的地方都转换为单元格列表
* 输入:函数及其参数字符串
* 输出:将单元格区域转换之后的函数及其参数字符串
* 编辑:wang_jp
* 时间:2020年5月29日
*******************************************************************************/
func (f *Formula) DecomposeFuncPars(formula string) string {
	regstr := fmt.Sprintf("%s|%s|%s", _cellZone, _cell, _number) //Excel函数参数
	reg := regexp.MustCompile(regstr)
	strs := reg.FindAllString(formula, -1) //提取函数参数
	indexs := reg.FindAllStringIndex(formula, -1)

	reg = regexp.MustCompile(_cellZone) //单元格区域
	var cells string                    //单元格或者值
	num := len(strs)
	for i, str := range strs { //遍历提取到的函数参数
		if reg.MatchString(str) { //匹配单元格区域
			cs := f.DecomposeCellZone(str)
			n := len(cs)
			for j, c := range cs {
				cells += c
				if j < n-1 {
					cells += ","
				}
			}
		} else {
			cells += str
		}
		if i < num-1 {
			cells += ","
		}
	}
	formula = formula[:indexs[0][0]] + cells + formula[indexs[len(indexs)-1][1]:]
	return formula
}

/*******************************************************************************
* 功能:提取公式中的计算元组
* 输入:Excel公式(不含=)
* 输出:[bool] 如果有未提取到的字符，输出false
* 编辑:wang_jp
* 时间:2020年5月29日
*******************************************************************************/
func (f *Formula) GetTokens(formula string) bool {
	upformula := strings.ToUpper(formula)
	reg := regexp.MustCompile(f.getRegexpString())
	if !reg.MatchString(upformula) { //没有匹配项
		return false
	}
	indexs := reg.FindAllStringIndex(upformula, -1)

	if len(indexs) == 1 { //只有一个元组的时候
		if indexs[0][1] != len(formula) { //坐标与长度不匹配
			return false
		}
	}
	for i, idex := range indexs {
		if i > 0 {
			if idex[0]-indexs[i-1][1] > 0 {
				return false
			}
		} else {
			if idex[0] != 0 { //第一个坐标不为0
				return false
			}
		}
		if i == len(indexs)-1 {
			if idex[1] != len(formula) { //最后一个坐标与长度不匹配
				return false
			}
		}
	}
	f.TokenStr = reg.FindAllString(upformula, -1) //提取元组
	reg = regexp.MustCompile(_FUNC)
	for i, tok := range f.TokenStr {
		if reg.MatchString(tok) {
			f.TokenStr[i] = f.DecomposeFuncPars(tok)
		}
	}
	return true
}
