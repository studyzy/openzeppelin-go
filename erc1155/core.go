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
	"encoding/json"
	"errors"

	"github.com/studyzy/openzeppelin-go/common"
)

// Option 初始化ERC20合约的选项
type Option struct {
	// BeforeTransfer 在转账前执行的逻辑
	BeforeTransfer func(operator, from, to common.Account, ids, amounts []*common.SafeUint256, data []byte) error
	// AfterTransfer 在转账成功后执行的逻辑
	AfterTransfer func(operator, from, to common.Account, ids, amounts []*common.SafeUint256, data []byte) error
}

func asSingletonArray(element *common.SafeUint256) []*common.SafeUint256 {
	array := make([]*common.SafeUint256, 1)
	array[0] = element

	return array
}

/**
 * @dev Transfers `amount` tokens of token type `id` from `from` to `to`.
 *
 * Emits a {TransferSingle} event.
 *
 * Requirements:
 *
 * - `to` cannot be the zero address.
 * - `from` must have a balance of tokens of type `id` of at least `amount`.
 * - If `to` refers to a smart contract, it must implement {IERC1155Receiver-onERC1155Received} and return the
 * acceptance magic value.
 */
func (c *ERC1155Contract) baseSafeTransferFrom(from, to common.Account, id, amount *common.SafeUint256, data []byte) error {
	err := common.Require(!to.IsZero(), "ERC1155: transfer to the zero address")
	if err != nil {
		return err
	}
	sender, err := c.sdk.GetTxSender()
	if err != nil {
		return err
	}
	ids := asSingletonArray(id)
	amounts := asSingletonArray(amount)

	if c.option.BeforeTransfer != nil {
		err = c.option.BeforeTransfer(sender, from, to, ids, amounts, data)
		if err != nil {
			return err
		}
	}
	//update from balance
	fromBalance, err := c.dal.GetBalance(id, from)
	err = common.Require(fromBalance.GTE(amount), "ERC1155: insufficient balance for transfer")
	if err != nil {
		return err
	}
	fromBalance, _ = common.SafeSub(fromBalance, amount)
	err = c.dal.SetBalance(id, from, fromBalance)
	if err != nil {
		return err
	}
	//update to balance
	toBalance, err := c.dal.GetBalance(id, to)
	if err != nil {
		return err
	}
	toBalance, _ = common.SafeAdd(toBalance, amount)
	err = c.dal.SetBalance(id, to, toBalance)
	if err != nil {
		return err
	}
	err = c.sdk.EmitEvent("transferSingle", sender.ToString(), from.ToString(), to.ToString(), id.ToString(), amount.ToString())
	if err != nil {
		return err
	}
	if c.option.AfterTransfer != nil {
		err = c.option.AfterTransfer(sender, from, to, ids, amounts, data)
		if err != nil {
			return err
		}
	}
	return c.doSafeTransferAcceptanceCheck(sender, from, to, id, amount, data)
}

/**
 * @dev xref:ROOT:erc1155.adoc#batch-operations[Batched] version of {_safeTransferFrom}.
 *
 * Emits a {TransferBatch} event.
 *
 * Requirements:
 *
 * - If `to` refers to a smart contract, it must implement {IERC1155Receiver-onERC1155BatchReceived} and return the
 * acceptance magic value.
 */
func (c *ERC1155Contract) baseSafeBatchTransferFrom(from, to common.Account, ids, amounts []*common.SafeUint256, data []byte) error {
	err := common.Require(len(ids) == len(amounts), "ERC1155: ids and amounts length mismatch")
	if err != nil {
		return err
	}

	err = common.Require(!to.IsZero(), "ERC1155: transfer to the zero address")
	if err != nil {
		return err
	}
	sender, err := c.sdk.GetTxSender()
	if err != nil {
		return err
	}
	if c.option.BeforeTransfer != nil {
		err = c.option.BeforeTransfer(sender, from, to, ids, amounts, data)
		if err != nil {
			return err
		}
	}
	for i := 0; i < len(ids); i++ {
		id := ids[i]
		amount := amounts[i]
		//update from balance
		fromBalance, err := c.dal.GetBalance(id, from)
		if err != nil {
			return err
		}
		err = common.Require(fromBalance.GTE(amount), "ERC1155: insufficient balance for transfer")
		if err != nil {
			return err
		}
		fromBalance, _ = common.SafeSub(fromBalance, amount)
		err = c.dal.SetBalance(id, from, fromBalance)
		if err != nil {
			return err
		}
		//update to balance
		toBalance, err := c.dal.GetBalance(id, to)
		toBalance, _ = common.SafeAdd(toBalance, amount)
		err = c.dal.SetBalance(id, to, toBalance)
		if err != nil {
			return err
		}
	}
	err = c.sdk.EmitEvent("transferBatch", sender.ToString(), from.ToString(), to.ToString(), uint256sToString(ids), uint256sToString(amounts))
	if err != nil {
		return err
	}
	if c.option.AfterTransfer != nil {
		err = c.option.AfterTransfer(sender, from, to, ids, amounts, data)
		if err != nil {
			return err
		}
	}
	return c.doSafeBatchTransferAcceptanceCheck(sender, from, to, ids, amounts, data)
}

/**
 * @dev Creates `amount` tokens of token type `id`, and assigns them to `to`.
 *
 * Emits a {TransferSingle} event.
 *
 * Requirements:
 *
 * - `to` cannot be the zero address.
 * - If `to` refers to a smart contract, it must implement {IERC1155Receiver-onERC1155Received} and return the
 * acceptance magic value.
 */
func (c *ERC1155Contract) baseMint(to common.Account, id, amount *common.SafeUint256, data []byte) error {
	if err := common.Require(!to.IsZero(), "ERC1155: mint to the zero address"); err != nil {
		return err
	}
	operator, err := c.sdk.GetTxSender()
	if err != nil {
		return err
	}
	ids := asSingletonArray(id)
	amounts := asSingletonArray(amount)
	from := c.sdk.NewZeroAccount()
	if c.option.BeforeTransfer != nil {
		err = c.option.BeforeTransfer(operator, from, to, ids, amounts, data)
		if err != nil {
			return err
		}
	}
	//update to balance
	toBalance, err := c.dal.GetBalance(id, to)
	if err != nil {
		return err
	}
	toBalance, _ = common.SafeAdd(toBalance, amount)
	err = c.dal.SetBalance(id, to, toBalance)
	if err != nil {
		return err
	}
	err = c.sdk.EmitEvent("transferSingle", operator.ToString(), from.ToString(), to.ToString(), id.ToString(), amount.ToString())
	if err != nil {
		return err
	}
	if c.option.AfterTransfer != nil {
		err = c.option.AfterTransfer(operator, from, to, ids, amounts, data)
		if err != nil {
			return err
		}
	}
	return c.doSafeTransferAcceptanceCheck(operator, from, to, id, amount, data)
}

/**
 * @dev xref:ROOT:erc1155.adoc#batch-operations[Batched] version of {_mint}.
 *
 * Emits a {TransferBatch} event.
 *
 * Requirements:
 *
 * - `ids` and `amounts` must have the same length.
 * - If `to` refers to a smart contract, it must implement {IERC1155Receiver-onERC1155BatchReceived} and return the
 * acceptance magic value.
 */
func (c *ERC1155Contract) baseMintBatch(to common.Account, ids, amounts []*common.SafeUint256, data []byte) error {
	if err := common.Require(!to.IsZero(), "ERC1155: mint to the zero address"); err != nil {
		return err
	}
	err := common.Require(len(ids) == len(amounts), "ERC1155: ids and amounts length mismatch")
	if err != nil {
		return err
	}
	operator, err := c.sdk.GetTxSender()
	if err != nil {
		return err
	}
	from := c.sdk.NewZeroAccount()
	if c.option.BeforeTransfer != nil {
		err = c.option.BeforeTransfer(operator, from, to, ids, amounts, data)
		if err != nil {
			return err
		}
	}
	for i := 0; i < len(ids); i++ {
		id := ids[i]
		amount := amounts[i]
		//update to balance
		toBalance, err := c.dal.GetBalance(id, to)
		toBalance, _ = common.SafeAdd(toBalance, amount)
		err = c.dal.SetBalance(id, to, toBalance)
		if err != nil {
			return err
		}
	}

	err = c.sdk.EmitEvent("transferBatch", operator.ToString(), from.ToString(), to.ToString(), uint256sToString(ids), uint256sToString(amounts))
	if err != nil {
		return err
	}
	if c.option.AfterTransfer != nil {
		err = c.option.AfterTransfer(operator, from, to, ids, amounts, data)
		if err != nil {
			return err
		}
	}
	return c.doSafeBatchTransferAcceptanceCheck(operator, from, to, ids, amounts, data)
}

/**
 * @dev Destroys `amount` tokens of token type `id` from `from`
 *
 * Emits a {TransferSingle} event.
 *
 * Requirements:
 *
 * - `from` cannot be the zero address.
 * - `from` must have at least `amount` tokens of token type `id`.
 */
func (c *ERC1155Contract) baseBurn(from common.Account, id, amount *common.SafeUint256) error {
	if err := common.Require(!from.IsZero(), "ERC1155: burn from the zero address"); err != nil {
		return err
	}
	operator, err := c.sdk.GetTxSender()
	if err != nil {
		return err
	}
	ids := asSingletonArray(id)
	amounts := asSingletonArray(amount)
	to := c.sdk.NewZeroAccount()
	if c.option.BeforeTransfer != nil {
		err = c.option.BeforeTransfer(operator, from, to, ids, amounts, nil)
		if err != nil {
			return err
		}
	}
	//update from balance
	fromBalance, err := c.dal.GetBalance(id, from)
	if err != nil {
		return err
	}
	err = common.Require(fromBalance.GTE(amount), "ERC1155: burn amount exceeds balance")
	if err != nil {
		return err
	}
	fromBalance, _ = common.SafeSub(fromBalance, amount)
	err = c.dal.SetBalance(id, from, fromBalance)
	if err != nil {
		return err
	}
	err = c.sdk.EmitEvent("transferSingle", operator.ToString(), from.ToString(), to.ToString(), id.ToString(), amount.ToString())
	if err != nil {
		return err
	}
	if c.option.AfterTransfer != nil {
		err = c.option.AfterTransfer(operator, from, to, ids, amounts, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
 * @dev xref:ROOT:erc1155.adoc#batch-operations[Batched] version of {_burn}.
 *
 * Emits a {TransferBatch} event.
 *
 * Requirements:
 *
 * - `ids` and `amounts` must have the same length.
 */
func (c *ERC1155Contract) baseBurnBatch(from common.Account, ids, amounts []*common.SafeUint256) error {
	if err := common.Require(!from.IsZero(), "ERC1155: burn from the zero address"); err != nil {
		return err
	}
	err := common.Require(len(ids) == len(amounts), "ERC1155: ids and amounts length mismatch")
	if err != nil {
		return err
	}
	operator, err := c.sdk.GetTxSender()
	if err != nil {
		return err
	}
	to := c.sdk.NewZeroAccount()
	if c.option.BeforeTransfer != nil {
		err = c.option.BeforeTransfer(operator, from, to, ids, amounts, nil)
		if err != nil {
			return err
		}
	}
	for i := 0; i < len(ids); i++ {
		id := ids[i]
		amount := amounts[i]
		//update from balance
		fromBalance, err := c.dal.GetBalance(id, from)
		if err != nil {
			return err
		}
		err = common.Require(fromBalance.GTE(amount), "ERC1155: burn amount exceeds balance")
		if err != nil {
			return err
		}
		fromBalance, _ = common.SafeSub(fromBalance, amount)
		err = c.dal.SetBalance(id, from, fromBalance)
		if err != nil {
			return err
		}
	}

	err = c.sdk.EmitEvent("transferBatch", operator.ToString(), from.ToString(), to.ToString(), uint256sToString(ids), uint256sToString(amounts))
	if err != nil {
		return err
	}
	if c.option.AfterTransfer != nil {
		err = c.option.AfterTransfer(operator, from, to, ids, amounts, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
 * @dev Approve `operator` to operate on all of `owner` tokens
 *
 * Emits an {ApprovalForAll} event.
 */
func (c *ERC1155Contract) baseSetApprovalForAll(owner, operator common.Account, approved bool) error {
	err := common.Require(owner != operator, "ERC1155: setting approval status for self")
	if err != nil {
		return err
	}
	err = c.dal.SetOperatorApproval(owner, operator, approved)
	if err != nil {
		return err
	}
	//_operatorApprovals[owner][operator] = approved;
	a := "false"
	if approved {
		a = "true"
	}
	return c.sdk.EmitEvent("approvalForAll", owner.ToString(), operator.ToString(), a)
	//emit ApprovalForAll(owner, operator, approved);
}

func (c *ERC1155Contract) baseCheckOnERC721Received(from, to common.Account, tokenId *common.SafeUint256, data []byte) (bool, error) {
	if c.sdk.IsContract(to) {
		sender, err := c.sdk.GetTxSender()
		if err != nil {
			return false, err
		}
		args := make([]common.KeyValue, 4)
		args[0] = common.KeyValue{
			Key:   "operator",
			Value: sender.Bytes(),
		}
		args[1] = common.KeyValue{
			Key:   "from",
			Value: from.Bytes(),
		}
		args[2] = common.KeyValue{
			Key:   "tokenId",
			Value: []byte(tokenId.ToString()),
		}
		args[3] = common.KeyValue{
			Key:   "data",
			Value: data,
		}
		response := c.sdk.CallContract(to, "onERC721Received", args)
		if response.Status == common.OK {
			return true, nil
		}
		if len(response.Message) == 0 {
			return false, errors.New("ERC721: transfer to non ERC721Receiver implementer")
		}
		return false, errors.New(response.Message)
	} else {
		return true, nil
	}
}

func (c *ERC1155Contract) doSafeTransferAcceptanceCheck(operator, from, to common.Account,
	id, amount *common.SafeUint256, data []byte) error {
	if c.sdk.IsContract(to) {
		args := make([]common.KeyValue, 4)
		args[0] = common.KeyValue{
			Key:   "operator",
			Value: operator.Bytes(),
		}
		args[1] = common.KeyValue{
			Key:   "from",
			Value: from.Bytes(),
		}
		args[2] = common.KeyValue{
			Key:   "to",
			Value: to.Bytes(),
		}
		args[3] = common.KeyValue{
			Key:   "id",
			Value: []byte(id.ToString()),
		}
		args[4] = common.KeyValue{
			Key:   "amount",
			Value: []byte(amount.ToString()),
		}
		args[5] = common.KeyValue{
			Key:   "data",
			Value: data,
		}
		response := c.sdk.CallContract(to, "onERC1155Received", args)
		if response.Status == common.OK {
			return nil
		}
		if len(response.Message) == 0 {
			return errors.New("ERC1155: transfer to non-ERC1155Receiver implementer")
		}
		return errors.New(response.Message)
	} else {
		return nil
	}

}

func (c *ERC1155Contract) doSafeBatchTransferAcceptanceCheck(sender common.Account, from common.Account, to common.Account, ids []*common.SafeUint256, amounts []*common.SafeUint256, data []byte) error {
	if c.sdk.IsContract(to) {
		args := make([]common.KeyValue, 4)
		args[0] = common.KeyValue{
			Key:   "operator",
			Value: sender.Bytes(),
		}
		args[1] = common.KeyValue{
			Key:   "from",
			Value: from.Bytes(),
		}
		args[2] = common.KeyValue{
			Key:   "to",
			Value: to.Bytes(),
		}
		args[3] = common.KeyValue{
			Key:   "ids",
			Value: []byte(uint256sToString(ids)),
		}
		args[4] = common.KeyValue{
			Key:   "amounts",
			Value: []byte(uint256sToString(amounts)),
		}
		args[5] = common.KeyValue{
			Key:   "data",
			Value: data,
		}
		response := c.sdk.CallContract(to, "onERC1155BatchReceived", args)
		if response.Status == common.OK {
			return nil
		}
		if len(response.Message) == 0 {
			return errors.New("ERC1155: transfer to non-ERC1155Receiver implementer")
		}
		return errors.New(response.Message)
	} else {
		return nil
	}
}
func uint256sToString(ids []*common.SafeUint256) string {
	data, _ := json.Marshal(ids)
	return string(data)
}
