package vkutil

import (
	"log"
	"net/url"

	"github.com/zhuharev/vk"
)

func (api *Api) Execute(code string) ([]byte, error) {
	log.Println(code)
	return api.VkApi.Request(vk.METHOD_EXECUTE, url.Values{
		"code": {code},
	})
}
