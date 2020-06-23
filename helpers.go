// +build js

package main

import (
	"errors"
	"fmt"
	"reflect"
	"syscall/js"
)

// funcToJs automatically wraps fn in js.Func type.
//
// It checks fn to ensure that only accepts parameters of types that are
// convertible from Js to Go, it returns 0 or one value. It optionally supports
// to receive the 'this' Js parameter, with the condition that fn has to receive
// as a first parameter a value of the type js.Value.
// If any of these checks fails, it returns an error with a specific message.
//
// The js.Func returned automatically checks on each Js function call that the
// expected number of parameters and types are passed otherwise it panics with
// an specific message.
func funcToJs(fn interface{}) (js.Func, error) {
	// First of all, check that fn fulfills the conditions for being callable from
	// the Js world
	fnType := reflect.TypeOf(fn)
	if fnType == nil || fnType.Kind() != reflect.Func {
		return js.Func{}, errors.New("fn argument isn't a function")
	}

	var (
		useThis       = false
		fnNumParamsIn = fnType.NumIn()
	)

	if fnNumParamsIn > 0 {
		i := 0
		// check if the first fn argument is of the type jsThis
		if fnType.In(0).Name() == reflect.TypeOf(js.Value{}).Name() {
			useThis = true
			i++
		}

		for ; i < fnNumParamsIn; i++ {
			switch fnType.In(i).Kind() {
			case reflect.Bool:
			case reflect.Float32, reflect.Float64:
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			case reflect.String:
			default:
				return js.Func{}, fmt.Errorf(
					"'fn' %d parameter is of a type not convertible from Js", i,
				)
			}
		}
	}

	if fnType.NumOut() > 1 {
		return js.Func{}, errors.New(
			"fn cannot return more than 1 parameter because Js can only returns 0 or 1",
		)
	}

	if fnType.NumOut() == 1 {
		switch ot := fnType.Out(0); ot.Kind() {
		case reflect.Bool:
		case reflect.Float32, reflect.Float64:
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		case reflect.String:
		case reflect.Map:
			if ot.Key().Kind() != reflect.String {
				return js.Func{}, errors.New("'fn' returns a map with not keys String type, Js only supports map[string]interface{}")
			}
			if el := ot.Elem(); el.Kind() != reflect.Interface || el.NumMethod() != 0 {
				return js.Func{}, errors.New("'fn' returns a map of values of not empty interface type, Js only supports map[string]interface{}")
			}
		case reflect.Slice:
			if el := ot.Elem(); el.Kind() != reflect.Interface || el.NumMethod() != 0 {
				return js.Func{}, errors.New("'fn' returns a slice of values of not empty interface type, Js only supports []interface{}")
			}
		default:
			return js.Func{}, errors.New("'fn' return parameter is of a type not convertible to Js")
		}
	}

	// create Js function
	fnVal := reflect.ValueOf(fn)
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var in []reflect.Value
		if (useThis && len(args) != (fnNumParamsIn-1)) ||
			(!useThis && len(args) != fnNumParamsIn) {
			panic("called a Go function with a different number of expected parameters")
		}

		if fnNumParamsIn > 0 {
			in = make([]reflect.Value, 0, fnNumParamsIn)
		}

		if useThis {
			in = append(in, reflect.ValueOf(this))
		}

		// start fn params at index 1 if useThis is true otherwise at index 0
		for i, j := 0, len(in); i < len(args); i, j = i+1, j+1 {
			switch fnType.In(j).Kind() {
			case reflect.Bool:
				if args[i].Type() != js.TypeBoolean {
					panic(fmt.Sprintf("argument %d must be of type bool", i))
				}
				in = append(in, reflect.ValueOf(args[i].Bool()))
			case reflect.Float32:
				if args[i].Type() != js.TypeNumber {
					panic(fmt.Sprintf("argument %d must be of type number", i))
				}
				in = append(in, reflect.ValueOf(float32(args[i].Float())))
			case reflect.Float64:
				if args[i].Type() != js.TypeNumber {
					panic(fmt.Sprintf("argument %d must be of type number", i))
				}
				in = append(in, reflect.ValueOf(args[i].Float()))
			case reflect.Int:
				if args[i].Type() != js.TypeNumber {
					panic(fmt.Sprintf("argument %d must be of type number", i))
				}
				in = append(in, reflect.ValueOf(args[i].Int()))
			case reflect.Int8:
				if args[i].Type() != js.TypeNumber {
					panic(fmt.Sprintf("argument %d must be of type number", i))
				}
				in = append(in, reflect.ValueOf(int8(args[i].Int())))
			case reflect.Int16:
				if args[i].Type() != js.TypeNumber {
					panic(fmt.Sprintf("argument %d must be of type number", i))
				}
				in = append(in, reflect.ValueOf(int16(args[i].Int())))
			case reflect.Int32:
				if args[i].Type() != js.TypeNumber {
					panic(fmt.Sprintf("argument %d must be of type number", i))
				}
				in = append(in, reflect.ValueOf(int32(args[i].Int())))
			case reflect.Int64:
				if args[i].Type() != js.TypeNumber {
					panic(fmt.Sprintf("argument %d must be of type number", i))
				}
				in = append(in, reflect.ValueOf(int64(args[i].Int())))
			case reflect.Uint:
				if args[i].Type() != js.TypeNumber {
					panic(fmt.Sprintf("argument %d must be of type number", i))
				}
				in = append(in, reflect.ValueOf(uint(args[i].Int())))
			case reflect.Uint8:
				if args[i].Type() != js.TypeNumber {
					panic(fmt.Sprintf("argument %d must be of type number", i))
				}
				in = append(in, reflect.ValueOf(uint8(args[i].Int())))
			case reflect.Uint16:
				if args[i].Type() != js.TypeNumber {
					panic(fmt.Sprintf("argument %d must be of type number", i))
				}
				in = append(in, reflect.ValueOf(uint16(args[i].Int())))
			case reflect.Uint32:
				if args[i].Type() != js.TypeNumber {
					panic(fmt.Sprintf("argument %d must be of type number", i))
				}
				in = append(in, reflect.ValueOf(uint32(args[i].Int())))
			case reflect.Uint64:
				if args[i].Type() != js.TypeNumber {
					panic(fmt.Sprintf("argument %d must be of type number", i))
				}
				in = append(in, reflect.ValueOf(uint64(args[i].Int())))
			case reflect.String:
				if args[i].Type() != js.TypeString {
					panic(fmt.Sprintf("argument %d must be of type string", i))
				}
				in = append(in, reflect.ValueOf(args[i].String()))
			}
		}

		outVals := fnVal.Call(in)
		if len(outVals) > 0 {
			switch out := outVals[0]; out.Kind() {
			case reflect.Bool:
				return out.Bool()
			case reflect.Float32, reflect.Float64:
				return out.Float()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return out.Int()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return out.Uint()
			case reflect.String:
				return out.String()
			default:
				return out.Interface()
			}
		}

		return nil
	}), nil
}
