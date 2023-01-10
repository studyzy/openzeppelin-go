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
	"fmt"

	"github.com/studyzy/openzeppelin-go/common"
)

var _ IERC1155 = (*ERC1155Contract)(nil)

type ERC1155Contract struct {
	option Option
	dal    *ERC1155Dal
	sdk    common.ContractSDK
}

func (c *ERC1155Contract) SupportsInterface(interfaceId string) bool {
	return interfaceId == "ERC1155" || interfaceId == "ERC1155Metadata" || interfaceId == "ERC165"
}

func (c *ERC1155Contract) BalanceOf(account common.Account, id *common.SafeUint256) (*common.SafeUint256, error) {
	err := common.Require(!account.IsZero(), "ERC1155: address zero is not a valid owner")
	if err != nil {
		return nil, err
	}
	return c.dal.GetBalance(id, account)
}

func (c *ERC1155Contract) BalanceOfBatch(accounts []common.Account, ids []*common.SafeUint256) ([]*common.SafeUint256, error) {
	err := common.Require(len(accounts) == len(ids), "ERC1155: accounts and ids length mismatch")
	if err != nil {
		return nil, err
	}
	batchBalances := make([]*common.SafeUint256, len(accounts))

	for i := 0; i < len(accounts); i++ {
		batchBalances[i], err = c.dal.GetBalance(ids[i], accounts[i])
		if err != nil {
			return nil, err
		}
	}

	return batchBalances, nil
}

func (c *ERC1155Contract) SetApprovalForAll(operator common.Account, approved bool) error {
	sender, err := c.sdk.GetTxSender()
	if err != nil {
		return err
	}
	return c.dal.SetOperatorApproval(sender, operator, approved)
}

func (c *ERC1155Contract) IsApprovedForAll(account common.Account, operator common.Account) (bool, error) {
	return c.dal.GetOperatorApproval(account, operator)
}

func (c *ERC1155Contract) SafeTransferFrom(from, to common.Account, id, amount *common.SafeUint256, data []byte) error {
	sender, err := c.sdk.GetTxSender()
	if err != nil {
		return err
	}
	isApproved, err := c.IsApprovedForAll(from, sender)
	if err != nil {
		return err
	}
	err = common.Require(from.Equal(sender) || isApproved, "ERC1155: caller is not token owner or approved")
	if err != nil {
		return err
	}
	return c.baseSafeTransferFrom(from, to, id, amount, data)
}

func (c *ERC1155Contract) SafeBatchTransferFrom(from, to common.Account, ids, amounts []*common.SafeUint256, data []byte) error {
	sender, err := c.sdk.GetTxSender()
	if err != nil {
		return err
	}
	isApproved, err := c.IsApprovedForAll(from, sender)
	if err != nil {
		return err
	}
	err = common.Require(from.Equal(sender) || isApproved, "ERC1155: caller is not token owner or approved")
	if err != nil {
		return err
	}
	return c.baseSafeBatchTransferFrom(from, to, ids, amounts, data)
}

func (c *ERC1155Contract) Uri(id *common.SafeUint256) (string, error) {
	uri, err := c.dal.GetUri()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(uri, id.ToString()), nil
}
