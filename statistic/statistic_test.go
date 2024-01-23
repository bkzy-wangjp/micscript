package statistic

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestBaseStatistics(t *testing.T) {
	tsdstr, _ := GetDataFromJSON("testdata/datas.json")
	tsdt, _ := Tsds2Tsdt(tsdstr)
	//tsds, err := tsdt.Extract("2024-01-15 06:00:00", "2024-01-15 07:00:00")

	tests := []struct {
		tsds    []TimeSeriesData
		advance int
		group   int
	}{
		{tsdt, 1, 0},
	}
	for _, tt := range tests {
		data := Tsds(tt.tsds)
		res := data.Statistics(tt.advance, tt.group)
		res.Increment = nil
		StructFormatPrint(res)
	}
}

func TestBaseStatistics2(t *testing.T) {
	tsdsi, err := GetItDataFromJSON("testdata/gddata2.json")
	if err != nil {
		t.Log(err.Error())
	}
	tsdt := Tsdi2Tsdt(tsdsi)
	//tsds, err := tsdt.Extract("2024-01-15 06:00:00", "2024-01-15 07:00:00")
	//StructFormatPrint(tsdt)
	tests := []struct {
		tsds    []TimeSeriesData
		advance int
		group   int
	}{
		{tsdt, 1, 0},
	}
	for _, tt := range tests {
		data := Tsds(tt.tsds)
		res := data.Statistics(tt.advance, tt.group)
		res.Increment = nil
		StructFormatPrint(res)
	}
}

func TestPeakValleySelector(t *testing.T) {
	tsds, _ := GetDataFromJSON("testdata/pvdata.json")
	tsdt, _ := Tsds2Tsdt(tsds)
	tests := []struct {
		tsdt   Tsds    //数据
		stv    float64 //稳态判据值
		inin   float64 //拐点增量
		cp     int     //连续稳定的数据点数
		nz     int     //如果为0,保留负数;如果为1,将负数作为0处理
		butter bool    //是否启用滤波
	}{
		{tsdt, 0.03, 2.0, 3, 1, true},
	}

	for _, tt := range tests {
		var pvs PeakValleySelector
		if tt.butter {
			err := tt.tsdt.ButterFilter(1, 0.125, "lp")
			if err != nil {
				t.Error(err.Error())
			}
		}
		tt.tsdt.ReplaceLowValue(0.5, 0.0)
		pvs.New(tt.stv, tt.inin, 0, 0, tt.cp, tt.nz)
		pvs.DataFillter(tt.tsdt)
		fmt.Printf("峰值和:%f,谷之和:%f,峰谷差之和:%f,周期数:%d\n", pvs.PeakSum, pvs.ValleySum, pvs.PVDiffSum, pvs.PeriodCnt)
		for _, tsd := range pvs.PvDatas {
			StructFormatPrint(tsd)
		}
	}
}

func TestSdt(t *testing.T) {
	var tsdts Tsds
	now := time.Now()
	for i := 0; i < 360; i++ {
		var tsd TimeSeriesData
		tsd.Time = now.Add(time.Duration(i) * time.Second)
		tsd.Value = math.Sin(math.Pi / 180 * float64(i))
		tsdts = append(tsdts, tsd)
		fmt.Printf("%+v\n", tsd)
	}

	tsdt := tsdts.SdtFillter(0.0001)
	fmt.Print("--------------------------------\n")
	for _, ts := range tsdt {
		fmt.Printf("%+v\n", ts)
	}
}
