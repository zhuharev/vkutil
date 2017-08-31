package vkutil

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/zhuharev/vk"
)

// Topic board item
type Topic struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Created   EpochTime `json:"created"`
	CreatedBy int       `json:"created_by"`
	Updated   EpochTime `json:"upadted"`
	UpdatedBy int       `json:"updated_by"`
	IsClosed  int       `json:"is_closed"`
	Comments  int       `json:"comments"`
}

// ResponseTopics represent topic items
type ResponseTopics struct {
	Resp struct {
		Count int     `json:"count"`
		Items []Topic `json:"items"`
	} `json:"response"`
}

// BoardGetCommetns returns board items
func (api *Api) BoardGetCommetns(groupID, topicID int, args ...url.Values) ([]Comment, error) {
	params := setToUrlValues("group_id", groupID, args...)
	params = setToUrlValues("topic_id", topicID, params)
	//params = setToUrlValues("sort", "desc", params)

	bts, err := api.VkApi.Request(vk.METHOD_BOARD_GET_COMMENTS, params)
	if err != nil {
		return nil, err
	}
	log.Println(string(bts))
	var r ResponseComments
	err = json.Unmarshal(bts, &r)
	if err != nil {
		return nil, err
	}
	return r.Response.Items, nil
}
