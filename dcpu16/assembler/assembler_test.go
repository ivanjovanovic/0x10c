package assembler

import "testing"
import "os"
import "bufio"
import "io"
import "io/ioutil"
import "github.com/ivanjovanovic/0x10c/dcpu16"
import "fmt"

func readLines(filename string) []string {
  file, _ := os.Open(filename)
  reader := bufio.NewReader(file)
  lines := []string{}

  for {
    line, _, _ := reader.ReadLine()
    if _, err := reader.Peek(1); err == io.EOF {
      break
    }
    lines = append(lines, string(line))
  }
  return lines
}

func readFile(filename string) string {
  file, _ := os.Open(filename)
  reader := bufio.NewReader(file)

  contents, err := ioutil.ReadAll(reader)

  if err != nil {
    panic(err.Error())
  }

  return string(contents)
}

func assertSame(specified, assembled []string) bool {
  for i, _ := range specified {
    if specified[i] != assembled[i] {
      return false
    }
  }
  return true
}

func castMemoryContent(memory []dcpu16.Word) []string {
  memString := []string{}

  for i, _ := range memory {
    memString = append(memString, fmt.Sprintf("%04x", int(memory[i])))
  }
  return memString
}

// Print nice table with specified and comparable memory
func printMemoryComparison(spec, assembled []string) {
  fmt.Print("Specified\tAssembled\t\n")
  for i, val := range spec {
    fmt.Printf("   %s\t\t  %s\n", val, assembled[i])
  }
}

// Tests for official specification given on DCPU16 spec page
// http://0x10c.com/doc/dcpu-16.txt
func TestAssemblerSpecification(t *testing.T) {
  program := readFile("test/spec.dasm")
  memorySpec := readLines("test/spec.mem")

  memoryImage := Assemble(program)

  comparableMemoryImage := castMemoryContent(memoryImage)[0:len(memorySpec)]

  same := assertSame(memorySpec, comparableMemoryImage)

  if !same {
    printMemoryComparison(memorySpec, comparableMemoryImage)
    t.Error("Returned memory different from specified")
  }
}
