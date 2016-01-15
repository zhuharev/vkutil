package vkutil

import (
	"encoding/json"
)

type Response int

const (
	RESPONSE_IDS Response = iota
)

type ResponseIds struct {
	Resp IdsResp `json:"response"`
	ResponseError
}

type ResponseError struct {
	Error Err `json:"error"`
}

type Err struct {
	Code int    `json:"error_code"`
	Msg  string `json:"error_msg"`
}

type IdsResp struct {
	Count int   `json:"count"`
	Items []int `json:"items"`
}

/* likes */
type ResponseLikes struct {
}

type ResponsePosts struct {
	Response RespPosts `json:"response"`
	ResponseError
}

type RespPosts struct {
	Count int `count`
	Items []Post
}

type Post struct {
	Id        int    `post_id`
	AccessKey string `access_key`
	Likes     struct {
		Count int `count`
	}
	Text     string `text`
	PostType string `post_type`
}

type ResponseInt struct {
	Response int `json:"response"`
	ResponseError
}

/* messages */

type ResponseMessages struct {
	Response RespMessages `json:"response"`
	ResponseError
}

//todo create RespCount struct
type RespMessages struct {
	Count int       `json:"count"`
	Items []Message `json:"items"`
}

type RespFriendsGetMutual struct {
	Id            int   `json:"id"`
	CommonFriends []int `json:"common_friends"`
	CommonCount   int   `json:"common_count"`
}

type ResponseFriendsGetMutual struct {
	Response []RespFriendsGetMutual `json:"response"`
	ResponseError
}

type RespUnreadMessages struct {
	Count         int `json:"count"`
	UnreadDialogs int `json:"unread_dialogs"`
	Items         []struct {
		Unread  int     `json:"unread"`
		Message Message `json:"message"`
	}
}

type Message struct {
	Id        int    `json:"id"`
	Date      int    `json:"date"`
	Out       int    `json:"out"`
	UserId    int    `json:"user_id"`
	ReadState int    `json:"read_state"`
	Title     string `json:"title"`
	Body      string `json:"body"`
}

type RespUserWithCount struct {
	Count int    `json:"count"`
	Items []User `json:"items"`
}

type ResponseUserWithCount struct {
	Response RespUserWithCount `json:"response"`
	ResponseError
}

type RespPhotos struct {
	Count int     `json:"count"`
	Items []Photo `json:"items"`
}

type ResponsePhotos struct {
	Response RespPhotos `json:"response"`
	ResponseError
}

type City struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Area   string `json:"area,omitempty"`
	Region string `json:"region"`
}

type ResponseCities struct {
	Response struct {
		Count int    `json:"count"`
		Items []City `json:"items"`
	} `json:"response"`
	ResponseError
}

func ParseIdsResponse(data []byte) (count int, ids []int, err error) {
	var r ResponseIds
	err = json.Unmarshal(data, &r)
	if err != nil {
		return 0, []int{}, err
	}
	return r.Resp.Count, r.Resp.Items, nil
}

func ParseUsersResponse(data []byte) (users []User, err error) {
	var r ResponseUsers
	err = json.Unmarshal(data, &r)
	if err != nil {
		return []User{}, err
	}
	return r.Response, nil
}

func ParseMessagesResponse(data []byte) (ms []Message, err error) {
	var r ResponseMessages
	err = json.Unmarshal(data, &r)
	if err != nil {
		return []Message{}, err
	}
	return r.Response.Items, nil
}

func ParseUnreadMessagesResponse(data []byte) (ms []Message, err error) {
	var r struct {
		RespUnreadMessages `json:"response"`
		ResponseError
	}
	err = json.Unmarshal(data, &r)
	if err != nil {
		return []Message{}, err
	}
	for i := range r.Items {
		ms = append(ms, r.Items[i].Message)
	}
	return ms, nil
}

func ParseResponseUserWithCount(data []byte) (users []User, cnt int, e error) {
	var r ResponseUserWithCount
	e = json.Unmarshal(data, &r)
	return r.Response.Items, r.Response.Count, e
}
