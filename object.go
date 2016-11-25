package vkutil

type ObjectType string

//go:generate stringer -type=ObjectType
const (
	OBJECT_USER          ObjectType = "user"
	OBJECT_POST                     = "post"          // — запись на стене пользователя или группы;
	OBJECT_COMMENT                  = "comment"       // — комментарий к записи на стене;
	OBJECT_PHOTO                    = "photo"         //— фотография;
	OBJECT_AUDIO                    = "audio"         // — аудиозапись;
	OBJECT_VIDEO                    = "video"         // — видеозапись;
	OBJECT_NOTE                     = "note"          // — заметка;
	OBJECT_PHOTO_COMMENT            = "photo_comment" // — комментарий к фотографии;
	OBJECT_VIDEO_COMMENT            = "video_comment" // — комментарий к видеозаписи;
	OBJECT_TOPIC_COMMENT            = "topic_comment" // — комментарий в обсуждении;
	OBJECT_GROUP                    = "group"
	OBJECT_APPLICATION              = "application"
	OBJECT_PAGE                     = "page"
)

func (o *ObjectType) UnmarshalJSON(in []byte) error {
	*o = ObjectType(string(in[1 : len(in)-1]))
	return nil
}
