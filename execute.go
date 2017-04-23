package vkutil

import (
	"github.com/zhuharev/vk"
	"log"
	"net/url"
)

func (api *Api) Execute(code string) ([]byte, error) {
	if vk.DEBUG {
		log.Println(code)
	}
	return api.VkApi.Request(vk.METHOD_EXECUTE, url.Values{
		"code": {code},
	})
}
