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
	"github.com/studyzy/openzeppelin-go/erc721"
	"github.com/studyzy/openzeppelin-go/fabric"
)

// SmartContract provides functions for transferring tokens between accounts
type SmartContract struct {
	contractapi.Contract
	erc721Contract erc721.ERC721Contract
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
type approvalForAll struct {
	owner    string
	spender  string
	approved bool
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
	} else if topic == "approveForAll" {
		approved := data[2] == "true"
		approvalEvent := approvalForAll{data[0], data[1], approved}
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
func (s *SmartContract) Mint(ctx contractapi.TransactionContextInterface, recipient string, tokenId int) error {
	s.erc721Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	//TODO check recipient is valid
	account := fabric.NewMspUser(recipient)
	if tokenId <= 0 {
		return fmt.Errorf("mint tokenId must be a positive integer")
	}
	tokenId256 := common.NewSafeUint256(uint64(tokenId))
	success, err := s.erc721Contract.Mint(account, tokenId256)
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
func (s *SmartContract) Burn(ctx contractapi.TransactionContextInterface, tokenId int) error {
	s.erc721Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	if tokenId <= 0 {
		return fmt.Errorf("mint tokenId must be a positive integer")
	}
	tokenId256 := common.NewSafeUint256(uint64(tokenId))
	return s.erc721Contract.Burn(tokenId256)

}

// BalanceOf returns the balance of the given account
func (s *SmartContract) BalanceOf(ctx contractapi.TransactionContextInterface, account string) (int, error) {
	s.erc721Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	acc := fabric.NewMspUser(account)
	bal, err := s.erc721Contract.BalanceOf(acc)
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

func (s *SmartContract) OwnerOf(ctx contractapi.TransactionContextInterface, tokenId int) (int, error) {
	s.erc721Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	tokenIdNum := common.NewSafeUint256(uint64(tokenId))

	bal, err := s.erc721Contract.OwnerOf(tokenIdNum)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(bal.ToString())
}

// SafeTransferFrom transfers the value amount from the "from" address to the "to" address
// This function triggers a Transfer event
func (s *SmartContract) SafeTransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, tokenId int) error {
	s.erc721Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	//TODO check recipient is valid
	fromAcc := fabric.NewMspUser(from)
	toAcc := fabric.NewMspUser(to)
	if tokenId <= 0 {
		return fmt.Errorf("mint amount must be a positive integer")
	}
	tokenNum := common.NewSafeUint256(uint64(tokenId))
	return s.erc721Contract.SafeTransferFrom(fromAcc, toAcc, tokenNum)
}

// SafeTransferFrom2 transfers the value amount from the "from" address to the "to" address
// This function triggers a Transfer event
func (s *SmartContract) SafeTransferFrom2(ctx contractapi.TransactionContextInterface, from string, to string, tokenId int, data []byte) error {
	s.erc721Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	//TODO check recipient is valid
	fromAcc := fabric.NewMspUser(from)
	toAcc := fabric.NewMspUser(to)
	if tokenId <= 0 {
		return fmt.Errorf("mint amount must be a positive integer")
	}
	tokenNum := common.NewSafeUint256(uint64(tokenId))
	return s.erc721Contract.SafeTransferFrom2(fromAcc, toAcc, tokenNum, data)
}

// TransferFrom transfers the value amount from the "from" address to the "to" address
// This function triggers a Transfer event
func (s *SmartContract) TransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, tokenId int) error {
	s.erc721Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	//TODO check recipient is valid
	fromAcc := fabric.NewMspUser(from)
	toAcc := fabric.NewMspUser(to)
	if tokenId <= 0 {
		return fmt.Errorf("mint amount must be a positive integer")
	}
	tokenNum := common.NewSafeUint256(uint64(tokenId))
	return s.erc721Contract.TransferFrom(fromAcc, toAcc, tokenNum)
}

// Approve allows the spender to withdraw from the calling client's token account
// The spender can withdraw multiple times if necessary, up to the value amount
// This function triggers an Approval event
func (s *SmartContract) Approve(ctx contractapi.TransactionContextInterface, spender string, tokenId int) error {
	s.erc721Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	//TODO check recipient is valid
	account := fabric.NewMspUser(spender)
	if tokenId <= 0 {
		return fmt.Errorf("mint amount must be a positive integer")
	}
	tokenId256 := common.NewSafeUint256(uint64(tokenId))
	return s.erc721Contract.Approve(account, tokenId256)
}

func (s *SmartContract) SetApprovalForAll(ctx contractapi.TransactionContextInterface, spender string, approved bool) error {
	s.erc721Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	//TODO check recipient is valid
	account := fabric.NewMspUser(spender)
	return s.erc721Contract.SetApprovalForAll(account, approved)
}

func (s *SmartContract) GetApproved(ctx contractapi.TransactionContextInterface, tokenId int) (string, error) {
	s.erc721Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	tokenId256 := common.NewSafeUint256(uint64(tokenId))
	operator, err := s.erc721Contract.GetApproved(tokenId256)
	if err != nil {
		return "", err
	}
	return operator.ToString(), nil
}

func (s *SmartContract) IsApprovedForAll(ctx contractapi.TransactionContextInterface, owner string, spender string) (bool, error) {
	s.erc721Contract.SetSDK(fabric.NewSDkAdapter(ctx, encodeEvent, func(contractName string) (bool, error) {
		//TODO
		return false, nil
	}))
	//TODO check recipient is valid
	spenderAcc := fabric.NewMspUser(spender)
	ownerAcc := fabric.NewMspUser(owner)
	return s.erc721Contract.IsApprovedForAll(ownerAcc, spenderAcc)

}
