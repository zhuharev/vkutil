package vkutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhuharev/vk"
	//log "gopkg.in/inconshreveable/log15.v2"
)

var (
	DEBUG = false
)

type Api struct {
	vkApi *vk.Api
}

func NewUtils(api *vk.Api) *Api {
	return &Api{vkApi: api}
}

func (api *Api) GroupsGetAllMembers(groupId int) ([]int, error) {
	var result []int
	count, err := api.GroupsGetMembersCount(groupId)
	if err != nil {
		return []int{}, err
	}
	for i := 0; i < count; i = i + 1000 {
		resp, err := api.vkApi.Request(vk.METHOD_GROUPS_GET_MEMBERS, map[string][]string{
			"group_id": {fmt.Sprint(groupId)},
			"count":    {fmt.Sprint(1000)},
			"offset":   {fmt.Sprint(i)},
		})
		if err != nil {
			return []int{}, err
		}
		_, ids, err := ParseIdsResponse(resp)
		if err != nil {
			return []int{}, err
		}
		result = append(result, ids...)
	}
	return result, nil
}

func (api *Api) GroupsGetMembersCount(groupId int) (int, error) {
	resp, err := api.vkApi.Request(vk.METHOD_GROUPS_GET_MEMBERS, map[string][]string{
		"group_id": {fmt.Sprint(groupId)},
		"count":    {fmt.Sprint(0)},
	})
	if err != nil {
		return 0, err
	}
	r := ResponseIds{}
	err = json.Unmarshal([]byte(resp), &r)
	if err != nil {
		return 0, err
	}
	if r.Error.Code != 0 {
		return 0, errors.New(r.Error.Msg)
	}
	return r.Resp.Count, nil
}
