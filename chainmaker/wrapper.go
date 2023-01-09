package chainmaker

import (
	"encoding/json"
	"strconv"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"github.com/studyzy/openzeppelin-go/common"
)

// ReturnUint256 封装返回SafeUint256类型为Response，如果有error则忽略num，封装error
// @param num
// @param err
// @return Response
func ReturnUint256(num *common.SafeUint256, err error) protogo.Response {
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte(num.ToString()))
}

// ReturnString 封装返回string类型为Response，如果有error则忽略str，封装error
// @param str
// @param err
// @return Response
func ReturnString(str string, err error) protogo.Response {
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte(str))
}

// ReturnBool 封装返回Bool类型为Response，如果有error则忽略bool，封装error
// @param b
// @param err
// @return Response
func ReturnBool(b bool, err error) protogo.Response {
	if err != nil {
		return sdk.Error(err.Error())
	}
	if b {
		return sdk.Success([]byte("true"))
	}
	return sdk.Success([]byte("false"))
}

// ReturnUint8 封装返回uint8类型为Response，如果有error则忽略num，封装error
// @param num
// @param err
// @return Response
func ReturnUint8(num uint8, err error) protogo.Response {
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte(strconv.Itoa(int(num))))
}

// ReturnJson 封装返回对象类型为json格式到Response，如果有error则忽略对象，封装error
// @param obj
// @param err
// @return Response
func ReturnJson(obj interface{}, err error) protogo.Response {
	if err != nil {
		return sdk.Error(err.Error())
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success(data)
}
