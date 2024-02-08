package handlers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Split(t *testing.T) {
	reqUri := "/data/v1/test"
	pathSplits := strings.Split(reqUri, "/")
	assert.Equal(t, pathSplits[1], "data")
	assert.Equal(t, fmt.Sprintf("/%s", strings.Join(pathSplits[2:], "/")), "/v1/test")
}
