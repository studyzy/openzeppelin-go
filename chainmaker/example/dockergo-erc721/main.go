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
	"errors"
	"fmt"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"github.com/studyzy/openzeppelin-go/chainmaker"
	"github.com/studyzy/openzeppelin-go/common"
	"github.com/studyzy/openzeppelin-go/erc721"
)

type ERC721DockerGo struct {
	supper  *erc721.ERC721Contract
	methods map[string]func() protogo.Response
	adapter *chainmaker.SdkAdapter
}

func NewERC721DockerGo() *ERC721DockerGo {
	erc721Option := erc721.Option{
		BeforeTransfer: nil,
		AfterTransfer:  nil,
		Burnable:       true,
		Minable:        true,
	}
	adapter := chainmaker.NewSdkAdapter(sdk.Instance)
	contract := &ERC721DockerGo{methods: make(map[string]func() protogo.Response), adapter: adapter}
	contract.supper = erc721.NewERC721Contract(erc721Option, "TestNFT", "TNT", adapter)
	contract.registerMethods(erc721Option)
	return contract
}
func (erc721 *ERC721DockerGo) registerMethods(option erc721.Option) {

	erc721.RegisterMethod("name", erc721.name)
	erc721.RegisterMethod("symbol", erc721.symbol)
	erc721.RegisterMethod("tokenURI", erc721.tokenURI)
	erc721.RegisterMethod("balanceOf", erc721.balanceOf)
	erc721.RegisterMethod("ownerOf", erc721.ownerOf)
	erc721.RegisterMethod("safeTransferFrom", erc721.safeTransferFrom)
	erc721.RegisterMethod("transferFrom", erc721.transferFrom)
	erc721.RegisterMethod("approve", erc721.approve)
	erc721.RegisterMethod("setApprovalForAll", erc721.setApprovalForAll)
	erc721.RegisterMethod("getApproved", erc721.getApproved)
	erc721.RegisterMethod("isApprovedForAll", erc721.isApprovedForAll)
	if option.Minable {
		erc721.RegisterMethod("mint", erc721.mint)
	}
	if option.Burnable {
		erc721.RegisterMethod("burn", erc721.burn)
	}
}
func (erc721 *ERC721DockerGo) RegisterMethod(methodName string, fun func() protogo.Response) {
	erc721.methods[methodName] = fun
}
func (erc721 *ERC721DockerGo) InitContract() protogo.Response {
	err := erc721.updateErc20Info()
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte("Init contract success"))
}

// UpgradeContract used to upgrade contract
func (erc721 *ERC721DockerGo) UpgradeContract() protogo.Response {
	err := erc721.updateErc20Info()
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte("Upgrade contract success"))
}

// UpgradeContract upgrade contract func
func (erc721 *ERC721DockerGo) updateErc20Info() error {
	args := sdk.Instance.GetArgs()
	// name, symbol and decimal are optional
	name := string(args["name"])
	symbol := string(args["symbol"])

	admin, err := sdk.Instance.Sender()
	if err != nil {
		return fmt.Errorf("get sender failed, err:%s", err)
	}
	adminAccount, err := erc721.adapter.NewAccountFromString(admin)
	if err != nil {
		return fmt.Errorf("get sender failed, err:%s", err)
	}
	//此处支持在安装合约的时候指定name,symbol
	//如果没有参数指定，那么就使用NewERC20Contract构造的时候的值
	err = erc721.supper.InitERC721(name, symbol, adminAccount)
	if err != nil {
		return fmt.Errorf("set admin failed, err:%s", err)
	}
	return nil
}

// InvokeContract used to invoke user contract
func (erc721 *ERC721DockerGo) InvokeContract(method string) protogo.Response {
	if len(method) == 0 {
		return sdk.Error("method of param should not be empty")
	}
	if fun, ok := erc721.methods[method]; ok {
		return fun()
	}
	return sdk.Error("Invalid method")
}
func (erc721 *ERC721DockerGo) requireBool(key string) (bool, error) {
	args := sdk.Instance.GetArgs()
	acc, ok := args[key]
	if !ok {
		return false, errors.New("require bool:" + key)
	}
	return string(acc) == "true", nil
}
func (erc721 *ERC721DockerGo) requireAccount(key string) (common.Account, error) {
	args := sdk.Instance.GetArgs()
	acc, ok := args[key]
	if !ok {
		return nil, errors.New("require account:" + key)
	}
	return erc721.adapter.NewAccountFromString(string(acc))
}
func (erc721 *ERC721DockerGo) requireTokenId(key string) (*common.SafeUint256, error) {
	args := sdk.Instance.GetArgs()
	tokenId, ok := args[key]
	if !ok {
		return nil, errors.New("require account:" + key)
	}
	num, ok := common.ParseSafeUint256(string(tokenId))
	if !ok {
		return nil, errors.New("invalid uint256")
	}
	return num, nil
}

func (erc721 *ERC721DockerGo) name() protogo.Response {
	return chainmaker.ReturnString(erc721.supper.Name())
}

func (erc721 *ERC721DockerGo) symbol() protogo.Response {
	return chainmaker.ReturnString(erc721.supper.Symbol())
}

func (erc721 *ERC721DockerGo) balanceOf() protogo.Response {
	account, err := erc721.requireAccount("account")
	if err != nil {
		return sdk.Error(err.Error())
	}
	return chainmaker.ReturnUint256(erc721.supper.BalanceOf(account))
}
func (erc721 *ERC721DockerGo) ownerOf() protogo.Response {
	tokenId, err := erc721.requireTokenId("tokenId")
	if err != nil {
		return sdk.Error(err.Error())
	}
	return chainmaker.ReturnAccount(erc721.supper.OwnerOf(tokenId))
}

//    function safeTransferFrom(address from, address to, uint256 tokenId, bytes calldata data) external;
func (erc721 *ERC721DockerGo) safeTransferFrom() protogo.Response {
	from, err := erc721.requireAccount("from")
	if err != nil {
		return sdk.Error(err.Error())
	}
	to, err := erc721.requireAccount("to")
	if err != nil {
		return sdk.Error(err.Error())
	}
	tokenId, err := erc721.requireTokenId("tokenId")
	if err != nil {
		return sdk.Error(err.Error())
	}
	args := sdk.Instance.GetArgs()
	data := args["data"]
	if len(data) > 0 {
		return chainmaker.Return(erc721.supper.SafeTransferFrom2(from, to, tokenId, data))
	}
	return chainmaker.Return(erc721.supper.SafeTransferFrom(from, to, tokenId))
}

func (erc721 *ERC721DockerGo) transferFrom() protogo.Response {
	from, err := erc721.requireAccount("from")
	if err != nil {
		return sdk.Error(err.Error())
	}
	to, err := erc721.requireAccount("to")
	if err != nil {
		return sdk.Error(err.Error())
	}
	tokenId, err := erc721.requireTokenId("tokenId")
	if err != nil {
		return sdk.Error(err.Error())
	}

	return chainmaker.Return(erc721.supper.TransferFrom(from, to, tokenId))
}

func (erc721 *ERC721DockerGo) approve() protogo.Response {
	to, err := erc721.requireAccount("to")
	if err != nil {
		return sdk.Error(err.Error())
	}
	tokenId, err := erc721.requireTokenId("tokenId")
	if err != nil {
		return sdk.Error(err.Error())
	}
	return chainmaker.Return(erc721.supper.Approve(to, tokenId))
}

//     function setApprovalForAll(address operator, bool approved) external;
func (erc721 *ERC721DockerGo) setApprovalForAll() protogo.Response {
	operator, err := erc721.requireAccount("operator")
	if err != nil {
		return sdk.Error(err.Error())
	}
	approved, err := erc721.requireBool("approved")
	if err != nil {
		return sdk.Error(err.Error())
	}
	return chainmaker.Return(erc721.supper.SetApprovalForAll(operator, approved))

}

//    function getApproved(uint256 tokenId) external view returns (address operator);
func (erc721 *ERC721DockerGo) getApproved() protogo.Response {
	tokenId, err := erc721.requireTokenId("tokenId")
	if err != nil {
		return sdk.Error(err.Error())
	}
	return chainmaker.ReturnAccount(erc721.supper.GetApproved(tokenId))
}

//    function isApprovedForAll(address owner, address operator) external view returns (bool);
func (erc721 *ERC721DockerGo) isApprovedForAll() protogo.Response {
	owner, err := erc721.requireAccount("owner")
	if err != nil {
		return sdk.Error(err.Error())
	}
	operator, err := erc721.requireAccount("operator")
	if err != nil {
		return sdk.Error(err.Error())
	}

	return chainmaker.ReturnBool(erc721.supper.IsApprovedForAll(owner, operator))
}

//safeMint(address to, uint256 tokenId)
func (erc721 *ERC721DockerGo) mint() protogo.Response {
	to, err := erc721.requireAccount("to")
	if err != nil {
		return sdk.Error(err.Error())
	}
	tokenId, err := erc721.requireTokenId("tokenId")
	if err != nil {
		return sdk.Error(err.Error())
	}
	return chainmaker.Return(erc721.supper.Mint(to, tokenId))
}

//    function burn(uint256 tokenId) public virtual
func (erc721 *ERC721DockerGo) burn() protogo.Response {
	tokenId, err := erc721.requireTokenId("tokenId")
	if err != nil {
		return sdk.Error(err.Error())
	}
	return chainmaker.Return(erc721.supper.Burn(tokenId))
}

//    function tokenURI(uint256 tokenId) external view returns (string memory);
func (erc721 *ERC721DockerGo) tokenURI() protogo.Response {
	tokenId, err := erc721.requireTokenId("tokenId")
	if err != nil {
		return sdk.Error(err.Error())
	}
	return chainmaker.ReturnString(erc721.supper.TokenURI(tokenId))
}

func main() {
	erc20 := NewERC721DockerGo()
	err := sandbox.Start(erc20)
	if err != nil {
		sdk.Instance.Errorf(err.Error())
	}
}
