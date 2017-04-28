package vkutil

import (
	"strconv"
)

// UserPermission used in auth
type UserPermission int

const (
	// UPNotify notify comment
	UPNotify UserPermission = 1 << iota
	UPFriends
	UPPhotos
	UPAudio
	UPVideo
	_
	_

	UPPages
	_
	_
	_
	_
	_
	_
	_
	_
	UPOffline
)

const (
	// UPAll represent max permission
	UPAll = UPNotify | UPFriends | UPPhotos | UPAudio | UPOffline
)

func (u UserPermission) String() string {
	return strconv.Itoa(int(u))
}
