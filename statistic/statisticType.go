package statistic

import (
	"time"
)

const (
	_TIME_LYOUT = "2006-01-02 15:04:05"
)

type TimeSeriesDataS struct { //时间序列数据结构
	Time  string  //时间
	Value float64 //数值
}
type TimeSeriesDataI struct { //时间序列数据结构
	Time  int64   //时间,Unix毫秒数
	Value float64 //数值
}

type TimeSeriesData struct { //时间序列数据结构
	Time  time.Time //时间
	Value float64   //数值
}

// 时间序列数据数组
type Tsds []TimeSeriesData

type StatisticData struct { //统计数据结构
	Min         float64            //最小值(基本)
	Max         float64            //最大值(基本)
	Range       float64            //数据范围(Max-Min)(基本)
	Total       float64            //表示统计时间段内的累计值，结果的单位为标签点的工程单位(面积,值*时间(s))(基本)
	Sum         float64            //统计时间段内的算术累积值(值相加)(基本)
	Mean        float64            //统计时间段内的算术平均值(Mean = Sum/PointCnt)(基本)
	PowerAvg    float64            //统计时间段内的加权平均值,对BOOL量而言是ON率（Total/Duration）(基本)
	Diff        float64            //统计时间段内的差值(最后一个值减去第一个值)(基本)
	PlusDiff    float64            //正差值,用于累计值求差,可以削除清零对值的影响(统计周期内只可以有一次清零动作)
	Duration    int64              //统计时间段内的秒数(EndTime - BeginTime)(基本)
	PointCnt    int                //统计时间段内的数据点数(基本)
	RisingCnt   int                //统计时间段内数据上升的次数(基本)
	FallingCnt  int                //统计时间段内数据下降的次数(基本)
	LtzCnt      int                //小于0的次数
	GtzCnt      int                //大于0的次数
	EzCnt       int                //等于0的次数
	OutliersCnt int                //离群异常点数(小于Lower和大于Upper的点数)
	Lower       float64            //箱型图(盒须图)下边界,lower: Q1 - 1.5 * Qd
	Q1          float64            //下四分位
	Q3          float64            //上四分位
	Qd          float64            //四分位差,亦写为IQR
	Upper       float64            //箱型图(盒须图)上边界,upper: Q3 + 1.5 * Qd
	BeginTime   string             //开始时间(基本)
	EndTime     string             //结束时间(基本)
	SD          float64            //总体标准偏差(高级)
	STDDEV      float64            //样本标准偏差(高级)
	SE          float64            //标准误差(SE = STDDEV / PointCnt)(高级)
	Ske         float64            //偏度(高级)
	Kur         float64            //峰度(高级)
	Mode        float64            //众数(高级)
	Median      float64            //中位数(高级)
	GroupDist   float64            //组距GroupDistance(高级),DataGroup中两组数之间的距离
	DataGroup   map[string]float64 //数据分布组(高级),绘制出图形后是近似正态分布的草帽图
	Increment   map[string]float64 //相邻量数之间的增量,string是时间(基本)
	RawData     []TimeSeriesDataS  //原始数据
}

type PeakValleySelector struct {
	InflectionIncrement float64 //拐点增量
	SteadyValue         float64 //稳态判据值
	ContinuePoint       int     //连续稳定的数据点数
	NegativeAsZero      int     //如果为0,保留负数;如果为1,将负数作为0处理
	MaxValley           float64 //最大谷值:谷值必须小于该值
	MiniPeek            float64 //最小峰值:峰值必须大于该值
	PeekFirst           bool    //先有峰后有谷

	PeakSum            float64                 //峰之和
	ValleySum          float64                 //谷之和
	PVDiffSum          float64                 //峰谷差之和
	PeriodCnt          int                     //周期数
	PvDatas            []PeakValleyPeriodValue //峰谷周期值
	dataState          []int                   //数据状态,数组长度由ContinuePoint决定.1=升,0=平,-1=降
	processState       int                     //过程状态,1=升,0=平,-1=降
	processStateChange int                     //过程状态发生改变的值,1=升->平(峰成),2=升->降(峰成),3=降->平(谷成),4=降->升(谷成),5=降->升(谷成,由单个增量判断),6=升->降(峰成,由单个增量判断),0=其他
}

type PeakValleyPeriodValue struct { //峰谷周期值
	Peak   TimeSeriesData //峰值数据
	Valley TimeSeriesData //谷值数据
	PVDiff float64        //峰谷值之差
}
