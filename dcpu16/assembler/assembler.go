package assembler

import "strings"
import "github.com/ivanjovanovic/0x10c/dcpu16"

var programReader *strings.Reader
var byteCode []dcpu16.Word

/* var labels map[string]word*/
/* var labelPlaceholders map[string][]word*/


func Assemble(program string) []dcpu16.Word {

  programReader = strings.NewReader(program)
  byteCode := make([]dcpu16.Word, dcpu16.MEMORY_SIZE)

  /* labels := make(map[string]word)*/
  /* labelPlaceholders := make(map[string][]word)*/

  for token := nextToken() ; token != EOF ; {
    byteCode[i] = dcpu16.Word(i)
  }

  return byteCode
}
