package phrase

const DefaultBlockSize = 1024

type Pool struct {
	block []Phrase
	off   int
}

func NewPool(blockSize int) *Pool {
	return &Pool{
		block: make([]Phrase, blockSize),
	}
}

func (p *Pool) Get() *Phrase {
	if len(p.block) == 0 {
		return nil
	}

	if len(p.block) == p.off {
		p.block = make([]Phrase, len(p.block))
		p.off = 0
	}

	p.off++

	return &p.block[p.off-1]
}
