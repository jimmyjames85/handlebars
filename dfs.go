package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
)

type dfs struct {

	// set before traversing the tree
	maxDepth int
	maxWidth int

	// bookkeeping
	curDepth int
	curWidth map[int]int

	curArrByteOffset int
	arrWidth         map[int]int

	// If any time we get an error we set this variable and stop
	// traversal
	err error
}

const tabIndent = "   "

// this signature is defined in burger/jsonparser
func (d *dfs) arrayEach(value []byte, dataType jsonparser.ValueType, offset int, err error) {

	if d.err != nil {
		return
	}

	d.arrWidth[d.curArrByteOffset]++
	if d.arrWidth[d.curArrByteOffset] > d.maxWidth {
		d.err = fmt.Errorf("maximum width of %d was exceeded: %q", d.maxWidth, value)
		return
	}

	spacer := strings.Repeat(tabIndent, d.curDepth+1)
	if dataType == jsonparser.Object {
		fmt.Printf("%s(object)=>\n", spacer)
		err := jsonparser.ObjectEach(value, d.objectEach)
		if err != nil {
			if d.err == nil {
				d.err = err
			}
			return
		}
	} else {
		fmt.Printf("%s(%s)%s\n", spacer, dataType, value)
	}

}

// this signature is defined in burger/jsonparser
func (d *dfs) objectEach(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {

	if d.err != nil {
		return d.err
	}

	d.curDepth++

	spacer := strings.Repeat(tabIndent, d.curDepth)
	defer func() { d.curDepth-- }()
	d.curWidth[d.curDepth]++

	if d.curDepth > d.maxDepth {
		d.err = fmt.Errorf("maximum depth of %d was exceeded: %q=%q", d.maxDepth, string(key), string(value))
		return d.err
	}

	if dataType == jsonparser.Object {
		fmt.Printf("%s(object)%s=>\n", spacer, key)

		err := jsonparser.ObjectEach(value, d.objectEach)
		if err != nil {
			d.err = err
			return err
		}
	} else if dataType == jsonparser.Array {
		fmt.Printf("%s(array)%s=>\n", spacer, key)
		curOffset := d.curArrByteOffset
		d.curArrByteOffset = offset
		_, err := jsonparser.ArrayEach(value, d.arrayEach)
		d.curArrByteOffset = curOffset
		if err != nil {
			d.err = err
			return err
		}
	} else {
		fmt.Printf("%s%s=>%s\n", spacer, key, value)
	}
	// d.err could have been populated during the traversal
	return d.err
}

func New(maxDepth, maxWidth int) *dfs {
	ret := &dfs{maxDepth: maxDepth, maxWidth: maxWidth}
	ret.curWidth = make(map[int]int)
	ret.arrWidth = make(map[int]int)
	return ret
}

func handleDFSDemo() {

	d := New(5, 7)

	if len(os.Args) > 2 {

		maxDepth, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("Couldn't parse maxDepth\n%s maxDepth maxWidth", os.Args[0])
		}
		maxWidth, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Couldn't parse maxWidth\n%s maxDepth maxWidth", os.Args[0])
		}

		d = New(maxDepth, maxWidth)
	}

	err := jsonparser.ObjectEach([]byte(jsonPayload), d.objectEach)
	if err != nil {
		log.Fatalf("%v", err)
	}

}
