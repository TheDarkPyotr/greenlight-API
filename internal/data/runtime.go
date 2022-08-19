package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

//to allow input as "102 mins" and enable conversion from JSON to struct
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {

	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unquotedJSONValue, " ") //isolate the number
	if len(parts) != 2 || parts[1] != "mins" {     //verify format
		return ErrInvalidRuntimeFormat
	}

	// 32 bit size
	number, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(number)
	return nil

}

//Custom representation for Runtime field during json.Marshal
func (r Runtime) MarshalJSON() ([]byte, error) {

	jsonVal := fmt.Sprintf("%d mins", r)
	//Use quote to represent as JSON string
	quotedJSONValue := strconv.Quote(jsonVal)
	return []byte(quotedJSONValue), nil

}
