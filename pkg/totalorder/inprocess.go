// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package totalorder

import (
	"context"

	"github.com/pkg/errors"
)

type InProcess struct {
	queue chan TXIDAndHMAC
	store *Store
	doneC chan struct{}
}

func NewInProcess(store *Store) *InProcess {
	ip := &InProcess{
		doneC: make(chan struct{}),
		queue: make(chan TXIDAndHMAC),
		store: store,
	}
	go ip.run()
	return ip
}

func (ip *InProcess) Stop() {
	close(ip.doneC)
}

func (ip *InProcess) run() {
	for {
		select {
		case t := <-ip.queue:
			ip.store.Append(t)
		case <-ip.doneC:
			return
		}
	}
}

func (ip *InProcess) Broadcast(ctx context.Context, t TXIDAndHMAC) error {
	select {
	case ip.queue <- t:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-ip.doneC:
		return errors.Errorf("told to exit") // TODO, turn into a sentinal error
	}
}

func (ip *InProcess) Deliver(ctx context.Context, seq uint64) (TXIDAndHMAC, error) {
	return ip.store.Get(ctx, seq)
}
