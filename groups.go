package vkutil

import (
	"encoding/json"
	"fmt"
	"github.com/zhuharev/vk"
	"log"
	"net/url"
)

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
			Response []int `response`
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
