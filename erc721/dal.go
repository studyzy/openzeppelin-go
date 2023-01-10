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

package erc721

import (
	"bytes"
	"errors"

	"github.com/studyzy/openzeppelin-go/common"
)

const (
	balanceKey          = "b"
	tokenApprovalKey    = "a"
	operatorApprovalKey = "o"
	nameKey             = "name"
	symbolKey           = "symbol"
	adminKey            = "admin"
	tokenOwnerKey       = "t"
	baseURIKey          = "uri"
)

type ERC721DAL struct {
	sdk common.StateOperator
}

func NewERC20ContractDAL(sdk common.StateOperator) *ERC721DAL {
	return &ERC721DAL{sdk: sdk}
}

func (c *ERC721DAL) GetUint256(key string) (*common.SafeUint256, error) {
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

func (c *ERC721DAL) GetBalance(account common.Account) (*common.SafeUint256, error) {
	return c.GetUint256(balanceKey + account.ToString())
}
func (c *ERC721DAL) SetBalance(account common.Account, amount *common.SafeUint256) error {
	return c.sdk.PutState(balanceKey+account.ToString(), []byte(amount.ToString()))
}
func (c *ERC721DAL) IncreaseBalance(account common.Account) error {
	bal, err := c.GetBalance(account)
	if err != nil {
		bal = common.NewSafeUint256(0)
	}
	bal, _ = common.SafeAdd(bal, common.NewSafeUint256(1))
	err = c.SetBalance(account, bal)
	return err
}
func (c *ERC721DAL) DecreaseBalance(account common.Account) error {
	bal, err := c.GetBalance(account)
	if err != nil {
		bal = common.NewSafeUint256(0)
	}
	bal, _ = common.SafeSub(bal, common.NewSafeUint256(1))
	err = c.SetBalance(account, bal)
	return err
}
func (c *ERC721DAL) SetTokenApproval(tokenId *common.SafeUint256, spender common.Account) error {
	return c.sdk.PutState(tokenApprovalKey+tokenId.ToString(), spender.Bytes())
}
func (c *ERC721DAL) GetTokenApproval(tokenId *common.SafeUint256) (common.Account, error) {
	b, err := c.sdk.GetState(tokenApprovalKey + tokenId.ToString())
	if err != nil {
		return nil, err
	}
	return c.sdk.NewAccountFromBytes(b)
}
func (c *ERC721DAL) DeleteTokenApproval(tokenId *common.SafeUint256) error {
	return c.sdk.DelState(tokenApprovalKey + tokenId.ToString())
}

func (c *ERC721DAL) GetName() (string, error) {
	return bytes2String(c.sdk.GetState(nameKey))
}
func (c *ERC721DAL) SetName(name string) error {
	return c.sdk.PutState(nameKey, []byte(name))
}
func (c *ERC721DAL) GetSymbol() (string, error) {
	return bytes2String(c.sdk.GetState(symbolKey))
}
func (c *ERC721DAL) SetSymbol(symbol string) error {
	return c.sdk.PutState(symbolKey, []byte(symbol))
}

func (c *ERC721DAL) GetAdmin() (common.Account, error) {
	b, err := c.sdk.GetState(adminKey)
	if err != nil {
		return nil, err
	}
	return c.sdk.NewAccountFromBytes(b)
}
func (c *ERC721DAL) SetAdmin(admin common.Account) error {
	return c.sdk.PutState(adminKey, admin.Bytes())
}
func bytes2String(b []byte, err error) (string, error) {
	return string(b), err
}

func (c *ERC721DAL) SetTokenOwner(tokenId *common.SafeUint256, owner common.Account) error {
	return c.sdk.PutState(tokenOwnerKey+tokenId.ToString(), owner.Bytes())
}
func (c *ERC721DAL) GetTokenOwner(tokenId *common.SafeUint256) (common.Account, error) {
	b, err := c.sdk.GetState(tokenOwnerKey + tokenId.ToString())
	if err != nil {
		return nil, err
	}
	return c.sdk.NewAccountFromBytes(b)
}
func (c *ERC721DAL) DeleteTokenOwner(tokenId *common.SafeUint256) error {
	return c.sdk.DelState(tokenOwnerKey + tokenId.ToString())
}
func (c *ERC721DAL) SetOperatorApproval(owner common.Account, operator common.Account, approved bool) error {
	value := []byte("false")
	if approved {
		value = []byte("true")
	}
	key, _ := c.sdk.CreateCompositeKey(operatorApprovalKey, owner.ToString(), operator.ToString())
	return c.sdk.PutState(key, value)
}
func (c *ERC721DAL) GetOperatorApproval(owner common.Account, operator common.Account) (bool, error) {
	key, _ := c.sdk.CreateCompositeKey(operatorApprovalKey, owner.ToString(), operator.ToString())

	b, err := c.sdk.GetState(key)
	if err != nil {
		return false, err
	}
	return bytes.Equal(b, []byte("true")), nil
}
func (c *ERC721DAL) GetBaseURI() (string, error) {
	return bytes2String(c.sdk.GetState(baseURIKey))
}
func (c *ERC721DAL) SetBaseURI(name string) error {
	return c.sdk.PutState(baseURIKey, []byte(name))
}
