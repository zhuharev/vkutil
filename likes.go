package vkutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhuharev/vk"
	"net/url"
)

func (api *Api) LikesAdd(likeType string, ownerId, itemId int,
	args ...url.Values) error {
	resp, err := api.vkApi.Request(vk.METHOD_LIKES_ADD, url.Values{
		"type":     {likeType},
		"item_id":  {fmt.Sprint(itemId)},
		"owner_id": {fmt.Sprint(ownerId)},
	})
	if err != nil {
		return err
	}
	fmt.Println(string(resp))
	var r ResponseLikes
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return err
	}
	return nil
}

func (api *Api) LikesIsLiked(likeType string, ownerId, itemId int,
	args ...url.Values) (liked bool, copied bool, err error) {
	resp, err := api.vkApi.Request(vk.METHOD_LIKES_IS_LIKED, url.Values{
		"type":     {likeType},
		"item_id":  {fmt.Sprint(itemId)},
		"owner_id": {fmt.Sprint(ownerId)},
	})
	if err != nil {
		return
	}
	type ResponseLikesIsliked struct {
		Response struct {
			Liked  int `json:"liked"`
			Copied int `json:"copied"`
		} `json:"response"`
		ResponseError
	}
	var r ResponseLikesIsliked
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return
	}
	if r.Response.Liked == 1 {
		liked = true
	}
	if r.Response.Copied == 1 {
		copied = true
	}
	if r.Error.Msg != "" {
		err = errors.New(r.Error.Msg)
	}

	return
}

func (api *Api) LikesGetList(likeType string, ownerId, itemId int,
	args ...url.Values) (likes []int, err error) {
	resp, err := api.vkApi.Request(vk.METHOD_LIKES_IS_LIKED, url.Values{
		"type":     {likeType},
		"item_id":  {fmt.Sprint(itemId)},
		"owner_id": {fmt.Sprint(ownerId)},
	})
	if err != nil {
		return
	}

	var r IdsResp
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return
	}

	likes = r.Items

	return
}

func (api *Api) Get25KLikes(likeType string, ownerId, itemId, offset int,
	args ...url.Values) (count int, likes []int, err error) {
	fmtcode := `var i=%d,lim=i+25, a=API.likes.getList({type:"%[2]s",owner_id:%[3]d,item_id:%[4]d,count:1000,offset:i*1000});i=i+1;
while (i<lim){
	a.items=a.items+API.likes.getList({type:"%[2]s",owner_id:%[3]d,item_id:%[4]d,count:1000,offset:i*1000}).items;
	i=i+1;
};
return a;`
	code := fmt.Sprintf(fmtcode, offset, likeType, ownerId, itemId)
	bts, e := api.Execute(code)
	if e != nil {
		return
	}
	var resp ResponseIds
	e = json.Unmarshal(bts, &resp)
	if e != nil {
		return
	}
	if resp.Error.Msg != "" {
		e = errors.New(resp.Error.Msg)
		return
	}
	return resp.Resp.Count, resp.Resp.Items, nil
}

func (api *Api) LikesGetAll(likeType string, ownerId, itemId int,
	args ...url.Values) (likes []int, e error) {

	i := 0
	tmplikes := []int{}
	cnt, likes, e := api.Get25KLikes(likeType, ownerId, itemId, i, args...)
	if e != nil {
		return
	}
	i += 25

	for len(likes) < cnt {
		_, tmplikes, e = api.Get25KLikes(likeType, ownerId, itemId, i, args...)
		if e != nil {
			return
		}
		likes = append(likes, tmplikes...)
		i += 25
	}

	return
}
