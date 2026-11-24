//go:build linux || ignore || !windows
// +build linux ignore !windows

/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package title
// 模块名:
// 模块功能简介:
package title

import (
	"fmt"
)

var titleBase = `
unitengine
%s | %s: %s
`

var bakUrl = "https://patorjk.com/software/taag/#p=display&f=3D%20Diagonal" // 可以在这里做新的title

func EchoByeBye() {
	fmt.Println("bye!")
}
