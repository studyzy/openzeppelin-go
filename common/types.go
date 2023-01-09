package common

import (
	"math/big"
)

type Account interface {
	ToString() string
	Bytes() []byte
	Equal(Account) bool
	IsZero() bool
}

type SafeUint256 big.Int

func NewSafeUint256(i uint64) *SafeUint256 {
	z := big.NewInt(int64(i))
	return (*SafeUint256)(z)
}

var (
	// SafeUintOne is 1
	SafeUintOne    = (*SafeUint256)(big.NewInt(1))
	maxSafeUint256 *SafeUint256
	minSafeUint256 *SafeUint256
)

func init() {
	x := big.NewInt(1)
	x = x.Lsh(x, 256).Sub(x, big.NewInt(1))
	maxSafeUint256 = (*SafeUint256)(x)
	minSafeUint256 = (*SafeUint256)(big.NewInt(0))
}

// ToString get str from uint256
func (x *SafeUint256) ToString() string {
	return (*big.Int)(x).String()
}

// GTE if x>=y return true
func (x *SafeUint256) GTE(y *SafeUint256) bool {
	return (*big.Int)(x).Cmp((*big.Int)(y)) >= 0
}

// ParseSafeUint256 get uint256 obj from str
func ParseSafeUint256(x string) (*SafeUint256, bool) {
	z := big.NewInt(0)
	if x == "" {
		return (*SafeUint256)(z), true
	}
	z, ok := z.SetString(x, 10)
	if !ok || z.Cmp((*big.Int)(maxSafeUint256)) > 0 || z.Cmp((*big.Int)(minSafeUint256)) < 0 {
		return nil, false
	}
	return (*SafeUint256)(z), true
}

// SafeAdd sets z to the sum x+y and returns z.
func SafeAdd(x, y *SafeUint256) (*SafeUint256, bool) {
	z := big.NewInt(0)
	z = z.Add((*big.Int)(x), (*big.Int)(y))
	if z.Cmp((*big.Int)(maxSafeUint256)) > 0 {
		return nil, false
	}
	return (*SafeUint256)(z), true
}

// SafeSub sets z to the difference x-y and returns z.
func SafeSub(x, y *SafeUint256) (*SafeUint256, bool) {
	if (*big.Int)(x).Cmp((*big.Int)(y)) < 0 {
		return nil, false
	}
	return (*SafeUint256)((*big.Int)(x).Sub((*big.Int)(x), (*big.Int)(y))), true
}

// SafeMul sets z to the product x*y and returns z.
func SafeMul(x, y *SafeUint256) (*SafeUint256, bool) {
	z := (*big.Int)(x).Mul((*big.Int)(x), (*big.Int)(y))
	if z.Cmp((*big.Int)(maxSafeUint256)) > 0 || z.Cmp((*big.Int)(minSafeUint256)) < 0 {
		return nil, false
	}
	return (*SafeUint256)(z), true
}

// SafeDiv sets z to the quotient x/y for y != 0 and returns z.
// If y == 0, a division-by-zero run-time panic occurs.
// Div implements Euclidean division (unlike Go); see DivMod for more details.
func SafeDiv(x, y *SafeUint256) *SafeUint256 {
	return (*SafeUint256)((*big.Int)(x).Div((*big.Int)(x), (*big.Int)(y)))
}

type KeyValue struct {
	Key   string
	Value []byte
}

type Response struct {
	// A status code that should follow the HTTP status codes.
	Status int32
	// A message associated with the response code. error has message
	Message string
	// A payload that can be used to include metadata with this response. success with payload
	Payload []byte
}

const (
	// OK constant - status code less than 400, endorser will endorse it.
	// OK means init or invoke successfully.
	OK = 0

	// ERROR constant - default error value
	ERROR = 1
)

//// Success ...
//func Success(payload []byte) Response {
//	return Response{
//		Status:  OK,
//		Payload: payload,
//	}
//}
//
//// Error ...
//func Error(msg string) Response {
//	return Response{
//		Status:  ERROR,
//		Message: msg,
//	}
//}
