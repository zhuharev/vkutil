package vkutil

//"fmt"

type Attachment struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Photo struct {
		Photo130 string `json:"photo_130"`
		Photo604 string `json:"photo_604"`
		Photo807 string `json:"photo_807"`
	} `json:"photo"`
}

func (a Attachment) ToPhoto() (*Photo, error) {
	/*if typ, okconv := a["photo"].(string); !okconv || typ != "photo" {
		return nil, fmt.Errorf("Not photo")
	}*/
	ph := new(Photo)
	return ph, nil
}

func (a Attachment) IsPhoto() bool {
	return a.Type == OBJECT_PHOTO && a.Photo.Photo807 != ""
}
