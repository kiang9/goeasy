package assert

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// Equal -
func Equal(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n\t   %#v (expected)\n\n\t!= %#v (actual)\n\n",
			filepath.Base(file), line, expected, actual)
		t.FailNow()
	}
}

// NotEqual -
func NotEqual(t *testing.T, expected, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n\tnexp: %#v\n\n\tgot:  %#v\n\n",
			filepath.Base(file), line, expected, actual)
		t.FailNow()
	}
}

// NotNil -
func NotNil(t *testing.T, obj interface{}) {
	if isNil(obj) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n\tExpected value not to be <nil>\n\n",
			filepath.Base(file), line)
		t.FailNow()
	}
}

// Nil -
func Nil(t *testing.T, obj interface{}) {
	if !isNil(obj) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n\t   <nil> (expected)\n\n\t!= %#v (actual)\n\n",
			filepath.Base(file), line, obj)
		t.FailNow()
	}
}

func isNil(obj interface{}) bool {
	if obj == nil {
		return true
	}

	value := reflect.ValueOf(obj)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}
	return false
}

// Len -
func Len(t *testing.T, obj interface{}, length int) {
	ok, l := getLen(obj)
	if !ok {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n\t   can not get length of %#v\n\n\t \n\n",
			filepath.Base(file), line, obj)
		t.FailNow()
	}
	if l != length {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n\t   %#v (expected)\n\n\t!= %#v (actual)\n\n",
			filepath.Base(file), line, length, l)
		t.FailNow()
	}
}

//getLen try to get length of obj
//return (false, 0) if impossible
func getLen(x interface{}) (ok bool, length int) {
	v := reflect.ValueOf(x)
	defer func() {
		if err := recover(); err != nil {
			ok = false
		}
	}()
	return true, v.Len()
}

// True -
func True(t *testing.T, value bool) {
	if !value {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n\t   true (expected)\n\n\t!= false (actual)\n\n",
			filepath.Base(file), line)
		t.FailNow()
	}
}

// False -
func False(t *testing.T, value bool) {
	if value {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n\t   false (expected)\n\n\t!= true (actual)\n\n",
			filepath.Base(file), line)
		t.FailNow()
	}
}

// Error -
func Error(t *testing.T, err error) bool {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n\t   error (expected)\n\n\t!= <nil> (actual)\n\n",
			filepath.Base(file), line)
		t.FailNow()
		return false
	}
	return true
}

// NoError -
func NoError(t *testing.T, err error) bool {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n\t   <nil> (expected)\n\n\t!= error:%v (actual)\n\n",
			filepath.Base(file), line, err)
		t.FailNow()
		return false
	}
	return true
}
