package enforcer

import (
	"crypto/sha256"

	"github.com/rss3-network/global-indexer/internal/distributor"
)

const (
	requestUnit        = 1
	invalidRequestUnit = 1
)

// updateRequestsBasedOnDataCompare updates the requests based on the data comparison responses.
func updateRequestsBasedOnDataCompare(responses []distributor.DataResponse) {
	errResponseCount := markErrorResponsesAndCount(responses)

	if len(responses) == distributor.DefaultNodeCount-1 {
		handleTwoResponses(responses)
	} else if len(responses) == distributor.DefaultNodeCount-2 {
		handleSingleResponse(responses)
	} else {
		handleFullResponses(responses, errResponseCount)
	}
}

// compareData compares two byte slices and returns true if they are equal.
func compareData(src, des []byte) bool {
	if src == nil || des == nil {
		return false
	}

	srcHash, destHash := sha256.Sum256(src), sha256.Sum256(des)

	return string(srcHash[:]) == string(destHash[:])
}

// markErrorResponsesAndCount marks the error results and returns the count of error results.
func markErrorResponsesAndCount(responses []distributor.DataResponse) (errResponseCount int) {
	for i := range responses {
		if responses[i].Err != nil {
			responses[i].InvalidRequest = invalidRequestUnit
			errResponseCount++
		}
	}

	return errResponseCount
}

// handleTwoResponses handles the case when there are two results.
func handleTwoResponses(responses []distributor.DataResponse) {
	if responses[0].Err == nil && responses[1].Err == nil {
		updateRequestBasedOnComparison(responses)
	} else {
		markErrorResponse(responses...)
	}
}

// handleSingleResponse handles the case when there is only one response.
func handleSingleResponse(responses []distributor.DataResponse) {
	if responses[0].Err == nil {
		responses[0].Request = requestUnit
	} else {
		responses[0].InvalidRequest = invalidRequestUnit
	}
}

// handleFullResponses handles the case when there are more than two results.
func handleFullResponses(responses []distributor.DataResponse, errResponseCount int) {
	if errResponseCount < len(responses) {
		compareAndAssignRequests(responses, errResponseCount)
	}
}

func updateRequestBasedOnComparison(responses []distributor.DataResponse) {
	if compareData(responses[0].Data, responses[1].Data) {
		responses[0].Request = 2 * requestUnit
		responses[1].Request = requestUnit
	} else {
		responses[0].Request = requestUnit
	}
}

func markErrorResponse(responses ...distributor.DataResponse) {
	for i, result := range responses {
		if result.Err != nil {
			responses[i].InvalidRequest = invalidRequestUnit
		} else {
			responses[i].Request = requestUnit
		}
	}
}

func compareAndAssignRequests(responses []distributor.DataResponse, errResponseCount int) {
	d0, d1, d2 := responses[0].Data, responses[1].Data, responses[2].Data
	diff01, diff02, diff12 := compareData(d0, d1), compareData(d0, d2), compareData(d1, d2)

	switch errResponseCount {
	// responses contain 2 errors
	case len(responses) - 1:
		responses[0].Request = requestUnit
	// responses contain 1 error
	case len(responses) - 2:
		if diff01 {
			responses[0].Request = 2 * requestUnit
			responses[1].Request = requestUnit
		} else {
			responses[0].Request = requestUnit
		}
	// responses contain no error
	case len(responses) - 3:
		if diff01 && diff02 && diff12 {
			responses[0].Request = 2 * requestUnit
			responses[1].Request = requestUnit
			responses[2].Request = requestUnit
		} else if diff01 && !diff02 {
			responses[0].Request = 2 * requestUnit
			responses[1].Request = requestUnit
			responses[2].InvalidRequest = invalidRequestUnit
		} else if diff02 && !diff01 {
			responses[0].Request = 2 * requestUnit
			responses[1].InvalidRequest = invalidRequestUnit
			responses[2].Request = requestUnit
		} else if diff12 {
			// if the second response is non-null
			if responses[1].Valid {
				responses[0].InvalidRequest = invalidRequestUnit
				responses[1].Request = requestUnit
				responses[2].Request = requestUnit
			} else {
				// it means the last two responses are the null data
				responses[0].Request = requestUnit
			}
		} else if !diff01 {
			responses[0].Request = requestUnit
		}
	}
}
