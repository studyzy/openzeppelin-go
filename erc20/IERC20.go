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

package erc20

import (
	"github.com/studyzy/openzeppelin-go/common"
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
