package vkutil

type Relation int

const (
	RNo = iota
	RSingle
	RInRelationship
	REngaged
	Rmarried
	RItsComplicated
	RActivelySearching
	RInLove
)

func (r Relation) String() string {
	switch r {
	case RSingle:
		return "single"
	case RInRelationship:
		return "in relationship"
	case REngaged:
		return "engaged"
	case Rmarried:
		return "married"
	case RItsComplicated:
		return "its complicated"
	case RActivelySearching:
		return "actively searching"
	case RInLove:
		return "in love"
	}
	return ""
}
