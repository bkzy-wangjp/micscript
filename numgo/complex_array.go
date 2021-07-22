// numgo project numgo.go
package numgo

import (
	"fmt"
	"math"
	"math/cmplx"
)

type ComplexArray []complex128 //复数数组

/***********************************************
功能:实数数组转换为复数数组
输入:
	arr:[]float64 实数
输出:
	ComplexArray
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func RealArr2ComplexsArr(arr Array) ComplexArray {
	var cmps ComplexArray
	for _, v := range arr {
		cmps = append(cmps, Real2Complex(v))
	}
	return cmps
}

/***********************************************
功能:复制一个新的数组
输入:无
输出:ComplexArray;新的复数数组
说明:
编辑:wang_jp
时间:2020年10月20日
***********************************************/
func (cmps ComplexArray) Copy() ComplexArray {
	var res ComplexArray
	for _, c := range cmps {
		res = append(res, c)
	}
	return res
}

/***********************************************
功能:实数除以一个复数数组中的每一个元素
输入:
	r:float64:作为被除数的实数
	inplace:bool: 是否改变原始数据,true是,false否
输出:
	ComplexArray;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func (cmps ComplexArray) DivByReal(r float64, inplace ...bool) ComplexArray {
	res := cmps.Copy()
	if len(inplace) > 0 {
		if inplace[0] {
			res = cmps
		}
	}
	for i, c := range cmps {
		res[i] = ComplexDivByReal(r, c)
	}
	return res
}

/***********************************************
功能:实数与一个复数数组中的每一个元素相减
输入:
	r:float64:作为除数的实数
	inplace:bool: 是否改变原始数据,true是,false否
输出:
	ComplexArray;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func (cmps ComplexArray) SubByReal(r float64, inplace ...bool) ComplexArray {
	res := cmps.Copy()
	if len(inplace) > 0 {
		if inplace[0] {
			res = cmps
		}
	}
	for i, c := range cmps {
		res[i] = ComplexSubByReal(r, c)
	}
	return res
}

/***********************************************
功能:实数与一个复数数组中的每一个元素相加
输入:
	r:float64:实数
	cmps:ComplexArray:复数
输出:
	ComplexArray;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func (cmps ComplexArray) AddReal(r float64, inplace ...bool) ComplexArray {
	res := cmps.Copy()
	if len(inplace) > 0 {
		if inplace[0] {
			res = cmps
		}
	}
	for i, c := range cmps {
		res[i] = ComplexAddReal(r, c)
	}
	return res
}

/***********************************************
功能:两个复数数组中的每一个元素相除
输入:
	b:ComplexArray:复数
	inplace:bool: 是否改变原始数据,true是,false否
输出:
	ComplexArray;复数
说明:如果被除数组的实部和虚部都是0，则使结果元素等于被除数
	 如果两个数组的长度不相同，将短的扩展为长的
		扩展被除数时,新扩展的元素为:0+0i
		扩展除数时,新扩展的元素为:1+0i
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func (cmps ComplexArray) Div(b ComplexArray, inplace ...bool) ComplexArray {
	res := cmps.Copy()
	nb := b.Copy()
	if len(inplace) > 0 {
		if inplace[0] {
			res = cmps
		}
	}
	alen := len(cmps)
	blen := len(b)
	if alen != blen { //长度不同
		if alen < blen {
			res = append(res, make(ComplexArray, blen-alen)...) //扩展被除数
		} else {
			n := make(ComplexArray, alen-blen)
			for i, _ := range n {
				n[i] = 1 + 0i //新的除数为1+0i
			}
			nb = append(nb, n...) //扩展除数
		}
	}
	for i, c := range nb { //执行除法运算
		if real(c) == 0 && imag(c) == 0 { //实部和虚部都是零
			res[i] = c
		} else {
			res[i] = res[i] / c
		}
	}
	return res
}

/***********************************************
功能:复数数组元素同时乘以一个实数
输入:
	r:float64;实数
	inplace:bool: 是否改变原始数据,true是,false否
输出:
	ComplexArray;复数数组
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func (cmps ComplexArray) MulReal(r float64, inplace ...bool) ComplexArray {
	res := cmps.Copy()
	if len(inplace) > 0 {
		if inplace[0] {
			res = cmps
		}
	}
	for i, cv := range cmps {
		res[i] = ComplexMulReal(r, cv)
	}
	return res
}

/***********************************************
功能:复数数组的乘积
输入:
	cmps:ComplexArray:复数数组切片
输出:
	complex128;复数
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func (cmps ComplexArray) Prod() complex128 {
	var res complex128 = 1 + 0i
	if len(cmps) < 1 {
		return res
	}
	for _, cv := range cmps {
		res *= cv
	}
	return res
}

/***********************************************
功能:比较两个复数数组是否相等
输入:
	a,b:ComplexArray:复数数组切片
输出:
	如果两个数组长度相等且每个元素的实部和虚部都相等,返回true,否则返回false
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func (a ComplexArray) IsEqual(b ComplexArray) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if !ComplexsIsEqual(v, b[i]) {
			return false
		}
	}
	return true
}

/***********************************************
功能:数组中的元素是否全为零
输入:
输出:true or false
说明:数组中的元素全为0输出true，否则为false
	 空数组当成全为0的数组
编辑:wang_jp
时间:2020年10月20日
***********************************************/
func (a ComplexArray) IsAllZero() bool {
	if len(a) == 0 {
		return true
	}
	for _, v := range a {
		if v != 0 {
			return false
		}
	}
	return true
}

/***********************************************
功能:将数组作为多项式的系数列出多项式的行列式
输入:
输出:Matrix,N x N+1的矩阵
说明:a[0] * x**(N) + a[1] * x**(N-1) + ... + a[N-1] * x + a[N]=0
	上式中,c[0]始终为1，则可以转换为：
	a[1] * x**(N-1) + ... + a[N-1] * x + a[N]= - a[0] * x**(N)
	输出的c矩阵为:
		x[0]**(N-1)    x[0]**(N-2) ...  x[0]**(1)   x[0]**(0)   -x[0]**(N)
		x[1]**(N-1)    x[1]**(N-2) ...  x[1]**(1)   x[1]**(0)   -x[1]**(N)
		...       ...           ...     ...     ...
		x[N-1]**(N-1) x[N-1]**(N-2) ... x[N-1]**(1) x[N-1]**(0) -x[N-1]**(N)
	N 是输入数组的长度
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func (X ComplexArray) Polynominal() ComplexMatrix {
	n := len(X)
	var cmat ComplexMatrix //行列式
	for _, x := range X {
		var mr ComplexArray //行列式的行
		for i := n - 1; i >= 0; i-- {
			mr = append(mr, cmplx.Pow(x, Real2Complex(float64(i))))
		}
		mr = append(mr, -cmplx.Pow(x, Real2Complex(float64(n)))) //
		cmat = append(cmat, mr)
	}
	return cmat
}

/***********************************************
功能:用给定的根序列求多项式的系数。
输入:
输出:c:ComplexArray,长度为N+1的矩阵,c[0]始终为1.0
说明:c[0] * x**(N) + c[1] * x**(N-1) + ... + c[N-1] * x + c[N]
	对于给定的零序列，返回其前导系数为1的多项式的系数（序列中必须包含多个根是
	其重数的许多倍）
编辑:wang_jp
时间:2020年10月20日
***********************************************/
func (seq_of_zeros ComplexArray) Poly() (ComplexArray, error) {
	var c ComplexArray
	c = append(c, 1.0)
	if seq_of_zeros.IsAllZero() { //空数组或者全为0的数组
		return append(c, seq_of_zeros...), nil
	}
	mat := seq_of_zeros.Polynominal()       //多项式的增广矩阵
	_, err := mat.LeftSquare2Identity(true) //左侧方阵转换为单位矩阵
	if err != nil {
		return nil, err
	}

	for _, mr := range mat {
		c = append(c, mr[len(mr)-1]) //提取矩阵最后一列的数据作为系数
	}
	return c, nil
}

/***********************************************
功能:将数组中的所有元素取共轭
输入:inplace bool:是否更改原始数组
输出:所有元素共轭后的数组
说明:复数 a+ib 的共轭复数是 a-ib
编辑:wang_jp
时间:2020年10月21日
***********************************************/
func (cmps ComplexArray) Conjugate(inplace ...bool) ComplexArray {
	res := cmps.Copy()
	if len(inplace) > 0 {
		if inplace[0] {
			res = cmps
		}
	}
	for i, v := range res {
		res[i] = ComplexConjugate(v)
	}
	return res
}

/***********************************************
功能:判断是否全为实数，即每个元素的虚部全为0
输入:
输出:是true，否false
编辑:wang_jp
时间:2020年10月21日
***********************************************/
func (cmps ComplexArray) IsAllReal() bool {
	for _, v := range cmps {
		if imag(v) != 0.0 {
			return false
		}
	}
	return true
}

/***********************************************
功能:获取复数数组各个元素实部的副本
输入:
输出:实数数组
编辑:wang_jp
时间:2020年10月21日
***********************************************/
func (cmps ComplexArray) RealCopy() Array {
	var res Array
	for _, v := range cmps {
		res = append(res, real(v))
	}
	return res
}

/***********************************************
功能:获取复数数组各个元素的模(绝对值)
输入:
输出:实数数组
编辑:wang_jp
时间:2020年10月21日
***********************************************/
func (cmps ComplexArray) Abs() Array {
	var res Array
	for _, v := range cmps {
		res = append(res, cmplx.Abs(v))
	}
	return res
}

/***********************************************
功能:获取复数数组各个元素的角度(主值)
输入:
输出:实数数组
说明:复数Z=a+bi化成三角式r(cosθ+jsinθ),可简写作r∠θ,
	其中模(绝对值)r=√(a²+b²)；复角θ由tanθ=b/a解出并在0≤θ＜360°范围内取值（主值）
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func (cmps ComplexArray) Angle() Array {
	var res Array
	for _, v := range cmps {
		res = append(res, ComplexAngle(v))
	}
	return res
}

/***********************************************
功能:交换两个元素
输入:k1,k2 int:需要交换的两个元素下标
输出:
说明:
编辑:wang_jp
时间:2020年11月09日
***********************************************/
func (a ComplexArray) Swap(k1, k2 int) {
	al := len(a)
	if k1 >= al || k2 >= al {
		return
	}
	temp := a[k1]
	a[k1] = a[k2]
	a[k2] = temp
}

/***********************************************
功能:傅里叶变换
输入:args ...bool:可配置参数,
		args[0]:inplace,true:改变源数组,false:不改变源数组.不填写时不改变源数组
		args[1]:false:傅里叶变换,true:傅里叶逆变换.不填写时为傅里叶变换
输出:
说明:
编辑:wang_jp
时间:2020年11月09日
***********************************************/
func (cmps ComplexArray) Fft(args ...bool) ComplexArray {
	res := cmps.Copy()
	inv := 1.0 //傅里叶变换
	if len(args) > 0 {
		if args[0] {
			res = cmps
		}
		if len(args) > 1 {
			if args[1] { //傅里叶逆变换
				inv = -1.0
			}
		}
	}

	bit := 0
	n := len(cmps) //数组长度

	for { //循环计算bit取值
		if 1<<bit < n {
			bit += 1
		} else {
			break
		}
	}
	p := int(math.Pow(2, float64(bit)))
	if n < p {
		resadd := make(ComplexArray, p-n)
		res = append(res, resadd...)
		n = p
	}
	rev := make([]int, n)
	for i, _ := range rev { //初始化rev
		rev[i] = i
	}
	for i := 0; i < n-1; i++ {
		rev[i] = (rev[i>>1] >> 1) | ((i & 1) << (bit - 1))
		if i < rev[i] { //不加这条if会交换两次（就是没交换）
			res.Swap(i, rev[i])
		}
	}
	for mid := 1; mid < n; mid *= 2 { //mid是准备合并序列的长度的二分之一
		theta := PI / float64(mid)
		unitroot := complex(math.Cos(theta), inv*math.Sin(theta)) //单位根，pi的系数2已经约掉了

		for i := 0; i < n; i += mid * 2 { //mid*2是准备合并序列的长度，i是合并到了哪一位
			omega := 1. + 0i
			fmt.Println("-------------", i, "--------------")
			for j := 0; j < mid; j++ {
				fmt.Println(mid, (i + j), (i + j + mid)) //===========
				x := res[i+j]
				y := omega * res[i+j+mid]
				//蝴蝶变换
				res[i+j] = x + y
				res[i+j+mid] = x - y

				omega *= unitroot
			}
		}
	}

	return res[:len(cmps)]
}
