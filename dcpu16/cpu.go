package dcpu16

type dcpu16 struct {
	registers [8]Word
	O         Word
	PC        Word
	SP        Word
	literals  [0x20]Word
	skipMap   map[Word]Word

	// CPU needs reference to memory to operate
	m         memory
}

func newCPU(m memory) *dcpu16 {
	cpu := new(dcpu16)
	cpu.m = m
	return cpu
}

func (cpu *dcpu16) reset() {
	// fill registers
	for i, _ := range cpu.registers {
		cpu.registers[i] = 0
	}

	for i, _ := range cpu.literals {
		cpu.literals[i] = Word(i)
	}

	cpu.skipMap = map[Word]Word{
		0x00: 0, 0x01: 0, 0x02: 0, 0x03: 0, 0x04: 0,
		0x05: 0, 0x06: 0, 0x07: 0, 0x08: 0, 0x09: 0,
		0x0a: 0, 0x0b: 0, 0x0c: 0, 0x0d: 0, 0x0e: 0,
		0x0f: 0,

		0x10: 1, 0x11: 1, 0x12: 1, 0x13: 1, 0x14: 1,
		0x15: 1, 0x16: 1, 0x17: 1,

		0x18: 0, 0x19: 0, 0x1a: 0, 0x1b: 0, 0x1c: 0,
		0x1d: 0,

		0x1E: 1, 0x1F: 1}

	cpu.SP = 0
	cpu.PC = 0
	cpu.O = 0
}

func (cpu *dcpu16) step() {

	// fetch
	instruction := cpu.nextWord()

	// decode
	basic, opcode, a, b := cpu.decodeInstruction(instruction)

	// execute
	if basic {
		aOp := cpu.resolveOperand(a)
		bOp := cpu.resolveOperand(b)
		cpu.executeBasic(opcode, a, aOp, bOp)
	} else {
		aOp := cpu.resolveOperand(a)
		cpu.executeExtended(opcode, aOp)
	}
}

func (cpu *dcpu16) decodeInstruction(instruction Word) (bool, Word, Word, Word) {
	// decode operands
	opcode := instruction & 0x000f
	a := (instruction & 0x03f0) >> 4
	b := (instruction & 0xfc00) >> 10
	basic := true

	// non-basic opcodes
	if opcode == 0 {
		basic = false
		opcode = a
		a = b
		b = 0
	}

	return basic, opcode, a, b
}

func (cpu *dcpu16) resolveOperand(val Word) *Word {
	switch val {
	case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07:
		return &cpu.registers[val]

	case 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f:
		return &cpu.m[cpu.registers[val-0x08]]

	case 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17:
		registerIndex := val & 0x000f
		return &cpu.m[cpu.nextWord()+cpu.registers[registerIndex]]
	case 0x18: // POP
		val := &cpu.m[cpu.SP]
		cpu.SP++
		return val
	case 0x19: // PEEK
		return &cpu.m[cpu.SP]
	case 0x1a: // PUSH
		cpu.SP -= 1
		return &cpu.m[cpu.SP]
	case 0x1b:
		return &cpu.SP
	case 0x1c:
		return &cpu.PC
	case 0x1d:
		return &cpu.O
	case 0x1e:
		return &cpu.m[cpu.nextWord()]
	case 0x1f:
		cpu.nextWord()
		return &cpu.m[cpu.PC]
	}

	return &cpu.literals[val-0x20]
}

func (cpu *dcpu16) executeBasic(opcode, destCode Word, aOp, bOp *Word) {

	a := *aOp
	b := *bOp
	var res Word

	switch opcode {
	// 0x1: SET a, b - sets a to b
	// 1 cycle plus cost of a and b
	case 0x1:
		res = b
	// 0x2: ADD a, b - sets a to a+b, sets O to 0x0001 if there's an overflow, 0x0 otherwise
	// 2 cycles plus cost of a and b
	case 0x2:
		var sum uint = uint(a) + uint(b)
		cpu.O = Word(sum >> 16)
		res = Word(sum)

	// 0x3: SUB a, b - sets a to a-b, sets O to 0xffff if there's an underflow, 0x0 otherwise
	// 2 cycles plus cost of a and b
	case 0x3:
		var diff uint = uint(a) - uint(b)

		if diff < 0 {
			cpu.O = 0xffff
		} else {
			cpu.O = 0x0
		}
		res = Word(diff)

	// 0x4: MUL a, b - sets a to a*b, sets O to ((a*b)>>16)&0xffff
	// 2 cycles plus cost of a and b
	case 0x4:
		res = a * b
		cpu.O = ((a * b) >> 16) & 0xffff

	// 0x5: DIV a, b - sets a to a/b, sets O to ((a<<16)/b)&0xffff. if b==0, sets a and O to 0 instead.
	// 3 cycles plus cost of a and b
	case 0x5:
		if b == 0 {
			cpu.O = 0
			res = 0
		} else {
			res = a / b
			cpu.O = ((a << 16) / b) & 0xffff
		}

	// 0x6: MOD a, b - sets a to a % b. if b==0, sets a to 0 instead.
	// 3 cycles plus cost of a and b
	case 0x6:
		if b == 0 {
			res = 0
		} else {
			res = a % b
		}

	// 0x7: SHL a, b - sets a to a<<b, sets O to ((a<<b)>>16)&0xffff
	// 2 cycles plus cost of a and b
	case 0x7:
		res = a << b
		cpu.O = ((a << b) >> 16) & 0xffff

	// 0x8: SHR a, b - sets a to a>>b, sets O to ((a<<16)>>b)&0xffff
	// 2 cycles plus cost of a and b
	case 0x8:
		res = a >> b
		cpu.O = ((a << 16) >> b) & 0xffff

	// 0x9: AND a, b - sets a to a & b
	// 1 cycle plus cost of a and b
	case 0x9:
		res = a & b

	// 0xa: BOR a, b - sets a to a | b
	// 1 cycle plus cost of a and b
	case 0xa:
		res = a | b

	// 0xb: XOR a, b - sets a to a^b
	// 1 cycle plus cost of a and b
	case 0xb:
		res = a ^ b

	// 0xc: IFE a, b - performs next instruction only if a==b
	// 2 cycles, plus the cost of a and b, plus 1 if the test fails
	case 0xc:
		if a != b {
			cpu.skipNextInstruction()
			return
		}

	// 0xd: IFN a, b - performs next instruction only if a!=b
	// 2 cycles, plus the cost of a and b, plus 1 if the test fails
	case 0xd:
		if a == b {
			cpu.skipNextInstruction()
			return
		}

	// 0xe: IFG a, b - performs next instruction only if a>b
	// 2 cycles, plus the cost of a and b, plus 1 if the test fails
	case 0xe:
		if a <= b {
			cpu.skipNextInstruction()
			return
		}

	// 0xf: IFB a, b - performs next instruction only if (a&b)!=0
	// 2 cycles, plus the cost of a and b, plus 1 if the test fails
	case 0xf:
		if (a & b) == 0 {
			cpu.skipNextInstruction()
			return
		}

	// otherwise just panic
	default:
		panic("Illegal instruction executed")
	}

	if destCode < 0x1f { // destination is not an immediate value
		*aOp = res // write to the location that operand points to
	}
}

func (cpu *dcpu16) executeExtended(opcode Word, aOp *Word) {
	switch opcode {
	// 0x01: JSR a - pushes the address of the next instruction to the stack, then sets PC to a
	// JSR takes 2 cycles, plus the cost of a
	case 0x01:
		cpu.SP -= 1 // PUSH
		cpu.m[cpu.SP] = cpu.PC
		cpu.PC = *aOp
	}
}

func (cpu *dcpu16) nextWord() Word {
	val := cpu.m[cpu.PC]
	cpu.PC++
	return val
}

func (cpu *dcpu16) skipNextInstruction() {
	instruction := cpu.nextWord()
	_, _, a, b := cpu.decodeInstruction(instruction)

	cpu.PC += cpu.skipMap[a] + cpu.skipMap[b]
}

// API function
func (cpu *dcpu16) Step() {
	cpu.step()
}