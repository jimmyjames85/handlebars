package dfs

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
)

type node map[string]interface{}

func newNode() node {
	return make(map[string]interface{})
}

func (n node) spawnArray(width int) {

	if _, ok := n["array"]; !ok {
		return
	}

	var arr []interface{}
	for i := 0; i < width; i++ {
		mod := i % 3
		if mod == 0 {
			arr = append(arr, i)
		} else if mod == 1 {
			arr = append(arr, fmt.Sprintf("str%d", i))
		} else {
			arr = append(arr, i%2 == 0)
		}

	}
	n["array"] = arr
}

// spawn populates the node's data, provided no data already exists
func (n node) spawn(width int) {

	if len(n) != 0 {
		return
	}

	if width > 3 {
		rand.Seed(time.Now().Unix())
		num := rand.Int() % width
		n["random"] = num
		n["isEven"] = num%2 == 0
		n["array"] = "uninitialized"
		width -= 3
	}

	for i := 0; i < width; i++ {
		n[fmt.Sprintf("%d", i)] = "leaf"
	}
}

func (n node) expandLeaves(width int) {
	for k, v := range n {
		if v == "leaf" {
			newn := newNode()
			newn.spawn(width)
			n[k] = newn
		}
		// only expand
	}
}

func (n node) expandTestJson(width, depth int) {

	if width <= 0 || depth <= 0 {
		return
	}

	n.spawn(width)

	if depth-1 > 0 {
		n.spawnArray(width)
		n.expandLeaves(width)
	}

	for _, v := range n {
		next, ok := v.(node)
		if !ok {
			continue
		}
		next.expandTestJson(width-1, depth-1)
	}
}

func CreateTestJson(width, depth int) ([]byte, error) {
	root := newNode()
	root.expandTestJson(width, depth)
	return json.Marshal(root)
}

type jsonArray []interface{}
type jsonObj map[string]interface{}

type dfs struct {

	// set before traversing the tree
	maxDepth int
	maxWidth int

	// for debuging
	Debug bool

	// bookkeeping
	curDepth int
	curNode  []string

	// If any time we get an error we set this variable and stop
	// traversal
	err error
}

func (d *dfs) debugf(format string, a ...interface{}) {
	if d.Debug {
		fmt.Printf(format, a...)
	}
}

const tabIndent = "   "

// this signature is defined in burger/jsonparser
func (d *dfs) arrayEach(value []byte, dataType jsonparser.ValueType, offset int, arrErr error) {

	if d.err != nil {
		return
	}
	err := d.checkWidthAndDepth([]byte("no key"), value, dataType)
	if err != nil {
		d.err = err
		// TODO this will continue to iterate over the
		// other elements and just return
		return
	}

	// spacer := strings.Repeat(tabIndent, d.curDepth)
	if dataType == jsonparser.Object {
		// d.debugf("%s(object)=>\n", spacer)

		d.curDepth++
		d.curNode = append(d.curNode, "[]")
		err := jsonparser.ObjectEach(value, d.objectEach)
		d.curNode = d.curNode[:len(d.curNode)-1]
		d.curDepth--

		if err != nil {
			d.err = err
			return
		}
	} else if dataType == jsonparser.Array {

		d.curDepth++
		d.curNode = append(d.curNode, "[]")
		_, err := jsonparser.ArrayEach(value, d.arrayEach)
		d.curNode = d.curNode[:len(d.curNode)-1]
		d.curDepth--

		if err != nil {
			d.err = err
			return
		}
	} else {
		// d.debugf("%s(%s)%s\n", spacer, dataType, value)
	}

}

// this signature is defined in burger/jsonparser
func (d *dfs) objectEach(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {

	if d.err != nil {
		return d.err
	}
	err := d.checkWidthAndDepth(key, value, dataType)
	if err != nil {
		d.err = err
		return err
	}

	// spacer := strings.Repeat(tabIndent, d.curDepth)

	if dataType == jsonparser.Object {
		//d.debugf("%s(object)%s=>\n", spacer, key)

		d.curDepth++
		d.curNode = append(d.curNode, string(key))
		err := jsonparser.ObjectEach(value, d.objectEach)
		d.curNode = d.curNode[:len(d.curNode)-1]
		d.curDepth--
		if err != nil {
			d.err = err
			return err
		}
	} else if dataType == jsonparser.Array {
		// d.debugf("%s(array)%s=>\n", spacer, key)

		// jsonparser.ArrayEach does not allow for early exist
		// so we will iterate over every element, however,
		// we've already checked the width above

		d.curDepth++
		d.curNode = append(d.curNode, string(key))
		_, err := jsonparser.ArrayEach(value, d.arrayEach)
		d.curNode = d.curNode[:len(d.curNode)-1]
		d.curDepth--
		if err != nil {
			d.err = err
			return err
		}
	} else {
		// d.debugf("%s%s=>%s\n", spacer, key, value)
	}

	// it is possible that d.err is set during the traversel, so
	// we bubble it up here
	return d.err
}

func New(maxWidth, maxDepth int) *dfs {
	ret := &dfs{maxDepth: maxDepth, maxWidth: maxWidth}
	return ret
}

func (d *dfs) Validate(jsonPayload []byte) error {

	//start := time.Now()

	// l, err := d.lengthOf(jsonPayload, jsonparser.Object)
	// //fmt.Printf("top level time: %v\n", start.Sub(time.Now()))

	// if err != nil {
	// 	return errors.Wrap(err, "unable to parse json obj")
	// }
	// if l > d.maxWidth {
	// 	return fmt.Errorf("top level object exceeds max width")
	// }

	ret := jsonparser.ObjectEach(jsonPayload, d.objectEach)
	//fmt.Printf("total time: %v\n", start.Sub(time.Now()))
	return ret
}

func (d *dfs) lengthOf(value []byte, dataType jsonparser.ValueType) (int, error) {

	var retErr error
	var length int

	if dataType == jsonparser.Object {
		// var obj jsonObj
		// err := json.Unmarshal(value, &obj)
		// if err != nil {
		// 	retErr = errors.Wrap(err, "unable to determine length of obj")
		// 	length = -1
		// } else {
		// 	length = len(obj)
		// }

		// TODO is ^^ faster?

		jsonparser.ObjectEach([]byte(value),
			func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
				if length > d.maxWidth {
					retErr = fmt.Errorf("%s exceeds width", string(key))
				}
				length++
				return retErr
			})

	} else if dataType == jsonparser.Array {

		// var arr jsonArray
		// err := json.Unmarshal(value, &arr)
		// if err != nil {
		// 	retErr = errors.Wrap(err, "unable to determine length of array")
		// 	length = -1
		// 	//return -1, errors.Wrap(err, "unable to determine width")
		// } else {
		// 	length = len(arr)
		// }

		//jsonparser.ArrayEach does not allow for early exist so we can't do the following

		jsonparser.ArrayEach([]byte(value),
			func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				if err != nil {
					retErr = err
					return
				}
				if length > d.maxWidth {
					retErr = fmt.Errorf("arr: [] exceeds width: %d", length)
				}
				length++
			})

	}
	return length, retErr

}

func (d *dfs) spacer() string {
	return strings.Repeat(tabIndent, d.curDepth)
}

func (d *dfs) checkWidthAndDepth(key []byte, value []byte, dataType jsonparser.ValueType) error {

	return nil

	if d.curDepth >= d.maxDepth {
		return fmt.Errorf("max depth exceeded: %s", strings.Join(d.curNode, "."))
	}

	w, err := d.lengthOf(value, dataType)
	if err != nil {
		return errors.Wrap(err, strings.Join(d.curNode, "."))
	}
	if w > d.maxWidth {
		return fmt.Errorf("max width exceeded: %s", strings.Join(d.curNode, "."))
	}

	return nil

}

func whiteSpaceEnd(data []byte) int {
	if len(data) == 0 {
		return -1
	}
	for i, c := range data {
		if c == ' ' || c == '\t' || c == '\n' {
			continue
		} else {
			return i
		}
	}
	return -1
}

// findNextUnescapedQuote finds the first unescaped quote and returns
// it's index, otherwise it returns -1
func findNextUnescapedQuote(data []byte) int {
	for i := 0; i < len(data); i++ {
		if data[i] == '"' {
			return i
		} else if data[i] == '\\' {
			i++ // skip next 'escaped' char
		}
	}
	return -1
}

func CalculateJsonDepth(data []byte, depthLimit int) (int, error) {

	var stack []byte

	ws := whiteSpaceEnd(data)
	if ws == -1 || ws >= len(data) || data[ws] != '{' {
		return -1, fmt.Errorf("invalid json: unable to find begining of json object")
	}

	data = data[ws:]

	maxDepth := 1

	for i := 0; i < len(data); i++ {
		if maxDepth > depthLimit {
			return -1, fmt.Errorf("max depth reached at: %d", maxDepth)
		}
		b := data[i]
		switch b {
		case '{':
			stack = append(stack, '{')
			//fmt.Printf("stack: %s\n", string(stack))
			maxDepth = max(maxDepth, len(stack))
		case '}':
			if len(stack) == 0 {
				return -1, fmt.Errorf("unexpected close bracket: }")
			}
			lastBracket := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			//fmt.Printf("stack: %s%c\n", string(stack), b)
			if lastBracket != '{' {
				//fmt.Printf("%s^\n", strings.Repeat(" ", i))
				return -1, fmt.Errorf("mismatching brackets: } vs %c", lastBracket)
			}
		case '[':
			stack = append(stack, '[')
			//fmt.Printf("stack: %s\n", string(stack))
			maxDepth = max(maxDepth, len(stack))
		case ']':
			if len(stack) == 0 {
				return -1, fmt.Errorf("unexpected close bracket: ]")
			}
			lastBracket := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			//fmt.Printf("stack: %s%c\n", string(stack), b)
			if lastBracket != '[' {
				//fmt.Printf("%s^\n", strings.Repeat(" ", i))
				return -1, fmt.Errorf("mismatching brackets: ] vs %c", lastBracket)
			}
		case '"':
			startQuote := i
			endQuote := findNextUnescapedQuote(data[i+1:])
			if endQuote == -1 {
				return -1, fmt.Errorf("invalid string at %d", startQuote)
			}
			// fmt.Printf("\t%s\n", data[startQuote:i+endQuote+2])

			i += endQuote + 1
		}

	}
	// fmt.Printf("max depth = %d \n", maxDepth)
	// fmt.Printf("finalstack: %s\n", string(stack))
	return maxDepth, nil
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func KB(b []byte) float64 {
	return float64(len(b)) / 1024.0
}
