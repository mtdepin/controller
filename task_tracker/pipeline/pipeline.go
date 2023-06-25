package pipeline

//var Channel chan int

type PipeLine struct {
	dataChan chan interface{}
}

func (p *PipeLine) Init(nSize int) {

}