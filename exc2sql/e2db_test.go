package exc2sql

import "testing"

var config_cell = `
{
	"file_name":"testfile/2020年1月计划.xlsx"
	,"password":""
	,"sheets":["矿石自然类型明细总表"]
	,"sync_type":"cell"
	,"cells":[
		{
			"axismaps":{
				"g3":"col_f3"
				,"g4":"col_f4"
				,"g6":"col_f6"
			}
			,"db_dest":{
				"table_name":"dbtable_name"
				,"consts":{
					"常数项1":"const_1"
					,"常数项2":"const_2"
				}
				,"time_columns":{}
				,"wheres":[]
			}
		}
		,{
			"axismaps":{
				"g8":"col_g3"
				,"g9":"col_g9"
				,"g11":"col_G11"
			}
			,"db_dest":{
				"table_name":"dbtable_name"
				,"consts":{
					"常数项1":"const_1"
					,"常数项2":"const_2"
				}
				,"time_columns":{}
				,"wheres":[]
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
	,"sync_type":"row"
	,"cells":[]
	,"row":{
		"first_row":3
		,"colmaps":{
			"E":"ID"
			,"g":"矿石量(t)"
			,"h":"铜品位(%)"
			,"i":"铜金属量(t)"
			,"j":"硫品位(%)"
			,"k":"硫金属量(t)"
		}
		,"ignore_rows":[5,7,12,20,21,25,28,30,35,39,41]
		,"db_dest":{
			"table_name":"dbtable_name"
			,"consts":{
				"常数项1":"const_1"
				,"常数项2":"const_2"
			}
			,"time_columns":{}
			,"wheres":[]
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
	,"sync_type":"column"
	,"cells":[]
	,"row":{}
	,"column":{
		"first_col":"E"
		,"rowmaps":{
			"3":"r3"
			,"4":"r4"
			,"6":"r6"
			,"8":"r8"
			,"9":"r9"
			,"10":"r10"
		}
		,"ignore_cols":["F"]
		,"db_dest":{
			"table_name":"dbtable_name"
			,"consts":{
				"const_1":"12345"
				,"const_2":"54321"
			}
			,"time_columns":{}
			,"wheres":[]
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
		{config_col},
	}
	for _, tt := range tests {
		ef, err := newExcelFile(tt.cfg)
		if err != nil {
			t.Error(err)
		} else {
			err := ef.openFile()
			if err != nil {
				t.Error(err)
			} else {
				dbv, err := ef.getValues()
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
