package vkutil

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/zhuharev/vk"
	"net/url"
	"strings"
)

func (api *Api) GetCitiesById(id int, ids ...int) ([]City, error) {
	form := url.Values{"city_ids": {fmt.Sprint(id)}}
	if len(ids) > 0 {
		ids = append([]int{id}, ids...)
		form.Set("city_ids", strings.Join(arrIntToStr(ids), ","))
	}
	resp, e := api.VkApi.Request(vk.METHOD_DATABASE_GET_CITIES_BY_ID, form)
	if e != nil {
		return nil, e
	}
	type cities struct {
		Response []City `json:"response"`
	}
	var c cities
	e = json.Unmarshal(resp, &c)
	return c.Response, e
}

func (api *Api) GetCityById(id int) (c City, e error) {
	var (
		cts []City
	)
	cts, e = api.GetCitiesById(id)
	if e != nil {
		return
	}
	if len(cts) != 1 {
		e = fmt.Errorf("Unknown err expected %d, got %d", 1, len(cts))
		return
	}
	return cts[0], nil
}

func (api *Api) GetCityId(cname string) (int, error) {
	// russia only
	resp, e := api.VkApi.Request(vk.METHOD_DATABASE_GET_CITIES, url.Values{"q": {cname}, "country_id": {"1"}, "count": {"1"}})
	if e != nil {
		return 0, e
	}
	var rc ResponseCities
	e = json.Unmarshal(resp, &rc)
	if e != nil {
		return 0, e
	}
	color.Green("%v", rc.Response.Items)

	if rc.Response.Count == 0 {
		return 0, fmt.Errorf("City %s not found", cname)
	}
	return rc.Response.Items[0].Id, nil
}

func (api *Api) GetCities(q string) (r ResponseCities, e error) {
	r = ResponseCities{}
	var (
		resp []byte
	)
	resp, e = api.VkApi.Request(vk.METHOD_DATABASE_GET_CITIES, url.Values{"q": {q}, "country_id": {"1"}})
	if e != nil {
		return
	}
	e = json.Unmarshal(resp, &r)
	if e != nil {
		color.Green("%s", resp)
		return
	}
	return
}
