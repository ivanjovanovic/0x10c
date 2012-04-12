package main

import (
  "./dcpu16"
  "github.com/nsf/termbox-go"
  "fmt"
  "os"
)

func tbPrint(x, y int, fg, bg termbox.Attribute, msg string) {
  for _, c := range msg {
    termbox.SetCell(x, y, c, fg, bg)
    x++
  }
}

func tbPrintf(x, y int, fg, bg termbox.Attribute, format string, args ...interface{}) {
  s := fmt.Sprintf(format, args...)
  tbPrint(x, y, fg, bg, s)
}

func main() {

  // init termbox
  err := termbox.Init()
  if err != nil {
    panic(err)
  }
  defer termbox.Close()

  // define channels of communication
  c := make(chan int)
  sc := make(chan int)

  // run virtual machine
  go func() {
    dcpu16.VM.Reset()
    loadMemory()

    for {
      <-sc // wait for instruction to make a step
      dcpu16.VM.Cpu.Step()
      c <- 1
    }
  }()

  // run drawing
  go func() {
    for {
      termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

      registers := dcpu16.VM.Cpu.Registers()
      memY := 1
      for _, row := range registers {
        tbPrint(1, memY, termbox.ColorDefault, termbox.ColorDefault, row)
        memY++
      }
      memY++

      rows := dcpu16.VM.Memory.Rows()
      // bottom of memory
      for _, row := range rows {
        tbPrint(1, memY, termbox.ColorDefault, termbox.ColorDefault, row)
        memY++
      }

      termbox.Flush()
      <-c
    }
  }()

  // key input handler
  for {
    evt := termbox.PollEvent() // blocking call
    if evt.Type == termbox.EventKey {
      if evt.Key == termbox.KeyCtrlC || (evt.Mod == 0 && evt.Ch == 'q') {
        os.Exit(1)
      }

      if evt.Key == termbox.KeyArrowRight {
        sc <- 1 // we want to make a step
      }

    }
  }
}

func loadMemory() {
  mem := []int{
    0x8401, 0x8811, 0x0404, 0xc1e1, 0x1000, 0x780d, 0x1000, 0x7831, 0x1000,
    0x79a1, 0x1000, 0x6021}

  for i, val := range mem {
    dcpu16.VM.Memory.Set(uint16(i), uint16(val))
  }
}
