// Copyright 2023 Jetpack Technologies Inc and contributors. All rights reserved.
// Use of this source code is governed by the license in the LICENSE file.

package envir

import (
	"os"
	"strconv"
)

func IsCI() bool {
	ci, err := strconv.ParseBool(os.Getenv("CI"))
	return ci && err == nil
}

// GetValueOrDefault gets the value of an environment variable.
// If it's empty, it will return the given default value instead.
func GetValueOrDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		val = def
	}

	return val
}
