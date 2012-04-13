package assembler

import "regexp"
import "strings"
import "fmt"
import "dcpu16"

var byteCode []word
/* var labels map[string]word*/
/* var labelPlaceholders map[string][]word*/

func Assemble(program []string) []word {

  // initialize vars
  byteCode := make([]word, memory_size)
  /* labels := make(map[string]word)*/
  /* labelPlaceholders := make(map[string][]word)*/

  printProgram(program)

  // remove comments, comments at the end and empty lines
  // normalize labels so they are all on separate lines
  program = normalizeProgram(program)

  printProgram(program)

  for _, line := range program {
    switch {
    case isCode(line): decode(line)
    case isLabel(line): processLabel(line)
    }

  }

  return byteCode
}

func normalizeProgram(program []string) []string {
  trimmedProgram := make([]string,len(program))
  counter := 0
  whiteSpaceBegin, _ := regexp.Compile("^ +")
  multiWhiteSpace, _ := regexp.Compile(" +")
  whiteSpaceWithCommentsEnd, _ := regexp.Compile(" *;.*$")
  labelSingle, _ := regexp.Compile("^:[a-zA-Z0-9]+$")
  labelOnSameLine, _ := regexp.Compile("^(:[a-zA-Z0-9]+) +(.*)$")

  for _, line := range program {
    line = whiteSpaceBegin.ReplaceAllString(line, "")
    line = whiteSpaceWithCommentsEnd.ReplaceAllString(line, "")
    line = multiWhiteSpace.ReplaceAllString(line, " ")

    if len(line) == 0 { // discard empty lines
      continue
    }

    label := labelSingle.FindString(line)

    if len(label) > 0 { // we found label on single line
      trimmedProgram[counter] = line
      counter += 1
      continue
    }

    labelWithCode := labelOnSameLine.FindString(line)

    if len(labelWithCode) > 0 {
      elements := strings.Split(line, " ")
      trimmedProgram[counter] = elements[0] // put label inside
      counter += 1
      trimmedProgram[counter] = strings.Join(elements[1:len(elements)], " ") // put rest as command
      counter += 1
      continue
    }

    trimmedProgram[counter] = line
    counter += 1

  }

  return trimmedProgram
}

func isCode(line string) bool {
  return false
}

func isLabel(line string) bool {
  return false
}

func decode(line string) {

}

func processLabel(line string) {

}

func printProgram(program []string) {
  for _, line := range program {
    fmt.Printf("%s\n", line)
  }
}
