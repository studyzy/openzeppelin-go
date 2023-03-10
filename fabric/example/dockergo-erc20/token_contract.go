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

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/studyzy/openzeppelin-go/common"
	"github.com/studyzy/openzeppelin-go/erc20"
	"github.com/studyzy/openzeppelin-go/fabric"
)

// SmartContract provides functions for transferring tokens between accounts
type SmartContract struct {
	contractapi.Contract
	erc20Contract erc20.ERC20Contract
}

// event provides an organized struct for emitting events
type transfer struct {
	from  string
	to    string
	value int
}
type approval struct {
	owner   string
	spender string
	value   int
}

func encodeEvent(topic string, data ...string) ([]byte, error) {
	var payload []byte
	var err error
	if topic == "transfer" {
		val, _ := strconv.Atoi(data[2])
		transferEvent := transfer{data[0], data[1], val}
		payload, _ = json.Marshal(transferEvent)
	} else if topic == "approve" {
		val, _ := strconv.Atoi(data[2])
		approvalEvent := approval{data[0], data[1], val}
		payload, _ = json.Marshal(approvalEvent)
	} else {
		payload, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}
	return payload, err
}

// Mint creates new tokens and adds them to minter's account balance
// This function triggers a Transfer event
func (s *SmartContract) Mint(ctx contractapi.TransactionContextInterface, recipient string, amount int) error {
	s.erc20Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	//TODO check recipient is valid
	account := fabric.NewMspUser(recipient)
	if amount <= 0 {
		return fmt.Errorf("mint amount must be a positive integer")
	}
	amount256 := common.NewSafeUint256(uint64(amount))
	success, err := s.erc20Contract.Mint(account, amount256)
	if err != nil {
		return err
	}
	if success {
		return nil
	}
	return errors.New("mint fail")
}

// Burn redeems tokens the minter's account balance
// This function triggers a Transfer event
func (s *SmartContract) Burn(ctx contractapi.TransactionContextInterface, amount int) error {
	s.erc20Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	if amount <= 0 {
		return fmt.Errorf("mint amount must be a positive integer")
	}
	amount256 := common.NewSafeUint256(uint64(amount))
	success, err := s.erc20Contract.Burn(amount256)
	if err != nil {
		return err
	}
	if success {
		return nil
	}
	return errors.New("burn fail")
}

// Transfer transfers tokens from client account to recipient account
// recipient account must be a valid clientID as returned by the ClientID() function
// This function triggers a Transfer event
func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, recipient string, amount int) error {
	s.erc20Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	//TODO check recipient is valid
	account := fabric.NewMspUser(recipient)
	if amount <= 0 {
		return fmt.Errorf("mint amount must be a positive integer")
	}
	amount256 := common.NewSafeUint256(uint64(amount))
	success, err := s.erc20Contract.Transfer(account, amount256)
	if err != nil {
		return err
	}
	if success {
		return nil
	}
	return errors.New("transfer fail")
	//// Get ID of submitting client identity
	//clientID, err := ctx.GetClientIdentity().GetID()
	//if err != nil {
	//	return fmt.Errorf("failed to get client id: %v", err)
	//}
	//
	//err = transferHelper(ctx, clientID, recipient, amount)
	//if err != nil {
	//	return fmt.Errorf("failed to transfer: %v", err)
	//}
	//
	//// Emit the Transfer event
	//transferEvent := event{clientID, recipient, amount}
	//transferEventJSON, err := json.Marshal(transferEvent)
	//if err != nil {
	//	return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	//}
	//err = ctx.GetStub().SetEvent("Transfer", transferEventJSON)
	//if err != nil {
	//	return fmt.Errorf("failed to set event: %v", err)
	//}
	//
	//return nil
}

// BalanceOf returns the balance of the given account
func (s *SmartContract) BalanceOf(ctx contractapi.TransactionContextInterface, account string) (int, error) {
	s.erc20Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	acc := fabric.NewMspUser(account)
	bal, err := s.erc20Contract.BalanceOf(acc)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(bal.ToString())
}

// ClientAccountBalance returns the balance of the requesting client's account
func (s *SmartContract) ClientAccountBalance(ctx contractapi.TransactionContextInterface) (int, error) {

	// Get ID of submitting client identity
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return 0, fmt.Errorf("failed to get client id: %v", err)
	}
	return s.BalanceOf(ctx, clientID)
}

// ClientAccountID returns the id of the requesting client's account
// In this implementation, the client account ID is the clientId itself
// Users can use this function to get their own account id, which they can then give to others as the payment address
func (s *SmartContract) ClientAccountID(ctx contractapi.TransactionContextInterface) (string, error) {

	// Get ID of submitting client identity
	clientAccountID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to get client id: %v", err)
	}

	return clientAccountID, nil
}

// TotalSupply returns the total token supply
func (s *SmartContract) TotalSupply(ctx contractapi.TransactionContextInterface) (int, error) {
	s.erc20Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	num, err := s.erc20Contract.TotalSupply()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(num.ToString())
}

// Approve allows the spender to withdraw from the calling client's token account
// The spender can withdraw multiple times if necessary, up to the value amount
// This function triggers an Approval event
func (s *SmartContract) Approve(ctx contractapi.TransactionContextInterface, spender string, value int) error {
	s.erc20Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	//TODO check recipient is valid
	account := fabric.NewMspUser(spender)
	if value <= 0 {
		return fmt.Errorf("mint amount must be a positive integer")
	}
	amount256 := common.NewSafeUint256(uint64(value))
	success, err := s.erc20Contract.Approve(account, amount256)
	if err != nil {
		return err
	}
	if success {
		return nil
	}
	return errors.New("approve fail")
}

// Allowance returns the amount still available for the spender to withdraw from the owner
func (s *SmartContract) Allowance(ctx contractapi.TransactionContextInterface, owner string, spender string) (int, error) {
	s.erc20Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	//TODO check recipient is valid
	spenderAcc := fabric.NewMspUser(spender)
	ownerAcc := fabric.NewMspUser(owner)
	num, err := s.erc20Contract.Allowance(ownerAcc, spenderAcc)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(num.ToString())
}

// TransferFrom transfers the value amount from the "from" address to the "to" address
// This function triggers a Transfer event
func (s *SmartContract) TransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, value int) error {
	s.erc20Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	//TODO check recipient is valid
	fromAcc := fabric.NewMspUser(from)
	toAcc := fabric.NewMspUser(to)
	if value <= 0 {
		return fmt.Errorf("mint amount must be a positive integer")
	}
	amount256 := common.NewSafeUint256(uint64(value))
	success, err := s.erc20Contract.TransferFrom(fromAcc, toAcc, amount256)
	if err != nil {
		return err
	}
	if success {
		return nil
	}
	return errors.New("transferFrom fail")
}
