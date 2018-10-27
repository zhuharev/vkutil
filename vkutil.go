package vkutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/zhuharev/vk"
	"github.com/zhuharev/vkutil/structs"
	//log "gopkg.in/inconshreveable/log15.v2"
)

var (
	DEBUG = false
)

type Api struct {
	VkApi      *vk.Api
	StdinAllow bool
	debug      bool
	*structs.API
}

func New() *Api {
	va := new(vk.Api)
	return &Api{VkApi: va, API: &structs.API{VkAPI: va}}
}

func NewWithToken(token string) *Api {
	api := New()
	api.VkApi.AccessToken = token
	return api
}

// func NewWithTokenFile(tokenFilePath string) (*Api, error) {
// 	token := ""
//
// 	f, e := os.OpenFile(tokenFilePath, os.O_RDONLY, 0777)
// 	if e == nil {
// 		bts, e := ioutil.ReadAll(f)
// 		if e != nil {
// 			return nil, e
// 		}
// 		f.Close()
// 		token = string(bts)
// 	} else {
// 		return nil, e
// 	}
//
// 	return NewWithToken(token), nil
// }
//
// func NewUtils(api *vk.Api) *Api {
// 	return &Api{VkApi: api}
// }

// SetDebug enable loggin and set log flags longfile
func (api *Api) SetDebug(debug bool) {
	if debug {
		log.SetFlags(log.LstdFlags | log.Llongfile)
	}
	api.VkApi.SetDebug(debug)
	api.debug = debug
}

/*func (api *Api) GroupsGetAllMembers(groupId int) ([]int, error) {
	var result []int
	count, err := api.GroupsGetMembersCount(groupId)
	if err != nil {
		return []int{}, err
	}
	for i := 0; i < count; i = i + 1000 {
		resp, err := api.VkApi.Request(vk.METHOD_GROUPS_GET_MEMBERS, map[string][]string{
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
}*/

func (api *Api) GroupsGetMembersCount(groupId int) (int, error) {
	resp, err := api.VkApi.Request(vk.METHOD_GROUPS_GET_MEMBERS, map[string][]string{
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
