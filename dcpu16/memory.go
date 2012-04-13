package dcpu16

type Word uint16
type memory_t  []Word
const MEMORY_SIZE int = 65536

var memory = make(memory_t, MEMORY_SIZE)

func (m memory_t) clear() {
  for i := 0; i < MEMORY_SIZE; i++ {
    m[i] = 0
  }
}

func (m memory_t) Set(address, val uint16) {
  memory[Word(address)] = Word(val)
}
