package exc2sql

import "github.com/360EntSecGroup-Skylar/excelize/v2"

type ExcelFile struct {
	FileName string `json:"file_name"` //文件路径(含名称)
	Password string `json:"password"`  //文件密码(如无密码，请置空)
	//NamePrefix  string         `json:"name_prefix"`  //文件名前缀,可以为空
	Sheets []string `json:"sheets"` //Excel的工作表(Sheet)名称,如果为空则遍历每个Sheet
	//SheetPrefix string         `json:"sheet_prefix"` //工作表名前缀,可以为空
	SyncType string         `json:"sync_type"` //同步类型,{"cell":按单元格,"row":按行(固定列),"column":按列(固定行)}
	Cells    []*Cells       `json:"cells"`     //单元格参数(仅 SyncType=cell 时有效)
	Row      *Row           `json:"row"`       //定义行参数(仅 SyncType=row 时有效)
	Column   *Column        `json:"column"`    //定义列参数(仅 SyncType=column 时有效)
	exfile   *excelize.File //打开的excel文件指针
}

//定义单元格
type Cells struct {
	Axismaps map[string]string `json:"axismaps"` //Excel单元格坐标:数据库表列名
	DbDest   *DbDestination    `json:"db_dest"`  //目的数据库单元格
	axis     []string
}

//定义行
type Row struct {
	FirstRow   int               `json:"first_row"`   //第一行有效数据的行号
	Colmaps    map[string]string `json:"colmaps"`     //Excel列名:数据库表列名
	IgnoreRows []int             `json:"ignore_rows"` //忽略的行
	DbDest     *DbDestination    `json:"db_dest"`     //目的数据库单元格
	colnames   []string          //通过colmaps解析出来的excel列名
	colindexs  []int             //列索引(根据Colmaps计算所得,从0开始)
}

//定义列
type Column struct {
	FirstCol        string            `json:"first_col"`   //第一列有效数据的列名
	Rowmaps         map[string]string `json:"rowmaps"`     //Excel行号:数据库表列名
	IgnoreCols      []string          `json:"ignore_cols"` //忽略的列
	DbDest          *DbDestination    `json:"db_dest"`     //目的数据库单元格
	rows            []int             //从Rowmaps中解析出来的行列表
	firstcolindex   int               //第一列的索引(从1开始)
	ignorecolindexs []int             //忽略的列的索引
}

type DbDestination struct {
	TableName   string            `json:"table_name"`   //数据库表名称
	Consts      map[string]string `json:"consts"`       //常数项
	TimeColumns *TimeColumns      `json:"time_columns"` //时间列的获取方法,可以为空
	//条件,相互之间为Or的关系.
	//可以为空,为空时INSERT数据；不为空时UPDATE数据
	Wheres   []*Wheres `json:"wheres"`
	colnames []string
}

type DbValues struct {
	TableName string
	ColNames  []string
	Values    [][]string
}

//查询条件,相互之间为And的关系
type Wheres struct {
	ColNames []string `json:"col_names"` //列名
	Operator []string `json:"operator"`  //操作符{=,<>,>,<,>=,<=,LIKE,IN}
	//数值
	//  ""
	Values []*GetValueMethod `json:"values"` //值
}

//时间列的取值方式
type TimeColumns struct {
	ColNames []string `json:"col_names"` //列名
	//时间值获取方法:
	//  "now":当前时间
	//  "cell:CellAxis":从文件单元格获取,冒号':'后为单元格的坐标,
	//  "filesuffix:TimeFormat":文件名后缀,冒号':'后为时间字符串的格式
	//  "sheetsuffix:TimeFormat":工作表后缀,冒号':'后为时间字符串的格式
	GetFrom []*GetValueMethod `json:"get_from"`
	//目标数据库时间格式
	//  "localunix":Unix格式的秒(本地时区)
	//  "localunixms":Unix格式的毫秒(本地时区)
	//  "localunixmicro":Unix格式的微秒(本地时区)
	//  "localunixnano":Unix格式的纳秒(本地时区)
	//  "unix":Unix格式的秒(UTC时区)
	//  "unixms":Unix格式的毫秒(UTC时区)
	//  "unixmicro":Unix格式的微秒(UTC时区)
	//  "unixnano":Unix格式的纳秒(UTC时区)
	//  其他字符串:自定义时间格式字符串
	TimeFormat []string `json:"time_format"`
}

type GetValueMethod struct {
	//获取值的方法
	//  "now":获取当前时间,Format中填写时间格式
	//  "order":顺序号. 对于cell,column,Format为数组下标号加1;对于row,Format为行号
	//  "translate":对于当前cell的平移, Format为"X,Y",X为横坐标数字,Y为纵坐标数字
	Method string `json:"method"` //获取值的方法
	Format string `json:"format"` //获取值的格式
}

type DataBaseSet struct {
	DbType   string `json:"db_type"`   //{mysql,mssql,sqlite,influxdb}
	Host     string `json:"host"`      //访问地址
	Port     int    `json:"port"`      //端口号
	DbName   string `json:"db_name"`   //数据库名
	UserName string `json:"user_name"` //用户名
	Password string `json:"password"`  //密码
}
