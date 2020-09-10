// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import "sync"

var (
	hLock    sync.RWMutex
	handlers []func()

	once sync.Once
)

// RegisterExitHandler prepends a handler func to the list of exit
// handlers to run on Exit().
func RegisterExitHandler(handler func()) {
	hLock.Lock()
	defer hLock.Unlock()

	handlers = append([]func(){handler}, handlers...)
}

// Exit runs all exit handlers in a last-in-first-out order to mimick
// how golang processes deferred statements.
func Exit() {
	hLock.RLock()
	defer hLock.RUnlock()

	once.Do(func() {
		for _, h := range handlers {
			h()
		}
	})
}
