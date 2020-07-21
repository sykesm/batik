// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package protomsg

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
)

func toMessageSlice(in ...interface{}) ([]proto.Message, error) {
	inval := reflect.ValueOf(in)
	if reflect.ValueOf(in[0]).Kind() == reflect.Slice && inval.Len() == 1 {
		inval = reflect.ValueOf(in[0])
	}

	messages := make([]proto.Message, inval.Len(), inval.Len())
	for i := 0; i < inval.Len(); i++ {
		elem := inval.Index(i).Interface()
		if elem == nil {
			messages[i] = nil
			break
		}
		m, ok := elem.(proto.Message)
		if !ok {
			return nil, fmt.Errorf("protomsg: index %d of type %T is not a proto.Message", i, elem)
		}
		messages[i] = m
	}
	return messages, nil
}
