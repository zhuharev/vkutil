package vkutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/zhuharev/vk"
)

func (api *Api) MessagesGet(params ...url.Values) ([]Message, error) {
	rparams := url.Values{}
	if len(params) == 1 {
		rparams = params[0]
	}
	resp, err := api.VkApi.Request(vk.METHOD_MESSAGES_GET, rparams)
	if err != nil {
		return nil, err
	}
	return ParseMessagesResponse(resp)
}

func (api *Api) MessagesSend(uid int, message string,
	params ...url.Values) (int, error) {
	rparams := url.Values{}
	if len(params) == 1 {
		rparams = params[0]
	}
	rparams.Set("user_id", fmt.Sprint(uid))
	rparams.Set("message", message)
	resp, err := api.VkApi.Request(vk.METHOD_MESSAGES_SEND, rparams)
	if err != nil {
		return 0, err
	}
	var r ResponseInt
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return 0, err
	}
	if r.Error.Code != 0 {
		return 0, errors.New(r.Error.Msg)
	}
	return r.Response, nil
}

func (api *Api) MessagesGetDialogs(params ...url.Values) ([]Dialog, error) {
	rparams := url.Values{}
	if len(params) == 1 {
		rparams = params[0]
	}
	resp, err := api.VkApi.Request(vk.METHOD_MESSAGES_GET_DIALOGS, rparams)
	if err != nil {
		return nil, err
	}
	return ParseDialogsResponse(resp)
}

func (api *Api) MessagesMarkAsRead(messageId int, messages ...int) error {
	var ids = []int{messageId}
	if len(messages) > 0 {
		ids = append(ids, messages...)
	}
	resp, err := api.VkApi.Request(vk.METHOD_MESSAGES_MARK_AS_READ, url.Values{
		"message_ids": {strings.Join(arrIntToStr(ids), ",")},
	})
	if err != nil {
		return err
	}
	var r ResponseInt
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return err
	}
	if r.Error.Code != 0 {
		return errors.New(r.Error.Msg)
	}
	return nil
}

func (api *Api) MessagesGetHistory(udserId int, params ...url.Values) ([]Message, error) {
	rparams := url.Values{}
	if len(params) == 1 {
		rparams = params[0]
	}
	rparams.Set("user_id", fmt.Sprint(udserId))
	resp, err := api.VkApi.Request(vk.METHOD_MESSAGES_GET_HISTORY, rparams)
	if err != nil {
		return nil, err
	}
	return ParseMessagesResponse(resp)
}

func (api *Api) MessagesSetActivity(udserID int, params ...url.Values) error {
	rparams := url.Values{}
	if len(params) == 1 {
		rparams = params[0]
	}
	rparams.Set("user_id", fmt.Sprint(udserID))
	rparams.Set("type", "typing")
	resp, err := api.VkApi.Request(vk.METHOD_MESSAGES_SET_ACTIVITY, rparams)
	if err != nil {
		return err
	}
	_, err = ParseIntResponse(resp)
	return err
}
