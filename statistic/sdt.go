package statistic

//旋转门压缩算法

//旋转门压缩结构体
type SdtDoor struct {
	DoorWidth      float64 //门初始宽度
	LastHis        TimeSeriesData
	LastHis2       TimeSeriesData
	LastReal       TimeSeriesData
	MaxIntervalSec int64   //数据存储最大间隔秒数
	isinit         bool    //初始化状态
	closed         bool    //初始关门状态
	uk             float64 //上斜率
	dk             float64 //下斜率
	ub             float64 //上零点
	db             float64 //下零点
	timeInterval   float64 //时间间隔,毫秒
}

/*
***********************************************************

	功能: 新建旋转门压缩实例
	输入:
		door_width:float64:旋转门的宽度
		maxinterval ...int64:最大存储间隔
	输出:
	说明:
	编辑: wangjp
	时间: 2023年12月10日

***********************************************************
*/
func NewSdtDoor(door_width float64, maxinterval ...int64) *SdtDoor {
	var interval int64 = 0
	if len(maxinterval) > 0 {
		interval = maxinterval[0]
	}
	sdt := &SdtDoor{DoorWidth: door_width, MaxIntervalSec: interval, closed: true, isinit: true}
	return sdt
}

/*
***********************************************************

	功能: 旋转门压缩过滤器
	输入:
		dpoint TimeSeriesData:时间序列数据点
	输出:
		status int:数据的趋势状态,
			0:上门斜率和下门斜率不一致;
			1:上门斜率和下门斜率均大于0;
			-1:上门斜率和下门斜率均小于0.
		savle bool:当前输入点是否是可保存的数据点
	说明:
	编辑: wangjp
	时间: 2023年12月10日

***********************************************************
*/
func (sdt *SdtDoor) Filter(dpoint TimeSeriesData) (status int, save bool) {
	if sdt.DoorWidth == 0 { //门宽度为零
		return //保存每一个数据
	}

	save = false    //过滤检查结果
	if sdt.isinit { //初始化状态
		sdt.LastReal = dpoint
		sdt.LastHis = dpoint
		sdt.isinit = false
		save = true
	} else {
		if sdt.timeInterval == 0 {
			interval := dpoint.Time.Sub(sdt.LastHis.Time).Milliseconds()
			if interval > 10000 {
				sdt.timeInterval = 10000.0
			} else if interval > 1000 {
				sdt.timeInterval = 1000.0
			} else if interval > 100 {
				sdt.timeInterval = 100.0
			} else {
				sdt.timeInterval = 1.0
			}
		}
		deltaT := float64(dpoint.Time.Sub(sdt.LastHis.Time).Milliseconds()) / sdt.timeInterval
		//fmt.Printf("时间差:%fms\n", deltaT)
		if sdt.closed { //开门第一个点
			if deltaT > 0 {
				sdt.closed = false
				sdt.uk = (dpoint.Value - sdt.ub) / deltaT
				sdt.dk = (dpoint.Value - sdt.db) / deltaT
			}
		} else {
			uk := (dpoint.Value - sdt.ub) / deltaT
			dk := (dpoint.Value - sdt.db) / deltaT
			if uk > sdt.uk { //上斜率只保存增大的
				sdt.uk = uk
			}
			if dk < sdt.dk { //下斜率只保存减小的
				sdt.dk = dk
			}
			if sdt.dk <= sdt.uk { //下斜率小于等于上斜率,触发保存
				save = true
			}
			//fmt.Printf("时间差:%f,上斜率:%f,下斜率:%f\n", deltaT, uk, dk)
			if !save && sdt.MaxIntervalSec > 0 {
				if deltaT/1000.0 > float64(sdt.MaxIntervalSec) { //已经长时间没有触发保存
					save = true
				}
			}
		}
	}
	if save {
		sdt.LastHis2 = sdt.LastHis
		sdt.LastHis = dpoint
		sdt.ub = sdt.LastHis.Value + sdt.DoorWidth
		sdt.db = sdt.LastHis.Value - sdt.DoorWidth
		sdt.closed = true
		//fmt.Println("保存数据:", sdt.LastHis.Time, sdt.LastHis.Value)
	}

	sdt.LastReal = dpoint

	if sdt.uk < 0 && sdt.dk < 0 {
		status = -1
	}
	if sdt.uk > 0 && sdt.dk > 0 {
		status = 1
	}

	return
}
