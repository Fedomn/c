package main

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func Equals(tb testing.TB, msg string, wat, got interface{}) {
	if !reflect.DeepEqual(wat, got) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: %s \n\n\twat: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, msg, wat, got)
		tb.FailNow()
	}
}
