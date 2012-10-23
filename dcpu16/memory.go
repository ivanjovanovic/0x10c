package dcpu16

type Word uint16
type memory  []Word
const MEMORY_SIZE int = 65536

func newMemory() memory {
	return make(memory, MEMORY_SIZE)
}

func (m memory) clear() {
	for i := 0; i < MEMORY_SIZE; i++ {
		m[i] = 0
	}
}

func (m memory) Set(address, val uint16) {
	m[Word(address)] = Word(val)
}