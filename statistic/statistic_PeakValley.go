package statistic

import (
	"math"
)

/*
***********************************************************

	功能:新建峰谷数据选择器
	输入:
			stdv float64:稳态判据值
			inflectionIncrement float64:拐点增量
			minipeek float64:最小峰值.如果不为0,则峰值必须大于该值
			maxvalley float64:最大谷值.如果不为0,则谷值必须大于该值
			cp int:连续稳定的数据点数
			negativeAsZero int64:如果为0,保留负数;如果为1,将负数作为0处理
			peekfists ...bool:是否必须先有峰
	输出：无
	时间：2020年2月14日
	编辑：wang_jp

***********************************************************
*/
func (pvs *PeakValleySelector) New(stdv, inflectionIncrement, minipeek, maxvalley float64, cp, negativeAsZero int, peekfists ...int) {
	if cp < 2 { //连续点不能小于2
		cp = 2
	}
	pvs.ContinuePoint = cp
	pvs.SteadyValue = stdv
	pvs.MiniPeek = minipeek
	pvs.MaxValley = maxvalley
	pvs.NegativeAsZero = negativeAsZero
	pvs.InflectionIncrement = inflectionIncrement

	if len(peekfists) > 0 {
		pvs.PeekFirst = peekfists[0]
	}

	pvs.dataState = make([]int, cp) //初始化数据状态空间长度
}

/*
***********************************************************

	功能:数据筛选
	输入:
		input []TimeSeriesData:输入的数据结构
	输出：无
	时间：2020年2月14日
	编辑：wang_jp

***********************************************************
*/
func (pvs *PeakValleySelector) DataFillter(input Tsds) {
	var v0 float64                //上一个循环的值
	var vt float64                //当前循环的值
	var dv float64                //增量值
	var ds int                    //数据状态值，1=升,0=平,-1=降
	var pvd PeakValleyPeriodValue //一个周期的峰谷值
	var havepeak, havevalley bool //已经获取到了峰值/谷值

	for i, data := range input { //遍历数据
		if pvs.NegativeAsZero > 0 && data.Value < 0 {
			data.Value = 0
			input[i].Value = 0
		}
		if i == 0 { //第一个值
			v0 = data.Value
			continue //直接循环
		}
		vt = data.Value //获取当前值
		dv = vt - v0    //增量值

		if math.Abs(dv) <= pvs.SteadyValue { //增量绝对值小于稳态判据
			ds = 0 //数据稳定状态
		} else {
			if dv > 0 { //增量大于0
				ds = 1 //数据上升状态
			} else { //增量小于0
				ds = -1 //数据下降状态
			}
		}
		pvs.dataStateChange(dv, ds) //数据状态改变
		if i >= pvs.ContinuePoint { //检查的数据已经大于了连续点数
			switch pvs.processStateChange {
			case 1: //升->平(峰成),取平尾为峰
				pvd.Peak = input[i-(pvs.ContinuePoint-1)]
				if (pvs.MiniPeek != 0 && pvd.Peak.Value > pvs.MiniPeek) || pvs.MiniPeek == 0 {
					if (pvs.PeekFirst < 0 && havevalley) || pvs.PeekFirst == 0 {
						havepeak = true
					}
				}
				//fmt.Printf("升->平(峰成):%d,%+v,%s\n", pvs.processStateChange, pvs.dataState, pvd.Peak.Time)
			case 2: //升->降(峰成),取尾去一为峰
				pvd.Peak = input[i-pvs.ContinuePoint]
				if (pvs.MiniPeek != 0 && pvd.Peak.Value > pvs.MiniPeek) || pvs.MiniPeek == 0 {
					if (pvs.PeekFirst < 0 && havevalley) || pvs.PeekFirst == 0 {
						havepeak = true
					}
				}
				//fmt.Printf("升->降(峰成):%d,%+v,%s\n", pvs.processStateChange, pvs.dataState, pvd.Peak.Time)
			case 3: //降->平(谷成),取平尾为谷
				pvd.Valley = input[i-(pvs.ContinuePoint-1)]
				if (pvs.MaxValley != 0 && pvd.Valley.Value < pvs.MaxValley) || pvs.MaxValley == 0 {
					if (pvs.PeekFirst > 0 && havepeak) || pvs.PeekFirst == 0 {
						havevalley = true
					}
				}
				//fmt.Printf("降->平(谷成):%d,%+v,%s\n", pvs.processStateChange, pvs.dataState, pvd.Valley.Time)
			case 4: //降->升(谷成),取尾去一为谷
				pvd.Valley = input[i-pvs.ContinuePoint]
				if (pvs.MaxValley != 0 && pvd.Valley.Value < pvs.MaxValley) || pvs.MaxValley == 0 {
					if (pvs.PeekFirst > 0 && havepeak) || pvs.PeekFirst == 0 {
						havevalley = true
					}
				}
				//fmt.Printf("降->升(谷成):%d,%+v,%s\n", pvs.processStateChange, pvs.dataState, pvd.Valley.Time)
			case 5: //降->升(谷成),去一为谷
				pvd.Valley = input[i-1]
				if (pvs.MaxValley != 0 && pvd.Valley.Value < pvs.MaxValley) || pvs.MaxValley == 0 {
					if (pvs.PeekFirst > 0 && havepeak) || pvs.PeekFirst == 0 {
						havevalley = true
					}
				}
				//fmt.Printf("降->升(谷成):%d,%+v,%s\n", pvs.processStateChange, pvs.dataState, pvd.Valley.Time)
			case 6: //升->降(峰成),去一为峰
				pvd.Peak = input[i-1]
				if (pvs.MiniPeek != 0 && pvd.Peak.Value > pvs.MiniPeek) || pvs.MiniPeek == 0 {
					if (pvs.PeekFirst < 0 && havevalley) || pvs.PeekFirst == 0 {
						havepeak = true
					}
				}
				//fmt.Printf("升->降(峰成):%d,%+v,%s\n", pvs.processStateChange, pvs.dataState, pvd.Peak.Time)
			}
			pvs.processStateChange = 0
			if havepeak && havevalley { //已经获取了峰值和谷值
				pvd.PVDiff = pvd.Peak.Value - pvd.Valley.Value //峰谷值之差
				havepeak = false                               //复位
				havevalley = false                             //复位
				pvs.PvDatas = append(pvs.PvDatas, pvd)         //保存峰谷值
				pvs.PeakSum += pvd.Peak.Value                  //峰值和
				pvs.ValleySum += pvd.Valley.Value              //谷之和
				pvs.PVDiffSum += pvd.PVDiff                    //峰谷差之和
				pvs.PeriodCnt += 1
			}
		}

		v0 = vt //当前值保存为上一个循环值
	}
}

/*
***********************************************************

	功能:数据状态改变判断
	输入:
			increment float64:增量值
			ds int:当前数据状态
	输出：无
	时间：2020年2月14日
	编辑：wang_jp

***********************************************************
*/
func (pvs *PeakValleySelector) dataStateChange(increment float64, ds int) {
	replacer := 10                               //用于代替-1，防止-1与1相加等于0
	var dssum int                                //数据状态和
	for i := pvs.ContinuePoint - 1; i > 0; i-- { //数据状态列队后移并求数据状态和(不含0位)
		pvs.dataState[i] = pvs.dataState[i-1]
		if pvs.dataState[i] == -1 {
			dssum += replacer
		} else {
			dssum += pvs.dataState[i]
		}
	}
	pvs.dataState[0] = ds       //保存最新的数据状态
	if pvs.dataState[0] == -1 { //0位的数据状态和
		dssum += replacer
	} else {
		dssum += pvs.dataState[0]
	}

	var ps int //过程状态
	//用大增量判断过程状态
	if math.Abs(increment) > pvs.InflectionIncrement && pvs.InflectionIncrement > 0 { //拐点增量不为0且增量绝对值大于拐点增量
		if increment > 0 { //增量大于0
			ps = 1                    //升态
			switch pvs.processState { //先前的过程状态
			case -1:
				pvs.processStateChange = 5 //降->升(谷成)
			default:
				pvs.processStateChange = 0
			}
			pvs.processState = ps //保存过程状态
		} else { //增量小于0
			ps = -1                   //降态
			switch pvs.processState { //先前的过程状态
			case 1:
				pvs.processStateChange = 6 //升->降(峰成)
			default:
				pvs.processStateChange = 0
			}
			pvs.processState = ps //保存过程状态
		}
	} else {
		switch dssum { //用数据状态判断过程状态
		case 0: //全部为0
			ps = 0                    //稳态
			switch pvs.processState { //先前的过程状态
			case 1:
				pvs.processStateChange = 1 //升->平(峰成)
			case -1:
				pvs.processStateChange = 3 //降->平(谷成)
			default:
				pvs.processStateChange = 0
			}
			pvs.processState = ps //保存过程状态
		case replacer * pvs.ContinuePoint: //全部为-1
			ps = -1                   //降态
			switch pvs.processState { //先前的过程状态
			case 1:
				pvs.processStateChange = 2 //升->降(峰成)
			default:
				pvs.processStateChange = 0
			}
			pvs.processState = ps //保存过程状态
		case pvs.ContinuePoint: //全部为1
			ps = 1                    //升态
			switch pvs.processState { //先前的过程状态
			case -1:
				pvs.processStateChange = 4 //降->升(谷成)
			default:
				pvs.processStateChange = 0
			}
			pvs.processState = ps //保存过程状态
		}
	}
}
