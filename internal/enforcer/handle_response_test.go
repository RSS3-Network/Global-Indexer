package enforcer

import (
	"errors"
	"testing"

	"github.com/rss3-network/global-indexer/internal/distributor"
	"github.com/stretchr/testify/assert"
)

func TestUpdateRequestsBasedOnDataCompare(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		responses       []distributor.DataResponse
		requests        []int
		invalidRequests []int
	}{
		{
			name: "one_error_response",
			responses: []distributor.DataResponse{
				{Err: errors.New("error")},
			},
			requests:        []int{0},
			invalidRequests: []int{1},
		},
		{
			name: "one_valid_response",
			responses: []distributor.DataResponse{
				{Data: []byte("data1"), Valid: true},
			},
			requests:        []int{1},
			invalidRequests: []int{0},
		},
		{
			name: "two_error_responses",
			responses: []distributor.DataResponse{
				{Err: errors.New("error1")},
				{Err: errors.New("error2")},
			},
			requests:        []int{0, 0},
			invalidRequests: []int{1, 1},
		},
		{
			name: "one_error_with_two_responses",
			responses: []distributor.DataResponse{
				{Data: []byte("data1"), Valid: true},
				{Err: errors.New("error")},
			},
			requests:        []int{1, 0},
			invalidRequests: []int{0, 1},
		},
		{
			name: "two_responses_with_different_data",
			responses: []distributor.DataResponse{
				{Data: []byte("data1"), Valid: true},
				{Data: []byte("data2"), Valid: true},
			},
			requests:        []int{1, 0},
			invalidRequests: []int{0, 0},
		},
		{
			name: "two_responses_with_same_data",
			responses: []distributor.DataResponse{
				{Data: []byte("data1"), Valid: true},
				{Data: []byte("data1"), Valid: true},
			},
			requests:        []int{2, 1},
			invalidRequests: []int{0, 0},
		},
		{
			name: "three_errors",
			responses: []distributor.DataResponse{
				{Err: errors.New("error1")},
				{Err: errors.New("error2")},
				{Err: errors.New("error3")},
			},
			requests:        []int{0, 0, 0},
			invalidRequests: []int{1, 1, 1},
		},
		{
			name: "two_errors",
			responses: []distributor.DataResponse{
				{Data: []byte("data1"), Valid: true},
				{Err: errors.New("error2")},
				{Err: errors.New("error3")},
			},
			requests:        []int{1, 0, 0},
			invalidRequests: []int{0, 1, 1},
		},
		{
			name: "one_error_with_same_data",
			responses: []distributor.DataResponse{
				{Data: []byte("data1")},
				{Data: []byte("data1")},
				{Err: errors.New("error3")},
			},
			requests:        []int{2, 1, 0},
			invalidRequests: []int{0, 0, 1},
		},
		{
			name: "one_error_with_different_data",
			responses: []distributor.DataResponse{
				{Data: []byte("data1")},
				{Data: []byte("data2")},
				{Err: errors.New("error3")},
			},
			requests:        []int{1, 0, 0},
			invalidRequests: []int{0, 0, 1},
		},
		{
			name: "three_same_data",
			responses: []distributor.DataResponse{
				{Data: []byte("data1")},
				{Data: []byte("data1")},
				{Data: []byte("data1")},
			},
			requests:        []int{2, 1, 1},
			invalidRequests: []int{0, 0, 0},
		},
		{
			name: "three_different_data",
			responses: []distributor.DataResponse{
				{Data: []byte("data1")},
				{Data: []byte("data2")},
				{Data: []byte("data3")},
			},
			requests:        []int{1, 0, 0},
			invalidRequests: []int{0, 0, 0},
		},
		{
			name: "two_different_data_01",
			responses: []distributor.DataResponse{
				{Data: []byte("data1")},
				{Data: []byte("data1")},
				{Data: []byte("data2")},
			},
			requests:        []int{2, 1, 0},
			invalidRequests: []int{0, 0, 1},
		},
		{
			name: "two_different_data_02",
			responses: []distributor.DataResponse{
				{Data: []byte("data1")},
				{Data: []byte("data2")},
				{Data: []byte("data1")},
			},
			requests:        []int{2, 0, 1},
			invalidRequests: []int{0, 1, 0},
		},
		{
			name: "two_different_data_12_with_valid",
			responses: []distributor.DataResponse{
				{Data: []byte("data0")},
				{Data: []byte("data1"), Valid: true},
				{Data: []byte("data1")},
			},
			requests:        []int{0, 1, 1},
			invalidRequests: []int{1, 0, 0},
		},
		{
			name: "two_different_data_12_with_invalid",
			responses: []distributor.DataResponse{
				{Data: []byte("data0")},
				{Data: []byte("data1")},
				{Data: []byte("data1")},
			},
			requests:        []int{1, 0, 0},
			invalidRequests: []int{0, 0, 0},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			updateRequestsBasedOnDataCompare(tc.responses)
			for i, result := range tc.responses {
				assert.Equal(t, tc.requests[i], result.Request)
				assert.Equal(t, tc.invalidRequests[i], result.InvalidRequest)
			}
		})
	}
}
