package dcpu16

import "fmt"

func (m memory_t) Rows() []string {
  tableWidth := 16
  rows := make([]string, MEMORY_SIZE / tableWidth)
  counter := 0
  for i := 0; i < MEMORY_SIZE / tableWidth; i++ {
    memoryLocation := fmt.Sprintf("%04x:", i*tableWidth)
    cells := ""
    cellSum := 0
    // loop over
    for j := 0; j < tableWidth; j++ {
      value := memory[Word(i*tableWidth + j)]
      cells = fmt.Sprintf("%s %04x", cells, value)
      cellSum += int(value)
    }

    // include in output only rows that contain something
    if cellSum > 0 {
      counter += 1
      rows[counter] = fmt.Sprintf("%s %s", memoryLocation, cells)
    }
  }

  return rows
}

func (c *dcpu16_t) Print() {
  registers := make([]string, 3)
  registers[0] = fmt.Sprint("  A     B     C     X     Y     Z     I     J     O     PC     SP")
  registers[1] = fmt.Sprint("----- ----- ----- ----- ----- ----- ----- ----- ----- ------ ------")
  registers[2] = fmt.Sprintf("%04x  %04x  %04x  %04x  %04x  %04x  %04x  %04x  %04x   %04x   %04x",
    c.registers[0],c.registers[1], c.registers[2], c.registers[3],
    c.registers[4], c.registers[5], c.registers[6], c.registers[7],
    c.O, c.PC, c.SP)

  for i, _ := range registers {
    fmt.Printf("%s\n", registers[i])
  }
}

func (c *dcpu16_t) Registers() []string {
  registers := make([]string, 3)
  registers[0] = fmt.Sprint("  A     B     C     X     Y     Z     I     J     O     PC     SP")
  registers[1] = fmt.Sprint("----- ----- ----- ----- ----- ----- ----- ----- ----- ------ ------")
  registers[2] = fmt.Sprintf("%04x  %04x  %04x  %04x  %04x  %04x  %04x  %04x  %04x   %04x   %04x",
    c.registers[0],c.registers[1], c.registers[2], c.registers[3],
    c.registers[4], c.registers[5], c.registers[6], c.registers[7],
    c.O, c.PC, c.SP)

    return registers
}
