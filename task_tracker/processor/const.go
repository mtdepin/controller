package processor

const (
	UPLOAD_FINISH = iota + 1
	FILE_NUM_NOT_EQUAL
	FILE_NOT_EXIST
	CID_NOT_EQUAL
)

const (
	EXTEND_SIZE    = 1024 * 100 //10w
	TIME_INTERAL   = 10
	FACTOR         = 1
	INTERNAL       = 100
	FID_EVENT_SIZE = 300
)
