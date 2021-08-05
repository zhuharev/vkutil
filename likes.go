package vkutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/zhuharev/vk"
)

func (api *Api) LikesAdd(likeType ObjectType, ownerId, itemId int,
	args ...url.Values) error {
	resp, err := api.VkApi.Request(vk.METHOD_LIKES_ADD, url.Values{
		"type":     {string(likeType)},
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

func (api *Api) LikesIsLiked(likeType ObjectType, ownerId, itemId int,
	args ...url.Values) (liked bool, copied bool, err error) {
	resp, err := api.VkApi.Request(vk.METHOD_LIKES_IS_LIKED, url.Values{
		"type":     {string(likeType)},
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

func (api *Api) LikesGetList(likeType ObjectType, ownerId, itemId int,
	args ...url.Values) (likers []int, err error) {
	param := setToUrlValues("type", string(likeType), args...)
	param.Set("item_id", fmt.Sprint(itemId))
	param = setToUrlValues("owner_id", ownerId, param)
	resp, err := api.VkApi.Request(vk.METHOD_LIKES_GET_LIST, param)
	if err != nil {
		return
	}

	var r ResponseIds
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return
	}

	likers = r.Resp.Items

	return
}

func (api *Api) LikesGetListAll(likeType ObjectType, ownerId, itemId int,
	args ...url.Values) (likers []int, err error) {
	param := setToUrlValues("type", string(likeType), args...)
	param.Set("item_id", fmt.Sprint(itemId))
	param = setToUrlValues("owner_id", ownerId, param)
	param = setToUrlValues("count", 1000, param)
	offset := 0
	for {
		param = setToUrlValues("offset", offset, param)
		var resp []byte
		resp, err = api.VkApi.Request(vk.METHOD_LIKES_GET_LIST, param)
		if err != nil {
			return
		}

		var r ResponseIds
		err = json.Unmarshal(resp, &r)
		if err != nil {
			return
		}

		likers = append(likers, r.Resp.Items...)
		if offset >= r.Resp.Count {
			break
		}
		offset += 1000
	}

	return
}

func (api *Api) Get25KLikes(likeType ObjectType, ownerId, itemId, offset int,
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
		err = errors.New(resp.Error.Msg)

		return
	}

	return resp.Resp.Count, resp.Resp.Items, nil
}

func (api *Api) LikesGetAll(likeType ObjectType, ownerId, itemId int,
	args ...url.Values) (likes []int, e error) {

	i := 0
	var tmplikes []int
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

func (a *Api) LikesDelete(ot ObjectType, ownerId int, objectId int, params ...url.Values) error {
	var (
		param = url.Values{}
	)
	if params != nil {
		param = params[0]
	}
	param.Set("type", string(ot))
	param.Set("owner_id", fmt.Sprint(ownerId))
	param.Set("item_id", fmt.Sprint(objectId))

	bts, e := a.VkApi.Request(vk.METHOD_LIKES_DELETE, param)
	if e != nil {
		return e
	}
	_, e = ParseIntResponse(bts)
	return e
}
