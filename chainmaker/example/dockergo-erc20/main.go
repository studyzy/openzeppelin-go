/*
  Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

  SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"errors"
	"fmt"
	"strconv"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"github.com/studyzy/openzeppelin-go/chainmaker"
	"github.com/studyzy/openzeppelin-go/common"
	"github.com/studyzy/openzeppelin-go/erc20"
)

type ERC20DockerGo struct {
	supper  *erc20.ERC20Contract
	methods map[string]func() protogo.Response
	adapter *chainmaker.SdkAdapter
}

func NewERC20DockerGo() *ERC20DockerGo {
	erc20ption := erc20.Option{
		BeforeTransfer: nil,
		AfterTransfer:  nil,
		Burnable:       true,
		Minable:        true,
	}
	adapter := chainmaker.NewSdkAdapter(sdk.Instance)
	contract := &ERC20DockerGo{methods: make(map[string]func() protogo.Response), adapter: adapter}
	contract.supper = erc20.NewERC20Contract(erc20ption, "TestToken", "TT", adapter)
	contract.registerMethods(erc20ption)
	return contract
}
func (erc20 *ERC20DockerGo) registerMethods(option erc20.Option) {

	erc20.RegisterMethod("name", erc20.name)
	erc20.RegisterMethod("symbol", erc20.symbol)
	erc20.RegisterMethod("decimals", erc20.decimals)
	erc20.RegisterMethod("totalSupply", erc20.totalSupply)
	erc20.RegisterMethod("balanceOf", erc20.balanceOf)
	erc20.RegisterMethod("transfer", erc20.transfer)
	erc20.RegisterMethod("allowance", erc20.allowance)
	erc20.RegisterMethod("approve", erc20.approve)
	erc20.RegisterMethod("transferFrom", erc20.transferFrom)
	if option.Minable {
		erc20.RegisterMethod("mint", erc20.mint)
	}
	if option.Burnable {
		erc20.RegisterMethod("burn", erc20.burn)
		erc20.RegisterMethod("burnFrom", erc20.burnFrom)
	}
}
func (c *ERC20DockerGo) RegisterMethod(methodName string, fun func() protogo.Response) {
	c.methods[methodName] = fun
}
func (c *ERC20DockerGo) InitContract() protogo.Response {
	err := c.updateErc20Info()
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte("Init contract success"))
}

// UpgradeContract used to upgrade contract
func (c *ERC20DockerGo) UpgradeContract() protogo.Response {
	err := c.updateErc20Info()
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte("Upgrade contract success"))
}

// UpgradeContract upgrade contract func
func (c *ERC20DockerGo) updateErc20Info() error {
	args := sdk.Instance.GetArgs()
	// name, symbol and decimal are optional
	name := string(args["name"])
	symbol := string(args["symbol"])
	decimalsStr := string(args["decimals"])
	decimal := uint8(18)
	//decimals default to 18
	if len(decimalsStr) > 0 {
		d, err := strconv.ParseUint(decimalsStr, 10, 8)
		if err != nil {
			return fmt.Errorf("param decimals err")
		}
		decimal = uint8(d)
	}
	totalSupplyStr := string(args["totalSupply"])
	totalSupplyValue := common.NewSafeUint256(0)
	if len(totalSupplyStr) > 0 {
		t, ok := common.ParseSafeUint256(totalSupplyStr)
		if !ok {
			return fmt.Errorf("param totalSupply err")
		}
		totalSupplyValue = t
	}
	admin, err := sdk.Instance.Sender()
	if err != nil {
		return fmt.Errorf("get sender failed, err:%s", err)
	}
	adminAccount, err := c.adapter.NewAccountFromString(admin)
	if err != nil {
		return fmt.Errorf("get sender failed, err:%s", err)
	}
	//此处支持在安装合约的时候指定name,symbol
	//如果没有参数指定，那么就使用NewERC20Contract构造的时候的值
	err = c.supper.InitERC20(name, symbol, decimal, totalSupplyValue, adminAccount)
	if err != nil {
		return fmt.Errorf("set admin failed, err:%s", err)
	}
	return nil
}

// InvokeContract used to invoke user contract
func (c *ERC20DockerGo) InvokeContract(method string) protogo.Response {
	if len(method) == 0 {
		return sdk.Error("method of param should not be empty")
	}
	if fun, ok := c.methods[method]; ok {
		return fun()
	}
	return sdk.Error("Invalid method")
}

func (erc20 *ERC20DockerGo) requireAccount(key string) (common.Account, error) {
	args := sdk.Instance.GetArgs()
	acc, ok := args[key]
	if !ok {
		return nil, errors.New("require account:" + key)
	}
	return erc20.adapter.NewAccountFromString(string(acc))
}
func (erc20 *ERC20DockerGo) requireAmount(key string) (*common.SafeUint256, error) {
	args := sdk.Instance.GetArgs()
	amt, ok := args[key]
	if !ok {
		return nil, errors.New("require account:" + key)
	}
	num, ok := common.ParseSafeUint256(string(amt))
	if !ok {
		return nil, errors.New("invalid uint256")
	}
	return num, nil
}

func (erc20 *ERC20DockerGo) name() protogo.Response {
	return chainmaker.ReturnString(erc20.supper.Name())
}

func (erc20 *ERC20DockerGo) symbol() protogo.Response {
	return chainmaker.ReturnString(erc20.supper.Symbol())
}

func (erc20 *ERC20DockerGo) decimals() protogo.Response {
	return chainmaker.ReturnUint8(erc20.supper.Decimals())
}

func (erc20 *ERC20DockerGo) totalSupply() protogo.Response {
	return chainmaker.ReturnUint256(erc20.supper.TotalSupply())
}

func (erc20 *ERC20DockerGo) balanceOf() protogo.Response {
	account, err := erc20.requireAccount("account")
	if err != nil {
		return sdk.Error(err.Error())
	}
	return chainmaker.ReturnUint256(erc20.supper.BalanceOf(account))
}

func (erc20 *ERC20DockerGo) transfer() protogo.Response {
	to, err := erc20.requireAccount("to")
	if err != nil {
		return sdk.Error(err.Error())
	}
	amt, err := erc20.requireAmount("amount")
	if err != nil {
		return sdk.Error(err.Error())
	}
	return chainmaker.ReturnBool(erc20.supper.Transfer(to, amt))
}

func (erc20 *ERC20DockerGo) allowance() protogo.Response {
	owner, err := erc20.requireAccount("owner")
	if err != nil {
		return sdk.Error(err.Error())
	}
	spender, err := erc20.requireAccount("spender")
	if err != nil {
		return sdk.Error(err.Error())
	}
	return chainmaker.ReturnUint256(erc20.supper.Allowance(owner, spender))

}

func (erc20 *ERC20DockerGo) approve() protogo.Response {
	spender, err := erc20.requireAccount("spender")
	if err != nil {
		return sdk.Error(err.Error())
	}
	amt, err := erc20.requireAmount("amount")
	if err != nil {
		return sdk.Error(err.Error())
	}
	return chainmaker.ReturnBool(erc20.supper.Approve(spender, amt))
}

func (erc20 *ERC20DockerGo) transferFrom() protogo.Response {
	from, err := erc20.requireAccount("from")
	if err != nil {
		return sdk.Error(err.Error())
	}
	to, err := erc20.requireAccount("to")
	if err != nil {
		return sdk.Error(err.Error())
	}
	amt, err := erc20.requireAmount("amount")
	if err != nil {
		return sdk.Error(err.Error())
	}
	return chainmaker.ReturnBool(erc20.supper.TransferFrom(from, to, amt))
}

func (erc20 *ERC20DockerGo) mint() protogo.Response {
	account, err := erc20.requireAccount("acount")
	if err != nil {
		return sdk.Error(err.Error())
	}
	amt, err := erc20.requireAmount("amount")
	if err != nil {
		return sdk.Error(err.Error())
	}
	return chainmaker.ReturnBool(erc20.supper.Mint(account, amt))
}

func (erc20 *ERC20DockerGo) burn() protogo.Response {
	amt, err := erc20.requireAmount("amount")
	if err != nil {
		return sdk.Error(err.Error())
	}
	return chainmaker.ReturnBool(erc20.supper.Burn(amt))
}

func (erc20 *ERC20DockerGo) burnFrom() protogo.Response {
	account, err := erc20.requireAccount("account")
	if err != nil {
		return sdk.Error(err.Error())
	}
	amt, err := erc20.requireAmount("amount")
	if err != nil {
		return sdk.Error(err.Error())
	}
	return chainmaker.ReturnBool(erc20.supper.BurnFrom(account, amt))
}

func main() {
	erc20 := NewERC20DockerGo()
	err := sandbox.Start(erc20)
	if err != nil {
		sdk.Instance.Errorf(err.Error())
	}
}
