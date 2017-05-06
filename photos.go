package vkutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/zhuharev/vk"
)

type PhotoType string

const (
	PHOTO_WALL    = "wall"
	PHOTO_PROFILE = "profile"
	PHOTO_SAVED   = "saved"
)

type Size struct {
	Source string `json:"src"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Type   string `json:"type"`
}

type Photo struct {
	Id    int `json:"id"`
	Likes struct {
		Count int `json:"count"`
	} `json:"likes"`

	Sizes []Size `json:"sizes"`

	AlbumId int       `json:"album_id"`
	OwnerId int       `json:"owner_id"`
	UserId  int       `json:"user_id"`
	Text    string    `json:"text"`
	Date    EpochTime `json:"date"`
}

func (api *Api) PhotosGet(ownerId int, albumId string, params ...url.Values) ([]Photo, error) {
	rparams := url.Values{}
	rparams.Set("rev", "1")
	if len(params) == 1 {
		rparams = params[0]
	}
	rparams.Set("owner_id", fmt.Sprint(ownerId))
	rparams.Set("album_id", albumId)
	resp, e := api.VkApi.Request(vk.METHOD_PHOTOS_GET, rparams)
	if e != nil {
		return nil, e
	}
	var r ResponsePhotos
	e = json.Unmarshal(resp, &r)
	return r.Response.Items, e
}

// PhotosGetAll run "photos.getAll" method
func (api *Api) PhotosGetAll(ownerID int, params ...url.Values) ([]Photo, int, error) {
	rparams := url.Values{}
	if len(params) == 1 {
		rparams = params[0]
	}
	rparams.Set("owner_id", fmt.Sprint(ownerID))
	rparams.Set("extended", "1")
	rparams.Set("photo_sizes", "1")
	resp, e := api.VkApi.Request(vk.METHOD_PHOTOS_GET_ALL, rparams)
	if e != nil {
		return nil, 0, e
	}
	var r ResponsePhotos
	e = json.Unmarshal(resp, &r)
	return r.Response.Items, r.Response.Count, e
}

type wallUploadServerResponse struct {
	Response struct {
		UploadURL string `json:"upload_url"`
		AlbumID   int    `json:"album_id"`
		UserID    int    `json:"user_id"`
	} `json:"response"`
}

type uploadResponse struct {
	Server int    `json:"server"`
	Photo  string `json:"photo"`
	Hash   string `json:"hash"`
}

type photoResponse struct {
	Photos []Photo `json:"response"`
}

func (api *Api) uploadPhoto(uploadURL string, photo io.Reader) (uploadResp uploadResponse, err error) {
	var (
		bst []byte
	)

	var (
		body = bytes.NewBuffer(nil)
		wr   = multipart.NewWriter(body)
		part io.Writer
	)

	part, err = wr.CreateFormFile("photo", "photo.jpg")
	if err != nil {
		if api.debug {
			log.Println(err)
		}
		return
	}
	_, err = io.Copy(part, photo)
	if err != nil {
		if api.debug {
			log.Println(err)
		}
		return
	}
	err = wr.Close()
	if err != nil {
		if api.debug {
			log.Println(err)
		}
		return
	}

	var (
		resp *http.Response
		req  *http.Request
	)

	req, err = http.NewRequest("POST", uploadURL, body)
	if err != nil {
		if api.debug {
			log.Println(err)
		}
		return
	}
	req.Header.Set("Content-Type", wr.FormDataContentType())
	client := http.DefaultClient
	resp, err = client.Do(req)
	if err != nil {
		if api.debug {
			log.Println(err)
		}
		return
	}
	defer resp.Body.Close()
	bst, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		if api.debug {
			log.Println(err)
		}
		return
	}
	if api.debug {
		log.Printf("Upload response: %s\n", bst)
	}
	err = json.Unmarshal(bst, &uploadResp)
	if err != nil {
		if api.debug {
			log.Println(err)
		}
		return
	}
	return
}

// PhotosUploadWall upload photo to group
func (api *Api) PhotosUploadWall(groupID int, photo io.Reader) (_ Photo, err error) {
	var (
		bts []byte

		uploadServerResp wallUploadServerResponse
	)
	bts, err = api.VkApi.Request(vk.METHOD_PHOTOS_GET_WALL_UPLOAD_SERVER, url.Values{"group_id": {fmt.Sprint(groupID)}})
	if err != nil {
		if api.debug {
			log.Println(err)
		}
		return
	}

	err = json.Unmarshal(bts, &uploadServerResp)
	if err != nil {
		if api.debug {
			log.Println(err)
		}
		return
	}

	var (
		uploadResp uploadResponse
	)

	uploadResp, err = api.uploadPhoto(uploadServerResp.Response.UploadURL, photo)
	if err != nil {
		if api.debug {
			log.Println(err)
		}
		return
	}

	// fix quoted json
	// uploadResp.Photo, err = strconv.Unquote(uploadResp.Photo)
	// if err != nil {
	// 	if api.debug {
	// 		log.Println(err)
	// 	}
	// 	return
	// }

	params := url.Values{}
	params.Set("group_id", fmt.Sprint(groupID))
	params.Set("photo", uploadResp.Photo)
	params.Set("server", fmt.Sprint(uploadResp.Server))
	params.Set("hash", uploadResp.Hash)

	if api.debug {
		log.Println(params.Encode())
	}
	bts, err = api.VkApi.Request(vk.METHOD_PHOTOS_SAVE_WALL_PHOTO, params)
	if err != nil {
		if api.debug {
			log.Println(err)
		}
		return
	}
	if api.debug {
		log.Printf("Upload save response^ %s\n", bts)
	}

	var photoResp photoResponse

	err = json.Unmarshal(bts, &photoResp)
	if err != nil {
		if api.debug {
			log.Println(err)
		}
		return
	}

	if len(photoResp.Photos) == 0 {
		// todo
		return
	}

	return photoResp.Photos[0], nil
}
