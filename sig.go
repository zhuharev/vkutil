package vkutil

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/url"
	"sort"
)

type Value struct {
	Key   string
	Value string
}

type Values []Value

func NewValuesFromParam(param url.Values) Values {
	vs := Values{}
	for k, v := range param {
		if len(v) > 0 {
			vs = append(vs, Value{k, v[0]})
		}
	}
	return vs
}

func (v Values) Less(i, j int) bool {
	return v[i].Key < v[j].Key
}

func (v Values) Len() int      { return len(v) }
func (v Values) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

func (v Values) Sig(method, secret string) string {
	res := ""
	sort.Sort(v)
	for _, v := range v {
		res += v.Key + "=" + v.Value
	}

	str := "/" + method + "?" + res

	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func Sig(param url.Values, method, secret string) string {
	str := "/" + method + "?" + "v=5.64&lang=ru&https=1&access_token=79b27d162a1fcf7ebd24f563ac4cffc9f772a95b7f8744c2417509e664526c02dad185a03c6d1694c42c1"
	log.Println(str)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
