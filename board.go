package vkutil

import (
	"context"
	"encoding/json"
	"log"
	"net/url"

	"github.com/zhuharev/vk"
)

const (
	MethodBoardGetTopics = "board.getTopics"
)

// Topic board item
type Topic struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Created   EpochTime `json:"created"`
	CreatedBy int       `json:"created_by"`
	Updated   EpochTime `json:"updated"`
	UpdatedBy int       `json:"updated_by"`
	IsClosed  int       `json:"is_closed"`
	Comments  int       `json:"comments"`
}

type RespTopic struct {
	Count int     `json:"count"`
	Items []Topic `json:"items"`
}

// ResponseTopics represent topic items
type ResponseTopics struct {
	Resp RespTopic `json:"response"`
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

func (api *Api) BoardGetAllTopics(ctx context.Context, groupID int) ([]Topic, error) {
	var (
		count  = 100
		offset = 0
		result []Topic
	)
	for {
		topics, err := api.BoardGetTopics(ctx, groupID, count, offset)
		if err != nil {
			return nil, err
		}
		offset += len(topics)
		result = append(result, topics...)

		if len(topics) < count {
			break
		}
	}
	return result, nil
}

func (api *Api) BoardGetTopics(ctx context.Context, groupID int, count int, offset int) ([]Topic, error) {
	params := setToUrlValues("group_id", groupID)
	params = setToUrlValues("count", count, params)
	params = setToUrlValues("offset", offset, params)

	bts, err := api.VkApi.Request(MethodBoardGetTopics, params)
	if err != nil {
		return nil, err
	}

	var r ResponseTopics
	err = json.Unmarshal(bts, &r)
	if err != nil {
		return nil, err
	}
	return r.Resp.Items, nil
}
