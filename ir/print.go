//  print.go -- print human-readable dump of intermediate code

package ir

import (
	"fmt"
	"strings"
)

const indentBy = "   " // increment for additional indentation labels

// Print(label, tree) -- print a tree of IR structs on stdout
func Print(label string, tree interface{}) {
	fmt.Printf("\n========== %s ==========\n", label)
	subprint("", tree)
	fmt.Println()
}

// subprint(indent, tree) -- print part of the IR tree
func subprint(indent string, tree interface{}) {
	switch t := tree.(type) {
	case nil:
		return
	case []interface{}:
		for _, v := range t {
			subprint(indent, v)
		}
	case []Ir_chunk:
		for _, v := range t {
			subprint(indent, v)
		}
	case Ir_Function:
		iplus := indent + indentBy
		fmt.Printf("\n%sproc %s {%v}  parent:%s  start:%v\n",
			indent, t.Name, t.Coord, t.Parent, t.CodeStart)
		fmt.Printf("%sparam %v", iplus, t.ParamList)
		if t.Accumulate != "" {
			fmt.Printf(" [accumulate]")
		}
		fmt.Printf("\n%slocal %v\n", iplus, t.LocalList)
		fmt.Printf("%sstatic %v\n", iplus, t.StaticList)
		fmt.Printf("%sunbound %v\n", iplus, t.UnboundList)
		subprint(indent, t.CodeList)
	case Ir_chunk:
		fmt.Printf("%s%s:\n", indent, t.Label)
		subprint(indent+indentBy, t.InsnList)
	default:
		s := fmt.Sprintf("%T %v", tree, tree)
		if strings.HasPrefix(s, "ir.Ir_") {
			s = s[6:]
		}
		fmt.Printf("%s%s\n", indent, s)
	}
}
