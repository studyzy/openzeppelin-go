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
	"errors"
	"fmt"

	"github.com/studyzy/openzeppelin-go/common"
)

var _ IERC721 = (*ERC721Contract)(nil)

type ERC721Contract struct {
	option  Option
	_name   string
	_symbol string
	dal     *ERC721DAL
	sdk     common.ContractSDK
}

func (c *ERC721Contract) BalanceOf(owner common.Account) (*common.SafeUint256, error) {
	err := common.Require(!owner.IsZero(), "ERC721: address zero is not a valid owner")
	if err != nil {
		return nil, err
	}
	return c.dal.GetBalance(owner)
}

func (c *ERC721Contract) OwnerOf(tokenId *common.SafeUint256) (common.Account, error) {
	//address owner = _ownerOf(tokenId);
	owner, err := c.dal.GetTokenOwner(tokenId)
	if err != nil {
		return nil, err
	}
	err = common.Require(!owner.IsZero(), "ERC721: invalid token ID")
	if err != nil {
		return nil, err
	}
	return owner, nil
}

func (c *ERC721Contract) SafeTransferFrom2(from, to common.Account, tokenId *common.SafeUint256, data []byte) error {
	sender, err := c.sdk.GetTxSender()
	if err != nil {
		return err
	}
	_isApprovedOrOwner, err := c.baseIsApprovedOrOwner(sender, tokenId)
	if err != nil {
		return err
	}
	err = common.Require(_isApprovedOrOwner, "ERC721: caller is not token owner or approved")
	if err != nil {
		return err
	}
	return c.baseSafeTransfer(from, to, tokenId, data)
}

func (c *ERC721Contract) SafeTransferFrom(from, to common.Account, tokenId *common.SafeUint256) error {
	return c.SafeTransferFrom2(from, to, tokenId, nil)
}

func (c *ERC721Contract) TransferFrom(from, to common.Account, tokenId *common.SafeUint256) error {
	sender, err := c.sdk.GetTxSender()
	if err != nil {
		return err
	}
	_isApprovedOrOwner, err := c.baseIsApprovedOrOwner(sender, tokenId)
	if err != nil {
		return err
	}
	err = common.Require(_isApprovedOrOwner, "ERC721: caller is not token owner or approved")
	if err != nil {
		return err
	}
	return c.baseTransfer(from, to, tokenId)
}

func (c *ERC721Contract) Approve(to common.Account, tokenId *common.SafeUint256) error {
	//address owner = ERC721.ownerOf(tokenId);
	owner, err := c.dal.GetTokenOwner(tokenId)
	if err != nil {
		return err
	}
	common.Require(!to.Equal(owner), "ERC721: approval to current owner")
	sender, err := c.sdk.GetTxSender()
	if err != nil {
		return err
	}
	isApprovedForAll, err := c.IsApprovedForAll(owner, sender)
	err = common.Require(sender.Equal(owner) || isApprovedForAll,
		"ERC721: approve caller is not token owner or approved for all")
	if err != nil {
		return err
	}
	return c.baseApprove(to, tokenId)
}

func (c *ERC721Contract) SetApprovalForAll(operator common.Account, approved bool) error {
	//_setApprovalForAll(_msgSender(), operator, approved);
	sender, err := c.sdk.GetTxSender()
	if err != nil {
		return err
	}
	return c.dal.SetOperatorApproval(sender, operator, approved)
}

func (c *ERC721Contract) GetApproved(tokenId *common.SafeUint256) (common.Account, error) {
	//_requireMinted(tokenId);
	err := common.Require(c.exists(tokenId), "ERC721: invalid token ID")
	if err != nil {
		return nil, err
	}
	return c.dal.GetTokenApproval(tokenId)
}

func (c *ERC721Contract) IsApprovedForAll(owner common.Account, operator common.Account) (bool, error) {
	//return _operatorApprovals[owner][operator];
	return c.dal.GetOperatorApproval(owner, operator)
}

func (c *ERC721Contract) SupportsInterface(interfaceId string) bool {
	//return
	//interfaceId == type(IERC721).interfaceId ||
	//	interfaceId == type(IERC721Metadata).interfaceId ||
	//	super.supportsInterface(interfaceId);
	return interfaceId == "ERC721" || interfaceId == "ERC721Metadata" || interfaceId == "ERC165"
}

func (c *ERC721Contract) Name() (string, error) {
	return c.dal.GetName()
}

func (c *ERC721Contract) Symbol() (string, error) {
	return c.dal.GetSymbol()
}

func (c *ERC721Contract) TokenURI(tokenId *common.SafeUint256) (string, error) {
	//_requireMinted(tokenId);
	err := common.Require(c.exists(tokenId), "ERC721: invalid token ID")
	if err != nil {
		return "", err
	}
	baseURI, err := c.dal.GetBaseURI()
	if err != nil {
		return "", err
	}
	if len(baseURI) > 0 {
		return fmt.Sprintf(baseURI, tokenId.ToString()), nil
	}
	return "", nil
}

func (c *ERC721Contract) Mint(account common.Account, tokenId *common.SafeUint256) (bool, error) {

	//check is admin
	sender, err := c.sdk.GetTxSender()
	if err != nil {
		return false, fmt.Errorf("Get sender address failed, err:%s", err)
	}

	admin, err := c.dal.GetAdmin()
	if err != nil {
		return false, err
	}
	if !sender.Equal(admin) {
		return false, errors.New("only admin can mint tokens")
	}
	//call base mint
	err = c.baseMint(account, tokenId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *ERC721Contract) Burn(tokenId *common.SafeUint256) error {
	sender, err := c.sdk.GetTxSender()
	if err != nil {
		return err
	}
	_isApprovedOrOwner, err := c.baseIsApprovedOrOwner(sender, tokenId)
	if err != nil {
		return err
	}
	err = common.Require(_isApprovedOrOwner, "ERC721: caller is not token owner or approved")
	if err != nil {
		return err
	}
	return c.baseBurn(tokenId)
}
