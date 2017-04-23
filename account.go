package vkutil

import (
	"github.com/zhuharev/vk"
)

func (api *Api) AccountSetOnline() error {
	bts, e := api.VkApi.Request(vk.METHOD_ACCOUNT_SET_ONLINE)
	if e != nil {
		return e
	}
	_, e = ParseIntResponse(bts)
	return e
}

func (api *Api) AccountSetOffline() error {
	bts, e := api.VkApi.Request(vk.METHOD_ACCOUNT_SET_OFFLINE)
	if e != nil {
		return e
	}
	_, e = ParseIntResponse(bts)
	return e
}
