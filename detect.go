package vkutil

import (
	"github.com/fatih/color"
	"net/url"
)

func (api *Api) detectCity(id int) (cid int, e error) {
	var (
		ids   []int
		users []User
		dup   = map[int]int{}
	)
	ids, _, e = api.UtilsFriendsGetOne(id)
	if e != nil {
		return
	}
	users, e = api.UsersGet(ids, url.Values{"fields": {"city"}})
	if e != nil {
		return
	}
	for _, u := range users {
		dup[u.City.ID]++
		if dup[u.City.ID] > dup[cid] {
			cid = u.City.ID
		}
	}
	for k, v := range dup {
		if v < 2 {
			continue
		}
		color.Green("%d - %v", k, v)
	}
	return
}

func (api *Api) DetectCity(uid int) (City, error) {
	cid, e := api.detectCity(uid)
	if e != nil {
		return City{}, e
	}
	return api.GetCityById(cid)
}
