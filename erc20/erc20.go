/*
ERC20 Token Standard:
https://eips.ethereum.org/EIPS/eip-20
*/

package erc20

import (
	"errors"
	"fmt"

	"github.com/studyzy/openzeppelin-go/common"
)

var _ IERC20 = (*ERC20Contract)(nil)

// ERC20Contract erc20 contract
type ERC20Contract struct {
	option  Option
	_name   string
	_symbol string
	dal     *ERC20ContractDAL
	sdk     common.ContractSDK
}

// NewERC20Contract ERC20Contract
// @param option
// @param name
// @param symbol
// @return *ERC20Contract
func NewERC20Contract(option Option, name, symbol string, sdk common.ContractSDK) *ERC20Contract {
	erc20 := &ERC20Contract{
		option:  option,
		_name:   name,
		_symbol: symbol,
		sdk:     sdk,
		dal:     NewERC20ContractDAL(sdk),
	}
	return erc20
}

func (c *ERC20Contract) InitERC20(name, symbol string, decimals uint8, totalSupply *common.SafeUint256, admin common.Account) error {

	//此处支持在安装合约的时候指定name,symbol
	//如果没有参数指定，那么就使用NewERC20Contract构造的时候的值
	if len(name) > 0 {
		c._name = name
	}
	if err := c.dal.SetName(c._name); err != nil {
		return err
	}
	if len(symbol) > 0 {
		c._symbol = symbol

	}
	if err := c.dal.SetSymbol(c._symbol); err != nil {
		return err
	}
	//通过安装合约时参数可以修改decimals，如果不指定那么就是18位小数
	if err := c.dal.SetDecimals(decimals); err != nil {
		return err
	}
	//通过安装合约时参数可以指定发行总量，如果不指定则发行量是0，后期再调用mint函数来铸币
	//total supply default to zero
	if err := c.dal.SetTotalSupply(totalSupply); err != nil {
		return err
	}
	//set Admin，方便后面mint的时候判断权限
	if err := c.dal.SetAdmin(admin); err != nil {
		return fmt.Errorf("set admin failed, err:%s", err)
	}
	return nil
}

func (c *ERC20Contract) TotalSupply() (*common.SafeUint256, error) {
	return c.dal.GetTotalSupply()
}

func (c *ERC20Contract) BalanceOf(account common.Account) (*common.SafeUint256, error) {
	return c.dal.GetBalance(account)
}

func (c *ERC20Contract) Transfer(to common.Account, amount *common.SafeUint256) (bool, error) {
	from, err := c.sdk.GetTxSender()
	if err != nil {
		return false, fmt.Errorf("Get sender address failed, err:%s", err)
	}

	err = c.baseTransfer(from, to, amount, c.option)
	if err != nil {
		return false, err
	}
	return true, nil
}

/**
 * @dev See {IERC20-transferFrom}.
 *
 * Emits an {Approval} event indicating the updated allowance. This is not
 * required by the EIP. See the note at the beginning of {ERC20}.
 *
 * NOTE: Does not update the allowance if the current allowance
 * is the maximum `uint256`.
 *
 * Requirements:
 *
 * - `from` and `to` cannot be the zero address.
 * - `from` must have a balance of at least `amount`.
 * - the caller must have allowance for ``from``'s tokens of at least
 * `amount`.
 */
func (c *ERC20Contract) TransferFrom(from, to common.Account, amount *common.SafeUint256) (bool, error) {

	sender, err := c.sdk.GetTxSender()
	if err != nil {
		return false, fmt.Errorf("Get sender address failed, err:%s", err)
	}

	err = c.baseSpendAllowance(from, sender, amount, c.option)
	if err != nil {
		return false, fmt.Errorf("spend allowance failed, err:%s", err)
	}
	err = c.baseTransfer(from, to, amount, c.option)
	if err != nil {
		return false, err
	}
	return true, nil
}

/**
 * @dev See {IERC20-approve}.
 *
 * NOTE: If `amount` is the maximum `uint256`, the allowance is not updated on
 * `transferFrom`. This is semantically equivalent to an infinite approval.
 *
 * Requirements:
 *
 * - `spender` cannot be the zero address.
 */
func (c *ERC20Contract) Approve(spender common.Account, amount *common.SafeUint256) (bool, error) {
	sender, err := c.sdk.GetTxSender()
	if err != nil {
		return false, fmt.Errorf("Get sender address failed, err:%s", err)
	}

	err = c.baseApprove(sender, spender, amount, c.option)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *ERC20Contract) Allowance(owner, spender common.Account) (*common.SafeUint256, error) {
	return c.dal.GetAllowance(owner, spender)
}

func (c *ERC20Contract) Name() (string, error) {
	return c.dal.GetName()
}
func (c *ERC20Contract) Symbol() (string, error) {
	return c.dal.GetSymbol()
}

func (c *ERC20Contract) Decimals() (uint8, error) {
	return c.dal.GetDecimals()
}

/** Creates `amount` tokens and assigns them to `account`, increasing
 * the total supply.
 *
 * Emits a {Transfer} event with `from` set to the zero address.
 *
 * Requirements:
 *
 * - `account` cannot be the zero address.
 */
func (c *ERC20Contract) Mint(account common.Account, amount *common.SafeUint256) (bool, error) {

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
	err = c.baseMint(account, amount, c.option)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *ERC20Contract) Burn(amount *common.SafeUint256) (bool, error) {
	spender, err := c.sdk.GetTxSender()
	if err != nil {
		return false, fmt.Errorf("Get sender address failed, err:%s", err)
	}
	//call base burn
	err = c.baseBurn(spender, amount, c.option)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *ERC20Contract) BurnFrom(account common.Account, amount *common.SafeUint256) (bool, error) {
	spender, err := c.sdk.GetTxSender()
	if err != nil {
		return false, fmt.Errorf("Get sender address failed, err:%s", err)
	}
	err = c.baseSpendAllowance(account, spender, amount, c.option)
	if err != nil {
		return false, err
	}
	//call base burn
	err = c.baseBurn(account, amount, c.option)
	if err != nil {
		return false, err
	}
	return true, nil
}
