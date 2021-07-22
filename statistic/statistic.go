// statistic
package statistic

import (
	"fmt"
	"strings"
	"time"

	"github.com/bkzy/micscript/filter"

	"github.com/bkzy/micscript/numgo"
)

/************************************************************
功能:对时序数据统计计算
输入:data []TimeSeriesData:历史数据数组
	cfg int:0=仅基础统计,1=需要高级统计,2=需要原始数据
	group int:高级计算时的分组数量,默认0=100
返回:统计计算结果
编辑:wang_jp
时间:2019年12月7日
************************************************************/
func (tsds Tsds) Statistics(cfg int, group ...int) StatisticData {
	var total, sum float64 = 0.0, 0.0
	var t0, st time.Time
	var v0, min, max, datarange, diff, plusdiff, cavg, pavg float64
	var count, rising, falling, ltz, gtz, ez, outlers int = 0, 0, 0, 0, 0, 0, 0
	var dataArr []float64
	var summary StatisticData
	var bgt, edt time.Time //最早和最后数据时间

	lower, q1, q2, q3, upper, qd := tsds.Quartiles()
	count = len(tsds)
	increment := make(map[string]float64, count)

	for k, dv := range tsds { //遍历数据
		t := dv.Time
		if dv.Value == 0 { //等于零计次
			ez += 1
		}
		if dv.Value > 0 { //大于零计次
			gtz += 1
		}
		if dv.Value < 0 { //小于零计次
			ltz += 1
		}
		if dv.Value < lower || dv.Value > upper { //离群点
			outlers += 1
		}
		if k > 0 { //非第一个数
			if dv.Time.UnixNano() > edt.UnixNano() {
				edt = dv.Time //最后数据时间
			}
			if dv.Time.UnixNano() < bgt.UnixNano() {
				bgt = dv.Time //最早数据时间
			}
			v := dv.Value
			total += (v0 + v) / 2 * float64(t.Unix()-t0.Unix())
			if v > v0 { //数据上升统计
				rising += 1
			}
			if v < v0 { //数据下降统计
				falling += 1
			}
			increment[dv.Time.Format(_TIME_LYOUT)] = v - v0

			v0 = v
		} else { //第一个数
			v0 = dv.Value
			min = v0
			max = v0
			diff = v0
			plusdiff = v0
			st = t
			sum = 0.0
			bgt = dv.Time
			edt = dv.Time
		}
		dataArr = append(dataArr, v0) //存入缓存数组
		t0 = t
		if v0 < min { //最小值
			min = v0
		}
		if v0 > max { //最大值
			max = v0
		}
		sum += v0 //算术累积值
	}
	diff = v0 - diff   //差值
	if plusdiff < v0 { //初始值小于结束值
		if count <= 2 { //如果原始数据小于等于2个
			plusdiff = v0 //等于最后值
		} else {
			plusdiff = max - plusdiff + v0 //最大值-第一个值+最后一个值
		}
	}
	if count > 0 { //算术平均值
		cavg = sum / float64(count)
	}
	dur := t0.Unix() - st.Unix() //持续时间
	if dur > 0 {                 //加权平均值
		pavg = total / float64(dur)
	}
	datarange = max - min

	summary.Diff = diff
	summary.PlusDiff = plusdiff
	summary.Duration = dur
	summary.FallingCnt = falling
	summary.Mean = cavg
	summary.Min = min
	summary.Max = max
	summary.PointCnt = count
	summary.PowerAvg = pavg
	summary.Range = datarange
	summary.RisingCnt = rising
	summary.Sum = sum
	summary.Total = total
	summary.Increment = increment
	summary.LtzCnt = ltz
	summary.GtzCnt = gtz
	summary.EzCnt = ez
	summary.OutliersCnt = outlers
	summary.Median = q2
	summary.Q1 = q1
	summary.Q3 = q3
	summary.Qd = qd
	summary.Upper = upper
	summary.Lower = lower
	summary.BeginTime = bgt.Format(_TIME_LYOUT)
	summary.EndTime = edt.Format(_TIME_LYOUT)

	if cfg > 0 { //是否需要高级统计
		var nga numgo.Array
		nga = dataArr
		mode, gd, dmp := nga.Mode(group...)                     //求众数
		_, sd, std, se, ske, kur := nga.CalcAdvangceStatistic() //求标准偏差等

		summary.Kur = kur
		summary.SD = sd
		summary.SE = se
		summary.Ske = ske
		summary.STDDEV = std
		summary.Mode = mode
		summary.DataGroup = dmp
		summary.GroupDist = gd
	}
	if cfg > 1 { //需要原始数据
		summary.RawData = tsds.ToTimeSeriesDataS()
	}

	return summary
}

/************************************************************
功能:时间参数格式化
输入:时间字符串,可选的时区
输出:格式化后的时间变量,错误信息
时间:2019年11月28日
编辑:wang_jp
************************************************************/
func TimeParse(s string, loc ...*time.Location) (time.Time, error) {
	location := time.Local
	if len(loc) > 0 {
		for _, v := range loc {
			location = v
		}
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006-01-02 15:04", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006-1-2 15:04:05", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006-1-2 15:04", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006/01/02 15:04:05", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006/01/02 15:04", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006/1/2 15:04:05", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006/1/2 15:04", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006-01-02T15:04:05Z", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006-01-02T15:04:05", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006-01-02T15:04", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006-1-2T15:04:05Z", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006-1-2T15:04", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("20060102150405", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("200601021504", s, location)
	if err == nil {
		return t, nil
	}
	return t, err
}

/************************************************************
功能:根据需要的数据指令从历史统计结构体中返回相应的数据
输入:历史统计结构体,数据类型指令
输出:浮点数
时间:2019年12月23日
编辑:wang_jp
************************************************************/
func SelectValueFromStatisticData(v StatisticData, keyword string) float64 {
	switch strings.ToLower(keyword) { //查询并返回结果
	case "min":
		return v.Min
	case "max":
		return v.Max
	case "range":
		return v.Range
	case "total":
		return v.Total
	case "sum":
		return v.Sum
	case "mean":
		return v.Mean
	case "poweravg":
		return v.PowerAvg
	case "diff":
		return v.Diff
	case "plusdiff":
		return v.PlusDiff
	case "duration", "dur", "times":
		return float64(v.Duration)
	case "pointcnt", "count", "pcnt":
		return float64(v.PointCnt)
	case "risingcnt", "rising", "rising_counts", "rcnt":
		return float64(v.RisingCnt)
	case "fallingcnt", "falling", "falling_counts", "fcnt":
		return float64(v.FallingCnt)
	case "ltzcnt", "ltz":
		return float64(v.LtzCnt)
	case "gtzcnt", "gtz":
		return float64(v.GtzCnt)
	case "ezcnt", "ez":
		return float64(v.EzCnt)
	case "sd":
		return v.SD
	case "stddev":
		return v.STDDEV
	case "se":
		return v.SE
	case "ske":
		return v.Ske
	case "kur":
		return v.Kur
	case "mode":
		return v.Mode
	case "q2", "median":
		return v.Median
	case "q1":
		return v.Q1
	case "q3":
		return v.Q3
	case "qd", "iqr":
		return v.Qd
	case "q4", "upper":
		return v.Upper
	case "q0", "lower":
		return v.Lower
	case "groupdist":
		return v.GroupDist
	default:
		return 0.0
	}
}

/************************************************************
功能:字符串时间格式的时序数据转换为时间类时间格式的时序数据
输入:[]TimeSeriesDataS
输出:[]TimeSeriesData,error
时间:2020年5月17日
编辑:wang_jp
************************************************************/
func Tsds2Tsdt(tsds []TimeSeriesDataS) (Tsds, error) {
	var tsdt []TimeSeriesData
	var err error
	for _, d := range tsds {
		var dt TimeSeriesData
		dt.Time, err = TimeParse(d.Time)
		if err != nil {
			return tsdt, err
		}
		dt.Value = d.Value
		tsdt = append(tsdt, dt)
	}
	return tsdt, nil
}

/************************************************************
功能:Unix毫秒时间格式的时序数据转换为时间类时间格式的时序数据
输入:[]TimeSeriesDataS
输出:[]TimeSeriesData,error
时间:2020年5月17日
编辑:wang_jp
************************************************************/
func Tsdi2Tsdt(tsdi []TimeSeriesDataI) Tsds {
	var tsdt []TimeSeriesData
	for _, d := range tsdi {
		var dt TimeSeriesData
		dt.Time = time.Unix(d.Time/1e3, d.Time%1e3*1e6)
		dt.Value = d.Value
		tsdt = append(tsdt, dt)
	}
	return tsdt
}

/************************************************************
功能:对时序数据按时间排序
输入:[desc] 顺序，false:从小到大(默认),true:从大到小
输出:无
时间:2020年8月26日
编辑:wang_jp
************************************************************/
func (tsds Tsds) SortByTime(desc ...bool) {
	flag := true
	vLen := len(tsds)
	des := false
	if len(desc) > 0 {
		des = desc[0]
	}
	if des == false {
		for i := 0; i < vLen-1; i++ {
			flag = true
			for j := 0; j < vLen-i-1; j++ {
				if tsds[j].Time.UnixNano() > tsds[j+1].Time.UnixNano() {
					tsds[j], tsds[j+1] = tsds[j+1], tsds[j]
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
				if tsds[j].Time.UnixNano() < tsds[j+1].Time.UnixNano() {
					tsds[j], tsds[j+1] = tsds[j+1], tsds[j]
					flag = false
					continue
				}
			}
			if flag {
				break
			}
		}
	}
}

/************************************************************
功能:获取时间序列数组中的数值
输入:无
输出:numgo.Array
时间:2020年8月26日
编辑:wang_jp
************************************************************/
func (tsds Tsds) GetDataArray() numgo.Array {
	var dataArr numgo.Array
	for _, tsd := range tsds {
		dataArr = append(dataArr, tsd.Value)
	}
	return dataArr
}

/************************************************************
功能:将数组传回时序数据
输入:numgo.Array
输出:
时间:2020年10月22日
编辑:wang_jp
************************************************************/
func (tsds Tsds) PutArray2Data(arr numgo.Array) error {
	if len(tsds) != len(arr) {
		return fmt.Errorf("PutArray2Data: Destination array and source array must be the same length.[目的数组和源数组长度必须相同]")
	}
	for i, tsd := range tsds {
		tsd.Value = arr[i]
		tsds[i] = tsd
	}
	return nil
}

/************************************************************
功能:获取四分位数(quartiles)
输入:[group ...int] 可选的分组数,不选时采用数组长度自然分组
输出:
	Lower:箱型图最小值,lower=Q1-1.5*Qd
	Q1:下四分位值
	Q2:等于Median
	Q3:上四分位值
	Upper:箱型图的最大值,upper=Q3+1.5*Qd
	IQR:四分位差,IQR=Q3-Q1,亦写为Qd
说明:
编辑:wang_jp
时间:2020年8月26日
************************************************************/
func (tsds Tsds) Quartiles(group ...int) (Lower, Q1, Q2, Q3, Upper, IQR float64) {
	var darr numgo.Array
	darr = tsds.GetDataArray()
	Lower, Q1, Q2, Q3, Upper, IQR = darr.Quartiles(group...)
	return
}

/************************************************************
功能:用中位数填充离群点
输入:无
输出:无
说明:超出Lower和Upper范围的数为离群点
时间:2020年8月26日
编辑:wang_jp
************************************************************/
func (tsds Tsds) FillOutliersByMidean() {
	Lower, _, Q2, _, Upper, _ := tsds.Quartiles()
	for i, tsd := range tsds {
		if tsd.Value < Lower || tsd.Value > Upper {
			tsd.Value = Q2
			tsds[i] = tsd
		}
	}
}

/************************************************************
功能:用四分位Lower或者Upper填充离群点
输入:无
输出:无
说明:超出Lower和Upper范围的数为离群点
时间:2020年8月26日
编辑:wang_jp
************************************************************/
func (tsds Tsds) FillOutliersByExtremum() {
	Lower, _, _, _, Upper, _ := tsds.Quartiles()
	for i, tsd := range tsds {
		if tsd.Value < Lower {
			tsd.Value = Lower
			tsds[i] = tsd
		} else if tsd.Value > Upper {
			tsd.Value = Upper
			tsds[i] = tsd
		}
	}
}

/************************************************************
功能:将时间序列数据中的时间值转换为字符串格式
输入:无
输出:[]TimeSeriesDataS
说明:
时间:2020年8月26日
编辑:wang_jp
************************************************************/
func (tsds Tsds) ToTimeSeriesDataS() []TimeSeriesDataS {
	var tsdss []TimeSeriesDataS
	for _, tsd := range tsds {
		var dss TimeSeriesDataS
		dss.Time = tsd.Time.Format(_TIME_LYOUT)
		dss.Value = tsd.Value
		tsdss = append(tsdss, dss)
	}
	return tsdss
}

/************************************************************
功能:将时间、数值键值对类型的时间序列数据转换为时序数据
输入:tvm map[string]float64:key为时间字符串，value为浮点型数值
输出:error:错误信息
说明:
时间:2020年8月26日
编辑:wang_jp
************************************************************/
func ParseFromTimeValueMap(tvm map[string]float64) (Tsds, error) {
	var datas Tsds
	for t, v := range tvm {
		var data TimeSeriesData
		tm, err := TimeParse(t)
		if err != nil {
			return datas, err
		}
		data.Time = tm
		data.Value = v
		datas = append(datas, data)
	}
	return datas, nil
}

/************************************************************
功能: 移动窗口滤波
输入: n:int,滤波窗口长度,1<n<len(data)
		fillvalue:string, 填充值方法{"mean","median","mode"},省略时为"mean"
输出:
说明: n<=1或者n>len(tsds)时,不进行滤波,直接返回原数据
		"median"时,n必须大于2，否则用"mean"方法
编辑: wangjp
时间: 2020年10月14日
************************************************************/
func (tsds Tsds) MoveWindowFilter(n int, fillvlue ...string) {
	var arr numgo.Array
	arr = tsds.GetDataArray()
	arr = arr.MoveWindowFilter(n, fillvlue...)
	for i, tsd := range tsds {
		tsd.Value = arr[i]
		tsds[i] = tsd
	}
}

/************************************************************
功能: 替换低于给定值的数据
输入: lowlimit:float64:底限值,小于等于该值的数据都将被替换
	  newvalue:float64:替换超限值的新值
输出:
说明:
编辑: wangjp
时间: 2020年10月14日
************************************************************/
func (tsds Tsds) ReplaceLowValue(lowlimit, newvalue float64) {
	for i, tsd := range tsds {
		if tsd.Value <= lowlimit {
			tsd.Value = newvalue
		}
		tsds[i] = tsd
	}
}

/************************************************************
功能: 复制一个相同的数组
输入: 无
输出: Tsds数组
说明:
编辑: wangjp
时间: 2020年10月14日
************************************************************/
func (tsds Tsds) Copy() Tsds {
	var res Tsds
	for _, tsd := range tsds {
		res = append(res, tsd)
	}
	return res
}

/************************************************************
功能: 对数据进行Butterworth滤波
输入: N int: 滤波阶数
	  Wn float/[2]float:自然频率数组，也称归一化的截止频率,与fs具有相同的单位
		fcf=截止频率*2/采样频率；
		如果是低通，高通滤波，fcf只有一个元素,范围:0.0-1.0
		如果是带通，带阻滤波，fcf数组有两个元素,分别是滤波带的起始频率和阶数频率,范围:0.0-fs/2.0
	  btype string:滤波类型，字符串；lp：表示低通；hp：表示高通；bp：表示带通；bs：表示带阻；
	  fs float64:数字系统的采样频率,可选。
		默认情况下,fs是2.0，表示每个周期采样2次
输出: 错误信息
说明:
编辑: wangjp
时间: 2020年10月22日
************************************************************/
func (tsds Tsds) ButterFilter(N int, Wn interface{}, btype string, fs ...float64) error {
	arr := tsds.GetDataArray()
	arr, err := filter.ButterFilter(arr, N, Wn, btype, fs...)
	if err != nil {
		return err
	}
	err = tsds.PutArray2Data(arr)
	return err
}
