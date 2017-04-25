package vkutil

import (
	"fmt"
	"net/url"
)

// NewCountOffsetParams return params with setted count and offset(optional)
func NewCountOffsetParams(count int, offsets ...int) url.Values {
	if count == 0 {
		count = 1
	}
	param := url.Values{
		"count": {fmt.Sprint(count)},
	}
	if len(offsets) > 0 && offsets[0] != 0 {
		param.Set("offset", fmt.Sprint(offsets[0]))
	}
	return param
}
