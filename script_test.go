package mics

import "testing"

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
