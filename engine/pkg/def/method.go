// Package def
// @Title  title
// @Description  desc
package def

import "context"

type MethodCallFunc func(ctx context.Context, req interface{}) (interface{}, error)
