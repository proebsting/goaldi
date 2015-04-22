//  import_test.go -- test importing Go values into Goaldi

package runtime

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

type extl struct{ i int }

func (*extl) GoaldiExternal() {}

type impr struct{ i int }

func (*impr) Import() Value { return NewNumber(4.713) }

func TestImport(t *testing.T) {

	// test Goaldi types
	testImp(t, NilValue, NilValue)
	testImp(t, ONE, ONE)
	testImp(t, EMPTY, EMPTY)
	testImp(t, "abc", NewString("abc"))
	testImp(t, STDIN, STDIN)

	// test nil flavors
	testImp(t, nil, NilValue)
	testImp(t, (*float64)(nil), NilValue)
	testImp(t, (*os.File)(nil), NilValue)

	// test simple types
	testImp(t, false, ZERO)
	testImp(t, true, ONE)
	testImp(t, 0, ZERO)
	testImp(t, 1, ONE)
	testImp(t, 0.0, ZERO)
	testImp(t, 1.0, ONE)
	testImp(t, uint16(1), ONE)
	testImp(t, "7.8", NewString("7.8"))

	// test file import
	i := fmt.Sprintf("%#v", Import(bufio.NewReader(os.Stdin)))
	o := fmt.Sprintf("%#v", Import(bufio.NewWriter(os.Stdout)))
	expect(t, "stdin", "file(*bufio.Reader,r)", i)
	expect(t, "stdout", "file(*bufio.Writer,w)", o)

	// test external imports
	testImp(t, &impr{1}, NewNumber(4.713))
	x := &extl{2}
	testImp(t, x, x)
	expect(t, "external", "&runtime.extl{i:2}", fmt.Sprintf("%#v", Import(x)))
	m := make(map[int]string)
	f := Import(&m)
	expect(t, "external", &m, f)
}

func testImp(t *testing.T, goval interface{}, expected Value) {
	imported := Import(goval)
	if Identical(imported, expected) != expected {
		t.Errorf("import(%T:%v) expected %v got %T:%v\n",
			goval, goval, expected, imported, imported)
	}
}
