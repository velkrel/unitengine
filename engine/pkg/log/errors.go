/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package log
// 模块名: 模块名
// 功能描述: 描述
package log

import "errors"

var (
	ErrRotationTime = errors.New("rotationTime must >= 1min and <= 24hour")
	ErrLevel        = errors.New("log level must <= 6")
)
