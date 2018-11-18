package vkutil

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/zhuharev/vk"
)

type Group struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	IsClosed   int    `json:"is_closed"`
	Type       string `json:"type"`
	IsAdmin    int    `json:"is_admin"`
	IsMember   int    `json:"is_member"`
	City       struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"city"`
	Country struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"country"`
	Place struct {
		ID         int     `json:"id"`
		Title      string  `json:"title"`
		Latitude   float64 `json:"latitude"`
		Longitude  float64 `json:"longitude"`
		Created    int     `json:"created"`
		Icon       string  `json:"icon"`
		GroupID    int     `json:"group_id"`
		GroupPhoto string  `json:"group_photo"`
		Checkins   int     `json:"checkins"`
		Updated    int     `json:"updated"`
		Type       int     `json:"type"`
		Country    int     `json:"country"`
		City       int     `json:"city"`
	} `json:"place"`
	Description  string `json:"description"`
	WikiPage     string `json:"wiki_page"`
	MembersCount int    `json:"members_count"`
	Counters     struct {
		Topics int `json:"topics"`
		Videos int `json:"videos"`
	} `json:"counters"`
	CanPost        int    `json:"can_post"`
	CanSeeAllPosts int    `json:"can_see_all_posts"`
	Activity       string `json:"activity"`
	Status         string `json:"status"`
	Contacts       []struct {
		UserID int    `json:"user_id"`
		Desc   string `json:"desc"`
	} `json:"contacts"`
	Links []struct {
		ID        int    `json:"id"`
		URL       string `json:"url"`
		Name      string `json:"name"`
		Desc      string `json:"desc"`
		Photo50   string `json:"photo_50"`
		Photo100  string `json:"photo_100"`
		EditTitle int    `json:"edit_title,omitempty"`
	} `json:"links"`
	Verified int    `json:"verified"`
	Site     string `json:"site"`
	Cover    struct {
		Enabled int `json:"enabled"`
		Images  []struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"images"`
	} `json:"cover"`
	Photo50  string `json:"photo_50"`
	Photo100 string `json:"photo_100"`
	Photo200 string `json:"photo_200"`
}

func (g Group) Closed() bool {
	return g.IsClosed != 0
}

func (api *Api) GroupsGet(args ...url.Values) ([]int, error) {
	params := url.Values{}
	if len(args) == 1 {
		params = args[0]
	}
	r, err := api.VkApi.Request(vk.METHOD_GROUPS_GET, params)
	if err != nil {
		return nil, err
	}
	_, ids, err := ParseIdsResponse(r)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// GroupsGetByID get group by id
func (api *Api) GroupsGetByID(ids interface{}, args ...url.Values) ([]Group, error) {
	params := url.Values{}
	if len(args) == 1 {
		params = args[0]
	}
	params = setToUrlValues("group_ids", ids, params)
	r, err := api.VkApi.Request(vk.METHOD_GROUPS_GET_BY_ID, params)
	if err != nil {
		return nil, err
	}
	type GroupsResponse struct {
		Groups []Group `json:"response"`
	}
	var gr GroupsResponse
	err = json.Unmarshal(r, &gr)
	if err != nil {
		return nil, err
	}
	return gr.Groups, nil
}

func (api *Api) GroupsJoin(groupId int, args ...url.Values) error {
	params := url.Values{}
	if len(args) == 1 {
		params = args[0]
	}
	params.Set("group_id", fmt.Sprint(groupId))
	r, err := api.VkApi.Request(vk.METHOD_GROUPS_JOIN, params)
	if err != nil {
		return err
	}
	_, e := ParseIntResponse(r)
	if e != nil {
		return e
	}
	return nil
}

func (api *Api) GroupsLeave(groupId int, args ...url.Values) error {
	params := url.Values{}
	if len(args) == 1 {
		params = args[0]
	}
	params.Set("group_id", fmt.Sprint(groupId))
	r, err := api.VkApi.Request(vk.METHOD_GROUPS_LEAVE, params)
	if err != nil {
		return err
	}
	_, e := ParseIntResponse(r)
	if e != nil {
		return e
	}
	return nil
}

func (api *Api) GroupsGetMembers(gid int, args ...url.Values) (count int, ids []int,
	err error) {
	params := url.Values{}
	if len(args) == 1 {
		params = args[0]
	}
	params.Set("group_id", fmt.Sprint(gid))
	r, err := api.VkApi.Request(vk.METHOD_GROUPS_GET_MEMBERS, params)
	if err != nil {
		return
	}
	count, ids, err = ParseIdsResponse(r)
	return
}

func (api *Api) GroupsGetAllMembers(gid int) ([]int, error) {
	var members []int = []int{}
	count, _, err := api.GroupsGetMembers(gid)
	if err != nil {
		return nil, err
	}
	var fmtcode = `
	var limiter = 0;
	var cnt = %[1]d;
	var c = %[3]d;
	var members = API.groups.getMembers({"group_id": %[2]d,"v":"5.32","offset" : %[3]d, "count": "1000"}).items;
var  offset = %[3]d+1000;
while (limiter < 24 && c < cnt)
{
	members = members + API.groups.getMembers({"group_id": %[2]d, "v": "5.32", "count": "1000", "offset": offset}).items;
	offset = offset + 1000;
	c = c + 1000;
	limiter = limiter+1;
};
return members;
`
	//offset := 0
	//limit := 1000
	for len(members) < count {
		code := fmt.Sprintf(fmtcode, count, gid, len(members))
		resp, err := api.Execute(code)
		if err != nil {
			log.Println(err)
		}
		var r struct {
			Response []int `json:"response"`
		}
		err = json.Unmarshal(resp, &r)
		if err != nil {
			fmt.Println(string(resp))
			return nil, err
		}
		//arr := arrStrToInt(strings.Split(r.Response, ","))
		//for _, j := range arr {
		//	if j != 0 {
		members = append(members, r.Response...)

		//	}
		//}
	}
	if len(members) > count {
		members = members[:count]
	}
	return members, nil
}

/*
function getMembers20k(group_id, members_count) {
    var code =  'var members = API.groups.getMembers({"group_id": ' + group_id + ', "v": "5.27", "sort": "id_asc", "count": "1000", "offset": ' + membersGroups.length + '}).items;' // делаем первый запрос и создаем массив
            +	'var offset = 1000;' // это сдвиг по участникам группы
            +	'while (offset < 25000 && (offset + ' + membersGroups.length + ') < ' + members_count + ')' // пока не получили 20000 и не прошлись по всем участникам
            +	'{'
                +	'members = members + "," + API.groups.getMembers({"group_id": ' + group_id + ', "v": "5.27", "sort": "id_asc", "count": "1000", "offset": (' + membersGroups.length + ' + offset)}).items;' // сдвиг участников на offset + мощность массива
                +	'offset = offset + 1000;' // увеличиваем сдвиг на 1000
            +	'};'
            +	'return members;'; // вернуть массив members

    VK.Api.call("execute", {code: code}, function(data) {
        if (data.response) {
            membersGroups = membersGroups.concat(JSON.parse("[" + data.response + "]")); // запишем это в массив
            $('.member_ids').html('Загрузка: ' + membersGroups.length + '/' + members_count);
            if (members_count >  membersGroups.length) // если еще не всех участников получили
                setTimeout(function() { getMembers20k(group_id, members_count); }, 333); // задержка 0.333 с. после чего запустим еще раз
            else // если конец то
                alert('Ура тест закончен! В массиве membersGroups теперь ' + membersGroups.length + ' элементов.');
        } else {
            alert(data.error.error_msg); // в случае ошибки выведем её
        }
    });
*/

func (a *Api) GroupsSearch(q string, params ...url.Values) ([]Group, error) {
	param := setToUrlValues("q", q, params...)
	bts, err := a.VkApi.Request(vk.METHOD_GROUPS_SEARCH, param)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Resp struct {
			Count int     `json:"count"`
			Items []Group `json:"items"`
		} `json:"response"`
	}

	var r Response

	err = json.Unmarshal(bts, &r)
	if err != nil {
		return nil, err
	}

	return r.Resp.Items, nil
}

func (api *Api) GroupIsMember(publicID, userID int) (bool, error) {
	param := setToUrlValues("group_id", publicID)
	param = setToUrlValues("user_id", userID, param)
	var r ResponseInt
	err := api.RequestTyped(&r, vk.METHOD_GROUPS_IS_MEMBER, param)
	if err != nil {
		return false, err
	}
	return r.Response == 1, nil
}
