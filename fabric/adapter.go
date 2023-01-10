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

package fabric

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/studyzy/openzeppelin-go/common"
)

var _ common.Account = (*MspUser)(nil)

type MspUser struct {
	ID string
}

func (m MspUser) IsZero() bool {
	return m.ID == "0x0"
}

func NewMspUser(id string) common.Account {
	return &MspUser{ID: id}
}
func (m MspUser) ToString() string {
	return m.ID
}

func (m MspUser) Bytes() []byte {
	return []byte(m.ID)
}

func (m MspUser) Equal(account common.Account) bool {
	return m.ID == account.ToString()
}

var _ common.ContractSDK = (*SdkAdapter)(nil)

type SdkAdapter struct {
	ctx           contractapi.TransactionContextInterface
	eventEncoder  func(string, ...string) ([]byte, error)
	contractExist func(string) (bool, error)
}

func (s SdkAdapter) NewAccountFromBytes(b []byte) (common.Account, error) {
	return NewMspUser(string(b)), nil
}

func (s SdkAdapter) NewAccountFromString(str string) (common.Account, error) {
	return NewMspUser(str), nil

}

func (s SdkAdapter) NewZeroAccount() common.Account {
	return NewMspUser("0x0")
}

func (s SdkAdapter) DelState(key string) error {
	return s.ctx.GetStub().DelState(key)
}

func (s SdkAdapter) CreateCompositeKey(prefix string, data ...string) (string, error) {
	return s.ctx.GetStub().CreateCompositeKey(prefix, data)
}

func (s SdkAdapter) IsContract(account common.Account) bool {
	exist, err := s.contractExist(account.ToString())
	if err != nil {
		return false
	}
	return exist
}

func (s SdkAdapter) CallContract(account common.Account, method string, args []common.KeyValue) common.Response {
	convert := func(args []common.KeyValue) [][]byte {
		result := make([][]byte, len(args))
		for i, arg := range args {
			result[i] = arg.Value
		}
		return result
	}
	response := s.ctx.GetStub().InvokeChaincode(account.ToString(), convert(args), "")
	return common.Response{
		Status:  response.Status,
		Message: response.Message,
		Payload: response.Payload,
	}
}

func NewSDkAdapter(ctx contractapi.TransactionContextInterface,
	eventEncoder func(string, ...string) ([]byte, error),
	contractExist func(string) (bool, error)) *SdkAdapter {
	return &SdkAdapter{ctx: ctx, eventEncoder: eventEncoder, contractExist: contractExist}
}
func (s SdkAdapter) GetState(key string) (value []byte, err error) {
	return s.ctx.GetStub().GetState(key)
}

func (s SdkAdapter) PutState(key string, value []byte) error {
	return s.ctx.GetStub().PutState(key, value)
}

func (s SdkAdapter) GetTxSender() (common.Account, error) {
	id, err := s.ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, err
	}
	return NewMspUser(id), nil
}

func (s SdkAdapter) EmitEvent(topic string, data ...string) error {
	payload, err := s.eventEncoder(topic, data...)
	if err != nil {
		return err
	}
	return s.ctx.GetStub().SetEvent(topic, payload)
}
