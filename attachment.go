package vkutil

//"fmt"

type Attachment struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Photo struct {
		Photo130  string `json:"photo_130"`
		Photo604  string `json:"photo_604"`
		Photo807  string `json:"photo_807"`
		Photo1280 string `json:"photo_1280"`
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
	return a.Type == OBJECT_PHOTO && a.PhotoURL() != ""
}

func (a Attachment) PhotoURL() string {
	if a.Photo.Photo1280 != "" {
		return a.Photo.Photo1280
	} else if a.Photo.Photo807 != "" {
		return a.Photo.Photo807
	} else {
		return a.Photo.Photo604
	}
}
