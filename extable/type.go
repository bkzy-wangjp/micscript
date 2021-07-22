// extable project extable.go
package extable

const (
//_WidthFactor = 16.8 //宽度系数
)

//Excel表格的行
type TableRow struct {
	Index   int     //行号
	Height  float64 //行高
	Visible bool    //可见性
}

//Excel表格的列
type TableCol struct {
	Axis    string  //列名
	Width   float64 //列宽
	Visible bool    //可见性
}

//Excel表格的单元格
type TableCell struct {
	Axis               string  //单元格坐标
	Value              string  //值
	Formula            string  //公式
	Comments           string  //批注
	StyleIndex         int     //样式索引
	Visible            bool    //所在列的可见性
	IsHead             bool    //是否表头
	IsInMerged         bool    //是否合并单元格的一部分
	IsMergeStart       bool    //是否合并单元格的左上角(仅IsMerged为true时有用)
	NeedCalcFormula    bool    //需要对公式进行计算
	CalcFormulaScanCnt int     //计算单元格公式的时候扫描的次数
	RowSpan            int     //合并的行数(仅IsMergeStart为true时有用)
	ColSpan            int     //合并的列数(仅IsMergeStart为true时有用)
	MergeStart         string  //是合并单元格一部分时,合并格的起始坐标(仅IsMerged为true时有用)
	MergeEnd           string  //合并单元格的结束坐标
	IndexC             int     //列号(从1开始)
	IndexR             int     //行号(从1开始)
	Width              float64 //列宽
	Html               string  //描述表格的HTML语句
	Status             int     //前端状态:默认=0,已编辑=1,已编辑但公式未校验=2
	Err                string  //计算错误信息
}

//Excel工作表
type Table struct {
	RowCnt int           //行数
	ColCnt int           //列数
	Rows   []TableRow    //行
	Cols   []TableCol    //列
	Cells  [][]TableCell //单元格
	Html   string        //描述表格的HTML语句
}
