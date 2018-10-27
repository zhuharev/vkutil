// Copyright 2018 Kirill Zhuharev

package wall

import (
	"net/url"

	"github.com/zhuharev/vk"
	"github.com/zhuharev/vkutil/structs"
	"github.com/zhuharev/vkutil/util"
)

func Edit(api *structs.API, ownerID int, postID int, message string, args ...url.Values) error {
	params := util.SetToUrlValues("owner_id", ownerID, args...)
	params = util.SetToUrlValues("post_id", postID, params)
	params = util.SetToUrlValues("message", message, params)

	_, err := api.VkAPI.Request(vk.METHOD_WALL_EDIT, params)
	if err != nil {
		return err
	}
	return nil
}
