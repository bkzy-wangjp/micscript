// filter project filter.go
package filter

import (
	"fmt"

	"github.com/bkzy/micscript/numgo"
)

/***********************************************
功能:使用N阶Butterworth(巴特沃斯)数字滤波器对数据进行滤波
输入:X numgo.Array:待滤波的原始数据
	N int: 滤波阶数
	Wn float/[2]float:自然频率数组，也称归一化的截止频率,与fs具有相同的单位
		fcf=截止频率*2/采样频率；
		如果是低通，高通滤波，fcf只有一个元素,范围:0.0-1.0
		如果是带通，带阻滤波，fcf数组有两个元素,分别是滤波带的起始频率和阶数频率,范围:0.0-fs/2.0
	btype string:滤波类型，字符串；lp：表示低通；hp：表示高通；bp：表示带通；bs：表示带阻；
	fs float64:数字系统的采样频率,可选。
		默认情况下,fs是2.0，表示每个周期采样2次
输出:滤波后的数据,错误信息
说明:
编辑:wang_jp
时间:2020年10月21日
***********************************************/
func ButterFilter(X numgo.Array, N int, Wn interface{}, btype string, fs ...float64) (numgo.Array, error) {
	b, a, err := Butter(N, Wn, btype, fs...)
	if err != nil {
		return nil, fmt.Errorf("filter parameter design error[滤波器参数设计错误]:%s", err.Error())
	}
	return Filtfilt(b, a, X)
}

/***********************************************
功能:输入参数归一化
输入:b参数和a参数
输出:
说明:a参数的第一个元素a[0]必须为1.0。如果不为1，则需要执行归一化操作
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func normalization(b, a numgo.Array) (numgo.Array, numgo.Array, error) {
checka:
	if len(a) > 1 && a[0] == 0.0 { //如果a[0]为0,则从抛弃之,循环判断执行
		a = a[1:]
		goto checka
	}

	if len(a) < 1 {
		return b, a, fmt.Errorf("there must be at least one nonzero `a` coefficient.[至少要有一个有效的A参数]")
	}

	//把a[0]标准化为1.0
	firsta := a[0]
	if firsta != 1.0 {
		b.DivScalar(firsta)
		a.DivScalar(firsta)
	}
	//取a\b长度的最大值为长度
	n := len(a)
	if len(b) > n {
		n = len(b)
	}
	//长度不足的用0补齐，使a \ b长度相同
	if len(a) < n {
		add := make(numgo.Array, n-len(a))
		a = append(a, add...)
	}
	if len(b) < n {
		add := make(numgo.Array, n-len(b))
		b = append(b, add...)
	}
	return b, a, nil
}

/***********************************************
功能:过滤
输入:b,a:差分方程的系数
	 X:待滤波的原始数据
	 Zi:初始状态
输出:
说明:
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func filter(b, a, X, zi numgo.Array) {
	//标准化参数
	//normalization(b, a)
	//取a\b长度的最大值为长度
	nfilt := len(a)
	if len(b) > nfilt {
		nfilt = len(b)
	}

	a[0] = 0.0 //设置a0
	y_tmp := numgo.NewArray(nfilt)
	for i := 0; i < nfilt-1; i++ {
		for j := 0; j <= i; j++ {
			y_tmp[i] += (b[j]*X[i-j] - a[j]*y_tmp[i-j])
		}
		y_tmp[i] += zi[i]
	}
	xlen := len(X)
	i := 0
	for i = nfilt - 1; i < xlen; i++ {
		y_tmp[nfilt-1] = b[0] * X[i]
		for j := 1; j < nfilt; j++ {
			y_tmp[nfilt-1] += (b[j]*X[i-j] - a[j]*y_tmp[nfilt-1-j])
		}

		X[i-nfilt+1] = y_tmp[0]

		if nfilt > 1 {
			for i, y := range y_tmp[1:] {
				y_tmp[i] = y
			}
		}
	}
	for j := 0; j < nfilt-1; j++ {
		X[i-nfilt+1] = y_tmp[j]
		i++
	}

	a[0] = 1.0 //恢复a0
}

/***********************************************
功能: 初始化临时数据
输入: X:待滤波的原始数据
	  nflit:A\B参数的长度
输出: 临时存储数组,长度为:6*nflit+len(X)
说明:
编辑: wang_jp
时间: 2020年10月15日
***********************************************/
func initTx(X numgo.Array, nfilt int) numgo.Array {
	xlen := len(X)
	nfact := 3 * nfilt //length of edge transients(边缘效应的长度)
	//tlen := 2*nfact + xlen //临时X区的长度

	tmp := X[0]
	var tx numgo.Array
	//初始化tx[:nfact]
	for i := nfact; i > 0; i-- {
		tx = append(tx, 2.0*tmp-X[i])
	}
	//初始化tx[nfact:nfact+xlen]
	tx = append(tx, X...)
	//初始化tx[nfact+xlen:]
	tmp = X[xlen-1]
	for i := 0; i < nfact; i++ {
		tx = append(tx, 2.0*tmp-X[xlen-(i+2)])
	}
	return tx
}

/***********************************************
功能: 初始化zi数组
输入: b,a:差分方程的系数
输出: zi数组,长度：nfilt-1
说明: 参照python scipy.linalg.lfilter_zi
编辑: wang_jp
时间: 2020年10月15日
***********************************************/
func initZi(b, a numgo.Array) (numgo.Array, error) {
	//取a\b长度的最大值为长度
	nfilt := len(a)
	if len(b) > nfilt {
		nfilt = len(b)
	}
	eye := numgo.InitAsIdentity(nfilt - 1) //nfilt-1 x nfilt-1矩阵
	comp, err := a.Companion()
	if err != nil {
		return nil, err
	}
	ct := comp.Transpose()

	IminusA, err := eye.Sub(ct) //矩阵差
	if err != nil {
		return nil, err
	}
	var B numgo.Array
	for i, bv := range b[1:] {
		B = append(B, bv-a[i+1]*b[0])
	}

	var iminusa0 numgo.Array
	iat := IminusA.Transpose()
	iminusa0 = iat[0]
	a0sum := iminusa0.Sum() //IminusA第一列的和
	if a0sum == 0.0 {
		return nil, fmt.Errorf("IminusA[0].Sum must not be zero.[IminusA[0].Sum 不能为0]")
	}

	var zi numgo.Array
	zi = append(zi, B.Sum()/a0sum)
	asum := 1.0
	csum := 0.0
	for i := 1; i < nfilt-1; i++ { //填充sp矩阵
		asum += a[i]
		csum += b[i] - a[i]*b[0]
		zi = append(zi, asum*zi[0]-csum)
	}

	return zi, nil
}

/***********************************************
功能:滤波
输入:b,a:差分方程的系数
	 X:待滤波的原始数据
输出:滤波后的数据,错误信息
说明:
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func Filtfilt(b, a, X numgo.Array) (numgo.Array, error) {
	//标准化参数
	normalization(b, a)
	//取a\b长度的最大值为长度
	nfilt := len(a)
	if len(b) > nfilt {
		nfilt = len(b)
	}

	xlen := len(X)     //X数组的长度
	nfact := 3 * nfilt //length of edge transients(边缘效应的长度)
	if xlen <= nfact || nfilt < 2 {
		if xlen <= nfact {
			return nil, fmt.Errorf("the length of the input x must be more than three times the filter order, defined as max(length(b)-1,length(a)-1)[输入数据X的长度必须大于3倍的A或者B参数的长度]")
		} else {
			return nil, fmt.Errorf("the length of A and B must greater than 1[A参数或者B参数的长度必须大于1]")
		}
	}

	tx_arr := initTx(X, nfilt)
	zi_arr, err := initZi(b, a) //zi矩阵,nfilt-1 x 1型
	if err != nil {
		return nil, err
	}

	zi_data := make(numgo.Array, len(zi_arr))
	copy(zi_data, zi_arr) //获取初始化好的zi
	tx0 := tx_arr[0]
	for i, zv := range zi_data {
		zi_data[i] = zv * tx0
	}
	//第一次滤波
	filter(b, a, tx_arr, zi_data)
	//反转tx_arr
	tx_arr.Reverse()

	copy(zi_data, zi_arr) //获取初始化好的zi
	tx0 = tx_arr[0]
	for i, zv := range zi_data {
		zi_data[i] = zv * tx0
	}
	//第二次滤波
	filter(b, a, tx_arr, zi_data)
	//再反转tx_arr
	tx_arr.Reverse()
	return tx_arr[nfact : len(tx_arr)-nfact], nil
}
