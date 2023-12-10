package statistic

/*
***********************************************************

	功能:新建峰谷数据选择器
	输入:
			doorwidth float64:稳态判据值(门宽度)
			minipvd float64:最小峰谷差
			minipeek float64:最小峰值.如果不为0,则峰值必须大于该值
			maxvalley float64:最大谷值.如果不为0,则谷值必须大于该值
			negativeAsZero int64:如果为0,保留负数;如果为1,将负数作为0处理
			peekfists ...int:是否必须先有峰,-1:必须先有谷,0:随机,1:必须先有峰
	输出：无
	时间：2020年2月14日
	编辑：wang_jp

***********************************************************
*/
func (pvs *PeakValleySelector) New(doorwidth, minipvd, minipeek, maxvalley float64, negativeAsZero int, peekfists ...int) {

	pvs.SteadyValue = doorwidth
	pvs.MiniPeek = minipeek
	pvs.MaxValley = maxvalley
	pvs.NegativeAsZero = negativeAsZero
	pvs.MiniPvd = minipvd

	if len(peekfists) > 0 {
		pvs.PeekFirst = peekfists[0]
	}
	pvs.sdt = NewSdtDoor(doorwidth)
}

/*
***********************************************************

	功能:数据筛选(旋转门算法)
	输入:
		input []TimeSeriesData:输入的数据结构
	输出：无
	时间：2020年2月14日
	编辑：wang_jp

***********************************************************
*/
func (pvs *PeakValleySelector) DataFillter(input Tsds) {
	var ds int                    //数据状态值，1=升,0=平,-1=降
	var save bool                 //是否保存
	var pvd PeakValleyPeriodValue //一个周期的峰谷值
	var havepeak, havevalley bool //已经获取到了峰值/谷值
	//fmt.Println("开始旋转门过滤")
	for i, data := range input { //遍历数据
		//去除负数
		if pvs.NegativeAsZero > 0 && data.Value < 0 {
			data.Value = 0
			input[i] = data
		}
		ds, save = pvs.sdt.Filter(data)
		if !save {
			input[i] = pvs.sdt.LastHis
		}
		//fmt.Printf("状态:%d, 保存:%t, 数据:%v\n", ds, save, data)

		if save { //检查的数据已经大于了连续点数
			pvs.processStatus(ds) //数据状态改变
			//fmt.Printf("数据点切换: %v,过程状态:%d,%d\n", data, pvs._processStateChange, pvs.processStateChange)
			switch pvs.processStateChange {
			case -2: //谷成
				if pvs._processStateChange == -2 {
					if pvs.sdt.LastHis2.Value < pvd.Valley.Value {
						pvd.Valley = pvs.sdt.LastHis2
					}
				} else {
					pvd.Valley = pvs.sdt.LastHis2
				}
				if (pvs.MaxValley != 0 && pvd.Valley.Value < pvs.MaxValley) || pvs.MaxValley == 0 {
					if (pvs.PeekFirst < 0 && !havepeak) ||
						((pvs.PeekFirst > 0 && havepeak) && (pvd.Peak.Value-pvd.Valley.Value > pvs.MiniPvd)) ||
						pvs.PeekFirst == 0 {
						havevalley = true
						//fmt.Printf("真谷成:%d,%v\n", pvs.processStateChange, pvd.Valley)
					}
				}
				//fmt.Printf("谷成:%d,%v;MaxValley:%f,PeekFirst:%d,HavePeek:%t,HaveValley:%t\n",
				//	pvs.processStateChange, pvd.Valley, pvs.MaxValley, pvs.PeekFirst, havepeak, havevalley)
			case -1: //降
				pvd.Valley = data
				if pvs._processStateChange == 1 {
					pvd.Peak = pvs.sdt.LastHis
					//fmt.Printf("峰成:%d,%v\n", pvs.processStateChange, pvd.Peak)
					if (pvs.MiniPeek != 0 && pvd.Peak.Value > pvs.MiniPeek) || pvs.MiniPeek == 0 {
						if (pvs.PeekFirst > 0 && !havevalley) ||
							((pvs.PeekFirst < 0 && havevalley) && (pvd.Peak.Value-pvd.Valley.Value > pvs.MiniPvd)) ||
							pvs.PeekFirst == 0 {
							havepeak = true
							//fmt.Printf("真峰成:%d,%v\n", pvs.processStateChange, pvd.Peak)
						}
					}
				}
				//fmt.Printf("降:%d,%v\n", pvs.processStateChange, pvd.Peak)
			//case 0: //平
			//fmt.Printf("平:%d,%v\n", pvs.processStateChange, data)
			case 1: //升
				pvd.Peak = data
				if pvs._processStateChange == -1 {
					pvd.Valley = pvs.sdt.LastHis
					//fmt.Printf("谷成:%d,%v\n", pvs.processStateChange, pvd.Valley)
					if (pvs.MaxValley != 0 && pvd.Valley.Value < pvs.MaxValley) || pvs.MaxValley == 0 {
						if (pvs.PeekFirst < 0 && !havepeak) ||
							((pvs.PeekFirst > 0 && havepeak) && (pvd.Peak.Value-pvd.Valley.Value > pvs.MiniPvd)) ||
							pvs.PeekFirst == 0 {
							havevalley = true
							//fmt.Printf("真谷成:%d,%v\n", pvs.processStateChange, pvd.Valley)
						}
					}
				}
				//fmt.Printf("升:%d,%v\n", pvs.processStateChange, pvd.Valley)
			case 2: //峰成
				if pvs._processStateChange == 2 {
					if pvs.sdt.LastHis2.Value > pvd.Peak.Value {
						pvd.Peak = pvs.sdt.LastHis2
					}
				} else {
					pvd.Peak = pvs.sdt.LastHis2
				}
				if (pvs.MiniPeek != 0 && pvd.Peak.Value > pvs.MiniPeek) || pvs.MiniPeek == 0 {
					if (pvs.PeekFirst > 0 && !havevalley) ||
						((pvs.PeekFirst < 0 && havevalley) && (pvd.Peak.Value-pvd.Valley.Value > pvs.MiniPvd)) ||
						pvs.PeekFirst == 0 {
						havepeak = true
						//fmt.Printf("真峰成:%d,%v\n", pvs.processStateChange, pvd.Peak)
					}
				}
				//fmt.Printf("峰成:%d,%v;MaxValley:%f,PeekFirst:%d,HavePeek:%t,HaveValley:%t\n",
				//	pvs.processStateChange, pvd.Peak, pvs.MiniPeek, pvs.PeekFirst, havepeak, havevalley)
			default:
				//fmt.Printf("未定义状态:%d\n", pvs.processStateChange)
			}
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
	}
	if pvs.PeekFirst > 0 && havepeak { //如果先有峰,且已经有峰,取最后一点为谷
		pvd.Valley = input[len(input)-1]
		if pvs.MaxValley > 0 && pvd.Valley.Value < pvs.MaxValley { //谷值小于最大估值
			pvd.PVDiff = pvd.Peak.Value - pvd.Valley.Value //峰谷值之差
			if pvd.PVDiff > pvs.MiniPvd {                  //峰谷差大于最小峰谷差
				pvs.PvDatas = append(pvs.PvDatas, pvd) //保存峰谷值
				pvs.PeakSum += pvd.Peak.Value          //峰值和
				pvs.ValleySum += pvd.Valley.Value      //谷之和
				pvs.PVDiffSum += pvd.PVDiff            //峰谷差之和
				pvs.PeriodCnt += 1
			}
		}
	}
	if pvs.PeekFirst < 0 && havevalley { //如果先有谷,且已经有谷,取最后一点为峰
		pvd.Peak = input[len(input)-1]
		if pvs.MiniPeek > 0 && pvd.Peak.Value > pvs.MiniPeek { //峰值大于最小峰值
			pvd.PVDiff = pvd.Peak.Value - pvd.Valley.Value //峰谷值之差
			if pvd.PVDiff > pvs.MiniPvd {                  //峰谷差大于最小峰谷差
				pvs.PvDatas = append(pvs.PvDatas, pvd) //保存峰谷值
				pvs.PeakSum += pvd.Peak.Value          //峰值和
				pvs.ValleySum += pvd.Valley.Value      //谷之和
				pvs.PVDiffSum += pvd.PVDiff            //峰谷差之和
				pvs.PeriodCnt += 1
			}
		}
	}
}

/*
***********************************************************

	功能:数据状态改变判断
	输入:
			ds int:当前数据状态;-1:数据在下降,0:数据为平,1:数据在上升
	输出：无
	时间：2020年2月14日
	编辑：wang_jp

***********************************************************
*/
func (pvs *PeakValleySelector) processStatus(ds int) {
	pvs._processStateChange = pvs.processStateChange
	switch pvs.processState {
	case -1:
		switch ds {
		case -1:
			pvs.processStateChange = -1 //降->降
		case 0:
			pvs.processStateChange = -2 //降->平,谷成
		case 1:
			pvs.processStateChange = -2 //降->升,谷成
		}
	case 0:
		switch ds {
		case -1:
			pvs.processStateChange = 2 //平->降,峰成
		case 0:
			pvs.processStateChange = 0 //平->平
		case 1:
			pvs.processStateChange = -2 //平->升,谷成
		}
	case 1:
		switch ds {
		case -1:
			pvs.processStateChange = 2 //升->降,峰成
		case 0:
			pvs.processStateChange = 2 //升->平,峰成
		case 1:
			pvs.processStateChange = 1 //升->升
		}
	}
	pvs.processState = ds //上一个过程状态
}
