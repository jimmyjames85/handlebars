package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/aymerick/raymond"
)

func main() {
	handleBarDemo()
	//handleDFSDemo()
}

var jsonPayload = `
{
    "people": [
        "marcel marcel marcel",
        "jean claud",
        "ronald McDonald",
	{"obj":true, "arr": ["bar", "baz", {"this":"map will be coerced into a string"}], "arrSize": 3},
	{"obj":true, "arr":[11,22,33,44], "arrSize": 4},
	{"obj":true, "arr":[55,66,77,88], "arrSize": 4, "subObj": {"foo":"bar"}},
	"Jon jacob jingleheimerschmidt"
    ]
}`

var htmlPayload = `
<ul class="people">
{{#each people}}
  {{#if obj}}
	{{#if (gt arrSize 0) }}
			=============================================
			Sub Items
			=============================================
		{{#each arr}}
			<li>{{this}}</li>
		{{/each}}
	{{/if}}
  {{else}}
	<li>{{uppercase this}}</li>
  {{/if}}
{{/each}}
</ul>
`

// stringInterface is the type we unmarshal raw json into. This works
// nicely with raymond.Exec
type stringInterface map[string]interface{}

func handleBarDemo() {

	// These helpers are implemented below
	raymond.RegisterHelper("gt", gt)
	raymond.RegisterHelper("uppercase", uppercase)

	var obj stringInterface

	err := json.Unmarshal([]byte(jsonPayload), &obj)
	if err != nil {
		log.Fatalf("%v", err)
	}

	tpl := raymond.MustParse(htmlPayload)

	result, err := tpl.Exec(obj)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	fmt.Printf("%s\n", result)
}

// I found some example helpers to act as a guide here:
// https://www.npmjs.com/package/just-handlebars-helpers

// gt will attempt to coerce both a and b into floats and return true
// if a > b, otherise false
func gt(a interface{}, b interface{}, options *raymond.Options) bool {

	if a == nil || b == nil {
		return false
	}

	var fa, fb float64
	typ := reflect.TypeOf(a)
	if typ.Kind() == reflect.Int {
		fa = float64(a.(int))
	} else if typ.Kind() == reflect.Float64 {
		fa = a.(float64)
	} else {
		return false
	}

	typ = reflect.TypeOf(b)
	if typ.Kind() == reflect.Int {
		fb = float64(b.(int))
	} else if typ.Kind() == reflect.Float64 {
		fb = b.(float64)
	} else {
		return false
	}

	return fa > fb
}

// uppercase will attempt to coerece i into a string, if successful it
// will convert the entire string to uppercase
func uppercase(i interface{}, options *raymond.Options) string {

	typ := reflect.TypeOf(i)
	if typ.Kind() != reflect.String {
		return fmt.Sprintf("%v", i)
	}
	str := i.(string)
	return strings.ToUpper(str)
}
