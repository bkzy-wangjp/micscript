package mics

import (
	"testing"
)

func TestSortStringSlice(t *testing.T) {
	tests := []struct {
		orderby string
		ipt     []string
		opt     []string
	}{
		{"desc", []string{"France", "Paris", "Italy", "Rome", "Japan", "Tokyo", "India", "New delhi"},
			[]string{"Tokyo", "Rome", "Paris", "New delhi", "Japan", "Italy", "India", "France"}},
		{"asc", []string{"France", "Paris", "Italy", "Rome", "Japan", "Tokyo", "India", "New delhi"},
			[]string{"France", "India", "Italy", "Japan", "New delhi", "Paris", "Rome", "Tokyo"}},
	}
	for _, tt := range tests {
		SortStringSlice(tt.ipt, tt.orderby)
		if !StringsCompiler(tt.ipt, tt.opt) {
			t.Errorf("\n 得到值:%v, \n 期望值:%v", tt.ipt, tt.opt)
		}
	}
}

func TestSortStringMap(t *testing.T) {
	tests := []struct {
		orderby string
		ipt     map[string]string
		opt     []string
	}{
		{"desc", map[string]string{"France": "Paris", "Italy": "Rome", "Japan": "Tokyo", "India": "New delhi"},
			[]string{"Japan", "Italy", "India", "France"}},
		{"asc", map[string]string{"France": "Paris", "Italy": "Rome", "Japan": "Tokyo", "India": "New delhi"},
			[]string{"France", "India", "Italy", "Japan"}},
	}
	for _, tt := range tests {
		opt := SortStringMap(tt.ipt, tt.orderby)
		if !StringsCompiler(opt, tt.opt) {
			t.Errorf("\n 得到值:%v, \n 期望值:%v", opt, tt.opt)
		}
	}
}

func TestMd5(t *testing.T) {
	str := `{"file_name":"testfile/2020年1月计划.xlsx","password":"","sheets":["矿石自然类型明细总表"],"cells":[{"axismaps":{"g3":"col_1","g4":"col_2","g6":"col_3"},"db_dest":{"table_name":"excel","consts":{"excel":"mul_cells_1"}}},{"axismaps":{"g8":"col_4","g9":"col_5","g11":"col_6"},"db_dest":{"table_name":"excel","consts":{"excel":"mul_cells_2"}}}],"row":{"first_row":3,"last_row":0,"load_rows":[13,14,15,16,17,18,19],"colmaps":{"E":"col_1","g":"col_2","h":"col_3","i":"col_4","j":"col_5","k":"col_6"},"ignore_rows":[5,7,12,20,21,25,28,30,35,39,41],"db_dest":{"table_name":"excel","consts":{"excel":"mul_row_load"}}},"column":{"first_col":"E","last_col":"","load_cols":["E","G","H","I","J","K"],"rowmaps":{"13":"col_1","14":"col_2","15":"col_3","16":"col_4","17":"col_5","18":"col_6"},"ignore_cols":["F"],"db_dest":{"table_name":"excel","consts":{"excel":"mul_column_load"}}}}true1628230101802admindf5cbcafc6a12759c6ac17e9f93e83516F7B`
	md5str := Md5str(str)
	t.Log(md5str)
}

func TestTimeParse(t *testing.T) {
	tests := []struct {
		timestr string
	}{
		{"2021-10-28T15:02:07.5622783+08:00"},
		{"2021-10-28T15:02:07.5622+08:00"},
		{"2021-10-28T15:02:07.56+08:00"},
		{"2021-10-28 15:02:07.56"},
		{"2021-10-28 15:02:07+0800"},
		{"2021-10-29T18:32:10.26367+08:00"},
	}
	for _, tt := range tests {
		tm, err := TimeParse(tt.timestr)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(tm)
		}
	}
}
