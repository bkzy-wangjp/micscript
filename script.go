// mics.go
package mics

import (
	"time"
)

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
	t, err = time.ParseInLocation("2006年01月02日 15:04:05", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006年01月02日 15:04", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006年01月02日 15点04分05秒", s, location)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006年01月02日 15点04分", s, location)
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
功能:比较两个切片组是否相同
输入:切片a和切片b
输出:相同时输出true，不同时输出false
时间:2021年8月2日
编辑:wang_jp
************************************************************/
func StringsCompiler(slices_a []string, slices_b []string) bool {
	if len(slices_a) != len(slices_b) {
		return false
	}
	for i, a := range slices_a {
		if a != slices_b[i] {
			return false
		}
	}
	return true
}
