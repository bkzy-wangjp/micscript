// numgo project numgo.go
package numgo

import (
	"fmt"
	"math"
)

type ComplexMatrix []ComplexArray //复数矩阵

/*********************************************
//功能:矩阵初等行变换（Elementary Transformation）-交换两行
//参数:两行的行号
//说明:
//作者:wangjp
//时间:2020年3月24日
**********************************************/
func (m ComplexMatrix) ElemTransRowSwap(r1, r2 int) {
	rows := len(m)
	if r1 < rows && r2 < rows && r1 != r2 {
		tmp := m[r1]
		m[r1] = m[r2]
		m[r2] = tmp
	}
}

/*********************************************
//功能:矩阵初等行变换（Elementary Transformation）-某行乘以常数
//参数: 行号r
		需要乘的参数k
//说明:
//作者:wangjp
//时间:2020年3月24日
**********************************************/
func (m ComplexMatrix) ElemTransRowMulK(r int, k complex128) {
	rows := len(m)
	if r < rows {
		for col, cval := range m[r] {
			m[r][col] = k * cval
		}
	}
}

/*********************************************
//功能:矩阵初等行变换（Elementary Transformation）-数k乘第 r2 行加到第 r1 行
//参数: 行号r1和r2
		需要乘的参数k
//说明:
//作者:wangjp
//时间:2020年3月24日
**********************************************/
func (m ComplexMatrix) ElemTransRowMulKAddToRow(r1, r2 int, k complex128) {
	rows := len(m)
	//fmt.Println("-------------------\nr1:", r1, "r2:", r2, "K:", k)
	//m.Print("原始矩阵")
	if r1 < rows && r2 < rows && r1 != r2 {
		for col, cval := range m[r2] {
			m[r1][col] += k * cval
			if math.Abs(real(m[r1][col])) < 1e-10 { //极小值设置为0
				m[r1][col] = complex(0, imag(m[r1][col]))
			}
			if math.Abs(imag(m[r1][col])) < 1e-10 { //极小值设置为0
				m[r1][col] = complex(real(m[r1][col]), 0)
			}
		}
	}
	//m.Print("行变换后的矩阵:")
}

/*********************************************
//功能: 初始化一个新矩阵
//参数: m为行数,n为列数
//返回: m x n的矩阵,矩阵的元素全为0
//说明:
//作者: wangjp
//时间: 2020年10月14日
**********************************************/
func InitComplexMatrix(m, n int) ComplexMatrix {
	var mat ComplexMatrix
	for i := 0; i < m; i++ {
		row := make([]complex128, n)
		mat = append(mat, row)
	}
	return mat
}

/*********************************************
//功能:初始化输出n阶单位矩阵
//参数: 单位矩阵的阶数n
//输出: n x n的单位矩阵
//说明:
//作者:wangjp
//时间:2020年3月24日
**********************************************/
func InitComplexMatrixAsIdentity(n int) ComplexMatrix {
	var mat ComplexMatrix
	for i := 0; i < n; i++ {
		row := make([]complex128, n)
		for k := range row {
			if i == k {
				row[k] = 1 + 0i
			}
		}
		mat = append(mat, row)
	}
	return mat
}

/*********************************************
//功能:将矩阵逐行打印输出,可选添加标题
//参数: 可选的标题
//说明:
//作者:wangjp
//时间:2020年3月24日
**********************************************/
func (m ComplexMatrix) Print(title ...string) {
	if len(title) > 0 {
		fmt.Println(title[0])
	}
	for _, row := range m {
		fmt.Println(row)
	}
}

/*********************************************
//功能:对矩阵的每行求和
//参数:无
//说明:
//作者:wangjp
//时间:2020年3月24日
**********************************************/
func (m ComplexMatrix) SumRow() []complex128 {
	var sums []complex128
	for _, mr := range m {
		var sum complex128
		for _, mv := range mr {
			sum += mv
		}
		sums = append(sums, sum)
	}
	return sums
}

/*********************************************
//功能:对矩阵的每行求平均
//参数:无
//说明:
//作者:wangjp
//时间:2020年3月24日
**********************************************/
func (m ComplexMatrix) MeanRow() []complex128 {
	var mean []complex128
	for _, mr := range m {
		var avg complex128
		if len(mr) > 0 {
			var sum complex128
			for _, mv := range mr {
				sum += mv
			}
			avg = sum / Real2Complex(float64(len(mr)))
		}
		mean = append(mean, avg)
	}
	return mean
}

/*********************************************
//功能:复制一个矩阵
//参数:
//输出:新的矩阵
//说明:
//作者:wangjp
//时间:2020年10月19日
**********************************************/
func (mat ComplexMatrix) Copy() ComplexMatrix {
	m, n := mat.RowsAndCols()
	newmat := InitComplexMatrix(m, n)
	for i, row := range mat {
		for j, v := range row {
			newmat[i][j] = v
		}
	}
	return newmat
}

/*********************************************
//功能:对矩阵进行初等行变换，对左侧方阵部分求上三角矩阵
//参数:confg[0]:true:替换原始矩阵,false:不替换原始矩阵
//	   confg[1]:true:矩阵不可逆时返回错误,false:矩阵不可逆时继续进行
//输出:上三角方阵
//说明:矩阵的行数不可大于列数
//	   做行变换将左侧方阵部分转换为上三角阵
//     主对角线以下都是零的方阵称为上三角矩阵。
//     上三角矩阵具有行列式为对角线元素相乘、上三角矩阵乘以系数后也是上三角矩阵、
//     上三角矩阵间的加减法和乘法运算的结果仍是上三角矩阵等性质。
//作者:wangjp
//时间:2020年10月19日
**********************************************/
func (m ComplexMatrix) LeftSquare2Triu(confg ...bool) (ComplexMatrix, error) {
	n, c := m.RowsAndCols() //矩阵的行数
	if c < n {
		return nil, fmt.Errorf("the number of rows of a matrix cannot be greater than the number of columns.[矩阵的行数不可大于列数]")
	}
	mat := m.Copy()
	return_err := false
	if len(confg) > 0 {
		if confg[0] {
			mat = m
		}
		if len(confg) > 1 {
			return_err = confg[1]
		}
	}
	for i, ri := range mat { //逐行进行初等行变换
		if ri[i] == 0 && i < n-1 { //如果第i行i列为0
			ok := false
			for j := i + 1; j < n; j++ { //循环遍历下三角元素,找到一个不为零的i列元素
				if mat[j][i] != 0 { //找一个i列元素非0行的行
					mat.ElemTransRowSwap(i, j) //将该j行与i行交换
					ok = true
					break //跳出j循环
				}
			}
			if !ok && return_err {
				return nil, fmt.Errorf("the matrix is not invertible and can not be converted to upper triangular matrix.[矩阵不可逆,不能转换为上三角阵]")
			}
		}
		for j := i + 1; j < n; j++ {
			if mat[j][i] != 0 { //i行i列下方的j行i列元素不为0
				k := -1 * mat[j][i] / mat[i][i] //准备系数,用以将j行i列变为0
				mat.ElemTransRowMulKAddToRow(j, i, k)
			}
		}
	}
	return mat, nil
}

/*********************************************
//功能:进行初等行变换将矩阵左侧方阵部分转换为单位矩阵
//参数:confg[0]:true:替换原始矩阵,false:不替换原始矩阵
//	   confg[1]:true:矩阵不可逆时返回错误,false:矩阵不可逆时继续进行
//输出:左侧方阵转换为单位方阵的矩阵
//说明:矩阵的行数不可大于列数
//作者:wangjp
//时间:2020年10月19日
**********************************************/
func (m ComplexMatrix) LeftSquare2Identity(confg ...bool) (ComplexMatrix, error) {
	n := m.Rows() //矩阵的行数
	mat := m.Copy()
	return_err := false
	if len(confg) > 0 {
		if confg[0] {
			mat = m
		}
		if len(confg) > 1 {
			return_err = confg[1]
		}
	}
	_, err := mat.LeftSquare2Triu(true, return_err)
	if err != nil {
		return nil, err
	}
	//化左侧成单位矩阵
	for i := n - 1; i >= 0; i-- {
		mat.ElemTransRowMulK(i, 1/mat[i][i]) //该行每个元素除以对角线上的元素，将对角线转换为1
		for j := i - 1; j >= 0; j-- {        //对角线上部逐行遍历
			for r := i; r < 2*n; r++ {
				k := -mat[j][i]
				mat.ElemTransRowMulKAddToRow(j, i, k)
			}
		}
	}
	return mat, nil
}

/*********************************************
//功能:求矩阵的逆矩阵
//参数:
//输出:矩阵的逆矩阵,状态(求逆矩阵成功true,不成功false)
//说明:矩阵必须为方阵
//作者:wangjp
//时间:2020年3月24日
**********************************************/
func (m ComplexMatrix) Inverse() (ComplexMatrix, error) {
	n := len(m) //矩阵的行数
	cCnt := 0
	if n > 0 {
		cCnt = len(m[0])            //矩阵的列数
		if cCnt == 0 || n != cCnt { //列数等于0,或者行列不相等
			return nil, fmt.Errorf("the matrix must be a square matrix.[矩阵必须是方阵]")
		}
	} else { //行数等于0
		return nil, fmt.Errorf("the number of rows of a matrix must be greater than zero.[矩阵的行数必须大于零]")
	}

	//建立n*2n的矩阵,左半部分为输入矩阵,右半部分为单位矩阵
	//单位矩阵是方阵,它从左上角到右下角的对角线（称为主对角线）上的元素均为1。除此以外全都为0
	var mat ComplexMatrix
	for i, mr := range m {
		row := make([]complex128, n*2)
		for j, mv := range mr {
			row[j] = mv
		}
		row[n+i] = 1
		mat = append(mat, row)
	}

	var err error
	_, err = mat.LeftSquare2Identity(true, true) //左侧化为单位矩阵
	if err != nil {
		return nil, err
	}

	var rp ComplexMatrix
	for _, m := range mat {
		rp = append(rp, m[n:])
	}
	return rp, nil
}

/*********************************************
//功能:求矩阵乘积
//参数:矩阵A和矩阵B
//输出:矩阵的乘积,错误信息
//说明:第一个矩阵a是m×s型的，第二个矩阵b是s×n型的,结果阵Result是m×n型的,如果a矩阵的列与b矩阵
//    的行不相等，则输出错误信息
//作者:wangjp
//时间:2020年3月24日
**********************************************/
func (ma ComplexMatrix) Mul(mb ComplexMatrix) (ComplexMatrix, error) { //求矩阵乘积
	arow := len(ma) //a矩阵行数
	acol := 0       //a矩阵列数
	if arow > 0 {
		acol = len(ma[0])
	}
	brow := len(mb) //b矩阵行数
	bcol := 0       //b矩阵列数
	if brow > 0 {
		bcol = len(mb[0])
	}
	if acol == brow {
		var res ComplexMatrix //结果矩阵
		for ar := 0; ar < arow; ar++ {
			var rr []complex128
			for bc := 0; bc < bcol; bc++ {
				var r complex128
				for ac := 0; ac < acol; ac++ {
					r += ma[ar][ac] * mb[ac][bc]
				}
				rr = append(rr, r)
			}
			res = append(res, rr)
		}
		return res, nil
	} else {
		return nil, fmt.Errorf("the number of columns of matrix A must be equal to the number of rows of matrix B[矩阵A的列数必须与矩阵B的行数相等]")
	}
}

/*********************************************
//功能:两个矩阵相加
//参数:矩阵B
//输出:矩阵的和,错误信息
//说明:两个矩阵都必须是m x n的矩阵,结果也是m x n的矩阵
//作者:wangjp
//时间:2020年10月14日
**********************************************/
func (ma ComplexMatrix) Add(mb ComplexMatrix) (ComplexMatrix, error) {
	arow := len(ma) //a矩阵行数
	acol := 0       //a矩阵列数
	if arow > 0 {
		acol = len(ma[0])
	}
	brow := len(mb) //b矩阵行数
	bcol := 0       //b矩阵列数
	if brow > 0 {
		bcol = len(mb[0])
	}
	if arow == brow && acol == bcol {
		var res ComplexMatrix //结果矩阵
		for i := 0; i < arow; i++ {
			var rr []complex128
			for j := 0; j < acol; j++ {
				rr = append(rr, ma[i][j]+mb[i][j])
			}
			res = append(res, rr)
		}
		return res, nil
	} else {
		return nil, fmt.Errorf("matrices A and B must have the same dimension[两个矩阵的维数必须相同]")
	}
}

/*********************************************
//功能:两个矩阵相减
//参数:矩阵B
//输出:矩阵的差,错误信息
//说明:两个矩阵都必须是m x n的矩阵,结果也是m x n的矩阵
//作者:wangjp
//时间:2020年10月14日
**********************************************/
func (ma ComplexMatrix) Sub(mb ComplexMatrix) (ComplexMatrix, error) {
	arow := len(ma) //a矩阵行数
	acol := 0       //a矩阵列数
	if arow > 0 {
		acol = len(ma[0])
	}
	brow := len(mb) //b矩阵行数
	bcol := 0       //b矩阵列数
	if brow > 0 {
		bcol = len(mb[0])
	}
	if arow == brow && acol == bcol {
		var res ComplexMatrix //结果矩阵
		for i := 0; i < arow; i++ {
			var rr []complex128
			for j := 0; j < acol; j++ {
				rr = append(rr, ma[i][j]-mb[i][j])
			}
			res = append(res, rr)
		}
		return res, nil
	} else {
		return nil, fmt.Errorf("matrices A and B must have the same dimension[两个矩阵的维数必须相同]")
	}
}

/*********************************************
//功能:用单元矩阵减去该矩阵
//参数:无
//说明:
//作者:wangjp
//时间:2020年10月14日
**********************************************/
func (m ComplexMatrix) SubFromIdentityMatrix() (ComplexMatrix, error) {
	if m.IsSquareMatix() {
		return InitComplexMatrixAsIdentity(len(m)).Sub(m)
	} else {
		return nil, fmt.Errorf("the matrix must be square matrix[必须是方阵才可以使用该方法]")
	}
}

/*********************************************
//功能: 获取矩阵的行数
//参数: 无
//返回: 矩阵行数
//说明:
//作者: wangjp
//时间: 2020年10月14日
**********************************************/
func (m ComplexMatrix) Rows() int {
	return len(m)
}

/*********************************************
//功能: 获取矩阵的列数
//参数: 无
//返回: 矩阵列数
//说明:
//作者: wangjp
//时间: 2020年10月14日
**********************************************/
func (m ComplexMatrix) Cols() int {
	if m.Rows() > 0 {
		return len(m[0])
	} else {
		return 0
	}
}

/*********************************************
//功能: 获取矩阵的行数和列数
//参数: 无
//返回: 矩阵行数和列数
//说明:
//作者: wangjp
//时间: 2020年10月14日
**********************************************/
func (m ComplexMatrix) RowsAndCols() (rows, cols int) {
	rows = m.Rows()
	cols = 0
	if rows > 0 {
		cols = len(m[0])
	}
	return
}

/*********************************************
//功能: 是否方阵
//参数: 无
//返回: 如果是方阵(行和列相同),返回true,否则false
//说明:
//作者: wangjp
//时间: 2020年10月14日
**********************************************/
func (m ComplexMatrix) IsSquareMatix() bool {
	rows := m.Rows()
	cols := 0
	if rows > 0 {
		cols = len(m[0])
	}
	return rows == cols
}

/*********************************************
//功能: 矩阵转置
//参数: 无
//返回: 转置后的矩阵
//说明:
//作者: wangjp
//时间: 2020年10月14日
**********************************************/
func (m ComplexMatrix) Transpose() ComplexMatrix {
	rows, cols := m.RowsAndCols()
	res := InitComplexMatrix(cols, rows)
	for i, mr := range m {
		for j, v := range mr {
			//注意不但要(i, j)--->(j, i)，且要a+bj--->a-bj(共轭运算)
			res[j][i] = complex(real(v), -imag(v))
		}
	}
	return res
}

/*********************************************
//功能: 矩阵乘以一个系数
//参数: 系数
//返回: 无
//说明:
//作者: wangjp
//时间: 2020年10月14日
**********************************************/
func (m ComplexMatrix) MulScalar(scalar complex128) {
	for i, mr := range m {
		for j := range mr {
			m[i][j] *= scalar
		}
	}
}

/*********************************************
//功能: 判断两个矩阵是否相等
//参数: 被比较的矩阵
//返回: 两个矩阵相同返回true,不同返回false
//说明:
//作者: wangjp
//时间: 2020年10月14日
**********************************************/
func (a ComplexMatrix) IsEqual(b ComplexMatrix) bool {
	ar, ac := a.RowsAndCols()
	br, bc := b.RowsAndCols()
	if ar == br && ac == bc { //矩阵维数相同
		for i, mar := range a {
			for j, av := range mar {
				if !ComplexsIsEqual(av, b[i][j]) { //有一个元素不同
					return false //则判断为不相等
				}
			}
		}
	} else {
		return false
	}
	return true
}

/***********************************************
功能:将矩阵中的所有元素取共轭
输入:inplace bool:是否更改原始矩阵
输出:所有元素共轭后的矩阵
说明:复数 a+ib 的共轭复数是 a-ib
编辑:wang_jp
时间:2020年10月21日
***********************************************/
func (cmpm ComplexMatrix) Conjugate(inplace ...bool) ComplexMatrix {
	res := cmpm.Copy()
	if len(inplace) > 0 {
		if inplace[0] {
			res = cmpm
		}
	}
	for i, arv := range res {
		for j, v := range arv {
			res[i][j] = ComplexConjugate(v)
		}
	}
	return res
}

/***********************************************
功能:获取矩阵中各个元素的模(绝对值)
输入:
输出:实数矩阵
编辑:wang_jp
时间:2020年10月21日
***********************************************/
func (cmpm ComplexMatrix) Abs() Matrix {
	var res Matrix
	for _, v := range cmpm {
		res = append(res, v.Abs())
	}
	return res
}

/***********************************************
功能:获取矩阵中各个元素的角度(主值)
输入:
输出:实数矩阵
说明:复数Z=a+bi化成三角式r(cosθ+jsinθ),可简写作r∠θ,
	其中模(绝对值)r=√(a²+b²)；复角θ由tanθ=b/a解出并在0≤θ＜360°范围内取值（主值）
编辑:wang_jp
时间:2020年10月16日
***********************************************/
func (cmpm ComplexMatrix) Angle() Matrix {
	var res Matrix
	for _, v := range cmpm {
		res = append(res, v.Angle())
	}
	return res
}
