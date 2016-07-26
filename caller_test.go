package caller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

//
// Testing targets
//

type testMessage struct {
	Body string `json:"body"`
}

const testPayload = `{"body":"Success!"}`

func testFun(m testMessage) {
	fmt.Print(m.Body)
}

func testFunSilent(_ testMessage) {}

//
// Tests
//

func TestNewCallerSuccess(t *testing.T) {
	c, err := New(testFun)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if c == nil {
		t.Error("Expected an instance of Caller, got nil")
	}
}

func TestNewCallerWithNonFunc(t *testing.T) {
	c, err := New(1)
	if err != ErrInvalidFunctionType {
		t.Errorf("Expected ErrInvalidFunctionType, got: %v", err)
	}
	if c != nil {
		t.Error("Expected nil, got an instance of Caller")
	}
}

func TestNewCallerWithFuncMultipleArgs(t *testing.T) {
	fun := func(a, b int) {}
	c, err := New(fun)
	if err != ErrInvalidFunctionInArguments {
		t.Errorf("Expected ErrInvalidFunctionInArguments, got: %v", err)
	}
	if c != nil {
		t.Error("Expected nil, got an instance of Caller")
	}
}

func TestNewCallerWithFuncReturnValue(t *testing.T) {
	fun := func(a int) int { return 0 }
	c, err := New(fun)
	if err != ErrInvalidFunctionOutArguments {
		t.Errorf("Expected ErrInvalidFunctionOutArguments, got: %v", err)
	}
	if c != nil {
		t.Error("Expected nil, got an instance of Caller")
	}
}

func TestCallSuccess(t *testing.T) {
	c, err := New(testFun)
	if err != nil {
		t.Fatal(err.Error())
	}

	out := captureStdoutAround(func() {
		if err := c.Call([]byte(testPayload)); err != nil {
			t.Fatal(err.Error())
		}
	})

	if string(out) != "Success!" {
		t.Errorf("Expected output to be %q, got %q", "Success!", out)
	}
}

func TestCallFalure(t *testing.T) {
	c, _ := New(testFunSilent)

	err := c.Call([]byte("{"))
	if err == nil {
		t.Error("Expected unmarshalling error, got nil")
	}
}

func TestUnmarshalSuccess(t *testing.T) {
	c, _ := New(testFunSilent)

	_, err := c.unmarshal([]byte(testPayload))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestUnmarshalFailure(t *testing.T) {
	c, _ := New(testFunSilent)

	_, err := c.unmarshal([]byte("{"))
	if err == nil {
		t.Error("Expected unmarshalling error, got nil")
	}
}

func captureStdoutAround(f func()) []byte {
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	out, err := ioutil.ReadAll(r)
	if err != nil {
		os.Stdout = origStdout
		panic(err)
	}
	r.Close()
	os.Stdout = origStdout

	return out
}

//
// Benchmarks
//

func BenchmarkCaller(b *testing.B) {
	c, _ := New(testFunSilent)

	for i := 0; i < b.N; i++ {
		c.Call([]byte(testPayload))
	}
}

func BenchmarkNoCaller(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var msg testMessage
		json.Unmarshal([]byte(testPayload), &msg)
		testFunSilent(msg)
	}
}

func BenchmarkDynamicNew(b *testing.B) {
	c, _ := New(testFunSilent)

	for i := 0; i < b.N; i++ {
		_ = c.newValue()
	}
}

func BenchmarkStaticNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = testMessage{}
	}
}

func BenchmarkDynamicCall(b *testing.B) {
	c, _ := New(testFunSilent)
	val, _ := c.unmarshal([]byte(testPayload))

	for i := 0; i < b.N; i++ {
		c.makeDynamicCall(val)
	}
}

func BenchmarkStaticCall(b *testing.B) {
	var msg testMessage

	for i := 0; i < b.N; i++ {
		testFunSilent(msg)
	}
}

func BenchmarkUnmarshalIntoInterface(b *testing.B) {
	c, _ := New(testFunSilent)
	val := c.newValue()

	for i := 0; i < b.N; i++ {
		json.Unmarshal([]byte(testPayload), val.Interface())
	}
}

func BenchmarkUnmarshalIntoTypedValue(b *testing.B) {
	var msg testMessage

	for i := 0; i < b.N; i++ {
		json.Unmarshal([]byte(testPayload), &msg)
	}
}
