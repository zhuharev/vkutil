package vkutil

import (
	"encoding/json"
	"fmt"
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
	Text string `text`
	Type string `post_type`
}

func ParseIdsResponse(data []byte) (count int, ids []int, err error) {
	var r ResponseIds
	err = json.Unmarshal(data, &r)
	if err != nil {
		return 0, []int{}, err
	}
	fmt.Println(r.Resp.Items)
	return r.Resp.Count, r.Resp.Items, nil
}
