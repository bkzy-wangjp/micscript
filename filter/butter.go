// filter project filter.go
package filter

import (
	"fmt"
	"math"
	"math/cmplx"
	"strings"

	"micscript/numgo"
)

const (
	pi = 3.141592653589793
)

/***********************************************
功能:返回n阶Butterworth(巴特沃斯) 滤波器初始参数
输入:
	n: 滤波阶数,n不可小于0,如果输入为负数自动转换为正数使用
输出:
	z:[]float64;零点参数
	p:[]complex128;极点参数
	k:float64;增益值
说明:
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func buttap(n int) (numgo.Array, numgo.ComplexArray, float64) {
	if n < 0 {
		n *= -1
	}

	var z, m numgo.Array
	for i := (-1*n + 1); i < n; i = i + 2 {
		m = append(m, float64(i))
	}
	var p numgo.ComplexArray
	for i := 0; i < n; i++ {
		tmp := pi * m[i] / (2 * float64(n))
		p = append(p, cmplx.Exp(complex(0, tmp))*-1)
	}
	k := 1.0
	return z, p, k
}

/***********************************************
功能:从零点和极点返回传递函数的相对阶数
输入:
	z:[]float64;零点参数
	p:[]complex128;极点参数
输出:
	degree:int;相对阶数
说明:Return relative degree of transfer function from zeros and poles
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func relativeDegree(z interface{}, p numgo.ComplexArray) (int, error) {
	zlen := 0
	zav, ok := z.(numgo.Array)
	if ok {
		zlen = len(zav)
	} else {
		zcv, ok := z.(numgo.ComplexArray)
		if ok {
			zlen = len(zcv)
		}
	}

	degree := len(p) - zlen
	if degree < 0 {
		return degree, fmt.Errorf("improper transfer function. Must have at least as many poles as zeros[传递函数不正确,极点的阶数不能小于零点的阶数]")
	}
	return degree, nil
}

/***********************************************
功能:将低通滤波器原型转换为不同的频率
输入:
	z:[]float64;零点参数,运行完毕后直接改变
	p:[]complex128;极点参数,运行完毕后直接改变
	k:float64;增益值,运行完毕后直接改变
	wo:float64 截止频率
输出:错误信息
说明:从具有单位截止频率的模拟低通滤波器原型返回截止频率为“wo”的模拟低通滤波器
	 使用零、极点和增益（'zpk'）表示
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func lp2lpZpk(z numgo.Array, p numgo.ComplexArray, k float64, wo float64) (numgo.Array, numgo.ComplexArray, float64, error) {
	degree, err := relativeDegree(z, p)
	if err != nil {
		return nil, nil, 0, err
	}
	//Scale all points radially from origin to shift cutoff frequency
	z_lp := z.Copy()
	z_lp.MulScalar(wo)
	p_lp := p.Copy()
	for i, pv := range p {
		p_lp[i] = numgo.ComplexMulReal(wo, pv)
	}
	//Each shifted pole decreases gain by wo, each shifted zero increases it.
	//Cancel out the net change to keep overall gain the same.
	k_lp := k * math.Pow(wo, float64(degree))
	return z_lp, p_lp, k_lp, nil
}

/***********************************************
功能:将低通滤波器原型转换为不同的频率
输入:
	z:[]float64;零点参数,运行完毕后直接改变
	p:[]complex128;极点参数,运行完毕后直接改变
	k:float64;增益值,运行完毕后直接改变
	wo:float64 截止频率
输出:错误信息
说明:从具有单位截止频率的模拟高通滤波器原型返回截止频率为“wo”的模拟低通滤波器
	 使用零、极点和增益（'zpk'）表示
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func lp2hpZpk(z numgo.Array, p numgo.ComplexArray, k float64, wo float64) (numgo.Array, numgo.ComplexArray, float64, error) {
	degree, err := relativeDegree(z, p)
	if err != nil {
		return nil, nil, 0, err
	}
	fz := z.Copy()
	fz.MulScalar(-1.0)
	fp := p.MulReal(-1.0)
	// Invert positions radially about unit circle to convert LPF to HPF
	// Scale all points radially from origin to shift cutoff frequency
	z_hp := z.Copy()
	z_hp.DivByScalar(wo)
	p_hp := p.Copy()
	for i, pv := range p {
		p_hp[i] = numgo.ComplexDivByReal(wo, pv)
	}
	// If lowpass had zeros at infinity, inverting moves them to origin.
	z_hp = append(z_hp, make(numgo.Array, degree)...)

	// Cancel out gain change caused by inversion
	k_hp := k * real(numgo.ComplexDivByReal(fz.Product(), fp.Prod()))

	return z_hp, p_hp, k_hp, nil
}

/***********************************************
功能:基带移位
输入:
	cmps:[]complex128;
	wo:float64 中心频率,默认为1
输出:[]complex128;
说明:np.concatenate((z_lp + np.sqrt(z_lp**2 - wo**2),
                  z_lp - np.sqrt(z_lp**2 - wo**2)))
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func shiftband(cmps numgo.ComplexArray, wo float64) numgo.ComplexArray {
	wo = wo * wo
	var a, b numgo.ComplexArray
	for _, v := range cmps {
		tmp := v * v
		wc := numgo.Real2Complex(wo)
		tmp -= wc
		sqt := cmplx.Sqrt(tmp)
		a = append(a, v+sqt)
		b = append(b, v-sqt)
	}
	return append(a, b...)
}

/***********************************************
功能:将低通滤波器原型转换为不同的频率
输入:
	z:[]complex128;零点参数,运行完毕后直接改变
	p:[]complex128;极点参数,运行完毕后直接改变
	k:float64;增益值,运行完毕后直接改变
	wo:float64 中心频率,默认为1
	bw:float64 带宽,默认为1
输出:错误信息
说明:从具有单位截止频率的模拟低通滤波器原型返回中心频率为“wo”且带宽为“bw”的模拟带通滤波器器
	 使用零、极点和增益（'zpk'）表示
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func lp2bpZpk(z, p numgo.ComplexArray, k, wo, bw float64) (numgo.ComplexArray, numgo.ComplexArray, float64, error) {
	degree, err := relativeDegree(z, p)
	if err != nil {
		return nil, nil, 0, err
	}
	//缩放零点和极点参数到所需的带宽
	z_lp := z.MulReal(bw / 2)
	p_lp := p.MulReal(bw / 2)
	// 复制极点和零，从基带移到+wo和-wo
	z_bp := shiftband(z_lp, wo)
	p_bp := shiftband(p_lp, wo)
	// Move degree zeros to origin, leaving degree zeros at infinity for BPF.
	z_bp = append(z_bp, make(numgo.ComplexArray, degree)...)
	// Cancel out gain change caused by inversion
	k_bp := k * math.Pow(bw, float64(degree))
	return z_bp, p_bp, k_bp, nil
}

/***********************************************
功能:将低通滤波器原型转换为不同的频率
输入:
	z:[]complex128;零点参数,运行完毕后直接改变
	p:[]complex128;极点参数,运行完毕后直接改变
	k:float64;增益值,运行完毕后直接改变
	wo:float64 中心频率,默认为1
	bw:float64 带宽,默认为1
输出:错误信息
说明:从具有单位截止频率的模拟低通滤波器原型返回中心频率为“wo”且阻带宽度为“bw”的模拟带阻滤波器
	 使用零、极点和增益（'zpk'）表示
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func lp2bsZpk(z, p numgo.ComplexArray, k, wo, bw float64) (numgo.ComplexArray, numgo.ComplexArray, float64, error) {
	degree, err := relativeDegree(z, p)
	if err != nil {
		return nil, nil, 0, err
	}

	fz := z.MulReal(-1.0)
	fp := p.MulReal(-1.0)
	//缩放零点和极点参数到所需的带宽

	z_hp := z.DivByReal(bw / 2)
	p_hp := p.DivByReal(bw / 2)
	// 复制极点和零，从基带移到+wo和-wo
	z_bs := shiftband(z_hp, wo)
	p_bs := shiftband(p_hp, wo)
	// Move degree zeros to origin, leaving degree zeros at infinity for BPF.
	tmp := make(numgo.ComplexArray, degree)
	for i := range tmp {
		tmp[i] = numgo.ComplexMulReal(wo, +1i)
	}
	z_bs = append(z_bs, tmp...)
	for i := range tmp {
		tmp[i] = numgo.ComplexMulReal(wo, -1i)
	}
	z_bs = append(z_bs, tmp...)

	// Cancel out gain change caused by inversion
	prd := fp.Prod()
	if prd == 0 {
		prd = 1
	}
	k_bs := k * real(fz.Prod()/prd)

	return z_bs, p_bs, k_bs, nil
}

/***********************************************
功能:使用双线性变换从模拟滤波器返回数字IIR滤波器
输入:
	z:[]complex128/[]float64;零点参数,运行完毕后直接改变
	p:[]complex128;极点参数,运行完毕后直接改变
	k:float64;增益值,运行完毕后直接改变
	fs:float64
输出:z,p,k
说明:
编辑:wang_jp
时间:2020年10月18日
***********************************************/
func bilinearZpk(z interface{}, p numgo.ComplexArray, k, fs float64) (interface{}, numgo.ComplexArray, float64, error) {
	degree, err := relativeDegree(z, p)
	if err != nil {
		return nil, nil, 0, err
	}
	fs2 := 2.0 * fs

	p_z := p.AddReal(fs2)     //fs2+p
	p_tmp := p.SubByReal(fs2) //fs2-p
	p_z.Div(p_tmp, true)      //(fs2+p)/(fs2-p)

	z_a, ok := z.(numgo.Array)
	if ok {
		z_tmp := z_a.Copy()
		z_tmp.SubByScalar(fs2)        //fs2-z
		z_a.AddScalar(fs2)            //fs2+z
		z_z, _ := z_a.DivArray(z_tmp) //(fs2+z)/(fs2-z)
		k_z_f := z_tmp.Product()
		k_z := k * real(numgo.ComplexDivByReal(k_z_f, p_tmp.Prod())) //k*real(prod(fs2-z)/prod(fs2-p))

		for i := 0; i < degree; i++ {
			z_z = append(z_z, -1)
		}
		return z_z, p_z, k_z, nil
	} else {
		z_c, _ := z.(numgo.ComplexArray)
		z_tmp := z_c.SubByReal(fs2) //fs2-z
		z_c = z_c.AddReal(fs2)      //fs2+z
		z_z := z_c.Div(z_tmp)       //(fs2+z)/(fs2-z)
		k_z_c := z_tmp.Prod()       //prod(fs2-z)
		prd := p_tmp.Prod()
		if prd == 0 {
			prd = 1
		}
		k_z := k * real(k_z_c/prd) //k*real(prod(fs2-z)/prod(fs2-p))

		for i := 0; i < degree; i++ {
			z_z = append(z_z, -1+0i)
		}
		return z_z, p_z, k_z, nil
	}
}

/***********************************************
功能:零极点的回归多项式传递函数表示
输入:
	z:[]complex128/[]float64;零点参数,运行完毕后直接改变
	p:[]complex128;极点参数,运行完毕后直接改变
	k:float64;增益值,运行完毕后直接改变
输出:b,a numgo.Array:b参数和a参数数组
	 err error:错误信息
说明:
编辑:wang_jp
时间:2020年10月21日
***********************************************/
func zpk2tf(z interface{}, p numgo.ComplexArray, k float64) (b, a numgo.Array, err error) {
	z_poly, er := poly(z)
	if er != nil {
		err = er
		return
	}
	b_c := z_poly.MulReal(k)
	a_c, er := poly(p)
	if er != nil {
		err = er
		return
	}
	b = b_c.RealCopy()
	a = a_c.RealCopy()
	return
}

/***********************************************
功能:用给定的根序列求多项式的系数。
输入:seq_of_zeros ComplexArray/Array:零点根序列,长度为N
输出:c:ComplexArray,长度为N+1的矩阵,c[0]始终为1.0
说明:c[0] * x**(N) + c[1] * x**(N-1) + ... + c[N-1] * x + c[N]
	对于给定的零序列，返回其前导系数为1的多项式的系数（序列中必须包含多个根是
	其重数的许多倍）
编辑:wang_jp
时间:2020年10月20日
***********************************************/
func poly(seq_of_zeros interface{}) (numgo.ComplexArray, error) {
	seq_flt, ok := seq_of_zeros.(numgo.Array)
	if ok { //是实数数组
		arr, err := seq_flt.Poly()
		if err != nil {
			return nil, err
		}
		arr_c := numgo.RealArr2ComplexsArr(arr)
		return arr_c, nil
	} else {
		seq_cmplx, ok := seq_of_zeros.(numgo.ComplexArray)
		if ok {
			return seq_cmplx.Poly()
		}
	}
	return nil, fmt.Errorf("data type error,must be []float64 or []complex128.[数据类型错误,必须是64位实数数组或者128位复数数组]")
}

/***********************************************
功能:N阶IIR数字滤波器设计,返回滤波器系数
输入:
	N int: 滤波阶数
	Wn float/[2]float:自然频率数组，也称归一化的截止频率,与fs具有相同的单位
		fcf=截止频率*2/采样频率；
		如果是低通，高通滤波，fcf只有一个元素,范围:0.0-1.0
		如果是带通，带阻滤波，fcf数组有两个元素,分别是滤波带的起始频率和阶数频率,范围:0.0-fs/2.0
	btype string:滤波类型，字符串；lp：表示低通；hp：表示高通；bp：表示带通；bs：表示带阻；
	fs float64:数字系统的采样频率,可选。
		默认情况下,fs是2.0，表示每个周期采样2次
输出:b系数和a系数
说明:
编辑:wang_jp
时间:
***********************************************/
func iirFilter(N int, Wn interface{}, btype string, fs ...float64) (b, a numgo.Array, err error) {
	z, p, k := buttap(N)

	var z_b numgo.ComplexArray //带宽滤波时候的z参数
	switch strings.ToLower(btype) {
	case "lowpass", "low", "lp", "highpass", "high", "hp":
		warped := 1.0
		wn, ok := Wn.(float64)
		if !ok {
			err = fmt.Errorf(" Wn should be of float64 type when 'lowpass' or 'highpass'.[当进行'低通'或者'高通'滤波时,参数 Wn 应该是 float64类型]")
			return
		}
		if len(fs) > 0 {
			if wn <= 0. || wn >= fs[0]/2 {
				err = fmt.Errorf("digital filter critical frequencies must be 0 < Wn < fs/2(fs=%.3f,fs/2=%.3f).[数字滤波器的临界频率'Wn'必须大0小于fs/2]", fs[0], fs[0]/2.0)
				return
			}
		} else if wn <= 0. || wn >= 1. {
			err = fmt.Errorf("digital filter critical frequencies must be 0.0 < Wn < 1.0.[数字滤波器的临界频率'Wn'必须大0.0小于1.0]")
			return
		}

		_fs := 2.0
		if len(fs) > 0 {
			if fs[0] > 0 {
				wn = 2 * wn / fs[0]
			}
		}
		warped = 2 * _fs * math.Tan(pi*wn/_fs)
		switch strings.ToLower(btype) {
		case "lowpass", "low", "lp":
			z, p, k, err = lp2lpZpk(z, p, k, warped)
			if err != nil {
				return
			}
		case "highpass", "high", "hp":
			z, p, k, err = lp2hpZpk(z, p, k, warped)
			if err != nil {
				return
			}
		}
		zo, po, ko, er := bilinearZpk(z, p, k, _fs)
		if er != nil {
			err = er
			return
		}

		b, a, err = zpk2tf(zo, po, ko)
	case "bandpass", "bp", "pass", "bandstop", "bs", "stop":
		wn, ok := Wn.(numgo.Array)
		if !ok {
			err = fmt.Errorf(" Wn should be of [2]float64 type array when 'bandpass' or 'bandstop'.[当‘带通’或者‘带阻’滤波时,参数 Wn 应该是 [2]float64类型的数组]")
			return
		}
		if len(wn) < 2 {
			err = fmt.Errorf(" Wn must specify start and stop frequencies for 'bandpass' or 'bandstop' filter.[设计'带通'或者'带阻'滤波器时,Wn 应该是开始频率和结束频率的数组]")
			return
		} else {

			if len(fs) > 0 {
				if wn[0] <= 0. || wn[0] >= fs[0]/2 || wn[1] <= 0. || wn[1] >= fs[0]/2 {
					err = fmt.Errorf("digital filter critical frequencies must be 0 < Wn < fs/2(fs=%.3f,fs/2=%.3f).[数字滤波器的临界频率'Wn'必须大0小于fs/2]", fs[0], fs[0]/2.0)
					return
				}
			} else if wn[0] <= 0. || wn[0] >= 1. || wn[1] <= 0. || wn[1] >= 1. {
				err = fmt.Errorf("digital filter critical frequencies must be 0.0 < Wn < 1.0.[数字滤波器的临界频率'Wn'必须大0.0小于1.0]")
				return
			}
		}
		_fs := 2.0
		if len(fs) > 0 {
			if fs[0] > 0 {
				//_fs = fs[0]
				for i, v := range wn {
					wn[i] = 2 * v / fs[0]
				}
			}
		}

		var warped numgo.Array
		for _, v := range wn {
			warped = append(warped, 2*_fs*math.Tan(pi*v/_fs))
		}
		bw := warped[1] - warped[0]
		wo := math.Sqrt(warped[0] * warped[1])

		switch strings.ToLower(btype) {
		case "bandpass", "bp", "pass":
			z_b, p, k, err = lp2bpZpk(numgo.RealArr2ComplexsArr(z), p, k, wo, bw)
			if err != nil {
				return
			}
		case "bandstop", "bs", "stop":
			z_b, p, k, err = lp2bsZpk(numgo.RealArr2ComplexsArr(z), p, k, wo, bw)
			if err != nil {
				return
			}
		}
		zo, po, ko, er := bilinearZpk(z_b, p, k, _fs)
		if er != nil {
			err = er
			return
		}
		b, a, err = zpk2tf(zo, po, ko)
	}

	return
}

/***********************************************
功能:N阶Butterworth(巴特沃斯)数字滤波器设计,返回滤波器系数
输入:
	N int: 滤波阶数
	Wn float/[2]float:自然频率数组，也称归一化的截止频率,与fs具有相同的单位
		fcf=截止频率*2/采样频率；
		如果是低通，高通滤波，fcf只有一个元素,范围:0.0-1.0
		如果是带通，带阻滤波，fcf数组有两个元素,分别是滤波带的起始频率和阶数频率,范围:0.0-fs/2.0
	btype string:滤波类型，字符串；lp：表示低通；hp：表示高通；bp：表示带通；bs：表示带阻；
	fs float64:数字系统的采样频率,可选。
		默认情况下,fs是2.0，表示每个周期采样2次
输出:b系数和a系数
说明:
编辑:wang_jp
时间:
***********************************************/
func Butter(N int, Wn interface{}, btype string, fs ...float64) (b, a numgo.Array, err error) {
	return iirFilter(N, Wn, btype, fs...)
}
