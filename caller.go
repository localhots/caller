// Package caller is used to dynamically call functions with data unmarshalled
// into the functions' first argument. Its main purpose is to hide common
// unmarshalling code from each function implementation thus reducing
// boilerplate and making package interaction code sexier.
package caller

import (
	"encoding/json"
	"errors"
	"reflect"
)

// Caller wraps a function and makes it ready to be dynamically called.
type Caller struct {
	// Unmarshaller is a BYOB unmarshaller function. By default it uses JSON.
	Unmarshaller func(data []byte, v interface{}) error
	fun          reflect.Value
	argtyp       reflect.Type
}

var (
	// ErrInvalidFunctionType is an error that is returned by the New function
	// when its argument is not a function.
	ErrInvalidFunctionType = errors.New("argument must be function")
	// ErrInvalidFunctionInArguments is an error that is returned by the New
	// function when its argument-function has a number of input arguments other
	// than 1.
	ErrInvalidFunctionInArguments = errors.New("function must have only one input argument")
	// ErrInvalidFunctionOutArguments is an error that is returned by the New
	// function when its argument-function returs any values.
	ErrInvalidFunctionOutArguments = errors.New("function must not have output arguments")
)

// New creates a new Caller instance using the function given as an argument.
// It returns the Caller instance and an error if something is wrong with the
// argument-function.
func New(fun interface{}) (c *Caller, err error) {
	fval := reflect.ValueOf(fun)
	ftyp := reflect.TypeOf(fun)
	if ftyp.Kind() != reflect.Func {
		return nil, ErrInvalidFunctionType
	}
	if ftyp.NumIn() != 1 {
		return nil, ErrInvalidFunctionInArguments
	}
	if ftyp.NumOut() != 0 {
		return nil, ErrInvalidFunctionOutArguments
	}

	c = &Caller{
		Unmarshaller: json.Unmarshal,
		fun:          fval,
		argtyp:       ftyp.In(0),
	}

	return c, nil
}

// Call creates an instance of the Caller function's argument type, unmarshalls
// the payload into it and dynamically calls the Caller function with this
// instance.
func (c *Caller) Call(data []byte) error {
	val, err := c.unmarshal(data)
	if err != nil {
		return err
	}

	c.makeDynamicCall(val)
	return nil
}

func (c *Caller) unmarshal(data []byte) (val reflect.Value, err error) {
	val = c.newValue()
	err = c.Unmarshaller(data, val.Interface())
	return
}

func (c *Caller) makeDynamicCall(val reflect.Value) {
	c.fun.Call([]reflect.Value{val.Elem()})
}

func (c *Caller) newValue() reflect.Value {
	return reflect.New(c.argtyp)
}
