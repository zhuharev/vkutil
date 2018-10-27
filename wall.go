package vkutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/fatih/color"
	vk "github.com/zhuharev/vk"
	"github.com/zhuharev/vkutil/wall"
)

func (api *Api) WallGet(ownerId int, filter ...url.Values) ([]*Post, error) {
	vals := url.Values{}
	if len(filter) == 1 {
		vals = filter[0]
	}
	vals.Set("owner_id", fmt.Sprint(ownerId))
	resp, err := api.VkApi.Request(vk.METHOD_WALL_GET, vals)
	if err != nil {
		return nil, err
	}
	var r ResponsePosts
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, err
	}
	if r.Error.Msg != "" {
		return nil, errors.New(r.Error.Msg)
	}
	return r.Response.Items, nil
}

func (api *Api) UtilsWallPostCount(ownerId int) (count int, e error) {
	fcode := `return API.wall.get({owner_id:%d,count:1}).count;`
	bts, e := api.Execute(fcode)
	count, e = ParseIntResponse(bts)
	return
}

func (api *Api) GoUtilsWallGet(ownerId int) (resp chan []*Post, done chan struct{}, errs chan error) {

	var (
		//	postCount int
		posts []*Post
		e     error
	)

	resp = make(chan []*Post)
	done = make(chan struct{})
	errs = make(chan error)

	posts, e = api.WallGet(ownerId, url.Values{"count": {"100"}})
	if e != nil {
		errs <- e
		done <- struct{}{}
		return
	}

	resp <- posts

	// return if post count < requested
	if len(posts) != 200 {
		done <- struct{}{}
		return
	}

	return
}

func (api *Api) UtilsWallGetAll(ownerId int) (posts []*Post, e error) {

	var (
		offset = 0
	)

	for {
		tposts, e := api.utilsWallGetTwoThousandPost(ownerId, offset)
		if e != nil {
			return nil, e
		}

		posts = append(posts, tposts...)

		color.Green("LEN %d", len(tposts))

		if len(tposts) < 2000 {
			break
		}

		offset += 2000
	}

	return

}

func (a *Api) utilsWallGetTwoThousandPost(ownerId int, offsets ...int) ([]*Post, error) {

	var (
		offset = 0
	)

	if len(offsets) > 0 {
		offset = offsets[0]
	}

	fcode := `var cnt = 1;
	var offs = %d;
	var own = %d;
var a = API.wall.get({owner_id:own,count:100,offset:offs,filter:"all"}).items;
while(cnt<20) {
a = a+ API.wall.get({owner_id:own,count:100,offset:offs+(100*cnt),filter:"all"}).items;
cnt = cnt+1;
}
return a;`

	type Resp struct {
		Items []*Post `json:"response"`
		ResponseError
	}

	var r Resp

	bts, e := a.Execute(fmt.Sprintf(fcode, offset, ownerId))
	if e != nil {
		return nil, e
	}
	e = json.Unmarshal(bts, &r)
	if e != nil {
		color.Yellow("%s", bts)
		return nil, e
	}
	if r.Error.Code != 0 {
		return nil, fmt.Errorf(r.Error.Msg)
	}
	return r.Items, nil
}

func (a *Api) UtilsGetTwoThousandComments(ownerId int, postId int, offsets ...int) ([]*Comment, error) {

	var (
		offset = 0
	)

	if len(offsets) > 0 {
		offset = offsets[0]
	}

	fcode := `var cnt = 1;
	var offs = %d;
	var own = %d;
	var post_id = %d;
var a = API.wall.getComments({post_id:post_id,owner_id:own,count:100,offset:offs}).items;
while(cnt<25) {
a = a+ API.wall.getComments({post_id:post_id,owner_id:own,count:100,offset:offs+(100*cnt)}).items;
cnt = cnt+1;
}
return a;`

	var r ResponseCommentsList

	bts, e := a.Execute(fmt.Sprintf(fcode, offset, ownerId, postId))
	if e != nil {
		return nil, e
	}
	e = json.Unmarshal(bts, &r)
	if e != nil {
		color.Yellow("%s", bts)
		return nil, e
	}
	if r.Error.Code != 0 {
		return nil, fmt.Errorf(r.Error.Msg)
	}
	return r.Response, nil
}

func (api *Api) uploadAttachments(groupID int, imageURLs []string) ([]string, error) {

	if groupID < 0 {
		groupID = -groupID
	}

	var (
		resAttachments []string

		attachments []io.Reader
	)

	for _, imgURL := range imageURLs {
		resp, err := http.Get(imgURL)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, resp.Body)
	}

	for _, img := range attachments {
		photo, err := api.PhotosUploadWall(groupID, img)
		if err != nil {
			log.Println(err)
			continue
		}

		// close if http connection or opened file
		if closer, ok := img.(io.ReadCloser); ok {
			closer.Close()
		}

		resAttachments = append(resAttachments, fmt.Sprintf("photo%d_%d", photo.OwnerId, photo.Id))
	}

	return resAttachments, nil
}

type OptsWallPost struct {
	OwnerId   int
	Message   string
	FromGroup bool
	ImageURLs []string
}

func (api *Api) WallPost(opts OptsWallPost, filter ...url.Values) (int, error) {
	vals := url.Values{}
	if len(filter) == 1 {
		vals = filter[0]
	}
	vals.Set("owner_id", fmt.Sprint(opts.OwnerId))
	vals.Set("message", opts.Message)
	if opts.FromGroup {
		vals.Set("from_group", "1")
	}

	if len(opts.ImageURLs) != 0 {
		attachments, err := api.uploadAttachments(opts.OwnerId, opts.ImageURLs)
		if err != nil {
			return 0, err
		}
		vals.Set("attachments", strings.Join(attachments, ","))
	}

	resp, err := api.VkApi.Request(vk.METHOD_WALL_POST, vals)
	if err != nil {
		return 0, err
	}

	type rsp struct {
		Response struct {
			PostId int `json:"post_id"`
		} `json:"response"`
		ResponseError
	}

	var r = rsp{}

	err = json.Unmarshal(resp, &r)
	if err != nil {
		return 0, err
	}
	if r.Error.Msg != "" {
		return 0, errors.New(r.Error.Msg)
	}
	return r.Response.PostId, nil
}

func (api *Api) WallRepost(ot ObjectType, ownerId int, objectId int, params ...url.Values) (int, error) {
	var (
		param = url.Values{}
		sId   = fmt.Sprintf("%s%d_%d", ot, ownerId, objectId)
	)
	if params != nil {
		param = params[0]
	}
	param.Set("object", sId)
	bts, e := api.VkApi.Request(vk.METHOD_WALL_REPOST, param)
	if e != nil {
		return 0, e
	}
	type wallPostResponse struct {
		Success int `json:"success"`
		PostId  int `json:"post_id"`
	}
	var wpr wallPostResponse
	e = json.Unmarshal(bts, &wpr)
	if e != nil {
		return 0, e
	}
	return wpr.PostId, nil
}

func (api *Api) WallDelete(ownerId int, postId int, params ...url.Values) error {
	var (
		param = url.Values{}
	)
	if params != nil {
		param = params[0]
	}
	param.Set("owner_id", fmt.Sprint(ownerId))
	param.Set("post_id", fmt.Sprint(postId))
	bts, e := api.VkApi.Request(vk.METHOD_WALL_DELETE, param)
	if e != nil {
		return e
	}
	_, e = ParseIntResponse(bts)
	if e != nil {
		return e
	}
	return nil
}

// ResponseRepost represent repost response
type ResponseRepost struct {
	Items    []Post  `json:"items"`
	Profiles []User  `json:"profiles"`
	Groups   []Group `json:"groups"`
}

// WallGetRepostsIds return ids about wall reposts
func (api *Api) WallGetRepostsIds(ownerID int, postID int, params ...url.Values) (rr ResponseRepost, err error) {
	type RRWrap struct {
		Resp ResponseRepost `json:"response"`
	}
	var rrw RRWrap
	param := setToUrlValues("owner_id", ownerID, params...)
	param = setToUrlValues("post_id", postID, param)
	bts, err := api.VkApi.Request(vk.METHOD_WALL_GET_REPOSTS, param)
	if err != nil {
		return
	}
	err = json.Unmarshal(bts, &rrw)
	if err != nil {
		return
	}
	return rrw.Resp, nil
}

// TODO
// WallGetReposts return info about wall reposts
func (api *Api) WallGetReposts(ownerID int, postID int, params ...url.Values) (rr ResponseRepost, err error) {
	type RRWrap struct {
		Resp ResponseRepost `json:"response"`
	}
	var rrw RRWrap
	param := setToUrlValues("owner_id", ownerID, params...)
	param = setToUrlValues("post_id", postID, param)
	bts, err := api.VkApi.Request(vk.METHOD_WALL_GET_REPOSTS, param)
	if err != nil {
		return
	}
	err = json.Unmarshal(bts, &rrw)
	if err != nil {
		return
	}
	return rrw.Resp, nil
}

func (api *Api) WallGetByID(postID string) (p Post, err error) {
	type Resp struct {
		Response []Post `json:"response"`
	}
	var r Resp

	param := setToUrlValues("posts", postID)
	bts, err := api.VkApi.Request(vk.METHOD_WALL_GET_BY_ID, param)
	if err != nil {
		return
	}
	err = json.Unmarshal(bts, &r)
	if err != nil {
		return
	}
	if len(r.Response) != 1 {
		return Post{}, fmt.Errorf("not found")
	}
	return r.Response[0], nil
}

func (api *Api) WallCreateComment(ownerID, postID int, fromGroup bool,
	message string, replyToComment int, attachments []string, params ...url.Values) (int, error) {
	param := setToUrlValues("owner_id", ownerID)
	param = setToUrlValues("post_id", postID, param)
	param = setToUrlValues("from_group", fromGroup, param)
	param = setToUrlValues("message", message, param)
	param = setToUrlValues("reply_to_comment", replyToComment, param)
	param = setToUrlValues("attachments", strings.Join(attachments, ","), param)

	bts, err := api.VkApi.Request(vk.METHOD_WALL_CREATE_COMMENT, param)
	if err != nil {
		return 0, err
	}
	log.Printf("%s", bts)
	type Resp struct {
		Response struct {
			CommentID int `json:"comment_id"`
		} `json:"response"`
	}
	var r Resp
	err = json.Unmarshal(bts, &r)
	if err != nil {
		return 0, err
	}

	return r.Response.CommentID, nil
}

func (api *Api) WallGetComments(ownerID, postID int, args ...url.Values) ([]Comment, error) {
	params := setToUrlValues("owner_id", ownerID, args...)
	params = setToUrlValues("post_id", postID, params)
	params = setToUrlValues("count", 100, params)

	bts, err := api.VkApi.Request(vk.METHOD_WALL_GET_COMMENTS, params)
	if err != nil {
		return nil, err
	}
	//	log.Println(string(bts))
	var r ResponseComments
	err = json.Unmarshal(bts, &r)
	if err != nil {
		return nil, err
	}
	return r.Response.Items, nil
}

func (api *Api) WallEdit(ownerID int, postID int, message string, args ...url.Values) error {
	return wall.Edit(api.API, ownerID, postID, message, args...)
}
