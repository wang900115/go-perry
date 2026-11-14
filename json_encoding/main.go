// Package main demonstrates encoding/json (v1): struct tags, omitempty vs
// omitzero, custom Marshaler/Unmarshaler, json.RawMessage for delayed
// decoding, strict decoding, and number handling with interface{} values.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ---------- Part 1: struct tags ----------

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"` // dropped when ""
	Password string `json:"-"`               // never encoded or decoded
	Age      int    `json:"age,string"`      // encoded as a JSON string: "30"
	internal string // unexported: invisible to encoding/json
}

// ---------- Part 2: omitempty vs omitzero (Go 1.24+) ----------

type Event struct {
	Name     string    `json:"name"`
	EmptyAt  time.Time `json:"empty_at,omitempty"` // struct is never "empty": always encoded
	ZeroAt   time.Time `json:"zero_at,omitzero"`   // uses IsZero(): dropped
	Tags     []string  `json:"tags,omitempty"`     // nil and len 0 both dropped
	ZeroTags []string  `json:"zero_tags,omitzero"` // only nil is zero; []string{} is kept
}

// ---------- Part 3: custom Marshaler / Unmarshaler ----------

// Duration is encoded as a human string like "1m30s" instead of nanoseconds.
type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("duration must be a string: %w", err)
	}
	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(v)
	return nil
}

type Job struct {
	Name    string   `json:"name"`
	Timeout Duration `json:"timeout"`
}

// ---------- Part 4: RawMessage for a tagged union ----------

type Envelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"` // kept as raw bytes, decoded later
}

type Click struct {
	X, Y int
}

type KeyPress struct {
	Key string `json:"key"`
}

func decodeEvent(raw []byte) (any, error) {
	var env Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, err
	}
	switch env.Type {
	case "click":
		var c Click
		err := json.Unmarshal(env.Data, &c)
		return c, err
	case "key":
		var k KeyPress
		err := json.Unmarshal(env.Data, &k)
		return k, err
	default:
		return nil, fmt.Errorf("unknown event type %q", env.Type)
	}
}

func main() {
	fmt.Println("== Part 1: struct tags ==")
	u := User{ID: 1, Name: "perry", Password: "secret", Age: 30, internal: "x"}
	b, _ := json.Marshal(u)
	fmt.Println(string(b))

	var u2 User
	in := `{"id":2,"NAME":"Bob","password":"leak","age":"41"}`
	err := json.Unmarshal([]byte(in), &u2)
	fmt.Printf("err=%v name=%q password=%q age=%d\n", err, u2.Name, u2.Password, u2.Age)

	fmt.Println("== Part 2: omitempty vs omitzero ==")
	e := Event{Name: "deploy", Tags: []string{}, ZeroTags: []string{}}
	b, _ = json.Marshal(e)
	fmt.Println(string(b))
	b, _ = json.Marshal(Event{Name: "nil slices"})
	fmt.Println(string(b))

	fmt.Println("== Part 3: custom Marshaler / Unmarshaler ==")
	b, _ = json.Marshal(Job{Name: "backup", Timeout: Duration(90 * time.Second)})
	fmt.Println(string(b))
	var j Job
	err = json.Unmarshal([]byte(`{"name":"sync","timeout":"2h15m"}`), &j)
	fmt.Printf("err=%v timeout=%v\n", err, time.Duration(j.Timeout))
	err = json.Unmarshal([]byte(`{"name":"bad","timeout":5}`), &j)
	fmt.Println("err:", err)

	fmt.Println("== Part 4: RawMessage tagged union ==")
	for _, s := range []string{
		`{"type":"click","data":{"X":10,"Y":20}}`,
		`{"type":"key","data":{"key":"Enter"}}`,
		`{"type":"scroll","data":{}}`,
	} {
		v, err := decodeEvent([]byte(s))
		fmt.Printf("%T %+v err=%v\n", v, v, err)
	}

	fmt.Println("== Part 5: strict decoding ==")
	body := `{"id":3,"name":"Eve","admin":true}`
	var u3 User
	fmt.Println("lenient err:", json.Unmarshal([]byte(body), &u3))
	dec := json.NewDecoder(strings.NewReader(body))
	dec.DisallowUnknownFields()
	fmt.Println("strict err: ", dec.Decode(&u3))

	var ute *json.UnmarshalTypeError
	err = json.Unmarshal([]byte(`{"id":"four"}`), &u3)
	if errors.As(err, &ute) {
		fmt.Printf("type error: field=%s value=%s want=%v\n", ute.Field, ute.Value, ute.Type)
	}

	fmt.Println("== Part 6: numbers in interface{} ==")
	big := `{"order_id":9007199254740993}`
	var m map[string]any
	_ = json.Unmarshal([]byte(big), &m)
	fmt.Printf("default: %T %.0f\n", m["order_id"], m["order_id"])

	dec = json.NewDecoder(bytes.NewReader([]byte(big)))
	dec.UseNumber()
	_ = dec.Decode(&m)
	n := m["order_id"].(json.Number)
	i, _ := n.Int64()
	fmt.Printf("UseNumber: %T %s -> int64 %d\n", m["order_id"], n, i)
}
