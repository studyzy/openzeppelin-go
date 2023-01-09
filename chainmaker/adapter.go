package chainmaker

import (
	"bytes"
	"strings"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"chainmaker.org/chainmaker/contract-utils/address"
	"github.com/studyzy/openzeppelin-go/common"
)

type SdkAdapter struct {
	cmsdk sdk.SDKInterface
}

func NewSdkAdapter(cmsdk sdk.SDKInterface) *SdkAdapter {
	return &SdkAdapter{cmsdk: cmsdk}
}

func (s SdkAdapter) NewAccountFromBytes(b []byte) (common.Account, error) {
	return s.NewAccountFromString(string(b))
}

func (s SdkAdapter) NewAccountFromString(str string) (common.Account, error) {
	addr, _ := address.ParseAddress(str)
	a := Address{addr: *addr}
	return &a, nil
}

func (s SdkAdapter) NewZeroAccount() common.Account {
	zero, _ := address.ParseAddress(address.ZeroAddr)
	return &Address{addr: *zero}
}

func (s SdkAdapter) DelState(key string) error {
	return s.cmsdk.DelState(key, "")
}

func (s SdkAdapter) CreateCompositeKey(prefix string, data ...string) (string, error) {
	return prefix + "_" + strings.Join(data, "_"), nil
}

func (s SdkAdapter) IsContract(account common.Account) bool {
	//TODO
	return false
}

func (s SdkAdapter) CallContract(account common.Account, method string, args []common.KeyValue) common.Response {
	response := s.cmsdk.CallContract(account.ToString(), method, args2Map(args))
	return common.Response{
		Status:  response.Status,
		Message: response.Message,
		Payload: response.Payload,
	}
}
func args2Map(args []common.KeyValue) map[string][]byte {
	m := make(map[string][]byte)
	for _, arg := range args {
		m[arg.Key] = arg.Value
	}
	return m
}

func (s SdkAdapter) GetState(key string) (value []byte, err error) {
	str, err := s.cmsdk.GetState(key, "")
	return []byte(str), err
}

func (s SdkAdapter) PutState(key string, value []byte) error {
	return s.cmsdk.PutState(key, "", string(value))
}

func (s SdkAdapter) GetTxSender() (common.Account, error) {
	sender, err := s.cmsdk.Sender()
	if err != nil {
		return nil, err
	}
	return s.NewAccountFromString(sender)
}

func (s SdkAdapter) EmitEvent(topic string, data ...string) error {
	s.cmsdk.EmitEvent(topic, data)
	return nil
}

var _ common.ContractSDK = (*SdkAdapter)(nil)

var _ common.Account = (*Address)(nil)

type Address struct {
	addr address.Address
}

func (a *Address) IsZero() bool {
	return address.IsZeroAddress(a.addr.ToString())
}

func (a *Address) ToString() string {
	return a.addr.ToString()
}
func (a *Address) Bytes() []byte {
	return []byte(a.addr.ToString())
}

func (a *Address) Equal(account common.Account) bool {
	return bytes.Equal(a.Bytes(), account.Bytes())
}