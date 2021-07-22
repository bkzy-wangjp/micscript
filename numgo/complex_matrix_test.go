// numgo project numgo.go
package numgo

import (
	"testing"
)

func TestCmplxInitComplexMatrixAsIdentity(t *testing.T) {
	tests := []struct {
		n   int
		res ComplexMatrix
	}{
		{0, ComplexMatrix{}},
		{1, ComplexMatrix{{1}}},
		{2, ComplexMatrix{{1, 0}, {0, 1}}},
		{3, ComplexMatrix{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}},
		{4, ComplexMatrix{{1, 0, 0, 0}, {0, 1, 0, 0}, {0, 0, 1, 0}, {0, 0, 0, 1}}},
	}

	for _, tt := range tests {
		res := InitComplexMatrixAsIdentity(tt.n)
		if res.IsEqual(tt.res) == false {
			t.Error("错误:期望值:", tt.res, "得到值:", res)
		}
	}
}

func TestCmplxInitComplexMatrix(t *testing.T) {
	tests := []struct {
		m   int
		n   int
		res ComplexMatrix
	}{
		{0, 0, ComplexMatrix{}},
		{1, 1, ComplexMatrix{{0}}},
		{2, 3, ComplexMatrix{{0, 0, 0}, {0, 0, 0}}},
		{3, 4, ComplexMatrix{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}},
	}

	for _, tt := range tests {
		res := InitComplexMatrix(tt.m, tt.n)
		if res.IsEqual(tt.res) == false {
			t.Error("错误:期望值:", tt.res, "得到值:", res)
		}
	}
}

func TestCmplxMatrixAdd(t *testing.T) {
	tests := []struct {
		a   ComplexMatrix
		b   ComplexMatrix
		res ComplexMatrix
	}{
		{ComplexMatrix{}, ComplexMatrix{}, ComplexMatrix{}},
		{ComplexMatrix{{0}}, ComplexMatrix{{1}}, ComplexMatrix{{1}}},
		{ComplexMatrix{{0i, 0, 0}, {1, 0, 0}},
			ComplexMatrix{{0, 2i, 0}, {0, 0, 3}},
			ComplexMatrix{{0, 2i, 0}, {1, 0, 3}}},
	}

	for _, tt := range tests {
		res, err := tt.a.Add(tt.b)
		if err != nil {
			t.Log(err)
		}
		if res.IsEqual(tt.res) == false {
			t.Error("错误:期望值:", tt.res, "得到值:", res)
		}
	}
}

func TestCmplxElemTransRowMulK(t *testing.T) {
	tests := []struct {
		r   int
		k   complex128
		a   ComplexMatrix
		res ComplexMatrix
	}{
		{0, 1 + 1i, ComplexMatrix{}, ComplexMatrix{}},
		{2, 2.5, ComplexMatrix{{0, 2i, 0}, {0, 0, 3}, {1, 1 + 1i, 1 + 2i}},
			ComplexMatrix{{0, 2i, 0}, {0, 0, 3}, {2.5, 2.5 + 2.5i, 2.5 + 5i}}},
	}

	for _, tt := range tests {
		tt.a.ElemTransRowMulK(tt.r, tt.k)
		if tt.a.IsEqual(tt.res) == false {
			t.Error("错误:期望值:", tt.res, "得到值:", tt.a)
		}
	}
}

func TestCmplxElemTransRowMulKAddToRow(t *testing.T) {
	tests := []struct {
		r1  int
		r2  int
		k   complex128
		a   ComplexMatrix
		res ComplexMatrix
	}{
		{0, 2, 2.5, ComplexMatrix{{0, 2i, 0}, {0, 0, 3}, {1, 1 + 1i, 1 + 2i}},
			ComplexMatrix{{2.5, 2.5 + 4.5i, 2.5 + 5i}, {0, 0, 3}, {1, 1 + 1i, 1 + 2i}}},
	}

	for _, tt := range tests {
		tt.a.ElemTransRowMulKAddToRow(tt.r1, tt.r2, tt.k)
		if tt.a.IsEqual(tt.res) == false {
			t.Error("错误:期望值:", tt.res, "得到值:", tt.a)
		}
	}
}

func TestCmplxElemTransRowSwap(t *testing.T) {
	tests := []struct {
		r1  int
		r2  int
		a   ComplexMatrix
		res ComplexMatrix
	}{
		{0, 2, ComplexMatrix{{0, 2i, 0}, {0, 0, 3}, {1, 1 + 1i, 1 + 2i}},
			ComplexMatrix{{1, 1 + 1i, 1 + 2i}, {0, 0, 3}, {0, 2i, 0}}},
	}

	for _, tt := range tests {
		tt.a.ElemTransRowSwap(tt.r1, tt.r2)
		if tt.a.IsEqual(tt.res) == false {
			t.Error("错误:期望值:", tt.res, "得到值:", tt.a)
		}
	}
}

func TestCmplxTranspose(t *testing.T) {
	tests := []struct {
		a   ComplexMatrix
		res ComplexMatrix
	}{
		{ComplexMatrix{{2 + 2i, 3 + 7i}, {3 - 6i, 6 + 2i}, {-7 + 7i, 6i}},
			ComplexMatrix{{2 - 2i, 3 + 6i, -7 - 7i}, {3 - 7i, 6 - 2i, -6i}}},
	}

	for _, tt := range tests {
		res := tt.a.Transpose()
		if res.IsEqual(tt.res) == false {
			t.Error("错误:\n期望值:", tt.res, "\n得到值:", res)
		}
	}
}

func TestCmplxLeftSquare2Triu(t *testing.T) {
	tests := []struct {
		a   ComplexMatrix
		res ComplexMatrix
	}{
		{ComplexMatrix{{3, 2, 1}, {2, 5, 6}, {7, 8, 9}},
			ComplexMatrix{{3, 2, 1}, {0, 3.66666667, 5.33333333}, {0, 0, 1.8181818181818183}}},
		{ComplexMatrix{{0, 2i, 5 + 0.78i}, {0, 0, 3}, {1, 1 + 1i, 1 + 2i}},
			ComplexMatrix{{1, 1 + 1i, 1 + 2i}, {0, 2i, 5 + 0.78i}, {0, 0, 3}}},
	}

	for _, tt := range tests {
		res, err := tt.a.LeftSquare2Triu()
		if err != nil {
			t.Log(err)
		}
		if res.IsEqual(tt.res) == false {
			t.Error("错误:\n期望值:", tt.res, "\n得到值:", res)
		}
	}
}

func TestCmplxLeftSquare2Identity(t *testing.T) {
	tests := []struct {
		a   ComplexMatrix
		res ComplexMatrix
	}{
		{ComplexMatrix{{3, 2, 1}, {2, 5, 6}, {7, 8, 9}},
			ComplexMatrix{{1, 0, 0}, {0, 1, 0}, {0, 0, 1.}}},
		{ComplexMatrix{{0, 2i, 5 + 0.78i}, {0, 0, 3}, {1, 1 + 1i, 1 + 2i}},
			ComplexMatrix{{1, 0, 0}, {0, 1, 0}, {0, 0, 1.}}},
		{ComplexMatrix{{2 + 2i, 3 + 7i, 3 - 6i, 1, 0, 0}, {6 + 2i, -7 + 7i, 6i, 0, 1, 0}, {6i, 5 - 4i, -2 - 1i, 0, 0, 1}},
			ComplexMatrix{{1, 0, 0, 0.04432952 - 0.02793007i, 0.02474261 - 0.02884646i, 0.01513244 - 0.10822204i},
				{0, 1, 0, -0.02949989 - 0.02599369i, 0.00519555 - 0.06322315i, 0.07998916 + 0.0251012i},
				{0, 0, 1, 0.01724413 + 0.11838205i, -0.0592229 - 0.0646097i, 0.03914194 - 0.07139897i}}},
	}

	for i, tt := range tests {
		res, err := tt.a.LeftSquare2Identity()
		if err != nil {
			t.Log(err)
		}
		if res.IsEqual(tt.res) == false {
			t.Error(i, "错误:\n期望值:", tt.res, "\n得到值:", res)
		}
	}
}

func TestCmplxInverse(t *testing.T) {
	tests := []struct {
		a   ComplexMatrix
		res ComplexMatrix
	}{
		{ComplexMatrix{{2 + 2i, 3 + 7i, 3 - 6i}, {6 + 2i, -7 + 7i, 6i}, {6i, 5 - 4i, -2 - 1i}},
			ComplexMatrix{{0.04432952 - 0.02793007i, 0.02474261 - 0.02884646i, 0.01513244 - 0.10822204i},
				{-0.02949989 - 0.02599369i, 0.00519555 - 0.06322315i, 0.07998916 + 0.0251012i},
				{0.01724413 + 0.11838205i, -0.0592229 - 0.0646097i, 0.03914194 - 0.07139897i}}},
		{ComplexMatrix{{3, 2, 1}, {2, 5, 6}, {7, 8, 9}},
			ComplexMatrix{{-0.15, -0.5, 0.35},
				{1.2, 1, -0.8},
				{-0.95, -0.5, 0.55}}},
	}

	for i, tt := range tests {
		res, err := tt.a.Inverse()
		if err != nil {
			t.Log(err)
		}
		if res.IsEqual(tt.res) == false {
			t.Error(i, "错误:\n期望值:", tt.res, "\n得到值:", res)
		}
	}
}
