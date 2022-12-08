package vkutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/zhuharev/vk"
	"github.com/zhuharev/vkutil/structs"
)

var (
	DEBUG = false
)

type Api struct {
	VkApi      *vk.Api
	StdinAllow bool
	debug      bool
	*structs.API

	Pools *poolsClient
}

func New(opts ...Opt) *Api {
	va := new(vk.Api)
	a := &Api{
		VkApi: va, API: &structs.API{VkAPI: va},
	}

	for _, opt := range opts {
		opt(a)
	}

	a.Pools = newPoolsClient(a)

	return a
}

type Opt func(*Api)

func Token(token string) Opt {
	return func(a *Api) {
		a.VkApi.AccessToken = token
	}
}

func ClientCred(id int, secret string) Opt {
	return func(a *Api) {
		a.VkApi.ClientId = id
		a.VkAPI.ClientSecret = secret
	}
}

func NewWithToken(token string) *Api {
	api := New()
	api.VkApi.AccessToken = token
	return api
}

func (api *Api) AccessTokenURL(redirectURI, code string) string {
	return fmt.Sprintf(
		"https://oauth.vk.com/access_token?client_id=%d&client_secret=%s&redirect_uri=%s&code=%s",
		api.VkApi.ClientId,
		api.VkApi.ClientSecret,
		redirectURI,
		code,
	)
}

func (api *Api) UsersGetOne(id int, args ...url.Values) (*User, error) {
	users, err := api.UsersGet(id, args...)
	if err != nil {
		return nil, err
	}
	if len(users) != 1 {
		return nil, errors.New("unexpected vk response")
	}
	return &users[0], nil
}

func (api *Api) UserByCode(redirectURI, code string) (*User, error) {
	cli := api.VkApi.HTTPClient()
	resp, err := cli.Get(api.AccessTokenURL(redirectURI, code))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var at struct {
		AccessToken      string `json:"access_token"`
		UserID           int    `json:"user_id"`
		Error            string
		ErrorDescription string `json:"error_description"`
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if api.debug {
		log.Printf("response access token url: %s", bts)
	}

	err = json.Unmarshal(bts, &at)
	if err != nil {
		return nil, err
	}

	if at.Error != "" {
		return nil, fmt.Errorf("err get access token: %s (%s)", at.Error, at.ErrorDescription)
	}

	if at.UserID == 0 {
		return nil, fmt.Errorf("empty user id in response")
	}

	return api.UsersGetOne(at.UserID, url.Values{"fields": {"photo"}})
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

func (api *Api) RequestTyped(resp interface{}, method string, args ...url.Values) error {
	data, err := api.VkAPI.Request(method, args...)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, resp)
}
