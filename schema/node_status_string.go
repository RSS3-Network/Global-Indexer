// Code generated by "enumer --values --type=Status --linecomment --output node_status_string.go --json --yaml --sql"; DO NOT EDIT.

package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const _StatusName = "onlineoffline"

var _StatusIndex = [...]uint8{0, 6, 13}

const _StatusLowerName = "onlineoffline"

func (i Status) String() string {
	if i < 0 || i >= Status(len(_StatusIndex)-1) {
		return fmt.Sprintf("Status(%d)", i)
	}
	return _StatusName[_StatusIndex[i]:_StatusIndex[i+1]]
}

func (Status) Values() []string {
	return StatusStrings()
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _StatusNoOp() {
	var x [1]struct{}
	_ = x[StatusOnline-(0)]
	_ = x[StatusOffline-(1)]
}

var _StatusValues = []Status{StatusOnline, StatusOffline}

var _StatusNameToValueMap = map[string]Status{
	_StatusName[0:6]:       StatusOnline,
	_StatusLowerName[0:6]:  StatusOnline,
	_StatusName[6:13]:      StatusOffline,
	_StatusLowerName[6:13]: StatusOffline,
}

var _StatusNames = []string{
	_StatusName[0:6],
	_StatusName[6:13],
}

// StatusString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func StatusString(s string) (Status, error) {
	if val, ok := _StatusNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _StatusNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Status values", s)
}

// StatusValues returns all values of the enum
func StatusValues() []Status {
	return _StatusValues
}

// StatusStrings returns a slice of all String values of the enum
func StatusStrings() []string {
	strs := make([]string, len(_StatusNames))
	copy(strs, _StatusNames)
	return strs
}

// IsAStatus returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Status) IsAStatus() bool {
	for _, v := range _StatusValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for Status
func (i Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Status
func (i *Status) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Status should be a string, got %s", data)
	}

	var err error
	*i, err = StatusString(s)
	return err
}

// MarshalYAML implements a YAML Marshaler for Status
func (i Status) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for Status
func (i *Status) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var err error
	*i, err = StatusString(s)
	return err
}

func (i Status) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *Status) Scan(value interface{}) error {
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
		return fmt.Errorf("invalid value of Status: %[1]T(%[1]v)", value)
	}

	val, err := StatusString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}
