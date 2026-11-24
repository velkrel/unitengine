// Package dto
// @Title  title
// @Description  desc
package dto

type ConcurrentTaskCallback struct {
	Callback func(err error, args ...interface{})
	Args     []interface{}
}
