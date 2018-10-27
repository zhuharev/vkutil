package util

import (
	"fmt"
	"net/url"
	"strings"
)

func SetToUrlValues(key string, val interface{}, params ...url.Values) url.Values {
	param := url.Values{}
	if len(params) > 0 {
		param = params[0]
	}
	switch i := val.(type) {
	case int, int64:
		if i != 0 {
			param.Set(key, fmt.Sprint(i))
		}
	case string:
		if i != "" {
			param.Set(key, i)
		}
	case []string:
		param.Set(key, strings.Join(i, ","))
	case []int:
		param.Set(key, strings.Join(ArrIntToStr(i), ","))
	case bool:
		fb := "1"
		if !i {
			fb = "0"
		}
		param.Set(key, fb)
	}
	return param
}

func ArrIntToStr(arr []int) (sarr []string) {
	for _, v := range arr {
		sarr = append(sarr, fmt.Sprint(v))
	}
	return
}
