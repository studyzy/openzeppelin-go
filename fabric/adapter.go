package fabric

import (
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/studyzy/token-go/common"
)

var _ common.Account = (*MspUser)(nil)

type MspUser struct {
	ID string
}

func NewMspUser(id string) common.Account {
	return &MspUser{ID: id}
}
func (m MspUser) ToString() string {
	return m.ID
}

func (m MspUser) Bytes() []byte {
	return []byte(m.ID)
}

func (m MspUser) Equal(account common.Account) bool {
	return m.ID == account.ToString()
}

var _ common.ContractSDK = (*SdkAdapter)(nil)

type SdkAdapter struct {
	ctx contractapi.TransactionContextInterface
}

func NewSDkAdapter(ctx contractapi.TransactionContextInterface) *SdkAdapter {
	return &SdkAdapter{ctx: ctx}
}
func (s SdkAdapter) GetState(key string) (value []byte, err error) {
	return s.ctx.GetStub().GetState(key)
}

func (s SdkAdapter) PutState(key string, value []byte) error {
	return s.ctx.GetStub().PutState(key, value)
}

func (s SdkAdapter) GetTxSender() (common.Account, error) {
	id, err := s.ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, err
	}
	return NewMspUser(id), nil
}

type event struct {
	from  string
	to    string
	value int
}

func (s SdkAdapter) EmitEvent(topic string, data ...string) error {
	var payload []byte
	var err error
	if topic == "transfer" {
		val, _ := strconv.Atoi(data[2])
		transferEvent := event{data[0], data[1], val}
		payload, _ = json.Marshal(transferEvent)
	} else if topic == "approve" {
		//TODO
	} else {
		payload, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}
	return s.ctx.GetStub().SetEvent(topic, payload)
}

func (s SdkAdapter) GetArgs() map[string][]byte {
	//TODO implement me
	panic("implement me")
}
