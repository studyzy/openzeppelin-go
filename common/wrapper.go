package common

import (
	"encoding/json"
	"strconv"
)

// ReturnUint256 封装返回SafeUint256类型为Response，如果有error则忽略num，封装error
// @param num
// @param err
// @return Response
func ReturnUint256(num *SafeUint256, err error) Response {
	if err != nil {
		return Error(err.Error())
	}
	return Success([]byte(num.ToString()))
}

// ReturnString 封装返回string类型为Response，如果有error则忽略str，封装error
// @param str
// @param err
// @return Response
func ReturnString(str string, err error) Response {
	if err != nil {
		return Error(err.Error())
	}
	return Success([]byte(str))
}

// ReturnBool 封装返回Bool类型为Response，如果有error则忽略bool，封装error
// @param b
// @param err
// @return Response
func ReturnBool(b bool, err error) Response {
	if err != nil {
		return Error(err.Error())
	}
	if b {
		return Success([]byte("true"))
	}
	return Success([]byte("false"))
}

// ReturnUint8 封装返回uint8类型为Response，如果有error则忽略num，封装error
// @param num
// @param err
// @return Response
func ReturnUint8(num uint8, err error) Response {
	if err != nil {
		return Error(err.Error())
	}
	return Success([]byte(strconv.Itoa(int(num))))
}

// ReturnJson 封装返回对象类型为json格式到Response，如果有error则忽略对象，封装error
// @param obj
// @param err
// @return Response
func ReturnJson(obj interface{}, err error) Response {
	if err != nil {
		return Error(err.Error())
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return Error(err.Error())
	}
	return Success(data)
}
