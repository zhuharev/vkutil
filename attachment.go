package vkutil

import (
//"fmt"
)

type Attachment map[string]interface{}

func (a Attachment) ToPhoto() (*Photo, error) {
	/*if typ, okconv := a["photo"].(string); !okconv || typ != "photo" {
		return nil, fmt.Errorf("Not photo")
	}*/
	ph := new(Photo)
	return ph, nil
}
