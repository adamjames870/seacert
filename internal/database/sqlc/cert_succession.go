package sqlc

import "errors"

type SuccessionReason string

const (
	SuccessionReplaced SuccessionReason = "replaced"
	SuccessionUpdated  SuccessionReason = "updated"
)

func (r SuccessionReason) String() string {
	return string(r)
}

func SuccessionReasonFromString(s string) (SuccessionReason, error) {
	switch s {
	case "replaced":

		return SuccessionReplaced, nil
	case "updated":
		return SuccessionUpdated, nil
	default:
		return "", errors.New("unknown succession reason")
	}
}
