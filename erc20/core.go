package erc20

import (
	"errors"

	"github.com/studyzy/openzeppelin-go/common"
)

// Option 初始化ERC20合约的选项
type Option struct {
	// BeforeTransfer 在转账前执行的逻辑
	BeforeTransfer func(from common.Account, to common.Account, amount *common.SafeUint256) error
	// AfterTransfer 在转账成功后执行的逻辑
	AfterTransfer func(from common.Account, to common.Account, amount *common.SafeUint256) error
	// Burnable 是否允许销毁
	Burnable bool
	// Minable 是否允许后续铸造
	Minable bool
}

func checkAccount(acct ...common.Account) error {
	for _, acc := range acct {
		if acc.IsZero() {
			return errors.New(" the zero address")
		}
	}
	return nil
}
func (c *ERC20Contract) SetSDK(sdk common.ContractSDK) {
	c.sdk = sdk
}
func (c *ERC20Contract) baseTransfer(from common.Account, to common.Account, amount *common.SafeUint256, option Option) error {
	//检查from和to的合法性
	err := checkAccount(from, to)
	if err != nil {
		return err
	}
	//触发用户自定义的BeforeTransfer
	if option.BeforeTransfer != nil {
		if err = option.BeforeTransfer(from, to, amount); err != nil {
			return err
		}
	}
	//检查from余额充足
	fromBalance, err := c.dal.GetBalance(from)
	if err != nil {
		return err
	}
	if !fromBalance.GTE(amount) {
		return errors.New("ERC20: transfer amount exceeds balance")
	}
	//更新from和to的余额
	fromNewBalance, _ := common.SafeSub(fromBalance, amount)
	err = c.dal.SetBalance(from, fromNewBalance)
	if err != nil {
		return err
	}
	toBalance, err := c.dal.GetBalance(to)
	if err != nil {
		return err
	}
	toNewBalance, ok := common.SafeAdd(toBalance, amount)
	if !ok {
		return errors.New("calculate new to balance error")
	}
	err = c.dal.SetBalance(to, toNewBalance)
	if err != nil {
		return err
	}
	//触发事件

	c.sdk.EmitEvent("transfer", from.ToString(), to.ToString(), amount.ToString())
	//触发用户自定义的afterTransfer
	if option.AfterTransfer != nil {
		return option.AfterTransfer(from, to, amount)
	}
	return nil
}

func (c *ERC20Contract) baseApprove(owner common.Account, spender common.Account, amount *common.SafeUint256, option Option) error {
	//检查from和to的合法性
	err := checkAccount(owner, spender)
	if err != nil {
		return err
	}
	//设置Allowance
	err = c.dal.SetAllowance(owner, spender, amount)
	if err != nil {
		return err
	}
	//触发事件Approval
	c.sdk.EmitEvent("approve", owner.ToString(), spender.ToString(), amount.ToString())
	return nil
}

func (c *ERC20Contract) baseSpendAllowance(owner common.Account, spender common.Account, amount *common.SafeUint256, option Option) error {
	//获得授权的额度
	currentAllowance, err := c.dal.GetAllowance(owner, spender)
	if err != nil {
		return err
	}
	//计算额度是否够用
	if !currentAllowance.GTE(amount) {
		return errors.New("ERC20: insufficient allowance")
	}
	//扣减授权额度
	newCurrentAllowance, ok := common.SafeSub(currentAllowance, amount)
	if !ok {
		return errors.New("spend allowance error")
	}
	return c.baseApprove(owner, spender, newCurrentAllowance, option)
}
func (c *ERC20Contract) baseMint(account common.Account, amount *common.SafeUint256, option Option) error {
	//检查account的合法性
	err := checkAccount(account)
	if err != nil {
		return err
	}
	from := c.sdk.NewZeroAccount()

	//触发用户自定义的BeforeTransfer
	if option.BeforeTransfer != nil {
		if err = option.BeforeTransfer(from, account, amount); err != nil {
			return err
		}
	}
	//更新TotalSupply
	totalSupply, err := c.dal.GetTotalSupply()
	if err != nil {
		return err
	}
	newTotal, ok := common.SafeAdd(totalSupply, amount)
	if !ok {
		return errors.New("calculate totalSupply failed")
	}
	err = c.dal.SetTotalSupply(newTotal)
	if err != nil {
		return err
	}
	//更新余额
	toBalance, err := c.dal.GetBalance(account)
	if err != nil {
		return err
	}
	toNewBalance, ok := common.SafeAdd(toBalance, amount)
	if !ok {
		return errors.New("calculate new to balance error")
	}
	err = c.dal.SetBalance(account, toNewBalance)
	if err != nil {
		return err
	}
	//触发事件
	c.sdk.EmitEvent("transfer", from.ToString(), account.ToString(), amount.ToString())
	//触发用户自定义的afterTransfer
	if option.AfterTransfer != nil {
		return option.AfterTransfer(from, account, amount)
	}
	return nil
}

func (c *ERC20Contract) baseBurn(account common.Account, amount *common.SafeUint256, option Option) error {
	//检查account的合法性
	err := checkAccount(account)
	if err != nil {
		return err
	}
	to := c.sdk.NewZeroAccount()
	//触发用户自定义的BeforeTransfer
	if option.BeforeTransfer != nil {
		if err = option.BeforeTransfer(account, to, amount); err != nil {
			return err
		}
	}
	//检查用户余额充足
	fromBalance, err := c.dal.GetBalance(account)
	if err != nil {
		return err
	}
	if !fromBalance.GTE(amount) {
		return errors.New("ERC20: burn amount exceeds balance")
	}
	//更新TotalSupply
	totalSupply, err := c.dal.GetTotalSupply()
	if err != nil {
		return err
	}
	newTotal, ok := common.SafeSub(totalSupply, amount)
	if !ok {
		return errors.New("calculate totalSupply failed")
	}
	err = c.dal.SetTotalSupply(newTotal)
	if err != nil {
		return err
	}
	//更新余额
	fromNewBalance, ok := common.SafeSub(fromBalance, amount)
	if !ok {
		return errors.New("calculate new to balance error")
	}
	err = c.dal.SetBalance(account, fromNewBalance)
	if err != nil {
		return err
	}
	//触发事件
	c.sdk.EmitEvent("transfer", account.ToString(), to.ToString(), amount.ToString())
	//触发用户自定义的afterTransfer
	if option.AfterTransfer != nil {
		return option.AfterTransfer(account, to, amount)
	}
	return nil
}
