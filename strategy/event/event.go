package event

const (
	FIDREP_UPSERT = iota + 1
	FIDREP_DELETE
)

type Event struct {
	Type int32
	Data interface{}
	Ret  chan error
}
