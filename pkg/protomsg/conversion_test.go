// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package protomsg

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestToMessageSlice(t *testing.T) {
	messages := []proto.Message{
		&emptypb.Empty{},
		&durationpb.Duration{Seconds: 999, Nanos: 888},
	}

	t.Run("ProtoSlice", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		result, err := toMessageSlice(messages)
		gt.Expect(err).NotTo(HaveOccurred())
		gt.Expect(result).To(Equal(messages))
	})

	t.Run("ProtoVariadic", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		result, err := toMessageSlice(messages[0], messages[1])
		gt.Expect(err).NotTo(HaveOccurred())
		gt.Expect(result).To(Equal(messages))
	})

	t.Run("ProtoElement", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		duration := &durationpb.Duration{Seconds: 987, Nanos: 123}

		result, err := toMessageSlice(duration)
		gt.Expect(err).NotTo(HaveOccurred())
		gt.Expect(result).To(Equal([]proto.Message{duration}))
	})

	t.Run("Nil", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		result, err := toMessageSlice(nil)
		gt.Expect(err).NotTo(HaveOccurred())
		gt.Expect(result).To(Equal([]proto.Message{nil}))
	})

	t.Run("BadTypes", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		_, err := toMessageSlice("bob")
		gt.Expect(err).To(MatchError("protomsg: index 0 of type string is not a proto.Message"))

		_, err = toMessageSlice([]interface{}{&emptypb.Empty{}, time.Now(), "fred"})
		gt.Expect(err).To(MatchError("protomsg: index 1 of type time.Time is not a proto.Message"))
	})
}
