package dcpu16

type vm_t struct {
  Cpu *dcpu16_t
  Memory memory_t
}

var VM *vm_t = new(vm_t)

func (vm *vm_t) Reset() {
  vm.Cpu = dcpu16
  vm.Cpu.reset()

  vm.Memory = memory
  vm.Memory.clear()
}
