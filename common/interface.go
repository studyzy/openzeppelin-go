package common

import "errors"

type StateOperator interface {
	GetState(key string) (value []byte, err error)
	PutState(key string, value []byte) error
}
type ContractSDK interface {
	StateOperator
	GetTxSender() (Account, error)
	EmitEvent(topic string, data ...string) error
	GetArgs() map[string][]byte
}

func Required(exp bool, msg string) error {
	if !exp {
		return errors.New(msg)
	}
	return nil
}
