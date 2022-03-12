package evaluate

import (
	"fmt"
	"testing"
)

func TestDumpCgroup(t *testing.T) {
	fmt.Printf("\n[Information Gathering - Cgroups]\n")
	DumpCgroup()
}
