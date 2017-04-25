package vkutil

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"

	"github.com/zhuharev/vk"
)

func (api *Api) UtilsResolveScreenName(name string) (ObjectType, int, error) {
	type Responseo struct {
		Respo struct {
			Type     ObjectType `json:"type"`
			ObjectId int        `json:"object_id"`
		} `json:"response"`
	}
	bts, e := api.VkApi.Request(vk.METHOD_UTILS_RESOLVE_SCREEN_NAME, url.Values{
		"screen_name": {name},
	})
	if e != nil {
		return "", 0, e
	}

	fmt.Println(string(bts))

	var res Responseo

	e = json.Unmarshal(bts, &res)
	return res.Respo.Type, res.Respo.ObjectId, e
}

// vkutil utils

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
		if i%700 == 0 {
			res = append(res, []int{})
			if i > 701 {
				//fmt.Println(res[len(res)-2][0], res[len(res)-2][len(res[len(res)-2])-1])

			}
		}
		res[len(res)-1] = append(res[len(res)-1], v)
	}
	return
}

func arrUniq(in []int) (out []int) {
	var (
		dup = map[int]struct{}{}
	)
	for _, v := range in {
		if _, ok := dup[v]; !ok {
			out = append(out, v)
			dup[v] = struct{}{}
		}
	}
	sort.Ints(out)
	return
}

func debug(f string, args ...interface{}) {
	if DEBUG {
		fmt.Printf(f+"\n", args...)
	}
}
