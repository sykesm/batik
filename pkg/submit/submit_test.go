// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package submit

//go:generate counterfeiter -o fakes/repository.go --fake-name Repository . fakeRepository
type fakeRepository Repository // private import to prevent import cycle in generated fake
