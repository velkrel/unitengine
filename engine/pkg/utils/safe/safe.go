// Package safe
// 模块名: 模块名
// 功能描述: 描述
package safe

import (
	"fmt"
	"runtime/debug"
)

func Do(f func() error) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("safe.Do panic, err:%v\ntrace:%s", r, debug.Stack())
		}
	}()
	err = f()
	return err
}
