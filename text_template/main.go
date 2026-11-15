// Package main demonstrates text/template and html/template:
// actions, pipelines, FuncMap, define/block, missingkey, partial output
// on errors, and contextual auto-escaping in html/template.
package main

import (
	"bytes"
	"errors"
	"fmt"
	htmltemplate "html/template"
	"os"
	"slices"
	"strings"
	"text/template"
)

type Item struct {
	Name  string
	Price int
	Tags  []string
}

type Order struct {
	ID       int
	Customer string
	Items    []Item
	Paid     bool
}

// Total is a method; templates can call methods with no args (or with args).
func (o Order) Total() int {
	sum := 0
	for _, it := range o.Items {
		sum += it.Price
	}
	return sum
}

// Discount returns (value, error). A non-nil error stops Execute.
func (o Order) Discount(code string) (int, error) {
	if code == "PERRY10" {
		return o.Total() / 10, nil
	}
	return 0, errors.New("invalid coupon " + code)
}

func part1() {
	fmt.Println("== Part 1: actions, range, if, whitespace trim ==")
	const src = `Order #{{.ID}} for {{.Customer}}
{{- range $i, $it := .Items}}
  {{$i}}. {{$it.Name}} ${{$it.Price}}
{{- else}}
  (no items)
{{- end}}
Total: ${{.Total}} {{if .Paid}}[PAID]{{else}}[UNPAID]{{end}}
`
	t := template.Must(template.New("order").Parse(src))
	o := Order{ID: 7, Customer: "Perry", Items: []Item{
		{Name: "keyboard", Price: 120}, {Name: "mouse", Price: 40},
	}}
	if err := t.Execute(os.Stdout, o); err != nil {
		fmt.Println("error:", err)
	}
	_ = t.Execute(os.Stdout, Order{ID: 8, Customer: "Ann", Paid: true})
}

func part2() {
	fmt.Println("== Part 2: FuncMap, pipelines, built-ins ==")
	funcs := template.FuncMap{
		"upper": strings.ToUpper,
		"join":  strings.Join,
	}
	// Funcs must be called BEFORE Parse, otherwise Parse fails with
	// `function "upper" not defined`.
	const src = `{{range .Items -}}
{{.Name | upper | printf "%-10s"}} tags={{join .Tags ","}} n={{len .Tags}}{{if gt .Price 100}} (expensive){{end}}
{{end -}}
first tag of first item: {{index (index .Items 0).Tags 0}}
range over int (Go 1.24+):{{range 3}} {{.}}{{end}}
discount: {{.Discount "PERRY10"}}
`
	t := template.Must(template.New("p2").Funcs(funcs).Parse(src))
	o := Order{Items: []Item{
		{Name: "keyboard", Price: 120, Tags: []string{"usb", "rgb"}},
		{Name: "mouse", Price: 40, Tags: []string{"wireless"}},
	}}
	if err := t.Execute(os.Stdout, o); err != nil {
		fmt.Println("error:", err)
	}

	_, err := template.New("bad").Parse(`{{.Name | upper}}`)
	fmt.Println("parse without Funcs:", err)
}

func part3() {
	fmt.Println("== Part 3: define / template / block ==")
	const layout = `{{define "page"}}<{{block "title" .}}Default title{{end}}> {{template "body" .}}{{end}}`
	base := template.Must(template.New("base").Parse(layout +
		`{{define "body"}}base body for {{.}}{{end}}`))

	// Clone, then override "title" in the copy only.
	custom := template.Must(template.Must(base.Clone()).Parse(`{{define "title"}}Custom title{{end}}`))

	_ = base.ExecuteTemplate(os.Stdout, "page", "perry")
	fmt.Println()
	_ = custom.ExecuteTemplate(os.Stdout, "page", "perry")
	fmt.Println()
	// Templates() comes from a map, so its order is random: sort it.
	var names []string
	for _, tmpl := range base.Templates() {
		names = append(names, tmpl.Name())
	}
	slices.Sort(names)
	fmt.Println("defined:", names)
}

func part4() {
	fmt.Println("== Part 4: missing keys and errors ==")
	data := map[string]string{"name": "perry"}

	t := template.Must(template.New("m").Parse(`hello {{.name}}, role={{.role}}`))
	_ = t.Execute(os.Stdout, data)
	fmt.Println()

	strict := template.Must(template.New("m").Option("missingkey=error").Parse(`hello {{.name}}, role={{.role}}`))
	var buf bytes.Buffer
	err := strict.Execute(&buf, data)
	fmt.Printf("partial output=%q\nerr=%v\n", buf.String(), err)

	// A method returning an error also aborts, after writing earlier text.
	buf.Reset()
	t2 := template.Must(template.New("d").Parse(`total={{.Total}} discount={{.Discount "BAD"}}`))
	err = t2.Execute(&buf, Order{Items: []Item{{Price: 50}}})
	fmt.Printf("partial output=%q\nerr=%v\n", buf.String(), err)

	// Struct fields that do not exist are always an error (missingkey only affects maps).
	err = template.Must(template.New("f").Parse(`{{.Nope}}`)).Execute(&buf, Order{})
	fmt.Println("unknown field err:", err)
}

func part5() {
	fmt.Println("== Part 5: text/template vs html/template ==")
	const src = `<a href="/u?q={{.}}" onclick="greet('{{.}}')">{{.}}</a>`
	input := `<b>"Tom" & 'Jerry'</b>`

	tt := template.Must(template.New("t").Parse(src))
	ht := htmltemplate.Must(htmltemplate.New("h").Parse(src))

	fmt.Print("text: ")
	_ = tt.Execute(os.Stdout, input)
	fmt.Println()
	fmt.Print("html: ")
	_ = ht.Execute(os.Stdout, input)
	fmt.Println()

	// Trusted content can be marked with typed strings (use with care).
	ht2 := htmltemplate.Must(htmltemplate.New("h2").Parse(`{{.}}`))
	fmt.Print("trusted: ")
	_ = ht2.Execute(os.Stdout, htmltemplate.HTML("<b>bold</b>"))
	fmt.Println()
}

func main() {
	part1()
	part2()
	part3()
	part4()
	part5()
}
