package regression

import (
	"testing"
	//"micscript/numgo"
)

func TestTableLookupF(t *testing.T) {
	tests := []struct {
		v1  int
		v2  int
		alp float64
		res float64
	}{
		{1, 1, 0.05, 161.0},
		{2, 2, 0.05, 19.0},
	}

	for i, tt := range tests {
		res := TableLookupF(tt.v1, tt.v2, tt.alp)
		if tt.res != res {
			t.Errorf("No.%d Wrong answer, got=%f, want=%f", i, res, tt.res)
		}
	}
}

func TestTableLookupT(t *testing.T) {
	tests := []struct {
		v   int
		alp float64
		res float64
	}{
		{1, 0.05, 6.31},
		{2, 0.05, 2.92},
	}

	for i, tt := range tests {
		res := TableLookupT(tt.v, tt.alp)
		if tt.res != res {
			t.Errorf("No.%d Wrong answer, got=%f, want=%f", i, res, tt.res)
		}
	}
}

/*
func TestRegression(t *testing.T) {
	x := numgo.Matrix{
		{41, 45, 51, 52, 59, 62, 69, 72, 78, 80, 90, 92, 98, 103},
		{49, 58, 62, 71, 62, 74, 71, 74, 79, 84, 85, 94, 91, 95},
	}
	y := numgo.Array{28, 39, 41, 44, 43, 50, 51, 57, 63, 66, 70, 76, 80, 81}

	x1 := numgo.Matrix{
		{10, 12, 14, 16, 19, 16, 7, 7, 12, 11, 12, 7, 11},
		{26, 26, 40, 32, 51, 33, 26, 25, 17, 24, 16, 16, 15},
		{0.2, -1.4, -0.8, 0.2, -1.4, 0.2, 2.7, 1.0, 2.2, -0.8, -0.5, 2.0, 1.1},
		{3.6, 4.4, 1.7, 1.4, 0.9, 2.1, 2.7, 4, 3.7, 3, 4.9, 4.1, 4.7},
	}
	y1 := numgo.Array{9, 17, 34, 42, 40, 27, 4, 27, 13, 56, 15, 8, 20}
	tests := []struct {
		x numgo.Matrix
		y numgo.Array
	}{
		{x, y},
		{x1, y1},
	}

	for _, tt := range tests {
		res := Regression(tt.x, tt.y)
		t.Logf("%+v", res)
	}
}
*/
