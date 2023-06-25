package processor

type RegionCid struct {
	CidMap map[string]bool
	Cids   string
}

const (
	WARN  = 1
	ERROR = 2
)
