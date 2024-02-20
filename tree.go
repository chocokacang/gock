package gock

import (
	"bytes"

	"github.com/chocokacang/gock/bytesconv"
)

var (
	strColon = []byte(":")
	strStar  = []byte("*")
	strSlash = []byte("/")
)

type trees map[string]*node

type nodeType uint8

type Param struct {
	Name  string
	Value string
}

type Params []Param

func (ps Params) Get(name string) (string, bool) {
	for _, entry := range ps {
		if entry.Name == name {
			return entry.Value, true
		}
	}
	return "", false
}

func (ps Params) ByName(name string) (value string) {
	value, _ = ps.Get(name)
	return
}

const (
	static nodeType = iota
	root
	param
	catchAll
)

type node struct {
}

type nodeValue struct {
}

type skippedNode struct {
}

func countParams(path string) uint16 {
	var n uint16
	s := bytesconv.StringToBytes(path)
	n += uint16(bytes.Count(s, strColon))
	n += uint16(bytes.Count(s, strStar))
	return n
}

func countSections(path string) uint16 {
	s := bytesconv.StringToBytes(path)
	return uint16(bytes.Count(s, strSlash))
}

func (n *node) addRoute(path string, handlres Handlers) {

}
