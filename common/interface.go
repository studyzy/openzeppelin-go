// Copyright 2023 studyzy Author
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import "errors"

type ChainBase interface {
	NewAccountFromBytes(b []byte) (Account, error)
	NewAccountFromString(str string) (Account, error)
	NewZeroAccount() Account
}
type StateOperator interface {
	ChainBase
	GetState(key string) (value []byte, err error)
	PutState(key string, value []byte) error
	DelState(key string) error
	CreateCompositeKey(prefix string, data ...string) (string, error)
}
type ContractSDK interface {
	StateOperator
	GetTxSender() (Account, error)
	EmitEvent(topic string, data ...string) error
	IsContract(account Account) bool
	CallContract(account Account, method string, args []KeyValue) Response
}

func Require(exp bool, msg string) error {
	if !exp {
		return errors.New(msg)
	}
	return nil
}
