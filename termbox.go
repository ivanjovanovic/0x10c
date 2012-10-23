package main

import (
  "github.com/nsf/termbox-go"
  "fmt"
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
