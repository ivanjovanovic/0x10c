package dcpu16

type word uint16
type memory_t  []word
const memory_size int = 65536

var memory = make(memory_t, memory_size)

func (m memory_t) clear() {
  for i := 0; i < memory_size; i++ {
    m[i] = 0
  }
}

func (m memory_t) Set(address, val uint16) {
  memory[word(address)] = word(val)
}
