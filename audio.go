package vkutil

import (
	"encoding/json"
	"fmt"
	"github.com/zhuharev/vk"
	"log"
	"net/url"
)

type Audio struct {
	ID                int    `json:"id"`
	OwnerID           int    `json:"owner_id"`
	Artist            string `json:"artist"`
	Title             string `json:"title"`
	Duration          int    `json:"duration"`
	Date              int    `json:"date"`
	URL               string `json:"url"`
	GenreID           int    `json:"genre_id,omitempty"`
	LyricsID          int    `json:"lyrics_id,omitempty"`
	NoSearch          int    `json:"no_search,omitempty"`
	ContentRestricted int    `json:"content_restricted,omitempty"`
	AlbumID           int    `json:"album_id,omitempty"`
}

type AudioResponse struct {
	Response struct {
		Count int     `json:"count"`
		Items []Audio `json:"items"`
	} `json:"response"`
}

func (a *Api) AudioGet(ownerIds ...int) ([]Audio, error) {
	var params = url.Values{}
	if len(ownerIds) > 0 {
		params.Set("owner_id", fmt.Sprint(ownerIds[0]))
	}
	bts, err := a.VkApi.Request(vk.METHOD_AUDIO_GET, params)
	if err != nil {
		return nil, err
	}

	var r AudioResponse

	err = json.Unmarshal(bts, &r)
	if err != nil {
		return nil, err
	}

	return r.Response.Items, nil
}

func (a *Api) AudioSetBroadcast(ownerId, audioId int) error {
	bts, err := a.VkApi.Request(vk.METHOD_AUDIO_SET_BROADCAST, url.Values{
		"audio": {fmt.Sprintf("%d_%d", ownerId, audioId)},
	})
	if err != nil {
		log.Println(err)
	}
	log.Println(string(bts))
	return nil
}
