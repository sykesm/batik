// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	validationv1 "github.com/sykesm/batik/pkg/pb/validation/v1"
)

func ToState(in *txv1.State, txID ID, index uint64) *State {
	if in == nil {
		return nil
	}
	return &State{
		ID:        StateID{OutputIndex: index, TxID: txID},
		StateInfo: ToStateInfo(in.Info),
		Data:      in.State,
	}
}

func ToStates(txID ID, in ...*txv1.State) []*State {
	var states []*State
	for i := range in {
		states = append(states, ToState(in[i], txID, uint64(i)))
	}
	return states
}

func FromState(in *State) *txv1.State {
	if in == nil {
		return nil
	}
	return &txv1.State{
		Info:  FromStateInfo(in.StateInfo),
		State: in.Data,
	}
}

func FromStates(in ...*State) []*txv1.State {
	var states []*txv1.State
	for i := range in {
		states = append(states, FromState(in[i]))
	}
	return states
}

func ToStateInfo(in *txv1.StateInfo) *StateInfo {
	if in == nil {
		return nil
	}
	return &StateInfo{
		Owners: ToParties(in.Owners...),
		Kind:   in.Kind,
	}
}

func FromStateInfo(in *StateInfo) *txv1.StateInfo {
	if in == nil {
		return nil
	}
	return &txv1.StateInfo{
		Owners: FromParties(in.Owners...),
		Kind:   in.Kind,
	}
}

func ToStateID(in *txv1.StateReference) *StateID {
	if in == nil {
		return nil
	}
	return &StateID{
		TxID:        NewID(in.Txid),
		OutputIndex: in.OutputIndex,
	}
}

func FromStateID(in *StateID) *txv1.StateReference {
	if in == nil {
		return nil
	}
	return &txv1.StateReference{
		Txid:        in.TxID.Bytes(),
		OutputIndex: in.OutputIndex,
	}
}

func ToStateIDs(in ...*txv1.StateReference) []*StateID {
	var ids []*StateID
	for i := range in {
		ids = append(ids, ToStateID(in[i]))
	}
	return ids
}

func FromStateIDs(in ...*StateID) []*txv1.StateReference {
	var ids []*txv1.StateReference
	for i := range in {
		ids = append(ids, FromStateID(in[i]))
	}
	return ids
}

func ToParty(in *txv1.Party) *Party {
	if in == nil {
		return nil
	}
	return &Party{
		PublicKey: in.PublicKey,
	}
}

func FromParty(in *Party) *txv1.Party {
	if in == nil {
		return nil
	}
	return &txv1.Party{
		PublicKey: in.PublicKey,
	}
}

func ToParties(in ...*txv1.Party) []*Party {
	var parties []*Party
	for i := range in {
		parties = append(parties, ToParty(in[i]))
	}
	return parties
}

func FromParties(in ...*Party) []*txv1.Party {
	var parties []*txv1.Party
	for i := range in {
		parties = append(parties, FromParty(in[i]))
	}
	return parties
}

func ToParameter(in *txv1.Parameter) *Parameter {
	if in == nil {
		return nil
	}
	return &Parameter{
		Name:  in.Name,
		Value: in.Value,
	}
}

func FromParameter(in *Parameter) *txv1.Parameter {
	if in == nil {
		return nil
	}
	return &txv1.Parameter{
		Name:  in.Name,
		Value: in.Value,
	}
}

func ToParameters(in ...*txv1.Parameter) []*Parameter {
	var parameters []*Parameter
	for i := range in {
		parameters = append(parameters, ToParameter(in[i]))
	}
	return parameters
}

func FromParameters(in ...*Parameter) []*txv1.Parameter {
	var parameters []*txv1.Parameter
	for i := range in {
		parameters = append(parameters, FromParameter(in[i]))
	}
	return parameters
}

func ToSignature(in *txv1.Signature) *Signature {
	if in == nil {
		return nil
	}
	return &Signature{
		PublicKey: in.PublicKey,
		Signature: in.Signature,
	}
}

func FromSignature(in *Signature) *txv1.Signature {
	if in == nil {
		return nil
	}
	return &txv1.Signature{
		PublicKey: in.PublicKey,
		Signature: in.Signature,
	}
}

func ToSignatures(in ...*txv1.Signature) []*Signature {
	var sigs []*Signature
	for i := range in {
		sigs = append(sigs, ToSignature(in[i]))
	}
	return sigs
}

func FromSignatures(in ...*Signature) []*txv1.Signature {
	var sigs []*txv1.Signature
	for i := range in {
		sigs = append(sigs, FromSignature(in[i]))
	}
	return sigs
}

func ResolvedToState(in *validationv1.ResolvedState) *State {
	if in == nil {
		return nil
	}
	return &State{
		ID:        *ToStateID(in.Reference),
		StateInfo: ToStateInfo(in.State.Info),
		Data:      in.State.State,
	}
}

func ResolvedToStates(in ...*validationv1.ResolvedState) []*State {
	var states []*State
	for i := range in {
		states = append(states, ResolvedToState(in[i]))
	}
	return states
}

func ResolvedFromState(in *State) *validationv1.ResolvedState {
	if in == nil {
		return nil
	}
	return &validationv1.ResolvedState{
		Reference: FromStateID(&in.ID),
		State:     FromState(in),
	}
}

func ResolvedFromStates(in ...*State) []*validationv1.ResolvedState {
	var resolved []*validationv1.ResolvedState
	for i := range in {
		resolved = append(resolved, ResolvedFromState(in[i]))
	}
	return resolved
}

func ToResolved(in *validationv1.ResolvedTransaction) *Resolved {
	if in == nil {
		return nil
	}
	return &Resolved{
		ID:              in.Txid,
		Inputs:          ResolvedToStates(in.Inputs...),
		References:      ResolvedToStates(in.References...),
		Outputs:         ToStates(in.Txid, in.Outputs...),
		Parameters:      ToParameters(in.Parameters...),
		RequiredSigners: ToParties(in.RequiredSigners...),
		Signatures:      ToSignatures(in.Signatures...),
	}
}

func FromResolved(in *Resolved) *validationv1.ResolvedTransaction {
	if in == nil {
		return nil
	}
	return &validationv1.ResolvedTransaction{
		Txid:            in.ID,
		Inputs:          ResolvedFromStates(in.Inputs...),
		References:      ResolvedFromStates(in.References...),
		Outputs:         FromStates(in.Outputs...),
		Parameters:      FromParameters(in.Parameters...),
		RequiredSigners: FromParties(in.RequiredSigners...),
		Signatures:      FromSignatures(in.Signatures...),
	}
}
