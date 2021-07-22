// regression project regression.go
package regression

import (
	"fmt"
	"math"

	"github.com/bkzy/micscript/numgo"
)

//********************************************************************************
//函数功能:回归分析函数
//X是因素矩阵,X矩阵是m行n列的矩阵,
//Y是指标数组,Y有n个元素
//cols:X矩阵的行数
//rows:X矩阵的列数
//********************************************************************************
func Regression(mat_x numgo.Matrix, arr_y numgo.Array) (Regre, error) {
	cols := len(arr_y) //列数
	rows := len(mat_x) //行数

	udf := rows            //回归平方和自由度
	qdf := cols - rows - 1 //残差平方和自由度
	tssdf := cols - 1      //总残差平方和自由度
	var regRes Regre
	var Sx_Matrix, Sy_Matrix numgo.Matrix //Sx,Sy矩阵
	//求Xi的平均值
	Xi_Aver := mat_x.MeanRow()
	//求Y的平均值
	Y_Aver := arr_y.Mean()

	for i, irow := range mat_x { //遍历每一行
		//求S矩阵S=(sij)m*m
		var dsums []float64
		for j, jrow := range mat_x { //遍历每一行
			var dsum float64
			for k, xjk := range jrow { //遍历行中的每一个元素
				dsum += (irow[k] - Xi_Aver[i]) * (xjk - Xi_Aver[j])
			}
			dsums = append(dsums, dsum)
		}
		Sx_Matrix = append(Sx_Matrix, dsums)

		//求Sy矩阵Sy=(s1y,s2y,...,smy)';
		var ysum float64
		var ysums []float64
		for k, xik := range irow {
			ysum += (xik - Xi_Aver[i]) * (arr_y[k] - Y_Aver)
		}
		ysums = append(ysums, ysum)
		Sy_Matrix = append(Sy_Matrix, ysums)
	}
	//Sx_Matrix.Print("Sx矩阵")
	//Sy_Matrix.Print("Sy矩阵")

	inversion, err := Sx_Matrix.Inverse() //求Sx的逆矩阵Sx^-1
	if err != nil {
		return regRes, fmt.Errorf("矩阵%+v不可逆:%s", Sx_Matrix, err.Error())
	}
	product, err := inversion.Mul(Sy_Matrix) //求Sx^-1与Sy的乘积
	if err != nil {
		return regRes, err
	}

	//inversion.Print("Sx的逆矩阵")
	//product.Print("Sx*Sy矩阵")

	//求回归系数 b'= Sx^-1*Sy,b0=Y_Aver-b1*Xi_Aver[1]-...-bm*Xi_Aver[m]
	var cii numgo.Array //S^-1矩阵对角线上的元素
	var regCoeff numgo.Array
	regCoeff = append(regCoeff, Y_Aver)
	for i, _ := range inversion {
		cii = append(cii, inversion[i][i]) //将Sx的逆阵对角线上的元素保存到cii数组中
		regCoeff[0] -= product[i][0] * Xi_Aver[i]
		regCoeff = append(regCoeff, product[i][0])
	}

	var regss, rss float64
	var yest numgo.Array     //样本估计值
	var reldev numgo.Array   //相对偏差
	var residual numgo.Array //残差
	var yscatter numgo.Matrix
	var ymin, ymax float64
	for n, y := range arr_y { //遍历每一列(每一个样本)
		//求Y的估计值
		ye := regCoeff[0]
		for m, _ := range mat_x { //遍历每一行
			ye += regCoeff[m+1] * mat_x[m][n]
		}
		yest = append(yest, ye)

		var scatter []float64
		scatter = append(scatter, y)
		scatter = append(scatter, ye)
		yscatter = append(yscatter, scatter)
		//寻找最大最小值
		if n == 0 {
			if y < ye {
				ymin = y
				ymax = ye
			} else {
				ymin = ye
				ymax = y
			}
		}
		if y > ymax {
			ymax = y
		}
		if ye > ymax {
			ymax = ye
		}
		if y < ymin {
			ymin = y
		}
		if ye < ymin {
			ymin = ye
		}

		//求回归平方和
		diff := ye - Y_Aver
		regss += diff * diff
		dev := 0.0
		if Y_Aver != 0.0 {
			dev = diff / Y_Aver * 100 //相对偏差
		}
		reldev = append(reldev, dev)
		//求残差
		residual = append(residual, y-ye)
		//求残差平方和
		rss += residual[n] * residual[n]
	}

	tss := regss + rss                                                       //求总残差平方和
	vr := math.Sqrt(regss / tss)                                             //求复相关系数
	vf := (regss / float64(udf)) / (rss / float64(qdf))                      //求F值
	sd := math.Sqrt(rss / float64(qdf))                                      //计算标准偏差
	vfa := TableLookupF(udf, qdf, 0.05)                                      //查表求检验水平为Alpha的F临界值
	vra := math.Sqrt(float64(udf) * vfa / (float64(qdf) + float64(udf)*vfa)) //求复相关系数R的临界值

	//求标准残差
	var stdres numgo.Array
	for _, res := range residual {
		stdres = append(stdres, res/sd)
	}

	//求各系数的T值和偏回归平方和
	var ts numgo.Array
	var vs numgo.Array
	for i, ci := range cii {
		t := (regCoeff[i+1] / math.Sqrt(ci)) / math.Sqrt(rss/float64(qdf)) //求T
		v := (regCoeff[i+1] * regCoeff[i+1]) / ci                          //求偏回归平方和V
		if t < 0 {
			t *= -1
		}
		ts = append(ts, t)
		vs = append(vs, v)
	}

	talpha := TableLookupT(qdf, 0.05) //查表求检验水平为Alpha的T临界值

	regRes.Coeff = regCoeff
	regRes.Xs = mat_x
	regRes.Ys = arr_y
	regRes.YEst = yest
	regRes.U = regss
	regRes.Q = rss
	regRes.TSS = tss
	regRes.Ta = talpha
	regRes.Ts = ts
	regRes.Vs = vs
	regRes.StdRes = stdres
	regRes.R = vr
	regRes.F = vf
	regRes.Ra = vra
	regRes.Fa = vfa
	regRes.SD = sd
	regRes.TssDf = tssdf
	regRes.Qdf = qdf
	regRes.Udf = udf
	regRes.QdQdf = rss / float64(qdf)
	regRes.UdUdf = regss / float64(udf)
	regRes.Residual = residual
	regRes.RelDev = reldev
	regRes.Ymin = ymin
	regRes.Ymax = ymax
	regRes.Yscatter = yscatter
	return regRes, nil
}

//**************************************************************
//函数功能:查F分布表
//参数:
//		v1:分子自由度v1,横坐标
//      v2:分母自由度v2,纵坐标
//      Alpha:检验水平,可取0.01和0.05
//返回值:从表中查出的F临界值
//**************************************************************
func TableLookupF(v1, v2 int, Alpha float64) float64 { //查F分布表

	var i, j int

	if v1 < 1 {
		//TRACE("查F分布表函数'Table_Lookup_F'的参数'v1'不合法:v1=%d",v1)
		return -1
	}
	if v2 < 1 {
		//TRACE("查F分布表函数'Table_Lookup_F'的参数'v2'不合法:v2=%d",v2)
		return -1
	}
	if Alpha != 0.01 && Alpha != 0.05 {
		//TRACE("查F分布表函数'Table_Lookup_F'的参数'Alpha'不合法:Alpha=%f",Alpha)
		return -1
	}

	if v1 < 11 {
		j = v1 - 1
	} else if v1 >= 11 && v1 < 14 {
		j = 10
	} else if v1 >= 14 && v1 < 17 {
		j = 11
	} else if v1 >= 17 && v1 < 22 {
		j = 12
	} else if v1 >= 22 && v1 < 27 {
		j = 13
	} else if v1 >= 27 && v1 < 35 {
		j = 14
	} else if v1 >= 35 && v1 < 50 {
		j = 15
	} else if v1 >= 50 && v1 < 90 {
		j = 16
	} else if v1 >= 90 && v1 <= 120 {
		j = 17
	} else {
		j = 18
	}
	if v2 < 31 {
		i = v2 - 1
	} else if v2 >= 31 && v2 < 50 {
		i = 30
	} else if v2 >= 50 && v2 < 90 {
		i = 31
	} else if v2 >= 90 && v2 <= 120 {
		i = 32
	} else {
		i = 33
	}
	if Alpha == 0.05 {
		return _F95_TABLE[i][j]
	} else {
		return _F99_TABLE[i][j]
	}
}

//**************************************************************
//函数功能:查T分布表
//参数:
//		v:自由度
//      Alpha:检验水平,可取0.005,0.01,0.025,0.05,0.1,0.2,0.25,0.3,0.4,0.45
//返回值:从表中查出的T临界值
//**************************************************************
func TableLookupT(v int, Alpha float64) float64 { //查T分布表
	var i, j int

	if v < 31 {
		i = v - 1
	} else if v >= 31 && v < 50 {
		i = 30
	} else if v >= 50 && v < 90 {
		i = 31
	} else if v >= 90 && v <= 120 {
		i = 32
	} else {
		i = 33
	}
	switch int(Alpha * 1000) {
	case 450:
		j = 0
		break
	case 400:
		j = 1
		break
	case 300:
		j = 2
		break
	case 250:
		j = 3
		break
	case 200:
		j = 4
		break
	case 100:
		j = 5
		break
	case 50:
		j = 6
		break
	case 25:
		j = 7
		break
	case 10:
		j = 8
		break
	case 5:
		j = 9
		break
	default:
		{
			//TRACE("查F分布表函数'Table_Lookup_T'的参数'Alpha'不合法:Alpha=%f",Alpha)
			return -1
			break
		}
	}
	return _T_TABLE[i][j]
}
