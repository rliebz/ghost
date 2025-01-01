// Package jsondiff prints a human-readable JSON diff.
package jsondiff

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Kind describes the result of the diff categorically.
type Kind int

// These must be kept in order from least severe to most severe.
const (
	Match Kind = iota
	NoMatch
	GotInvalid
	WantInvalid
	BothInvalid
)

func (k Kind) String() string {
	switch k {
	case Match:
		return "Match"
	case NoMatch:
		return "NoMatch"
	case GotInvalid:
		return "GotInvalid"
	case WantInvalid:
		return "WantInvalid"
	case BothInvalid:
		return "BothInvalid"
	default:
		return "InvalidKind"
	}
}

// Diff returns a pretty JSON diff of two inputs.
func Diff[T ~string | ~[]byte](got, want T) (string, Kind) {
	gotValue, gotErr := decode(got)
	wantValue, wantErr := decode(want)

	switch {
	case gotErr != nil && wantErr != nil:
		return "", BothInvalid
	case gotErr != nil:
		return "", GotInvalid
	case wantErr != nil:
		return "", WantInvalid
	}

	d := newDiffer()
	d.diffValues(gotValue, wantValue)
	return d.buf.String(), d.kind
}

func decode[T ~string | ~[]byte](v T) (any, error) {
	dec := json.NewDecoder(bytes.NewReader([]byte(v)))
	dec.UseNumber()

	var data any
	err := dec.Decode(&data)
	return data, err
}

type differ struct {
	buf    *bytes.Buffer
	kind   Kind
	level  int
	prefix byte
}

func newDiffer() *differ {
	return &differ{
		buf:    new(bytes.Buffer),
		level:  1, // start non-zero to make space for +/-/~
		prefix: ' ',
	}
}

func (d *differ) diffValues(got, want any) {
	switch {
	case got == nil && want == nil:
		d.writeValue(got)
		return
	case got == nil || want == nil:
		d.writeMismatch(got, want)
		return
	}

	gotType := reflect.TypeOf(got)
	wantType := reflect.TypeOf(want)

	// check types before converting to [reflect.Kind] to catch string/json.Number
	if gotType != wantType {
		d.writeMismatch(got, want)
		return
	}

	switch gotType.Kind() {
	case reflect.Bool:
		if got != want {
			d.writeMismatch(got, want)
			return
		}
	case reflect.String:
		if got != want {
			d.writeMismatch(got, want)
			return
		}
	case reflect.Map:
		d.diffMaps(got.(map[string]any), want.(map[string]any))
		return
	case reflect.Slice:
		d.diffSlices(got.([]any), want.([]any))
		return
	}

	d.writeValue(got)
}

func (d *differ) diffMaps(got, want map[string]any) {
	if len(got) == 0 && len(want) == 0 {
		d.buf.WriteString("{}")
		return
	}

	if d.buf.Len() == 0 {
		d.writeIndent()
	}

	d.buf.WriteString("{")
	d.level++

	keys := make(map[string]struct{})
	for k := range got {
		keys[k] = struct{}{}
	}
	for k := range want {
		keys[k] = struct{}{}
	}

	sortedKeys := make([]string, 0, len(keys))
	for k := range keys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	first := true
	for _, k := range sortedKeys {
		gotValue, gotOk := got[k]
		wantValue, wantOk := want[k]

		switch {
		case gotOk && wantOk:
			d.writeElementSeparator(first)
			fmt.Fprintf(d.buf, "%q: ", k)
			d.diffValues(gotValue, wantValue)
		case gotOk:
			d.failMatch()
			d.prefix = '+'
			d.writeElementSeparator(first)
			fmt.Fprintf(d.buf, "%q: ", k)
			d.writeValue(gotValue)
			d.prefix = ' '
		case wantOk:
			d.failMatch()
			d.prefix = '-'
			d.writeElementSeparator(first)
			fmt.Fprintf(d.buf, "%q: ", k)
			d.writeValue(wantValue)
			d.prefix = ' '
		default:
			panic(fmt.Sprintf("unexpected key %q in map; this is a ghost error", k))
		}

		first = false
	}
	d.level--
	d.writeNewlineIndent()
	d.buf.WriteByte('}')
}

func (d *differ) diffSlices(got, want []any) {
	if len(got) == 0 && len(want) == 0 {
		d.buf.WriteString("[]")
		return
	}

	if d.buf.Len() == 0 {
		d.writeIndent()
	}

	d.buf.WriteByte('[')
	d.level++
	for i := 0; i < len(got) || i < len(want); i++ {
		switch {
		case i >= len(got):
			d.failMatch()
			d.prefix = '-'
			d.writeExtraSliceElements(&i, want)
			d.prefix = ' '
		case i >= len(want):
			d.failMatch()
			d.prefix = '+'
			d.writeExtraSliceElements(&i, got)
			d.prefix = ' '
		default:
			d.writeElementSeparator(i == 0)
			d.diffValues(got[i], want[i])
		}
	}
	d.level--
	d.writeNewlineIndent()
	d.buf.WriteByte(']')
}

func (d *differ) writeExtraSliceElements(idx *int, s []any) {
	for ; *idx < len(s); *idx++ {
		d.writeElementSeparator(*idx == 0)
		d.writeValue(s[*idx])
	}
}

func (d *differ) writeValue(v any) {
	switch v := v.(type) {
	case bool:
		d.buf.WriteString(strconv.FormatBool(v))
	case json.Number:
		d.buf.WriteString(string(v))
	case string:
		d.buf.WriteString(strconv.Quote(v))
	case []any:
		d.writeSlice(v)
	case map[string]any:
		d.writeMap(v)
	default:
		if v != nil {
			panic(fmt.Sprintf("unexpected JSON type %T: %s; this is a ghost error", v, v))
		}
		d.buf.WriteString("null")
	}
}

// writeValueInline is like [writeValue], but with abbrevations for long
// strings, slices, and maps.
func (d *differ) writeValueInline(v any) {
	switch v := v.(type) {
	case []any:
		if len(v) == 0 {
			d.buf.WriteString("[]")
		} else {
			d.buf.WriteString("[...]")
		}
	case map[string]any:
		if len(v) == 0 {
			d.buf.WriteString("{}")
		} else {
			d.buf.WriteString("{...}")
		}
	default:
		d.writeValue(v)
	}
}

func (d *differ) writeSlice(s []any) {
	if len(s) == 0 {
		d.buf.WriteString("[]")
		return
	}
	d.buf.WriteByte('[')
	d.level++
	d.writeNewlineIndent()
	for i, vv := range s {
		if i != 0 {
			d.buf.WriteByte(',')
			d.writeNewlineIndent()
		}
		d.writeValue(vv)
	}
	d.level--
	d.writeNewlineIndent()
	d.buf.WriteByte(']')
}

func (d *differ) writeMap(m map[string]any) {
	if len(m) == 0 {
		d.buf.WriteString("{}")
		return
	}
	d.buf.WriteByte('{')
	d.level++

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	first := true
	for _, k := range keys {
		v := m[k]

		d.writeElementSeparator(first)

		fmt.Fprintf(d.buf, "%q: ", k)
		d.writeValue(v)

		first = false
	}
	d.level--
	d.writeNewlineIndent()
	d.buf.WriteByte('}')
}

func (d *differ) writeMismatch(got, want any) {
	// If we're in a multi-line diff, make sure it starts with [~]
	data := d.buf.Bytes()
	if idx := bytes.LastIndex(data, []byte("\n ")); idx != -1 {
		data[idx+1] = '~'
		d.buf.Reset()
		d.buf.Write(data)
	}

	d.writeValueInline(want)
	d.buf.WriteString(" => ")
	d.writeValueInline(got)
	d.failMatch()
}

// TODO: I would rather have the last element of EACH collection handle the
// middle commas correctly. Probably a hefty refactor, can save it for later.
func (d *differ) writeElementSeparator(first bool) {
	if !first {
		d.buf.WriteByte(',')
	}
	d.writeNewlineIndent()
}

func (d *differ) writeNewlineIndent() {
	d.buf.WriteByte('\n')
	d.writeIndent()
}

func (d *differ) writeIndent() {
	d.buf.WriteByte(d.prefix)
	d.buf.WriteString(strings.Repeat(" ", d.level*2-1))
}

func (d *differ) failMatch() {
	d.kind = NoMatch
}
