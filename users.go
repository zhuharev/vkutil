package vkutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhuharev/vk"
	"net/url"
	"strconv"
	"strings"
)

type ResponseUsers struct {
	Response []User `json:"response"`
	ResponseError
}

//structures from https://github.com/shalakhin/vk/blob/master/users.go
//MIT license
type (
	User struct {
		Id                     int          `json:"id"`
		FirstName              string       `json:"first_name"`
		LastName               string       `json:"last_name"`
		ScreenName             string       `json:"screen_name"`
		NickName               string       `json:"nickname"`
		Sex                    int          `json:"sex,omitempty"`
		Domain                 string       `json:"domain,omitempty"`
		Birthdate              Bdate        `json:"bdate,omitempty"`
		City                   GeoPlace     `json:"city,omitempty"`
		Country                GeoPlace     `json:"country,omitempty"`
		PhotoId                string       `json:"photo_id,omitempty"`
		Photo50                string       `json:"photo_50,omitempty"`
		Photo100               string       `json:"photo_100,omitempty"`
		Photo200               string       `json:"photo_200,omitempty"`
		PhotoMax               string       `json:"photo_max,omitempty"`
		Photo200Orig           string       `json:"photo_200_orig,omitempty"`
		PhotoMaxOrig           string       `json:"photo_max_orig,omitempty"`
		HasMobile              Bool         `json:"has_mobile,omitempty"`
		Online                 Bool         `json:"online,omitempty"`
		CanPost                Bool         `json:"can_post,omitempty"`
		CanSeeAllPosts         Bool         `json:"can_see_all_posts,omitempty"`
		CanSeeAudio            Bool         `json:"can_see_audio,omitempty"`
		CanWritePrivateMessage Bool         `json:"can_write_private_message,omitempty"`
		Site                   string       `json:"site,omitempty"`
		Status                 string       `json:"status,omitempty"`
		LastSeen               PlatformInfo `json:"last_seen,omitempty"`
		CommonCount            int          `json:"common_count,omitempty"`
		University             int          `json:"university,omitempty"`
		UniversityName         string       `json:"university_name,omitempty"`
		Faculty                int          `json:"faculty,omitempty"`
		FacultyName            string       `json:"faculty_name,omitempty"`
		Graduation             int          `json:"graduation,omitempty"`
		Relation               Relation     `json:"relation,omitempty"`
		Universities           []University `json:"universities,omitempty"`
		Schools                []School     `json:"schools,omitempty"`
		Relatives              []Relative   `json:"relatives,omitempty"`
		Deactivated            string       `json:"deactivated,omitempty"`
		Deleted                string       `json:"deleted,omitempty"`
		Banned                 string       `json:"banned,omitempty"`
		Counters               Counters     `json:"counters,omitempty"`
	}

	Counters struct {
		Albums        int `json:"albums"`
		Videos        int `json:"videos"`
		Audios        int `json:"audios"`
		Notes         int `json:"notes"`
		Friends       int `json:"friends"`
		Groups        int `json:"groups"`
		OnlineFriends int `json:"online_friends"`
		MutualFriends int `json:"mutual_friends"`
		UserVideos    int `json:"user_videos"`
		Followers     int `json:"followers"`
		//fields Returned Only For Desktop Applications:
		UserPhotos    int `json:"user_photos"`
		Subscriptions int `json:"subscriptions"`
	}

	GeoPlace struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}
	// PlatformInfo contains information about time and platform
	PlatformInfo struct {
		Time     EpochTime `json:"time"`
		Platform int       `json:"platform"`
	}
	// University contains information about the university
	University struct {
		ID              int    `json:"id"`
		Country         int    `json:"country"`
		City            int    `json:"city"`
		Name            string `json:"name"`
		Faculty         int    `json:"faculty"`
		FacultyName     string `json:"faculty_name"`
		Chair           int    `json:"chair"`
		ChairName       string `json:"chair_name"`
		Graduation      int    `json:"graduation"`
		EducationForm   string `json:"education_form"`
		EducationStatus string `json:"education_status"`
	}
	// School contains information about schools
	School struct {
		ID         int    `json:"id"`
		Country    int    `json:"country"`
		City       int    `json:"city"`
		Name       string `json:"name"`
		YearFrom   int    `json:"year_from"`
		YearTo     int    `json:"year_to"`
		Class      string `json:"class"`
		TypeStr    string `json:"type_str,omitempty"`
		Speciality string `json:"speciality,omitempty"`
	}
	// Relative contains information about relative to the user
	Relative struct {
		ID   int    `json:"id"`   // negative id describes non-existing users (possibly prepared id if they will register)
		Type string `json:"type"` // like `parent`, `grandparent`, `sibling`
		Name string `json:"name,omitempty"`
	}

	Bool  bool
	Bdate string
)

func (b Bdate) getPart(part int) (int, bool) {
	arr := strings.Split(string(b), ".")
	if len(arr) <= part {
		return 0, false
	}
	num, e := strconv.Atoi(arr[part])
	if e != nil {
		return 0, false
	}
	return num, true
}

func (b Bdate) Year() (int, bool) {
	return b.getPart(2)
}

func (b Bdate) Day() (int, bool) {
	return b.getPart(0)
}

func (b Bdate) Mounth() (int, bool) {
	return b.getPart(1)
}

func (bit Bool) UnmarshalJSON(data []byte) error {
	asString := string(data)
	if asString == "1" || asString == "true" {
		bit = true
	} else if asString == "0" || asString == "false" {
		bit = false
	} else {
		return errors.New(fmt.Sprintf("Boolean unmarshal error: invalid input %s", asString))
	}
	return nil
}

// todo get all friends
func (api *Api) UsersGet(idsi interface{}, args ...url.Values) (res []User, e error) {
	var (
		ids []string
	)
	switch idsi.(type) {
	case []string:
		ids = idsi.([]string)
	case []int:
		ids = arrIntToStr(idsi.([]int))
	default:
		return nil, errors.New("uncnow type")
	}

	var lim = 1000
	stop := false
	for i := 0; i < len(ids) && !stop; i += lim {

		if len(ids)-1 < i+lim {
			lim = (len(ids) - 1) - i
			stop = true
		}
		users, e := api.usersGet1K(ids[i:i+lim], args...)
		if e != nil {
			return nil, e
		}
		res = append(res, users...)
	}

	return
}

func (api *Api) usersGet1K(ids []string, args ...url.Values) ([]User, error) {
	params := url.Values{}
	if len(args) == 1 {
		params = args[0]
	}
	if len(ids) > 1000 {
		ids = ids[:1000]
	}
	params.Set("user_ids", strings.Join(ids, ","))
	resp, err := api.vkApi.Request(vk.METHOD_USERS_GET, params)
	if err != nil {
		return nil, err
	}
	var r ResponseUsers
	err = json.Unmarshal(resp, &r)
	return r.Response, err
}

func (api *Api) UtilsUsersGet(ids []int) (users []User, e error) {
	/*	var tcode = `var in = ["1,2,3","4,5,6","7,8,9"];
		var res = [];
		var i  = 0;
		while(i<in.length) {
		res = res + API.users.get({user_ids:in[i]});
		i=i+1;
		}
		return res;`

			offset := 0
			count := 25
			for i := 0; i < len(ids); i = i + count {
				resp, err := api.Execute(fmt.Sprintf(tcode, strings.Join(arrIntToStr(ids), ","),
					offset, offset+count))
				var r struct {
					Response string `response`
				}
				err = json.Unmarshal(resp, &r)
				if err != nil {
					fmt.Println(string(resp))
					return
				}
				arr := arrStrToInt(strings.Split(r.Response, ","))
				for _, j := range arr {
					res = append(res, j)
				}
			}*/

	return
}

func SplitArr(arr []int, topLevel, botLevel int) (r [][][]int) {
	var (
		res [][]int
		lim = botLevel * topLevel
	)

	for c := 0; c < botLevel; c++ {

		for i, v := range arr[c:lim] {
			if i%(topLevel) == 0 {
				res = append(res, []int{})
			}
			res[len(res)-1] = append(res[len(res)-1], v)
		}
	}

	return
}

func (api *Api) UsersGetFollowers(args ...url.Values) ([]int, error) {
	resp, err := api.vkApi.Request(vk.METHOD_USERS_GET_FOLLOWERS, args...)
	if err != nil {
		return nil, err
	}
	var r ResponseIds
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, err
	}
	return r.Resp.Items, nil
}

func (api *Api) UsersSearch(q string, args ...url.Values) ([]User, error) {
	params := url.Values{}
	if len(args) == 1 {
		params = args[0]
	}
	params.Set("q", q)

	resp, err := api.vkApi.Request(vk.METHOD_USERS_SEARCH, params)
	if err != nil {
		return nil, err
	}
	var r ResponseUserWithCount
	err = json.Unmarshal(resp, &r)
	return r.Response.Items, err
}

// utils
func (api *Api) UtilsGetProfilePhoto(ids ...int) ([]User, error) {
	return api.UsersGet(arrIntToStr(ids), url.Values{
		"fields": {"photo_id"},
	})
}

func (api *Api) UtilsUsersGetDomains(id int) (string, error) {
	resp, err := api.UsersGet(arrIntToStr([]int{id}), url.Values{
		"fields": {"domain"},
	})
	if err != nil {
		return "", err
	}
	if len(resp) == 1 {
		return resp[0].Domain, nil
	}
	return "", errors.New("Unknown error")
}

func (api *Api) UtilsUsersGetId(domain string) (int, error) {
	resp, err := api.UsersGet([]string{domain})
	if err != nil {
		return 0, err
	}
	if len(resp) == 1 {
		return resp[0].Id, nil
	}
	return 0, errors.New("Unknown error")
}
