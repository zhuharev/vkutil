package vkutil

import (
	"github.com/zhuharev/vk"
	"net/url"
	"strings"
)

func (a *Api) NewsFeedAddBan(ids []int) error {
	_, err := a.VkApi.Request(vk.METHOD_NEWSFEED_ADD_BAN, url.Values{
		"user_ids": {strings.Join(arrIntToStr(ids), ",")},
	})
	return err
}
