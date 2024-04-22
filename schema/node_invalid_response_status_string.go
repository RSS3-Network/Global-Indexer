// Code generated by "enumer --values --type=NodeInvalidResponseStatus --linecomment --output node_invalid_response_status_string.go --json --yaml --sql"; DO NOT EDIT.

package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const _NodeInvalidResponseStatusName = "challengeable"

var _NodeInvalidResponseStatusIndex = [...]uint8{0, 13}

const _NodeInvalidResponseStatusLowerName = "challengeable"

func (i NodeInvalidResponseStatus) String() string {
	if i < 0 || i >= NodeInvalidResponseStatus(len(_NodeInvalidResponseStatusIndex)-1) {
		return fmt.Sprintf("NodeInvalidResponseStatus(%d)", i)
	}
	return _NodeInvalidResponseStatusName[_NodeInvalidResponseStatusIndex[i]:_NodeInvalidResponseStatusIndex[i+1]]
}

func (NodeInvalidResponseStatus) Values() []string {
	return NodeInvalidResponseStatusStrings()
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _NodeInvalidResponseStatusNoOp() {
	var x [1]struct{}
	_ = x[NodeInvalidResponseStatusChallengeable-(0)]
}

var _NodeInvalidResponseStatusValues = []NodeInvalidResponseStatus{NodeInvalidResponseStatusChallengeable}

var _NodeInvalidResponseStatusNameToValueMap = map[string]NodeInvalidResponseStatus{
	_NodeInvalidResponseStatusName[0:13]:      NodeInvalidResponseStatusChallengeable,
	_NodeInvalidResponseStatusLowerName[0:13]: NodeInvalidResponseStatusChallengeable,
}

var _NodeInvalidResponseStatusNames = []string{
	_NodeInvalidResponseStatusName[0:13],
}

// NodeInvalidResponseStatusString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func NodeInvalidResponseStatusString(s string) (NodeInvalidResponseStatus, error) {
	if val, ok := _NodeInvalidResponseStatusNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _NodeInvalidResponseStatusNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to NodeInvalidResponseStatus values", s)
}

// NodeInvalidResponseStatusValues returns all values of the enum
func NodeInvalidResponseStatusValues() []NodeInvalidResponseStatus {
	return _NodeInvalidResponseStatusValues
}

// NodeInvalidResponseStatusStrings returns a slice of all String values of the enum
func NodeInvalidResponseStatusStrings() []string {
	strs := make([]string, len(_NodeInvalidResponseStatusNames))
	copy(strs, _NodeInvalidResponseStatusNames)
	return strs
}

// IsANodeInvalidResponseStatus returns "true" if the value is listed in the enum definition. "false" otherwise
func (i NodeInvalidResponseStatus) IsANodeInvalidResponseStatus() bool {
	for _, v := range _NodeInvalidResponseStatusValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for NodeInvalidResponseStatus
func (i NodeInvalidResponseStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for NodeInvalidResponseStatus
func (i *NodeInvalidResponseStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("NodeInvalidResponseStatus should be a string, got %s", data)
	}

	var err error
	*i, err = NodeInvalidResponseStatusString(s)
	return err
}

// MarshalYAML implements a YAML Marshaler for NodeInvalidResponseStatus
func (i NodeInvalidResponseStatus) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for NodeInvalidResponseStatus
func (i *NodeInvalidResponseStatus) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var err error
	*i, err = NodeInvalidResponseStatusString(s)
	return err
}

func (i NodeInvalidResponseStatus) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *NodeInvalidResponseStatus) Scan(value interface{}) error {
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
		return fmt.Errorf("invalid value of NodeInvalidResponseStatus: %[1]T(%[1]v)", value)
	}

	val, err := NodeInvalidResponseStatusString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}
