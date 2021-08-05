package exc2sql

import "testing"

var config_cell = `
{
	"file_name":"testfile/2020年1月计划.xlsx"
	,"password":""
	,"sheets":["矿石自然类型明细总表"]
	,"cells":[
		{
			"axismaps":{
				"g3":"col_1"
				,"g4":"col_2"
				,"g6":"col_3"
			}
			,"db_dest":{
				"table_name":"excel"
				,"consts":{
					"excel":"cells_1"
				}
			}
		}
		,{
			"axismaps":{
				"g8":"col_4"
				,"g9":"col_5"
				,"g11":"col_6"
			}
			,"db_dest":{
				"table_name":"excel"
				,"consts":{
					"excel":"cells_2"
				}
			}
		}
	]
	,"row":{}
	,"column":{}
}
`

var config_row = `
{
	"file_name":"testfile/2020年1月计划.xlsx"
	,"password":""
	,"sheets":["矿石自然类型明细总表"]
	,"cells":[]
	,"row":{
		"first_row":3
		,"last_row":0
		,"load_rows":[]
		,"colmaps":{
			"E":"col_1"
			,"g":"col_2"
			,"h":"col_3"
			,"i":"col_4"
			,"j":"col_5"
			,"k":"col_6"
		}
		,"ignore_rows":[5,7,12,20,21,25,28,30,35,39,41]
		,"db_dest":{
			"table_name":"excel"
			,"consts":{
				"excel":"row"
			}
		}
	}
	,"column":{}
}
`
var config_row_load = `
{
	"file_name":"testfile/2020年1月计划.xlsx"
	,"password":""
	,"sheets":["矿石自然类型明细总表"]
	,"cells":[]
	,"row":{
		"first_row":3
		,"last_row":0
		,"load_rows":[13,14,15,16,17,18,19]
		,"colmaps":{
			"E":"col_1"
			,"g":"col_2"
			,"h":"col_3"
			,"i":"col_4"
			,"j":"col_5"
			,"k":"col_6"
		}
		,"ignore_rows":[5,7,12,20,21,25,28,30,35,39,41]
		,"db_dest":{
			"table_name":"excel"
			,"consts":{
				"excel":"row_load"
			}
		}
	}
	,"column":{}
}
`

var config_col = `
{
	"file_name":"testfile/2020年1月计划.xlsx"
	,"password":""
	,"sheets":["矿石自然类型明细总表"]
	,"cells":[]
	,"row":{}
	,"column":{
		"first_col":"E"
		,"last_col":""
		,"load_cols":[]
		,"rowmaps":{
			"13":"col_1"
			,"14":"col_2"
			,"15":"col_3"
			,"16":"col_4"
			,"17":"col_5"
			,"18":"col_6"
		}
		,"ignore_cols":["F"]
		,"db_dest":{
			"table_name":"excel"
			,"consts":{
				"excel":"column"
			}
		}
	}
}
`
var config_col_load = `
{
	"file_name":"testfile/2020年1月计划.xlsx"
	,"password":""
	,"sheets":["矿石自然类型明细总表"]
	,"cells":[]
	,"row":{}
	,"column":{
		"first_col":"E"
		,"last_col":""
		,"load_cols":["E","G","H","M"]
		,"rowmaps":{
			"13":"col_1"
			,"14":"col_2"
			,"15":"col_3"
			,"16":"col_4"
			,"17":"col_5"
			,"18":"col_6"
		}
		,"ignore_cols":["F"]
		,"db_dest":{
			"table_name":"excel"
			,"consts":{
				"excel":"column_load"
			}
		}
	}
}
`

var config_mul = `
{
	"file_name":"testfile/2020年1月计划.xlsx"
	,"password":""
	,"sheets":["矿石自然类型明细总表"]
	,"cells":[
		{
			"axismaps":{
				"g3":"col_1"
				,"g4":"col_2"
				,"g6":"col_3"
			}
			,"db_dest":{
				"table_name":"excel"
				,"consts":{
					"excel":"mul_cells_1"
				}
			}
		}
		,{
			"axismaps":{
				"g8":"col_4"
				,"g9":"col_5"
				,"g11":"col_6"
			}
			,"db_dest":{
				"table_name":"excel"
				,"consts":{
					"excel":"mul_cells_2"
				}
			}
		}
	]
	,"row":{
		"first_row":3
		,"last_row":0
		,"load_rows":[13,14,15,16,17,18,19]
		,"colmaps":{
			"E":"col_1"
			,"g":"col_2"
			,"h":"col_3"
			,"i":"col_4"
			,"j":"col_5"
			,"k":"col_6"
		}
		,"ignore_rows":[5,7,12,20,21,25,28,30,35,39,41]
		,"db_dest":{
			"table_name":"excel"
			,"consts":{
				"excel":"mul_row_load"
			}
		}
	}
	,"column":{
		"first_col":"E"
		,"last_col":""
		,"load_cols":["E","G","H","I","J","K"]
		,"rowmaps":{
			"13":"col_1"
			,"14":"col_2"
			,"15":"col_3"
			,"16":"col_4"
			,"17":"col_5"
			,"18":"col_6"
		}
		,"ignore_cols":["F"]
		,"db_dest":{
			"table_name":"excel"
			,"consts":{
				"excel":"mul_column_load"
			}
		}
	}
}
`

func Test_OpenFile(t *testing.T) {
	tests := []struct {
		cfg string
	}{
		{config_cell},
		{config_row},
		{config_row_load},
		{config_col},
		{config_col_load},
		{config_mul},
	}
	for _, tt := range tests {
		ef, err := NewExcelFile(tt.cfg)
		if err != nil {
			t.Error(err)
		} else {
			err := ef.OpenFile()
			if err != nil {
				t.Error(err)
			} else {
				dbv, err := ef.GetValues()
				if err != nil {
					t.Error(err)
				} else {
					for _, v := range dbv {
						v.Print()
						t.Log(v.FormatInsertSql())
					}
				}
			}
		}
	}
}
