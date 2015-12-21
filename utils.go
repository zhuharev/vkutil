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
