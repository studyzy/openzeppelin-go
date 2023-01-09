package common

import "errors"

type ChainBase interface {
	NewAccountFromBytes(b []byte) (Account, error)
	NewAccountFromString(str string) (Account, error)
	NewZeroAccount() Account
}
type StateOperator interface {
	ChainBase
	GetState(key string) (value []byte, err error)
	PutState(key string, value []byte) error
	DelState(key string) error
	CreateCompositeKey(prefix string, data ...string) (string, error)
}
type ContractSDK interface {
	StateOperator
	GetTxSender() (Account, error)
	EmitEvent(topic string, data ...string) error
	IsContract(account Account) bool
	CallContract(account Account, method string, args []KeyValue) Response
}

func Require(exp bool, msg string) error {
	if !exp {
		return errors.New(msg)
	}
	return nil
}
