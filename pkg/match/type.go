package match

import (
	"github.com/google/uuid"
)

type gender uint

const (
	female gender = 0
	male   gender = 1
)

var GenderToValueMap = map[string]gender{
	"female": 0,
	"male":   1,
}

var GenderToStringMap = map[gender]string{
	0: "female",
	1: "male",
}

type SinglePerson struct {
	ID         uuid.UUID
	Name       string
	Height     uint
	Gender     gender
	WantedDate int
	// all match person
	MatchPerson map[uuid.UUID]*SinglePerson
	// new match person, does not display yet
	// a cache for temporary data, once display then move to MatchPerson
	NewMatchPerson map[uuid.UUID]*SinglePerson
}
