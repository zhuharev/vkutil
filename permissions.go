package vkutil

import (
	"strconv"
)

type UserPermission int

const (
	UP_Notify UserPermission = 1 << iota
	UP_Friends
	UP_Photos
	UP_Audio
	UP_Video
	UP_Pages
)

const (
	UP_All = UP_Notify | UP_Friends | UP_Photos | UP_Audio
)

func (u UserPermission) String() string {
	switch u {
	default:
		return strconv.Itoa(int(UP_All))
	}
}
