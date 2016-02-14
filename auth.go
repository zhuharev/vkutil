package vkutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/zhuharev/vk"
	"net/url"
)

func (api *Api) SignUp(firstName, lastName, phone string, args ...url.Values) (sid string, e error) {
	uv := url.Values{
		"first_name":    {firstName},
		"last_name":     {lastName},
		"phone":         {phone},
		"client_id":     {fmt.Sprint(api.vkApi.ClientId)},
		"client_secret": {api.vkApi.ClientSecret},
	}
	if len(args) > 0 {
		for k, v := range args[0] {
			if len(v) > 0 {
				uv.Set(k, v[0])

			}
		}
	}
	resp, e := api.vkApi.Request(vk.METHOD_AUTH_SIGNUP, uv)
	if e != nil {
		return "", e
	}
	color.Green("%s", resp)
	type Response struct {
		Resp struct {
			Sid string `json:"sid"`
		} `json:"response"`
		ResponseError
	}
	var r Response
	e = json.Unmarshal(resp, &r)
	if e != nil {
		return "", e
	}
	if r.Error.Code == 14 && api.StdinAllow {
		color.Red("Captcha")
		color.Yellow("Open in browser %s", r.Error.CaptchaImg)
		var captcha string
		fmt.Scanln(&captcha)
		var a = url.Values{}
		if len(args) > 0 {
			a = args[0]
		}
		a.Set("captcha_sid", r.Error.CaptchaSid)
		a.Set("captcha_key", captcha)
		return api.SignUp(firstName, lastName, phone, a)
	}
	if r.Error.Code != 0 {
		return "", errors.New(r.Error.Msg)
	}
	return r.Resp.Sid, nil
}

func (api *Api) Confirm(phone string, code string, password string) (sid string, e error) {
	resp, e := api.vkApi.Request(vk.METHOD_AUTH_CONFIRM, url.Values{
		"phone":         {phone},
		"client_id":     {fmt.Sprint(api.vkApi.ClientId)},
		"client_secret": {api.vkApi.ClientSecret},
		"code":          {code},
		"password":      {password},
	})
	color.Green("Confirm with code %s", code)
	if e != nil {
		return "", e
	}
	type Response struct {
		Resp struct {
			Sid string `json:"sid"`
		} `json:"response"`
		ResponseError
	}

	color.Green("%s", resp)

	var r Response
	e = json.Unmarshal(resp, &r)
	if e != nil {
		return "", e
	}
	if r.Error.Code == 14 && api.StdinAllow {
		color.Red("Captcha")
		color.Yellow("Open in browser %s", r.Error.CaptchaImg)
		var captcha string
		fmt.Scanln(&captcha)
		var a = url.Values{}

		a.Set("captcha_sid", r.Error.CaptchaSid)
		a.Set("captcha_key", captcha)
		return api.Confirm(phone, code, password)
	}
	if r.Error.Code != 0 {
		return "", errors.New(r.Error.Msg)
	}
	return r.Resp.Sid, nil
}
