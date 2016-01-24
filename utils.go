package vkutil

import (
	"fmt"
	"strconv"
)

func arrIntToStr(arr []int) (sarr []string) {
	for _, v := range arr {
		sarr = append(sarr, fmt.Sprint(v))
	}
	return
}

func arrStrToInt(arr []string) (iarr []int) {
	for _, v := range arr {
		i, _ := strconv.Atoi(v)
		if i != 0 {
			iarr = append(iarr, i)
		}
	}
	return
}

func uniqStrArr(arr []string) []string {
	var m = map[string]interface{}{}
	for _, v := range arr {
		m[v] = nil
	}
	arr = []string{}
	for k := range m {
		arr = append(arr, k)
	}
	return arr
}

func arrSplit1K(arr []int) (res [][]int) {
	//res = append(res, []int)
	for i, v := range arr {
		if i%1000 == 0 {
			res = append(res, []int{})
		}
		res[len(res)-1] = append(res[len(res)-1], v)
	}
	return
}
func debug(f string, args ...interface{}) {
	if DEBUG {
		fmt.Printf(f+"\n", args...)
	}
}
