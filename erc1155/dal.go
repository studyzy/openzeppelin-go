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

package erc1155

import (
	"bytes"
	"errors"

	"github.com/studyzy/openzeppelin-go/common"
)

const (
	balanceKey          = "b"
	operatorApprovalKey = "o"
	uri                 = "uri"
	adminKey            = "admin"
)

type ERC1155Dal struct {
	sdk common.StateOperator
}

func NewERC20ContractDAL(sdk common.StateOperator) *ERC1155Dal {
	return &ERC1155Dal{sdk: sdk}
}

func (c *ERC1155Dal) GetUint256(key string) (*common.SafeUint256, error) {
	fromBalStr, err := c.sdk.GetState(key)
	if err != nil {
		return nil, err
	}
	fromBalance, pass := common.ParseSafeUint256(string(fromBalStr))
	if !pass {
		return nil, errors.New("invalid uint256 data")
	}
	return fromBalance, nil
}

func (c *ERC1155Dal) GetBalance(token *common.SafeUint256, account common.Account) (*common.SafeUint256, error) {
	key, err := c.sdk.CreateCompositeKey(balanceKey, token.ToString(), account.ToString())
	if err != nil {
		return nil, err
	}
	return c.GetUint256(key)
}
func (c *ERC1155Dal) SetBalance(token *common.SafeUint256, account common.Account, amount *common.SafeUint256) error {
	key, err := c.sdk.CreateCompositeKey(balanceKey, token.ToString(), account.ToString())
	if err != nil {
		return err
	}
	return c.sdk.PutState(key, []byte(amount.ToString()))
}
func (c *ERC1155Dal) SetOperatorApproval(owner common.Account, operator common.Account, approved bool) error {
	value := []byte("false")
	if approved {
		value = []byte("true")
	}
	key, _ := c.sdk.CreateCompositeKey(operatorApprovalKey, owner.ToString(), operator.ToString())
	return c.sdk.PutState(key, value)
}
func (c *ERC1155Dal) GetOperatorApproval(owner common.Account, operator common.Account) (bool, error) {
	key, _ := c.sdk.CreateCompositeKey(operatorApprovalKey, owner.ToString(), operator.ToString())

	b, err := c.sdk.GetState(key)
	if err != nil {
		return false, err
	}
	return bytes.Equal(b, []byte("true")), nil
}

func (c *ERC1155Dal) GetAdmin() (common.Account, error) {
	b, err := c.sdk.GetState(adminKey)
	if err != nil {
		return nil, err
	}
	return c.sdk.NewAccountFromBytes(b)
}
func (c *ERC1155Dal) SetAdmin(admin common.Account) error {
	return c.sdk.PutState(adminKey, admin.Bytes())
}
func bytes2String(b []byte, err error) (string, error) {
	return string(b), err

}
func (c *ERC1155Dal) GetUri() (string, error) {
	return bytes2String(c.sdk.GetState(uri))
}
func (c *ERC1155Dal) SetUri(uri string) error {
	return c.sdk.PutState(uri, []byte(uri))
}
