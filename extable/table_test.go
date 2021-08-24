package extable

import (
	"testing"
)

func TestFomatToHtml(t *testing.T) {
	table := new(Table)
	err := table.OpenFile("testfile/MicForm-3.xlsx", true, 15, "Sheet1")
	if err != nil {
		t.Error(err.Error())
	} else {
		//t.Logf("%+v\n", table)
		//for i, cell := range table.Cells {
		//	t.Logf("第%d行------------------\n", i+1)
		//	for _, c := range cell {
		//		t.Logf("%s,%s,%s,%d,%s\n", c.Axis, c.Value, c.Formula, c.CalcFormulaScanCnt, c.Err)
		//	}
		//}
		//t.Log("--------------------------")
		t.Log(table.FomatToHtml())
	}
}

func TestDecomposeCellZone(t *testing.T) {
	tests := []struct {
		zone string
		res  []string
	}{
		{"E8:G10", []string{"E8", "F8", "G8", "E9", "F9", "G9", "E10", "F10", "G10"}},
		{"E10:G8", []string{"E8", "F8", "G8", "E9", "F9", "G9", "E10", "F10", "G10"}},
		{"G8:E10", []string{"E8", "F8", "G8", "E9", "F9", "G9", "E10", "F10", "G10"}},
		{"G10:E8", []string{"E8", "F8", "G8", "E9", "F9", "G9", "E10", "F10", "G10"}},
	}
	for i, tt := range tests {
		t.Logf("======第%d行=====", i)
		f := new(Formula)
		cells := f.DecomposeCellZone(tt.zone)
		t.Log(cells)
		for i, tok := range cells {
			if i < len(tt.res) {
				if tok != tt.res[i] {
					t.Errorf("错误:期望值:%s,得到值:%s", tt.res, cells)
					break
				}
			} else {
				t.Errorf("错误:期望值:%s,得到值:%s", tt.res, cells)
				break
			}
		}
	}
}

func TestDecomposeFuncPars(t *testing.T) {
	tests := []struct {
		foumula string
		res     string
	}{
		{"SUM(A1,b2,az32)", "SUM(A1,b2,az32)"},
		{"SUM(A1:B3)", "SUM(A1,B1,A2,B2,A3,B3)"},
		{"PRODUCT(A1:B3,b5:c8)", "PRODUCT(A1,B1,A2,B2,A3,B3,B5,C5,B6,C6,B7,C7,B8,C8)"},
		{"AVERAGE(E5,D6,C8,A1:B3,b5:c8)", "AVERAGE(E5,D6,C8,A1,B1,A2,B2,A3,B3,B5,C5,B6,C6,B7,C7,B8,C8)"},
		{"MEDIAN(E5,D6,C8,A1:B3,b5:c8,32.5,16.8)", "MEDIAN(E5,D6,C8,A1,B1,A2,B2,A3,B3,B5,C5,B6,C6,B7,C7,B8,C8,32.5,16.8)"},
	}
	for i, tt := range tests {
		t.Logf("======第%d行=====", i)
		f := new(Formula)
		str := f.DecomposeFuncPars(tt.foumula)
		t.Log(str)
		if str != tt.res {
			t.Errorf("错误:期望值:%s,得到值:%s", tt.res, str)
		}
	}
}

func TestSetCell(t *testing.T) {
	cell := new(TableCell)
	cell.Axis = "A11"
	cell.Value = "11"
	err := cell.SetCellValue("Book1.xlsx")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("ok")
	}
}
