// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package atexit

import "sync"

// AtExit maintains a stack of functions to call "at exit". This can be used to
// implement defer-like semantics across a user defined scope.
type AtExit struct {
	mutex    sync.Mutex
	handlers []func()
	once     sync.Once
}

// Create an instance of AtExit to orchestrate exit handlers.
func New() *AtExit {
	return &AtExit{}
}

// Register a handler to run on at Exit(). Handlers are exited LIFO order.
func (a *AtExit) Register(handler func()) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.handlers = append([]func(){handler}, a.handlers...)
}

// Exit runs exit handlers in a last-in-first-out order.
func (a *AtExit) Exit() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.once.Do(func() {
		for _, h := range a.handlers {
			h()
		}
	})
}
