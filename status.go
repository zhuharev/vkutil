package vkutil

import (
	"encoding/json"
	"fmt"
	"github.com/zhuharev/vk"
	"net/url"
)

func (a *Api) StatusSet(text string) error {
	bts, e := a.VkApi.Request(vk.METHOD_STATUS_SET, url.Values{"text": {text}})
	if e != nil {
		return e
	}
	var (
		resp ResponseInt
	)
	e = json.Unmarshal(bts, &resp)
	if e != nil {
		return e
	}
	if resp.ResponseError.Error.Code != 0 {
		return fmt.Errorf("%s", resp.Error.Msg)
	}
	if resp.Response != 1 {
		return fmt.Errorf("Unknown error, resp %d, need 1", resp.Response)
	}
	return nil
}
