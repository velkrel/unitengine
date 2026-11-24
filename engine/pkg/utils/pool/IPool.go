// Package pool
// @Title  title
// @Description  desc
package pool

type IPool[T any] interface {
	Get() T
	Put(obj T)
	Stats() *Stats
}
