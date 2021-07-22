package numgo

import (
	"testing"
)

func TestCmplxPoly(t *testing.T) {
	tests := []struct {
		arr ComplexArray
		res ComplexArray
	}{
		{ComplexArray{}, ComplexArray{1}},
		{ComplexArray{0, 0, 0}, ComplexArray{1, 0, 0, 0}},
		{ComplexArray{1, 2, 3}, ComplexArray{1, -6, 11, -6}},
		{ComplexArray{0.024, -2.483, -3.985}, ComplexArray{1., 6.444, 9.739523, -0.23747412}},
		{ComplexArray{-1. / 2, 0, 1. / 2}, ComplexArray{1., 0., -0.25, 0.}},
		{ComplexArray{(6 + 2i), (-7 + 7i), (0 + 6i)}, ComplexArray{1. + 0.i, 1. - 15.i, -110. + 22.i, 168. + 336.i}},
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

func TestUnitRoot(t *testing.T) {
	tests := []struct {
		n   int
		k   int
		res complex128
	}{
		{4, 0, 1 + 0i},
		{4, 1, 0 + 1i},
		{4, 2, -1 - 0i},
		{4, 3, 0 - 1i},
	}

	for _, tt := range tests {
		res := UnitRoot(tt.n, tt.k)
		if ComplexsIsEqual(res, tt.res) == false {
			t.Error("错误:期望值:", tt.res, "得到值:", res)
		}
	}
}
