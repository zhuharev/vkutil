package vkutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/zhuharev/vk"
)

type ResponseUsers struct {
	Response []User `json:"response"`
	ResponseError
}

// structures from https://github.com/shalakhin/vk/blob/master/users.go
// MIT license
type (
	User struct {
		Id                     int          `json:"id" xorm:"id pk"`
		FirstName              string       `json:"first_name"`
		LastName               string       `json:"last_name"`
		ScreenName             string       `json:"screen_name"`
		NickName               string       `json:"nickname"`
		Sex                    int          `json:"sex,omitempty"`
		Domain                 string       `json:"domain,omitempty"`
		Birthdate              Bdate        `json:"bdate,omitempty"`
		City                   GeoPlace     `json:"city,omitempty" xorm:"-"`
		Country                GeoPlace     `json:"country,omitempty" xorm:"-"`
		PhotoId                string       `json:"photo_id,omitempty"`
		Photo                  string       `json:"photo,omitempty"`
		Photo50                string       `json:"photo_50,omitempty"`
		Photo100               string       `json:"photo_100,omitempty"`
		Photo200               string       `json:"photo_200,omitempty"`
		PhotoMax               string       `json:"photo_max,omitempty"`
		Photo200Orig           string       `json:"photo_200_orig,omitempty"`
		PhotoMaxOrig           string       `json:"photo_max_orig,omitempty"`
		HasMobile              Bool         `json:"has_mobile,omitempty" xorm:"-"`
		Online                 Bool         `json:"online,omitempty" xorm:"-"`
		CanPost                Bool         `json:"can_post,omitempty" xorm:"-"`
		CanSeeAllPosts         Bool         `json:"can_see_all_posts,omitempty" xorm:"-"`
		CanSeeAudio            Bool         `json:"can_see_audio,omitempty" xorm:"-"`
		CanWritePrivateMessage Bool         `json:"can_write_private_message,omitempty" xorm:"-"`
		Site                   string       `json:"site,omitempty"`
		Status                 string       `json:"status,omitempty"`
		LastSeen               PlatformInfo `json:"last_seen,omitempty" xorm:"-"`
		CommonCount            int          `json:"common_count,omitempty"`
		University             int          `json:"university,omitempty"`
		UniversityName         string       `json:"university_name,omitempty"`
		Faculty                int          `json:"faculty,omitempty"`
		FacultyName            string       `json:"faculty_name,omitempty"`
		FollowersCount         int          `json:"followers_count"`
		Graduation             int          `json:"graduation,omitempty"`
		Relation               Relation     `json:"relation,omitempty"`
		Universities           []University `json:"universities,omitempty"`
		Schools                []School     `json:"schools,omitempty"`
		Relatives              []Relative   `json:"relatives,omitempty"`
		Deactivated            string       `json:"deactivated,omitempty"`
		Deleted                string       `json:"deleted,omitempty"`
		Banned                 string       `json:"banned,omitempty"`
		Counters               Counters     `json:"counters,omitempty" xorm:"-"`

		Instagram string `json:"instagram"`
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
		ID    int    `json:"id" xorm:"city_id"`
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

func (bit Bool) FromDB(data []byte) error {
	if string(data) == "1" {
		bit = true
	}
	return nil
}

func (u User) Field(fieldName string) interface{} {
	switch fieldName {
	case "first_name":
		return u.FirstName
	case "last_name":
		return u.LastName
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
	case int:
		ids = []string{fmt.Sprint(idsi.(int))}
	case map[int]struct{}:
		for k := range idsi.(map[int]struct{}) {
			ids = append(ids, fmt.Sprint(k))
		}
	default:
		ids = []string{}
	}

	ids = uniqStrArr(ids)

	if len(ids) == 0 {
		return api.usersGet1K(ids, args...)
	}

	var lim = 1000
	stop := false
	for i := 0; i < len(ids) && !stop; i += lim {

		if len(ids)-1 < i+lim {
			lim = (len(ids)) - i
			stop = true
		}
		//color.Cyan("%v %v", ids[i:i+lim], ids)
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
	if len(ids) > 0 {
		params.Set("user_ids", strings.Join(ids, ","))
	}
	resp, err := api.VkApi.Request(vk.METHOD_USERS_GET, params)
	if err != nil {
		return nil, err
	}
	var r ResponseUsers
	err = json.Unmarshal(resp, &r)
	if err != nil {
		fmt.Println(string(resp))
		return nil, err
	}
	return r.Response, err
}

//func (api *Api) UtilsUsersGet(ids []int) (users []User, e error) {
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

//	return
//}

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
	resp, err := api.VkApi.Request(vk.METHOD_USERS_GET_FOLLOWERS, args...)
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

func (api *Api) UsersSearch(q string, args ...url.Values) ([]User, int, error) {
	params := url.Values{}
	if len(args) == 1 {
		params = args[0]
	}
	params.Set("q", q)

	resp, err := api.VkApi.Request(vk.METHOD_USERS_SEARCH, params)
	if err != nil {
		return nil, 0, err
	}
	var r ResponseUserWithCount
	err = json.Unmarshal(resp, &r)
	return r.Response.Items, r.Response.Count, err
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

func (api *Api) UtilsUsersField(ids []int, filedname string) (map[int]interface{}, error) {
	ids = arrUniq(ids)

	var (
		m = map[int]interface{}{}
	)

	mc, done, errors := api.GoUtilsUsersField(ids, filedname)

	for {
		select {
		case nm := <-mc:
			color.Cyan("GET")
			for id, field := range nm {
				m[id] = field
			}
		case <-done:
			return m, nil
		case err := <-errors:
			color.Red("%s", err)
		}
	}
}

func (api *Api) GoUtilsUsersField(ids []int, filedname string) (chan map[int]interface{}, chan struct{}, chan error) {
	var (
		cnt   int
		ch    = make(chan map[int]interface{})
		done  = make(chan struct{})
		dones = make(chan struct{})
		errs  = make(chan error)

		jobs = make(chan []int, 500)
	)

	var run = func(jobs chan []int) {
		for arr := range jobs {

			m, e := api.usersField25K(arr, filedname)
			if e != nil {
				errs <- e
				dones <- struct{}{}
				continue
			}
			if len(m) != len(arr) {
				errs <- fmt.Errorf("[users.go:399] Len mishmatch got %d, need %d", len(m), len(arr))
				dones <- struct{}{}
				continue
			}
			ch <- m
			dones <- struct{}{}
		}
	}

	for i := 0; i < 15; i++ {
		go run(jobs)
	}

	{
		var currArr []int
		for _, id := range ids {

			if len(currArr) == 10000 {
				cnt++
				jobs <- currArr
				currArr = []int{}
			} else {
				currArr = append(currArr, id)
			}
		}

		if len(currArr) > 0 {
			cnt++
			jobs <- currArr
		}
	}

	go func(done chan struct{}) {
		for i := 0; i < cnt; i++ {
			<-dones
		}
		done <- struct{}{}
	}(done)

	return ch, done, nil
}

func (api *Api) usersField25K(ids []int, filedname string) (map[int]interface{}, error) {
	start := time.Now()
	ids = arrUniq(ids)
	color.Green("Got %d", len(ids))
	if len(ids) > 25000 {
		ids = ids[:25000]
	}
	firstK := ids
	if len(firstK) > 700 {
		firstK = firstK[:700]
	}
	head := fmt.Sprintf(`var a = API.users.get({user_ids:"%s",fields:"`+filedname+`"})@.`+filedname+`;`, strings.Join(arrIntToStr(firstK), ","))

	for _, arr := range arrSplit1K(ids[len(firstK):]) {
		head += fmt.Sprintf(`a=a+API.users.get({user_ids:"%s",fields:"`+filedname+`"})@.`+filedname+`;`, strings.Join(arrIntToStr(arr), ","))
	}
	head += "return a;"
	b, e := api.Execute(head)
	if e != nil {
		fmt.Println(string(b))
		return nil, e
	}

	type Autogenerated struct {
		Response []interface{} `json:"response"`
		ResponseError
	}
	var ag Autogenerated

	res := map[int]interface{}{}

	e = json.Unmarshal(b, &ag)
	if e != nil {
		color.Red("%s", b)
		return nil, e
	}
	if ag.Error.Code != 0 {
		return nil, fmt.Errorf("%s", ag.Error.Msg)
	}
	for i, v := range ag.Response {
		res[ids[i]] = v
	}

	if len(ids) != len(res) {
		color.Red("[473] Len mishmatch got %d, need %d", len(res), len(ids))
		users, e := api.UsersGet(ids, url.Values{"fields": {filedname}})
		if e != nil {
			color.Red("[484] %s", e)

		}
		d := map[int]struct{}{}
		for _, v := range users {
			res[v.Id] = v.Field(filedname)
			d[v.Id] = struct{}{}
		}
		var lose = []int{}
		for _, id := range ids {
			if _, has := d[id]; !has {
				lose = append(lose, id)
			}
		}

		/* fix losed users */
		users, e = api.UsersGet(lose, url.Values{"fields": {filedname}})
		if e != nil {
			color.Red("[484] %s", e)

		}
		for _, v := range users {
			res[v.Id] = v.Field(filedname)
			d[v.Id] = struct{}{}
		}

		lose = []int{}
		for _, id := range ids {
			if _, has := d[id]; !has {
				lose = append(lose, id)
			}
		}

		/* fix losed users */
		users, e = api.UsersGet(lose, url.Values{"fields": {filedname}})
		if e != nil {
			color.Red("[484] %s", e)

		}
		for _, v := range users {
			res[v.Id] = v.Field(filedname)
			d[v.Id] = struct{}{}
		}

		color.Green("Fixed [520] got %d, has %d", len(ids), len(res))
		//return nil, ErrArrayLengthMismatch
	}

	//if e != nil {

	//}
	color.Green("Getted %d for %s", len(ids), time.Since(start))
	return res, nil
}

func (api *Api) UtilsUsersGet(ids []int, fields []string) ([]User, error) {
	var (
		m []User
	)

	mc, done, _ := api.GoUtilsUsersGet(ids, fields)

	for {
		select {
		case nm := <-mc:
			m = append(m, nm...)
		case <-done:
			return m, nil
		}
	}
}

func (api *Api) GoUtilsUsersGet(ids []int, fields []string) (chan []User, chan struct{}, error) {
	var (
		wg    sync.WaitGroup
		cnt   int
		ch    = make(chan []User)
		done  = make(chan struct{})
		dones = make(chan struct{})
	)
	var currArr []int
	for _, id := range ids {
		if len(currArr) == 10000 {
			wg.Add(1)
			cnt++
			go func(ch chan []User, arr []int) {
				m, e := api.users25K(arr, fields)
				if e != nil {
					color.Red("%s", e)
				}
				ch <- m
				dones <- struct{}{}
			}(ch, currArr)
			currArr = []int{}
		} else {
			currArr = append(currArr, id)
		}
	}

	if len(currArr) > 0 {
		cnt++
		wg.Add(1)
		go func(ch chan []User, arr []int) {
			m, e := api.users25K(arr, fields)
			if e != nil {
				color.Red("%s", e)
			}
			ch <- m
			dones <- struct{}{}
			wg.Done()
		}(ch, currArr)
	}

	go func(done chan struct{}) {
		for i := 0; i < cnt; i++ {
			<-dones
		}
		done <- struct{}{}
	}(done)

	return ch, done, nil
}

func (api *Api) users25K(ids []int, fields []string) ([]User, error) {
	if len(ids) > 25000 {
		ids = ids[:25000]
	}
	firstK := ids
	if len(firstK) > 700 {
		firstK = firstK[:700]
	}
	head := fmt.Sprintf(`var a = API.users.get({user_ids:"%s",fields:"`+strings.Join(fields, ",")+`"});`, strings.Join(arrIntToStr(firstK), ","))

	for _, arr := range arrSplit1K(ids[len(firstK):]) {
		head += fmt.Sprintf(`a=a+API.users.get({user_ids:"%s",fields:"`+strings.Join(fields, ",")+`"});`, strings.Join(arrIntToStr(arr), ","))
	}
	head += "return a;"
	b, e := api.Execute(head)
	if e != nil {
		color.Red("%s", b)
		return nil, e
	}

	type Autogenerated struct {
		Response []User `json:"response"`
		ResponseError
	}
	var ag Autogenerated

	e = json.Unmarshal(b, &ag)
	if e != nil {
		color.Red("%s", b)
		return nil, e
	}
	if ag.Error.Code != 0 {
		return nil, fmt.Errorf("%s", ag.Error.Msg)
	}

	//if e != nil {

	//}
	return ag.Response, nil
}

func (api *Api) UsersGetSubscriptions() (users []int, groups []int, e error) {
	type SubscriptionsResponse struct {
		Response struct {
			Users struct {
				Count int   `json:"count"`
				Items []int `json:"items"`
			} `json:"users"`
			Groups struct {
				Count int   `json:"count"`
				Items []int `json:"items"`
			} `json:"groups"`
		} `json:"response"`
	}
	var r SubscriptionsResponse
	bts, e := api.VkApi.Request(vk.METHOD_USERS_GET_SUBSCRIPTIONS, url.Values{"count": {"200"}})
	if e != nil {
		return
	}
	e = json.Unmarshal(bts, &r)
	if e != nil {
		return
	}
	return r.Response.Users.Items, r.Response.Groups.Items, nil
}
