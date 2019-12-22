package ch12

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"text/scanner"
)

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Unmarshal(data []byte, out interface{}) error {
	scan := scanner.Scanner{Mode: scanner.GoTokens}
	scan.Init(bytes.NewReader(data))
	return nil
}

func encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Fprintf(buf, "%d", v.Uint())
	case reflect.String:
		fmt.Fprintf(buf, "%s", v.String())
	case reflect.Ptr:
		return encode(buf, v.Elem())
	case reflect.Slice, reflect.Array:
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := encode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(')')
	case reflect.Map:
		buf.WriteByte('(')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteByte('(')
			if err := encode(buf, key); err != nil { // key can be many types, not only string
				return err
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	case reflect.Struct:
		buf.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteByte('(')
			buf.WriteString(v.Type().Field(i).Name)
			buf.WriteByte(' ')
			if err := encode(buf, v.Field(i)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	// case reflect.Interface: ?!
	default: // float, complex, bool, chan, func, interface
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

func decode() {

}

// ========================

type lexer struct {
	scan  scanner.Scanner
	token rune // current token
}

func (lex *lexer) next() {
	lex.token = lex.scan.Scan()
}

func (lex *lexer) text() string {
	return lex.scan.TokenText()
}

func (lex *lexer) consume(want rune) {
	if lex.token != want {
		panic(fmt.Sprintf("got %q, want %q", lex.token, want))
	}
	lex.next()
}

func read(lex *lexer, v reflect.Value) { // v must be settable
	if lex.text() == "nil" {
		v.Set(reflect.Zero(v.Type()))
		lex.next()
		return
	}

	// driven by the token type, but not v.Kind()...
	switch lex.token {
	case scanner.Int:
		i, _ := strconv.Atoi(lex.text())
		v.SetInt(int64(i))
		return
	case scanner.String:
		s, _ := strconv.Unquote(lex.text())
		v.SetString(s)
		return
	case '(':
		lex.next()
		readList(lex, v)
		lex.consume(')')
		return
	}
	panic(fmt.Sprintf("unexpected token %q", lex.text()))
}

func readList(lex *lexer, v reflect.Value) { // readList: read array/slice/struct/map
	switch v.Kind() {
	case reflect.Array: // (item ...)
		for i := 0; !endList(lex); i++ {
			read(lex, v.Index(i)) // read into item
		}
	case reflect.Slice: // (item ...)
		for !endList(lex) {
			// Need to build a new slice:
			// - Elem Type: v.Type().Elem()
			// - New: reflect.New(...)
			// - Append: reflect.Append(v, ...)
			elemType := v.Type().Elem()
			item := reflect.New(elemType).Elem()
			read(lex, item) // read into item
			reflect.Append(v, item)
		}
	case reflect.Struct: // ((name value)) ...)
		for !endList(lex) {
			lex.consume('(')
			if lex.token != scanner.Ident {
				panic(fmt.Sprintf("got token %q, want field name", lex.text()))
			}
			name := lex.text()
			lex.next()
			read(lex, v.FieldByName(name))
			lex.consume(')')
		}
	case reflect.Map: // ((key value)) ...)
		// Need to build a new map:
		// - Make: reflect.MakeMap(...)
		// - Key Type: v.Type().Key()
		// - Value Type: v.Type().Elem()
		// - New: reflect.New(...)
		// - Set KV Pair: v.SetMapIndex(key, value)
		v.Set(reflect.MakeMap(v.Type()))
		for !endList(lex) {
			lex.consume('(')
			keyType := v.Type().Key()
			key := reflect.New(keyType).Elem()
			read(lex, key)
			valueType := v.Type().Elem()
			value := reflect.New(valueType).Elem()
			read(lex, value)
			v.SetMapIndex(key, value)
			lex.consume(')')
		}
	default:
		panic(fmt.Sprintf("cannot decode list into %v", v.Type()))
	}
}

func endList(lex *lexer) bool {
	switch lex.token {
	case scanner.EOF:
		panic("end of file")
	case ')':
		return true
	}
	return false
}

// ========================
