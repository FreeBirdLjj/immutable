package main

import (
	"fmt"
	"unsafe"

	immutable_map "github.com/freebirdljj/immutable/map"
)

type (
	Map[K any, V any] struct{}
)

func main() {
	m := immutable_map.Map[int, int]{}
	fmt.Println(unsafe.Sizeof(m))
}
