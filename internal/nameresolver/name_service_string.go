// Code generated by "enumer --values --type=NameService --linecomment --output name_service_string.go --json --sql"; DO NOT EDIT.

package nameresolver

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const _NameServiceName = "unknownethcsblensfc"

var _NameServiceIndex = [...]uint8{0, 7, 10, 13, 17, 19}

const _NameServiceLowerName = "unknownethcsblensfc"

func (i NameService) String() string {
	if i < 0 || i >= NameService(len(_NameServiceIndex)-1) {
		return fmt.Sprintf("NameService(%d)", i)
	}
	return _NameServiceName[_NameServiceIndex[i]:_NameServiceIndex[i+1]]
}

func (NameService) Values() []string {
	return NameServiceStrings()
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _NameServiceNoOp() {
	var x [1]struct{}
	_ = x[NameServiceUnknown-(0)]
	_ = x[NameServiceENS-(1)]
	_ = x[NameServiceCSB-(2)]
	_ = x[NameServiceLens-(3)]
	_ = x[NameServiceFarcaster-(4)]
}

var _NameServiceValues = []NameService{NameServiceUnknown, NameServiceENS, NameServiceCSB, NameServiceLens, NameServiceFarcaster}

var _NameServiceNameToValueMap = map[string]NameService{
	_NameServiceName[0:7]:        NameServiceUnknown,
	_NameServiceLowerName[0:7]:   NameServiceUnknown,
	_NameServiceName[7:10]:       NameServiceENS,
	_NameServiceLowerName[7:10]:  NameServiceENS,
	_NameServiceName[10:13]:      NameServiceCSB,
	_NameServiceLowerName[10:13]: NameServiceCSB,
	_NameServiceName[13:17]:      NameServiceLens,
	_NameServiceLowerName[13:17]: NameServiceLens,
	_NameServiceName[17:19]:      NameServiceFarcaster,
	_NameServiceLowerName[17:19]: NameServiceFarcaster,
}

var _NameServiceNames = []string{
	_NameServiceName[0:7],
	_NameServiceName[7:10],
	_NameServiceName[10:13],
	_NameServiceName[13:17],
	_NameServiceName[17:19],
}

// NameServiceString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func NameServiceString(s string) (NameService, error) {
	if val, ok := _NameServiceNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _NameServiceNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to NameService values", s)
}

// NameServiceValues returns all values of the enum
func NameServiceValues() []NameService {
	return _NameServiceValues
}

// NameServiceStrings returns a slice of all String values of the enum
func NameServiceStrings() []string {
	strs := make([]string, len(_NameServiceNames))
	copy(strs, _NameServiceNames)
	return strs
}

// IsANameService returns "true" if the value is listed in the enum definition. "false" otherwise
func (i NameService) IsANameService() bool {
	for _, v := range _NameServiceValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for NameService
func (i NameService) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for NameService
func (i *NameService) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("NameService should be a string, got %s", data)
	}

	var err error
	*i, err = NameServiceString(s)
	return err
}

func (i NameService) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *NameService) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of NameService: %[1]T(%[1]v)", value)
	}

	val, err := NameServiceString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}