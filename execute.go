package vkutil

import (
	"github.com/zhuharev/vk"
	"net/url"
)

func (api *Api) Execute(code string) ([]byte, error) {
	return api.vkApi.Request(vk.METHOD_EXECUTE, url.Values{
		"code": {code},
	})
}
