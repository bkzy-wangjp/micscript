package numgo

import (
	"math"
)

const (
	_RealDiff = 1e-6
	PI        = 3.141592653589793
)

/***********************************************
功能:实数除以一个复数
输入:
	r:float64:实数
	cmp:complex128:复数
输出:
	complex128;复数
说明:如果被除数为0,输出结果为0
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexDivByReal(r float64, cmp complex128) complex128 {
	// c_r := real(cmp)
	// c_i := imag(cmp)
	// return complex(c_r*r, -1*c_i*r) / (cmp * complex(c_r, -1*c_i))
	r_c := Real2Complex(r)
	if cmp != 0 {
		return r_c / cmp
	}
	return cmp
}

/***********************************************
功能:复数加上一个实数
输入:
	r:float64:实数
	cmp:complex128:复数
输出:
	complex128;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexAddReal(r float64, cmp complex128) complex128 {
	return complex(real(cmp)+r, imag(cmp))
}

/***********************************************
功能:复数减去一个实数
输入:
	r:float64:实数
	cmp:complex128:复数
输出:
	complex128;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexSubReal(r float64, cmp complex128) complex128 {
	return complex(real(cmp)-r, imag(cmp))
}

/***********************************************
功能:实数减去一个复数
输入:
	r:float64:实数
	cmp:complex128:复数
输出:
	complex128;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexSubByReal(r float64, cmp complex128) complex128 {
	//return complex(r-real(cmp), -imag(cmp))
	return Real2Complex(r) - cmp
}

/***********************************************
功能:实数乘以一个复数
输入:
	r:float64:实数
	cmp:complex128:复数
输出:
	complex128;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexMulReal(r float64, cmp complex128) complex128 {
	return complex(real(cmp)*r, imag(cmp)*r)
}

/***********************************************
功能:比较两个复数是否相等
输入:
	a,b:complex128:复数
输出:
	如果两个复数的实部和虚部都相等,返回true,否则返回false
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexsIsEqual(a, b complex128) bool {
	if math.Abs(real(a)-real(b)) > _RealDiff {
		return false
	}
	if math.Abs(imag(a)-imag(b)) > _RealDiff {
		return false
	}
	return true
}

/***********************************************
功能:实数转换为复数
输入:
	a:float64 实数
输出:
	complex128
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func Real2Complex(a float64) complex128 {
	return ComplexMulReal(a, 1+0i)
}

/***********************************************
功能:返回复数的共轭复数
输入:
	cmp:complex128 复数
输出:
	complex128
说明:复数 a+ib 的共轭复数是 a-ib
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexConjugate(cmp complex128) complex128 {
	return complex(real(cmp), -imag(cmp))
}

/***********************************************
功能:返回复数的角度(主值)
输入:
	cmp:complex128 复数
输出:
	float64,角度的主值
说明:复数Z=a+bi化成三角式r(cosθ+isinθ),可简写作r∠θ,
	其中模(绝对值)r=√(a²+b²)；复角θ由tanθ=b/a解出并在0≤θ＜360°范围内取值（主值）
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func ComplexAngle(cmp complex128) float64 {
	if real(cmp) == 0 {
		return PI / 2
	}
	return math.Atan(imag(cmp) / real(cmp))
}

/***********************************************
功能:单位根
输入:n,k int
	 inv int:可选,用于对单位根取共轭。-1或者1,不填写时默认为1
输出:
	complex128,单位根
说明:单位根指模为1的根，一般的x的n个单位根可以表示为:x=cos(2*pi*k/n)+sin(2*pi*k/n)i，
	其中：k=0,1,2,..,n-1,i是虚数单位
编辑:wang_jp
时间:2020年11月9日
***********************************************/
func UnitRoot(n, k int, inv ...int) complex128 {
	if n == 0 {
		return 0
	}
	fk := float64(k)
	fn := float64(n)
	var invs float64 = 1.0
	if len(inv) > 0 {
		if inv[0] < 0 {
			invs = -1
		}
	}
	theta := 2 * fk * PI / fn
	return complex(math.Cos(theta), invs*math.Sin(theta))
}
