// mics.go
package mics

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"strings"
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
	if strings.Contains(s, ".") {
		strs := strings.Split(s, ".") //用.切片
		dsec := "."                   //小数点后的秒格式
		tzone := "Z0700"
		havezone := false //是否包含时区
		if len(strs) > 1 {
			dseclen := 0                        //小数部分的长度
			if strings.Contains(strs[1], "+") { //包含时区
				havezone = true
				dotseconds := strings.Split(strs[1], "+") //分割小数的秒和时区
				dseclen = len(dotseconds[0])
				if len(dotseconds) > 1 { //时区切片
					if strings.Contains(dotseconds[1], ":") { //时区切片含有分号
						tzone = "Z07:00"
					}
				}
			} else {
				dseclen = len(strs[1])
			}

			for i := 0; i < dseclen; i++ {
				dsec += "0"
			}
		}
		layout := "2006-01-02 15:04:05"
		if strings.Contains(s, "T") {
			layout = "2006-01-02T15:04:05"
		}
		layout += dsec
		if havezone { //包含时区
			layout += tzone
			t, err := time.Parse(layout, s)
			if err == nil {
				return t, nil
			}
		} else {
			t, err := time.ParseInLocation(layout, s, location)
			if err == nil {
				return t, nil
			}
		}
	} else {
		if strings.Contains(s, "+") { //包含时区
			tzone := "Z0700"
			strs := strings.Split(s, "+") //时区
			if len(strs) > 1 {            //时区切片
				if strings.Contains(strs[1], ":") { //时区切片含有分号
					tzone = "Z07:00"
				}
			}
			layout := "2006-01-02 15:04:05"
			if strings.Contains(s, "T") {
				layout = "2006-01-02T15:04:05"
			}
			layout += tzone
			t, err := time.Parse(layout, s)
			if err == nil {
				return t, nil
			}
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
	return t, fmt.Errorf("The raw time string is:[%s],the error is:[%s]", s, err.Error())
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

/****************************************************
功能：判断元素在数组、Map中是否存在
输入：元素、数组或者Map、Slice
输出：存在输出true，不存在输出false
说明：对于数组、Slice，判断的是值是否存在，对于Map，判断的是Key是否存在
时间：2019年12月15日
编辑：wang_jp
****************************************************/
func IsExistItem(element interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == element {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(element)).IsValid() {
			return true
		}
	}
	return false
}

/*************************************************
功能:对字符串进行MD5加密
输入:待加密字符串
输出:加密后的字符串
说明:
编辑:wang_jp
时间:2020年3月12日
*************************************************/
func Md5str(strs ...string) string {
	var str string
	for _, s := range strs {
		str += s
	}
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

/*************************************************
功能:对key为string类型的map按key进行排序
输入:待排序的map,排序方式:{desc:按照[z-a]逆向排序,asc:按照[a-z]正向排序(默认)}
输出:排序后的key值切片
说明:
编辑:wang_jp
时间:2021年8月4日
*************************************************/
func SortStringMap(mp map[string]string, orderby ...string) []string {
	var keys []string
	for key := range mp {
		keys = append(keys, key)
	}
	SortStringSlice(keys, orderby...)
	return keys
}

/*************************************************
功能:对字符串切片进行排序
输入:待排序的字符串切片,排序方式:{desc:按照[z-a]逆向排序,asc:按照[a-z]正向排序}
输出:无
说明:排序后的结果存储在输入切片中
编辑:wang_jp
时间:2021年8月4日
*************************************************/
func SortStringSlice(values []string, orderby ...string) {
	oder := "asc"
	if len(orderby) > 0 {
		oder = orderby[0]
	}
	vLen := len(values)
	flag := true
	if oder == "asc" {
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
}
