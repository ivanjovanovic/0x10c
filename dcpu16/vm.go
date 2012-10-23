package dcpu16

type VM struct {
	CPU *dcpu16
	Memory memory
}

func (vm *VM) Reset() {
	vm.resetCpu()
	vm.clearMemory()
}

func (vm *VM) resetCpu() {
	vm.CPU.reset()
}

func (vm *VM) clearMemory() {
	vm.Memory.clear()
}

func NewVM() *VM {
	newVM := new(VM)
	newVM.Memory = newMemory()
        newVM.CPU = newCPU(newVM.Memory)
	return newVM
}