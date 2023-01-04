package erc20

import (
	"github.com/studyzy/token-go/common"
)

type IERC20 interface {
	Name() (string, error)
	Symbol() (string, error)
	Decimals() (uint8, error)
	TotalSupply() (*common.SafeUint256, error)
	BalanceOf(account common.Account) (*common.SafeUint256, error)
	Transfer(to common.Account, amount *common.SafeUint256) (bool, error)
	TransferFrom(from, to common.Account, amount *common.SafeUint256) (bool, error)
	Approve(spender common.Account, amount *common.SafeUint256) (bool, error)
	Allowance(owner, spender common.Account) (*common.SafeUint256, error)
}
type Mintable interface {
	Mint(account common.Account, amount *common.SafeUint256) (bool, error)
}
type Burnable interface {
	Burn(amount *common.SafeUint256) (bool, error)
	BurnFrom(account common.Account, amount *common.SafeUint256) (bool, error)
}
