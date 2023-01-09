package erc721

import (
	"errors"

	"github.com/studyzy/openzeppelin-go/common"
)

// Option 初始化ERC20合约的选项
type Option struct {
	// IsZeroAccount 判断一个账户是否为空
	IsZeroAccount func(account common.Account) bool
	// GetZeroAccount 获得一个空账户
	GetZeroAccount func() common.Account
	// IsValidAccount 判断一个账号是否有效
	IsValidAccount func(account common.Account) bool
	// BeforeTransfer 在转账前执行的逻辑
	BeforeTransfer func(from common.Account, to common.Account, tokenId *common.SafeUint256) error
	// AfterTransfer 在转账成功后执行的逻辑
	AfterTransfer func(from common.Account, to common.Account, tokenId *common.SafeUint256) error
	// Burnable 是否允许销毁
	Burnable bool
	// Minable 是否允许后续铸造
	Minable bool
}

/**
 * @dev Transfers `tokenId` from `from` to `to`.
 *  As opposed to {transferFrom}, this imposes no restrictions on msg.sender.
 *
 * Requirements:
 *
 * - `to` cannot be the zero address.
 * - `tokenId` token must be owned by `from`.
 *
 * Emits a {Transfer} event.
 */
func (c *ERC721Contract) baseTransfer(from, to common.Account, tokenId *common.SafeUint256) error {
	tokenOwner, err := c.dal.GetTokenOwner(tokenId)
	if err != nil {
		return err
	}
	err = common.Require(tokenOwner.Equal(from), "ERC721: transfer from incorrect owner")
	if err != nil {
		return err
	}
	if err = common.Require(to.IsZero(), "ERC721: transfer to the zero address"); err != nil {
		return err
	}

	if c.option.BeforeTransfer != nil {
		err = c.option.BeforeTransfer(from, to, tokenId)
	}
	// Clear approvals from the previous owner
	err = c.dal.DeleteTokenApproval(tokenId)
	fromBal, err := c.dal.GetBalance(from)
	newFromBal, _ := common.SafeSub(fromBal, common.NewSafeUint256(1))
	err = c.dal.SetBalance(from, newFromBal)
	toBal, err := c.dal.GetBalance(to)
	newToBal, _ := common.SafeAdd(toBal, common.NewSafeUint256(1))
	err = c.dal.SetBalance(to, newToBal)
	err = c.dal.SetTokenOwner(tokenId, to)
	c.sdk.EmitEvent("transfer", from.ToString(), to.ToString(), tokenId.ToString())

	if c.option.AfterTransfer != nil {
		err = c.option.AfterTransfer(from, to, tokenId)
	}
	return nil
}

/**
 * @dev Mints `tokenId` and transfers it to `to`.
 *
 * WARNING: Usage of this method is discouraged, use {_safeMint} whenever possible
 *
 * Requirements:
 *
 * - `tokenId` must not exist.
 * - `to` cannot be the zero address.
 *
 * Emits a {Transfer} event.
 */
func (c *ERC721Contract) baseMint(to common.Account, tokenId *common.SafeUint256) error {
	if err := common.Require(!to.IsZero(), "ERC721: mint to the zero address"); err != nil {
		return err
	}
	if err := common.Require(!c.exists(tokenId), "ERC721: token already minted"); err != nil {
		return err
	}
	if c.option.BeforeTransfer != nil {
		c.option.BeforeTransfer(c.option.GetZeroAccount(), to, tokenId)
	}
	c.dal.IncreaseBalance(to)
	//_balances[to] += 1;
	c.dal.SetTokenOwner(tokenId, to)
	//_owners[tokenId] = to;
	c.sdk.EmitEvent("transfer", c.option.GetZeroAccount().ToString(), to.ToString(), tokenId.ToString())
	//emit Transfer(address(0), to, tokenId);
	if c.option.AfterTransfer != nil {
		c.option.AfterTransfer(c.option.GetZeroAccount(), to, tokenId)
		//_afterTokenTransfer(address(0), to, tokenId);
	}
	return nil
}

/**
 * @dev Destroys `tokenId`.
 * The approval is cleared when the token is burned.
 *
 * Requirements:
 *
 * - `tokenId` must exist.
 *
 * Emits a {Transfer} event.
 */
func (c *ERC721Contract) baseBurn(tokenId *common.SafeUint256) error {
	//address owner = ERC721.ownerOf(tokenId);
	owner, err := c.dal.GetTokenOwner(tokenId)
	if err != nil {
		return err
	}
	//_beforeTokenTransfer(owner, address(0), tokenId);
	if c.option.BeforeTransfer != nil {
		if err = c.option.BeforeTransfer(owner, c.option.GetZeroAccount(), tokenId); err != nil {
			return err
		}
	}
	// Clear approvals
	//_approve(address(0), tokenId)
	err = c.baseApprove(c.option.GetZeroAccount(), tokenId)
	if err != nil {
		return err
	}
	err = c.dal.DecreaseBalance(owner)
	if err != nil {
		return err
	}
	//_balances[owner] -= 1;
	//delete _owners[tokenId];
	err = c.dal.DeleteTokenOwner(tokenId)
	if err != nil {
		return err
	}
	//emit Transfer(owner, address(0), tokenId);
	err = c.sdk.EmitEvent("transfer", owner.ToString(), c.option.GetZeroAccount().ToString(), tokenId.ToString())
	if err != nil {
		return err
	}
	if c.option.AfterTransfer != nil {
		err = c.option.AfterTransfer(owner, c.option.GetZeroAccount(), tokenId)
		if err != nil {
			return err
		}
		//_afterTokenTransfer(address(0), to, tokenId);
	}
	return nil
}

/**
 * @dev Approve `to` to operate on `tokenId`
 *
 * Emits an {Approval} event.
 */
func (c *ERC721Contract) baseApprove(to common.Account, tokenId *common.SafeUint256) error {
	//_tokenApprovals[tokenId] = to;
	err := c.dal.SetTokenApproval(tokenId, to)
	if err != nil {
		return err
	}
	//emit Approval(ERC721.ownerOf(tokenId), to, tokenId);
	owner, err := c.dal.GetTokenOwner(tokenId)
	if err != nil {
		return err
	}
	return c.sdk.EmitEvent("approval", owner.ToString(), to.ToString(), tokenId.ToString())
}

/**
 * @dev Approve `operator` to operate on all of `owner` tokens
 *
 * Emits an {ApprovalForAll} event.
 */
func (c *ERC721Contract) baseetApprovalForAll(owner, operator common.Account, approved bool) error {
	//require(owner != operator, "ERC721: approve to caller");
	err := common.Require(owner != operator, "ERC721: approve to caller")
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

/**
 * @dev Returns whether `tokenId` exists.
 *
 * Tokens can be managed by their owner or approved accounts via {approve} or {setApprovalForAll}.
 *
 * Tokens start existing when they are minted (`_mint`),
 * and stop existing when they are burned (`_burn`).
 */
func (c *ERC721Contract) exists(tokenId *common.SafeUint256) bool {
	owner, err := c.dal.GetTokenOwner(tokenId)
	if err != nil || owner == nil {
		return false
	}
	return !owner.IsZero()
}

func (c *ERC721Contract) baseCheckOnERC721Received(from, to common.Account, tokenId *common.SafeUint256, data []byte) (bool, error) {
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

/**
 * @dev Safely transfers `tokenId` token from `from` to `to`, checking first that contract recipients
 * are aware of the ERC721 protocol to prevent tokens from being forever locked.
 *
 * `data` is additional data, it has no specified format and it is sent in call to `to`.
 *
 * This internal function is equivalent to {safeTransferFrom}, and can be used to e.g.
 * implement alternative mechanisms to perform token transfer, such as signature-based.
 *
 * Requirements:
 *
 * - `from` cannot be the zero address.
 * - `to` cannot be the zero address.
 * - `tokenId` token must exist and be owned by `from`.
 * - If `to` refers to a smart contract, it must implement {IERC721Receiver-onERC721Received}, which is called upon a safe transfer.
 *
 * Emits a {Transfer} event.
 */
func (c *ERC721Contract) baseSafeTransfer(from, to common.Account, tokenId *common.SafeUint256, data []byte) error {
	err := c.baseTransfer(from, to, tokenId)
	if err != nil {
		return err
	}
	result, err := c.baseCheckOnERC721Received(from, to, tokenId, data)
	if err != nil {
		return err
	}
	return common.Require(result, "ERC721: transfer to non ERC721Receiver implementer")
}

/**
 * @dev Returns whether `spender` is allowed to manage `tokenId`.
 *
 * Requirements:
 *
 * - `tokenId` must exist.
 */
func (c *ERC721Contract) baseIsApprovedOrOwner(spender common.Account, tokenId *common.SafeUint256) (bool, error) {
	owner, err := c.dal.GetTokenOwner(tokenId)
	if err != nil {
		return false, err
	}
	approvedForAll, err := c.dal.GetOperatorApproval(owner, spender)
	if err != nil {
		return false, err
	}
	approveSpender, err := c.GetApproved(tokenId)
	if err != nil {
		return false, err
	}
	return spender.Equal(owner) || approvedForAll || approveSpender.Equal(spender), nil
}
