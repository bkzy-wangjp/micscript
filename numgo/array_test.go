package numgo

import (
	"testing"
)

func TestMul(t *testing.T) {
	a := Array{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	tests := []struct {
		arr Array
		n   int
		fv  string
		res Array
	}{
		{a, 1, "mean", a},
		{a, 4, "median", Array{2.5, 3.5, 4.5, 5.5, 6.5, 7.5, 8.5, 9.5, 10, 10.5, 11}},
	}

	for _, tt := range tests {
		res := tt.arr.MoveWindowFilter(tt.n, tt.fv)
		if res.IsEqual(tt.res) == false {
			t.Error("错误:期望值:", tt.res, "得到值:", res)
		}
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		arr Array
		res Array
	}{
		{Array{1, 2, 3, 4, 5, 6, 7, 8}, Array{8, 7, 6, 5, 4, 3, 2, 1}},
	}

	for _, tt := range tests {
		tt.arr.Reverse()
		if tt.arr.IsEqual(tt.res) == false {
			t.Error("错误:期望值:", tt.res, "得到值:", tt.arr)
		}
	}
}

func TestCompanion(t *testing.T) {
	tests := []struct {
		arr Array
		res Matrix
	}{
		{Array{1, 2}, Matrix{Array{-2}}},
		{Array{1, -10, 31, -30}, Matrix{Array{10., -31., 30.}, Array{1, 0, 0}, Array{0, 1, 0}}},
	}

	for _, tt := range tests {
		mt, err := tt.arr.Companion()
		if err != nil {
			t.Log(err)
		} else if mt.IsEqual(tt.res) == false {
			t.Error("错误:期望值:", tt.res, "得到值:", tt.arr)
		}
	}
}

func TestPolynominal(t *testing.T) {
	tests := []struct {
		arr Array
		res Matrix
	}{
		{Array{0, 0, 0}, Matrix{Array{0, 0, 1, 0}, Array{0, 0, 1, -0}, Array{0, 0, 1, -0}}},
		{Array{1, 2, 3}, Matrix{Array{1, 1, 1, -1}, Array{4, 2, 1, -8}, Array{9, 3, 1, -27}}},
	}

	for _, tt := range tests {
		mt := tt.arr.Polynominal()
		if mt.IsEqual(tt.res) == false {
			t.Error("错误:期望值:", tt.res, "得到值:", mt)
		}
	}
}

func TestPoly(t *testing.T) {
	tests := []struct {
		arr Array
		res Array
	}{
		{Array{}, Array{1}},
		{Array{0, 0, 0}, Array{1, 0, 0, 0}},
		{Array{1, 2, 3}, Array{1, -6, 11, -6}},
		{Array{0.024, -2.483, -3.985}, Array{1., 6.444, 9.739523, -0.23747412}},
		{Array{-1. / 2, 0, 1. / 2}, Array{1., 0., -0.25, 0.}},
	}

	for _, tt := range tests {
		res, err := tt.arr.Poly()
		if err != nil {
			t.Log(err)
		} else if res.IsEqual(tt.res) == false {
			t.Error("错误:期望值:", tt.res, "得到值:", res)
		}
	}
}

func TestGroup(t *testing.T) {
	X := Array{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	X1 := Array{0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08, 0.09, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16, 0.17, 0.18, 0.19, 0.20}
	tests := []struct {
		data Array
		min  float64
		max  float64
		grp  int
		gp   float64
		res  Array
		//gmp map[string]float64
	}{
		{X, 0, 0, 10, 1.9, Array{2, 2, 2, 2, 2, 2, 2, 2, 2, 2}},
		{X, 0, 0, 10, 1.9, Array{2, 2, 2, 2, 2, 2, 2, 2, 2, 2}},
		{X, 5, 20, 10, 1.5, Array{6, 1, 2, 1, 2, 1, 2, 1, 2, 1, 1}},
		{X, 5, 15, 10, 1, Array{5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 6}},
		{X, 1, 20, 10, 1.9, Array{2, 2, 2, 2, 2, 2, 2, 2, 2, 2}},
		{X, 0, 25, 10, 2.5, Array{2, 2, 3, 2, 3, 2, 3, 2, 1, 0}},
		{X1, 0, 0, 10, 0.019, Array{2, 2, 2, 2, 2, 2, 2, 2, 1, 3}},
	}

	for i, tt := range tests {
		gp, res, _ := tt.data.Group(tt.min, tt.max, tt.grp)
		if gp != tt.gp {
			t.Error(i, "组距错误:期望值:", tt.gp, "得到值:", gp)
		}
		if res.IsEqual(tt.res) == false {
			t.Error(i, "数组错误:\n期望值:", tt.res, "\n得到值:", res)
		}
	}
}

func TestMode(t *testing.T) {
	X := Array{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	tests := []struct {
		data Array
		grp  int
	}{
		{X, 10},
	}

	for i, tt := range tests {
		md, gp, mp := tt.data.Mode(tt.grp)
		t.Log(i, "Mode:", md, ":组距:", gp)
		t.Log(mp)
	}
}

func TestArraySwap(t *testing.T) {
	X := Array{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	tests := []struct {
		k1 int
		k2 int
	}{
		{1, 10},
	}
	for i, tt := range tests {
		x := X.Copy()
		x.Swap(tt.k1, tt.k2)
		if x[tt.k1] != X[tt.k2] || x[tt.k2] != X[tt.k1] {
			t.Error(i, "错误")
		}
	}
}

func TestArrayFft(t *testing.T) {
	tests := []struct {
		x Array
	}{
		{Array{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18,
			19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}},
		{Array{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}},
		{Array{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
		{Array{16, 17, 18, 19}},
	}
	for i, tt := range tests {
		res := tt.x.Fft()
		t.Log("==============", i, "===============")
		for i, r := range res {
			t.Log(i, r)
		}
	}
}
