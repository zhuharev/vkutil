package vkutil

import (
	"context"
	"net/url"

	"github.com/zhuharev/vk"
)

type poolsClient struct {
	client *Api
}

func newPoolsClient(client *Api) *poolsClient {
	return &poolsClient{
		client: client,
	}
}

func (p *poolsClient) GetByID(ctx context.Context, ownerID, poolID int) (*Poll, error) {
	params := url.Values{
		"owner_id": {toString(ownerID)},
		"poll_id":  {toString(poolID)},
	}

	var res PoolResponse

	err := p.client.VkApi.RequestTypedContext(ctx, vk.MethodPoolsGetByID, params, &res)
	if err != nil {
		return nil, err
	}

	return &res.Response, nil
}

func (p *poolsClient) GetVoters(ctx context.Context, ownerID, pollID int, answerIDs []int, offset, count int) (map[int][]int, error) {
	params := url.Values{
		"owner_id":   {toString(ownerID)},
		"poll_id":    {toString(pollID)},
		"answer_ids": {arrIntToString(answerIDs)},
		"offset":     {toString(offset)},
		"count":      {toString(count)},
	}

	var res PollVoters

	err := p.client.VkApi.RequestTypedContext(ctx, vk.MethodPoolsGetByID, params, &res)
	if err != nil {
		return nil, err
	}

	// convert to map
	m := make(map[int][]int)

	for _, r := range res.Response {
		m[r.AnswerID] = append(m[r.AnswerID], r.Users.Items...)
	}

	return m, nil
}

func (p *poolsClient) GetAllVoters(ctx context.Context, ownerID, poolID int, answerIDs []int) (map[int][]int, error) {
	var (
		offset int
		count  int = 1000

		res map[int][]int
	)

	for {
		m, err := p.GetVoters(ctx, ownerID, poolID, answerIDs, offset, count)
		if err != nil {
			return m, err
		}

		// isEmpty
		var mapHasValues bool
		for _, arr := range m {
			if len(arr) != 0 {
				mapHasValues = true

				break
			}
		}

		res = joinMaps(res, m)

		// если в ответе нет ID-шников пользователей, считаем что мы дошли до конца и возвращаем результат
		if !mapHasValues {
			return res, nil
		}

		offset += count
	}
}

type PoolResponse struct {
	Response Poll `json:"response"`
}

type Poll struct {
	Multiple      bool   `json:"multiple"`
	EndDate       int    `json:"end_date"`
	Closed        bool   `json:"closed"`
	IsBoard       bool   `json:"is_board"`
	CanEdit       bool   `json:"can_edit"`
	CanVote       bool   `json:"can_vote"`
	CanReport     bool   `json:"can_report"`
	CanShare      bool   `json:"can_share"`
	Created       int    `json:"created"`
	ID            int    `json:"id"`
	OwnerID       int    `json:"owner_id"`
	Question      string `json:"question"`
	Votes         int    `json:"votes"`
	DisableUnvote bool   `json:"disable_unvote"`
	Anonymous     bool   `json:"anonymous"`
	Friends       []struct {
		ID int `json:"id"`
	} `json:"friends"`
	AnswerIds []int  `json:"answer_ids"`
	EmbedHash string `json:"embed_hash"`
	Answers   []struct {
		ID    int     `json:"id"`
		Rate  float64 `json:"rate"`
		Text  string  `json:"text"`
		Votes int     `json:"votes"`
	} `json:"answers"`
	Background struct {
		Angle  int    `json:"angle"`
		Color  string `json:"color"`
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Points []struct {
			Color    string  `json:"color"`
			Position float64 `json:"position"`
		} `json:"points"`
		Type string `json:"type"`
	} `json:"background"`
}

type PollVoters struct {
	Response []struct {
		AnswerID int `json:"answer_id"`
		Users    struct {
			Count int   `json:"count"`
			Items []int `json:"items"`
		} `json:"users"`
	} `json:"response"`
}
