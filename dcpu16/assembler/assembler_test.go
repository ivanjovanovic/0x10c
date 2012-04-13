package assembler

import "testing"
import "os"
import "bufio"
import "io"
import "github.com/ivanjovanovic/0x10c/dcpu16"
import "fmt"

func readLines(filename string) []string {
  file, _ := os.Open(filename)
  reader := bufio.NewReader(file)
  lines := []string{}

  for {
    line, _, _ := reader.ReadLine()
    // is there a better way to get to EOF
    if _, err := reader.Peek(1); err == io.EOF {
      break
    }
    lines = append(lines, string(line))
  }
  return lines
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

func printMemory(spec, assembled []string) {
  fmt.Print("Specified:\tAssembled\t\n")
  for i, val := range spec {
    fmt.Printf("   %s\t\t  %s\n", val, assembled[i])
  }
}

func TestAssemblerSpecification(t *testing.T) {
  programLines := readLines("test/spec.dasm")
  memorySpec := readLines("test/spec.mem")

  memoryImage := Assemble(programLines)

  comparableMemoryImage := castMemoryContent(memoryImage)[0:len(memorySpec)]

  same := assertSame(memorySpec, comparableMemoryImage)

  if !same {
    printMemory(memorySpec, comparableMemoryImage)
    t.Error("Returned memory different from specified")
  }
}
