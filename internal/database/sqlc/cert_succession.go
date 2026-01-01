package sqlc

type SuccessionReason string

const (
	SuccessionReplaced SuccessionReason = "replaced"
	SuccessionUpdated  SuccessionReason = "updated"
)
