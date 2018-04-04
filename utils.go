package vkutil

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"

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
		if k == "" || k == "0" || IsZeroOfUnderlyingType(m) {
			continue
		}
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

func setToUrlValues(key string, val interface{}, params ...url.Values) url.Values {
	param := url.Values{}
	if len(params) > 0 {
		param = params[0]
	}
	switch i := val.(type) {
	case int, int64:
		param.Set(key, fmt.Sprint(i))
	case string:
		param.Set(key, i)
	case []string:
		param.Set(key, strings.Join(i, ","))
	case []int:
		param.Set(key, strings.Join(arrIntToStr(i), ","))
	}
	return param
}

func ParseCallbackURL(uri string) (token string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	vals, err := url.ParseQuery(u.Fragment)
	if err != nil {
		return
	}
	return vals.Get("access_token"), nil
}

func ParseDomain(uri string) (token string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	return strings.TrimPrefix(u.Path, "/"), nil
}

func ParseBoardURL(uri string) (groupID, topicID, postID int, err error) {
	uri = strings.TrimSuffix(uri, "?offset=last&scroll=1")
	if strings.Contains(uri, "?") {
		_, err = fmt.Sscanf(uri, "https://vk.com/topic-%d_%d?offset=%d", &groupID, &topicID, &postID)
		if err != nil {
			return
		}
		return
	}
	_, err = fmt.Sscanf(uri, "https://vk.com/topic-%d_%d", &groupID, &topicID)
	if err != nil {
		return
	}
	return
}

// IsZeroOfUnderlyingType check is x zero value
func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
