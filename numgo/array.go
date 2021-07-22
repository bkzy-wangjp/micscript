package numgo

import (
	"fmt"
	"math"
	"strings"
)

type Array []float64

/*********************************************
//功能:数组求平均
//参数:无
//说明:
//作者:wangjp
//时间:2020年3月24日
**********************************************/
func (arr Array) Mean() float64 {
	var sum float64
	if len(arr) > 0 {
		for _, v := range arr {
			sum += v
		}
		return sum / float64(len(arr))
	} else {
		return 0.0
	}
}

/*********************************************
//功能:数组求和
//参数:无
//说明:
//作者:wangjp
//时间:2020年3月24日
**********************************************/
func (arr Array) Sum() float64 {
	var sum float64
	for _, v := range arr {
		sum += v
	}
	return sum
}

/*********************************************
//功能:数组乘积
//参数:无
//说明:
//作者:wangjp
//时间:2020年5月29日
**********************************************/
func (arr Array) Product() float64 {
	var p float64 = 1.0
	for _, v := range arr {
		p *= v
		if v == 0 {
			break
		}
	}
	return p
}

/*********************************************
//功能:求最小值
//参数:无
//说明:
//作者:wangjp
//时间:2020年5月29日
**********************************************/
func (arr Array) Min() float64 {
	var min float64
	for i, v := range arr {
		if i == 0 {
			min = v
		} else {
			if v < min {
				min = v
			}
		}
	}
	return min
}

/*********************************************
//功能:求最大值
//参数:无
//说明:
//作者:wangjp
//时间:2020年5月29日
**********************************************/
func (arr Array) Max() float64 {
	var max float64
	for i, v := range arr {
		if i == 0 {
			max = v
		} else {
			if v > max {
				max = v
			}
		}
	}
	return max
}

/***********************************************
功能:计算众数
输入:待计算数组,分组数量(不设置时根据数组长度自动分为10组或者100组)
输出:众数、组距、分组Map
编辑:wang_jp
时间:2019年12月6日
***********************************************/
func (data Array) Mode(grp ...int) (float64, float64, map[string]float64) {
	group := 5
	if len(data) > 30 {
		group = 10
	}
	if len(data) > 300 {
		group = 100
	}
	if len(grp) > 0 {
		if grp[0] > 0 {
			group = grp[0]
		}
	}
	min := data.Min()
	max := data.Max()
	r := max - min                      //全距
	gd := math.Ceil(r) / float64(group) //组距

	dar := make(map[int]float64, group)
	dmp := make(map[string]float64, group+2)
	dmp[fmt.Sprintf("%.5f", min-gd)] = 0
	for i := 0; i < group; i++ { //初始化
		dar[i] = 0
		dmp[fmt.Sprintf("%.5f", float64(i)*gd+min)] = 0
	}
	dmp[fmt.Sprintf("%.5f", max+gd)] = 0

	nf := float64(len(data))
	if nf == 0 {
		nf = 1
	}
	max_sub := 0
	for _, v := range data {
		i := int(math.Floor((v - min) / gd))
		dar[i] += 1
		if dar[i] > dar[max_sub] {
			max_sub = i
		}
		dmp[fmt.Sprintf("%.5f", float64(i)*gd+min)] = math.Trunc(dar[i]/nf*1e5) / 1e3 //strconv.ParseFloat(fmt.Sprintf("%.3f", dar[i]/nf*100), 64)
	}
	L := min + float64(max_sub)*gd //组距下限

	mode := L + (dar[max_sub]-dar[max_sub-1])/((dar[max_sub]-dar[max_sub-1])*(dar[max_sub]-dar[max_sub+1]))*gd

	//gd, _, dmp = data.Group(0., 0., group)

	s := fmt.Sprint(mode)
	if strings.Contains(s, "Inf") {
		return 0.0, gd, dmp
	}

	return mode, gd, dmp
}

/***********************************************
功能:对数据进行分组
输入:umin, umax float64:用户自定义的最小和最大分组区间,
		如果不使用自定义,请全部输入0.0
	 grp int:可选的分组数。
	 如果不输入分组数,数据长度0~30分为5组,30~300分为10组,300以上分为100组
输出:组距、分组Map1,分组Map2
说明:Map1以自然数为下标,以组内的数据量为内容
	 Map2以数据分组界限为下标,以组内数据量占数据总数的百分比为内容
编辑:wang_jp
时间:2019年12月6日
***********************************************/
func (data Array) Group(umin, umax float64, grp ...int) (float64, Array, map[string]float64) {
	group := 5
	if len(data) > 30 {
		group = 10
	}
	if len(data) > 300 {
		group = 100
	}
	if len(grp) > 0 {
		group = grp[0]

	}

	min := data.Min()
	max := data.Max()
	group_p := 0
	if umin != 0.0 || umax != 0.0 {
		if min < umin {
			group_p += 1
		}
		min = umin
		max = umax
	}
	r := max - min //全距

	gd := r / float64(group) //组距
	var p float64 = 100.0
	dot := 0
	if gd < 1.0 {
		re := gd * 10
		for dot = 1; re < 1; dot++ {
			re *= 10
			p *= 10
		}
	}
	gd = math.Ceil(gd*p) / p //向上取整
	min = math.Floor((min+gd)*p)/p - gd
	dot_fomat := fmt.Sprintf("%%.%df", dot+2)

	dar := NewArray(group + group_p)
	dmp := make(map[string]float64, group+group_p)

	for i := 0; i < group+group_p; i++ { //初始化
		key := fmt.Sprintf(dot_fomat, float64(i)*gd+min)
		dmp[key] = 0
	}

	nf := float64(len(data))
	if nf == 0 {
		nf = 1
	}
	maxi := group + group_p
	for _, v := range data {
		i := int(math.Floor((v - min) / gd))
		if i >= maxi {
			i = maxi - 1
		}
		if i < 0 {
			i = 0
		}

		j := i
		key := fmt.Sprintf(dot_fomat, float64(j)*gd+min)

		dar[i] += 1
		dmp[key] = math.Trunc(dar[i]/nf*1e5) / 1e3
	}
	return gd, dar, dmp
}

/***********************************************
功能:计算中位数
输入:
输出:中位数
说明:特别适用于偏态分布，对于对称分布也可以应用。对于偏态分部，是概率分布图最高峰所在位置。
	 数组长度必须大于2
编辑:wang_jp
时间:2019年12月6日
***********************************************/
func (data Array) Median() float64 {
	data.BubbleSort() //冒泡排序
	var median float64
	n := len(data)
	if n > 2 {
		if n%2 == 0 {
			median = (data[n/2-1] + data[n/2]) / 2
		} else {
			median = data[n/2]
		}
	}
	return median
}

/***********************************************
功能:四分位数(quartiles)
输入:[group ...int] 可选的分组数,不选时采用数组长度自然分组
输出:
	Lower:箱型图最小值,lower=Q1-1.5*Qd
	Q1:下四分位值
	Q2:等于Median
	Q3:上四分位值
	Upper:箱型图的最大值,upper=Q3+1.5*Qd
	IQR:四分位差,Qd=Q3-Q1,亦写为Qd
说明:数组长度必须大于4
编辑:wang_jp
时间:2020年8月26日
***********************************************/
func (data Array) Quartiles(group ...int) (Lower, Q1, Q2, Q3, Upper, IQR float64) {
	n := len(data) //数据长度或者分组数
	if n < 4 {
		return
	}
	Q2 = data.Median()

	var vq1, vq3 float64
	data.BubbleSort() //冒泡排序

	q1 := (n + 1) / 4
	q3 := 3 * q1
	//限制范围
	if q1 < 1 {
		q1 = 1
	}
	if q1 >= n {
		q1 = n - 1
	}
	if q3 < 1 {
		q3 = 1
	}
	if q3 >= n {
		q3 = n - 1
	}
	if (n+1)%4 == 0 { //能被4 整除
		vq1 = data[q1-1]
		vq3 = data[q3-1]
	} else {
		vq1 = 0.75*data[q1-1] + 0.25*data[q1]
		vq3 = 0.25*data[q3-1] + 0.75*data[q3]
	}
	if len(group) > 0 { //定义了分组数
		n = group[0]
		q1 = (n + 1) / 4
		q3 = 3 * q1
		//限制范围
		if q1 < 1 {
			q1 = 1
		}
		if q1 >= n {
			q1 = n - 1
		}
		if q3 < 1 {
			q3 = 1
		}
		if q3 >= n {
			q3 = n - 1
		}
		min := data.Min()
		max := data.Max()
		dar := make([]float64, n)       //数据分组数
		gar := make([]float64, n)       //每组的累计数
		r := max - min                  //全距
		gd := math.Ceil(r) / float64(n) //组距

		for i := 0; i < n; i++ {
			dar[i] = min + float64(i)*gd
		}
		for _, v := range data {
			i := int(math.Floor((v - min) / gd))
			if i >= 0 && i < n {
				gar[i] += 1
			}
		}

		var sq1 float64 //Q1所在组以下的累计次数
		for i := 0; i < q1; i++ {
			sq1 += gar[i-1]
		}
		var sq3 float64 //Q3所在组以下的累计次数
		for i := 0; i < q3; i++ {
			sq3 += gar[i-1]
		}
		if gar[q1] != 0 {
			vq1 = dar[q1-1] + (float64(n)/4.0-sq1)/gar[q1-1]*gd
		}
		if gar[q3] != 0 {
			vq3 = dar[q3-1] + (float64(n)/4.0-sq3)/gar[q3-1]*gd
		}
	}
	IQR = vq3 - vq1       //Qd四分位差
	Lower = vq1 - 1.5*IQR //Lower
	Q1 = vq1              //Q1
	Q3 = vq3              //Q3
	Upper = vq3 + 1.5*IQR //Upper
	return
}

/***********************************************
功能:四分位数(quartile)
输入:[q int] 分位标记,可选-1、0、1、2、3、4,
		-1:Qd,四分位差,Qd=Q3-Q1,亦写为IQR
		0:Lower,箱型图最小值,lower=Q1-1.5*Qd
		1:Q1,下四分位值
		2:Q2,等于Median
		3:Q3,上四分位值
		4:Upper,箱型图的最大值,upper=Q3+1.5*Qd
	[group ...int] 可选的分组数,不选时采用数组长度自然分组
输出:四分位数
说明:将各个变量值按大小顺序排列，然后将此数列分成四等份
	 数组长度必须大于4
编辑:wang_jp
时间:2020年7月5日
***********************************************/
func (data Array) Quartile(q int, group ...int) float64 {
	Lower, Q1, Q2, Q3, Upper, IQR := data.Quartiles(group...)
	switch q {
	case -1:
		return IQR //Qd四分位差
	case 0:
		return Lower //Lower
	case 1:
		return Q1 //Q1
	case 2:
		return Q2 //Q2
	case 3:
		return Q3 //Q3
	case 4:
		return Upper //Upper
	default:
		return IQR //四分位差
	}
}

/***********************************************
功能:冒泡排序
输入:confg[0] 顺序，false:从小到大(默认),true:从大到小
	 confg[1] 是否替换原数组，false:不替换,true:替换
输出:无
编辑:wang_jp
时间:2019年12月6日
***********************************************/
func (arr Array) BubbleSort(confg ...bool) Array {
	flag := true
	vLen := len(arr)
	des := false     //排序方式
	inplace := false //是否替换
	n := len(confg)
	if n > 0 {
		des = confg[0]
		if n > 1 {
			inplace = confg[1]
		}
	}

	values := arr.Copy()
	if inplace {
		values = arr
	}

	if des == false {
		for i := 0; i < vLen-1; i++ {
			flag = true
			for j := 0; j < vLen-i-1; j++ {
				if values[j] > values[j+1] {
					values[j], values[j+1] = values[j+1], values[j]
					flag = false
					continue
				}
			}
			if flag {
				break
			}
		}
	} else {
		for i := 0; i < vLen-1; i++ {
			flag = true
			for j := 0; j < vLen-i-1; j++ {
				if values[j] < values[j+1] {
					values[j], values[j+1] = values[j+1], values[j]
					flag = false
					continue
				}
			}
			if flag {
				break
			}
		}
	}
	return values
}

/***********************************************
功能:复制一个新的数组
输入:无
输出:Array;新的数组
说明:
编辑:wang_jp
时间:2020年10月20日
***********************************************/
func (arr Array) Copy() Array {
	var res Array
	for _, c := range arr {
		res = append(res, c)
	}
	return res
}

/***********************************************
功能:高级统计功能
输入:待计算数组
输出:平均值、标准偏差、样本标准变差、标准差、偏度、峰度
编辑:wang_jp
时间:2019年12月6日
***********************************************/
func (data Array) CalcAdvangceStatistic() (avg, sd, std, se, ske, kur float64) {
	var sum float64 = 0
	n := len(data)
	nf := float64(n)
	if n > 1 {
		for _, v := range data {
			sum += v
		}
		avg = sum / nf //平均值

		var variance float64
		for _, v := range data {
			variance += math.Pow((v - avg), 2)
		}
		sd = math.Sqrt(variance / nf)        //标准偏差
		std = math.Sqrt(variance / (nf - 1)) //样本标准偏差
		se = sd / float64(n)                 //标准误差

		if n > 3 && std != 0 {
			var a, b float64
			for _, v := range data {
				a += math.Pow((v-avg)/std, 3)
				b += math.Pow((v-avg)/std, 4)
			}
			ske = (nf / ((nf - 1) * (nf - 2))) * a                                             //偏度
			kur = (nf*(nf-1)/((nf-1)*(nf-2)*(nf-3)))*b - 3*math.Pow((nf-2), 2)/((nf-2)*(nf-3)) //峰度
		}
	}
	return avg, sd, std, se, ske, kur
}

/***********************************************
功能:算术平均值
输入:无
输出:平均值
说明:当数据呈对称分布或接近对称分布时，均值、中位数、众数相等或接近相等，这时应选择均值作为集中
	趋势的代表值，因为均值包含了全部数据的信息
编辑:wang_jp
时间:2019年12月6日
***********************************************/
func (data Array) Average() float64 {
	var avg, sum float64 = 0, 0
	nf := float64(len(data))
	if nf > 1 {
		for _, v := range data {
			sum += v
		}
		avg = sum / nf //平均值
	}
	return avg
}

/***********************************************
功能:方差(样本值与平均值的差的平方)
输入:
输出:方差
说明:数列中每个元素与均值之差的平方和
编辑:wang_jp
时间:2019年12月6日
***********************************************/
func (data Array) Variance() float64 {
	avg := data.Average()
	var variance float64
	for _, v := range data {
		variance += math.Pow((v - avg), 2)
	}
	return variance
}

/***********************************************
功能:标准差(均方差)
输入:
输出:标准差
说明:标准差（Standard Deviation） ，中文环境中又常称均方差，是离均差平方的算术平均数的平方根，
	用σ表示。标准差是方差的算术平方根。标准差能反映一个数据集的离散程度。
编辑:wang_jp
时间:2019年12月6日
***********************************************/
func (data Array) Sd() (sd float64) {
	nf := float64(len(data))
	variance := data.Variance()
	sd = math.Sqrt(variance / nf) //标准偏差
	return sd
}

/***********************************************
功能:样本标准偏差
输入:
输出:样本标准偏差
编辑:wang_jp
时间:2019年12月6日
***********************************************/
func (data Array) Std() (std float64) {
	nf := float64(len(data))
	variance := data.Variance()
	//sd = math.Sqrt(variance / nf) //标准差
	std = math.Sqrt(variance / (nf - 1)) //样本标准偏差
	//se = sd / nf                         //标准误差
	return std
}

/***********************************************
功能:标准误差（Standard error）
输入:
输出:标准误差
编辑:wang_jp
时间:2019年12月6日
***********************************************/
func (data Array) Se() (se float64) {
	nf := float64(len(data))
	variance := data.Variance()
	sd := math.Sqrt(variance / nf) //标准差
	//std = math.Sqrt(variance / (nf - 1)) //样本标准偏差
	se = sd / nf //标准误差
	return se
}

/***********************************************
功能:偏度
输入:
输出:偏度
编辑:wang_jp
时间:2019年12月6日
***********************************************/
func (data Array) Ske() (ske float64) {
	nf := float64(len(data))
	avg := data.Average()
	std := data.Std()
	if nf > 3 && std != 0 {
		var a float64
		for _, v := range data {
			a += math.Pow((v-avg)/std, 3)
		}
		ske = (nf / ((nf - 1) * (nf - 2))) * a //偏度
	}
	return ske
}

/***********************************************
功能:峰度
输入:
输出:峰度
编辑:wang_jp
时间:2019年12月6日
***********************************************/
func (data Array) Kur() (kur float64) {
	nf := float64(len(data))
	avg := data.Average()
	std := data.Std()
	if nf > 3 && std != 0 {
		var b float64
		for _, v := range data {
			b += math.Pow((v-avg)/std, 4)
		}
		kur = (nf*(nf-1)/((nf-1)*(nf-2)*(nf-3)))*b - 3*math.Pow((nf-2), 2)/((nf-2)*(nf-3)) //峰度
	}
	return kur
}

/***********************************************
功能:变异系数
输入:
输出:变异系数
说明:变异系数，又称“离散系数”（英文：coefficient of variation），是概率分布离散程度的一个归一化量度，
	其定义为标准差与平均值之比。变异系数也被称为标准离差率或单位风险
编辑:wang_jp
时间:2019年12月6日
***********************************************/
func (data Array) Cv() (cv float64) {
	avg := data.Average()
	sd := data.Sd()
	if avg != 0 {
		cv = sd / avg * 100.0
	}
	return cv
}

/***********************************************
功能:数据范围（极差或全距）
输入:
输出:数据范围
说明:最大值与最小值之间的差值，用于描述X的数字分散程度，越小则数字之间越紧密
编辑:wang_jp
时间:2019年12月6日
**********************************************/
func (data Array) Range() float64 {
	return data.Max() - data.Min()
}

/***********************************************
功能:中程数
输入:
输出:中程数
说明:（最大值 + 最小值）/2
编辑:wang_jp
时间:2019年12月6日
***********************************************/
func (data Array) Midrange() float64 {
	return (data.Max() + data.Min()) / 2
}

/**********************************************
//功能: 移动窗口滤波
//输入: n:int,滤波窗口长度,1<n<len(data)
//		fillvalue:string, 填充值方法{"mean","median"},省略时为"mean"
//输出:
//说明: n<=1或者n>len(data)时,不进行滤波,直接返回原数据
//		"median"时,n必须大于2，否则用"mean"方法
//编辑: wangjp
//时间: 2020年10月14日
**********************************************/
func (data Array) MoveWindowFilter(n int, fillvlue ...string) Array {
	d_len := len(data)
	//数据长度小于滤波窗口长度时,不进行滤波
	if n <= 1 || d_len < n {
		return data
	}

	fv_method := "mean"
	if len(fillvlue) > 0 {
		fv_method = fillvlue[0]
	}

	var res Array
	for i, _ := range data {
		ed := n + i
		if ed > d_len {
			ed = d_len
		}
		var ar Array
		ar = data[i:ed]
		fv := ar.Mean()
		switch fv_method {
		case "mean":
			fv = ar.Mean()
		case "median":
			if len(ar) > 2 {
				fv = ar.Median()
			}
		default:
			fv = ar.Mean()
		}
		res = append(res, fv)
	}
	return res
}

/***********************************************
功能:每个元素都减去一个标量
输入:scalar:标量值
输出:
说明:
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func (data Array) SubScalar(scalar float64) {
	for i, v := range data {
		data[i] = v - scalar
	}
}

/***********************************************
功能:每个元素都加上一个标量
输入:scalar:标量值
输出:
说明:
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func (data Array) AddScalar(scalar float64) {
	for i, v := range data {
		data[i] = v + scalar
	}
}

/***********************************************
功能:每个元素都被一个标量减
输入:scalar:标量值
输出:
说明:
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func (data Array) SubByScalar(scalar float64) {
	for i, v := range data {
		data[i] = scalar - v
	}
}

/***********************************************
功能:每个元素都乘以一个标量
输入:scalar:标量值
输出:
说明:
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func (data Array) MulScalar(scalar float64) {
	for i, v := range data {
		data[i] = v * scalar
	}
}

/***********************************************
功能:每个(非零)元素都被一个标量除
输入:scalar:标量值
输出:
说明:
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func (data Array) DivByScalar(scalar float64) {
	for i, v := range data {
		if v != 0 {
			data[i] = scalar / v
		}
	}
}

/***********************************************
功能:每个元素都除以一个标量
输入:scalar:标量值
输出:
说明:
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func (data Array) DivScalar(scalar float64) {
	if scalar == 0 {
		return
	}
	for i, v := range data {
		data[i] = v / scalar
	}
}

/***********************************************
功能:两个数组是否相等
输入:数组B
输出:数组相等输出true，否则false
说明:
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func (a Array) IsEqual(b Array) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if math.Abs(b[i]-v) > 1e-5 { //浮点数不易直接比较
			return false
		}
	}
	return true
}

/***********************************************
功能:将数组反转
输入:
输出:
说明:
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func (a Array) Reverse() {
	alen := len(a)
	b := make(Array, alen)
	copy(b, a)
	alen--
	for i, v := range b {
		a[alen-i] = v
	}
}

/***********************************************
功能:初始化一个新数组
输入:n:数组的长度
输出:元素全为0的长度为n的数组
说明:
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func NewArray(n int) Array {
	return make(Array, n)
}

/***********************************************
功能: 生成伴随矩阵
输入: 无
输出: 数组的伴随矩阵(N-1 x N-1维度方阵)
	  伴随矩阵的第一行是:-a[1:]/a[0]，第二行开始的对角线都是1，其余为0
说明: a 数组的长度N必须大于等于2,a[0]必须非0
编辑: wang_jp
时间: 2020年10月15日
***********************************************/
func (a Array) Companion() (Matrix, error) {
	alen := len(a)
	if alen < 2 {
		return nil, fmt.Errorf("The length of `a` must be at least 2[原始A数组的长度不能小于2]")
	}
	if a[0] == 0 {
		return nil, fmt.Errorf("a[0] must not be zero[ a[0]的值不可为0 ]")
	}
	var c1 Array //伴随矩阵的第一行
	for _, av := range a[1:] {
		c1 = append(c1, av*-1.0/a[0])
	}
	var c Matrix
	c = append(c, c1)
	if alen == 2 { //a的长度为2时，直接返回
		return c, nil
	}
	tmp := InitAsIdentity(alen - 2)
	for _, tr := range tmp {
		tr = append(tr, 0.0)
		c = append(c, tr)
	}
	return c, nil
}

/***********************************************
功能:数组乘以另一个数组
输入:一个同维度的数组
输出:
说明:数组中的每个下标上对应的元素相乘
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func (a Array) MulArray(b Array) (Array, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("Both arrays must be the same length.[两个数组的长度必须相同]")
	}
	var c Array
	for i, v := range a {
		c = append(c, v*b[i])
	}
	return c, nil
}

/***********************************************
功能:数组除以另一个数组
输入:一个同维度的数组
输出:
说明:数组中的每个下标上对应的元素相除，如果被除数元素为0，则结果元素为0
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func (a Array) DivArray(b Array) (Array, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("Both arrays must be the same length.[两个数组的长度必须相同]")
	}
	var c Array
	for i, v := range b {
		if v == 0 {
			c = append(c, v)
		} else {
			c = append(c, a[i]/v)
		}
	}
	return c, nil
}

/***********************************************
功能:数组加上另一个数组
输入:一个同维度的数组
输出:
说明:数组中的每个下标上对应的元素相加
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func (a Array) AddArray(b Array) (Array, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("Both arrays must be the same length.[两个数组的长度必须相同]")
	}
	var c Array
	for i, v := range b {
		c = append(c, a[i]+v)
	}
	return c, nil
}

/***********************************************
功能:数组减去另一个数组
输入:一个同维度的数组
输出:
说明:数组中的每个下标上对应的元素相减
编辑:wang_jp
时间:2020年10月15日
***********************************************/
func (a Array) SubArray(b Array) (Array, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("Both arrays must be the same length.[两个数组的长度必须相同]")
	}
	var c Array
	for i, v := range b {
		c = append(c, a[i]-v)
	}
	return c, nil
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
func (a Array) IsAllZero() bool {
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
func (X Array) Polynominal() Matrix {
	n := len(X)
	var cmat Matrix //行列式
	for _, x := range X {
		var mr Array //行列式的行
		for i := n - 1; i >= 0; i-- {
			mr = append(mr, math.Pow(x, float64(i)))
		}
		mr = append(mr, -math.Pow(x, float64(n))) //
		cmat = append(cmat, mr)
	}
	return cmat
}

/***********************************************
功能:用给定的根序列求多项式的系数。
输入:
输出:c:Array,长度为N+1的矩阵,c[0]始终为1.0
说明:c[0] * x**(N) + c[1] * x**(N-1) + ... + c[N-1] * x + c[N]
	对于给定的零序列，返回其前导系数为1的多项式的系数（序列中必须包含多个根是
	其重数的许多倍）
编辑:wang_jp
时间:2020年10月20日
***********************************************/
func (seq_of_zeros Array) Poly() (Array, error) {
	var c Array
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
功能:交换两个元素
输入:k1,k2 int:需要交换的两个元素下标
输出:
说明:
编辑:wang_jp
时间:2020年11月09日
***********************************************/
func (a Array) Swap(k1, k2 int) {
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
输入:inv ...bool:省略或者false:傅里叶变换;true:傅里叶逆变换
输出:
说明:尚未完成
编辑:wang_jp
时间:2020年11月09日
***********************************************/
func (a Array) Fft(inv ...bool) Array {
	invs := false
	if len(inv) > 0 {
		invs = inv[0]
	}
	car := RealArr2ComplexsArr(a)
	car.Fft(true, invs)
	res := a.Copy()
	for i, _ := range res {
		if i == 0 {
			res[i] = real(car[i])
		} else {
			if i%2 > 0 {
				res[i] = real(car[(i+1)/2])
			} else {
				res[i] = imag(car[(i+1)/2])
			}
		}
	}
	return res
}
