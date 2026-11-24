// Package interfaces
// @Title  title
// @Description  desc
package interfaces

type IDataDef interface {
	IsRef() bool
	Ref()
	UnRef() bool
}

type IReset interface {
	Reset()
}
