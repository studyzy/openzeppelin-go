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

package erc20

import (
	"errors"
	"strconv"

	"github.com/studyzy/openzeppelin-go/common"
)

const (
	balanceKey     = "b"
	allowanceKey   = "a"
	totalSupplyKey = "totalSupplyKey"
	nameKey        = "name"
	symbolKey      = "symbol"
	decimalKey     = "decimal"
	adminKey       = "admin"
)

type ERC20ContractDAL struct {
	sdk common.StateOperator
}

func NewERC20ContractDAL(sdk common.StateOperator) *ERC20ContractDAL {
	return &ERC20ContractDAL{sdk: sdk}
}

func (c *ERC20ContractDAL) GetUint256(key string) (*common.SafeUint256, error) {
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

func (c *ERC20ContractDAL) GetBalance(account common.Account) (*common.SafeUint256, error) {
	return c.GetUint256(balanceKey + account.ToString())
}
func (c *ERC20ContractDAL) SetBalance(account common.Account, amount *common.SafeUint256) error {
	return c.sdk.PutState(balanceKey+account.ToString(), []byte(amount.ToString()))
}
func (c *ERC20ContractDAL) SetAllowance(owner common.Account, spender common.Account, amount *common.SafeUint256) error {
	key, _ := c.sdk.CreateCompositeKey(allowanceKey, owner.ToString(), spender.ToString())
	return c.sdk.PutState(key, []byte(amount.ToString()))
}
func (c *ERC20ContractDAL) GetAllowance(owner common.Account, spender common.Account) (*common.SafeUint256, error) {
	key, _ := c.sdk.CreateCompositeKey(allowanceKey, owner.ToString(), spender.ToString())
	return c.GetUint256(key)
}
func (c *ERC20ContractDAL) GetTotalSupply() (*common.SafeUint256, error) {
	return c.GetUint256(totalSupplyKey)
}
func (c *ERC20ContractDAL) SetTotalSupply(amount *common.SafeUint256) error {
	return c.sdk.PutState(totalSupplyKey, []byte(amount.ToString()))
}
func (c *ERC20ContractDAL) GetName() (string, error) {
	return bytes2String(c.sdk.GetState(nameKey))
}
func (c *ERC20ContractDAL) SetName(name string) error {
	return c.sdk.PutState(nameKey, []byte(name))
}
func (c *ERC20ContractDAL) GetSymbol() (string, error) {
	return bytes2String(c.sdk.GetState(symbolKey))
}
func (c *ERC20ContractDAL) SetSymbol(symbol string) error {
	return c.sdk.PutState(symbolKey, []byte(symbol))
}
func (c *ERC20ContractDAL) GetDecimals() (uint8, error) {
	d, err := c.sdk.GetState(decimalKey)
	if err != nil {
		return 0, err
	}
	decimal, err := strconv.ParseUint(string(d), 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(decimal), nil
}
func (c *ERC20ContractDAL) SetDecimals(decimal uint8) error {
	return c.sdk.PutState(decimalKey, []byte(strconv.Itoa(int(decimal))))
}

func (c *ERC20ContractDAL) GetAdmin() (common.Account, error) {
	b, err := c.sdk.GetState(adminKey)
	if err != nil {
		return nil, err
	}
	return c.sdk.NewAccountFromBytes(b)
}
func (c *ERC20ContractDAL) SetAdmin(admin common.Account) error {
	return c.sdk.PutState(adminKey, admin.Bytes())
}
func bytes2String(b []byte, err error) (string, error) {
	return string(b), err

}
