package enforcer

import (
	"crypto/sha256"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
)

const (
	requestUnit        = 1
	invalidRequestUnit = 1
)

// updateRequestsBasedOnDataCompare updates the requests based on the data comparison results.
func updateRequestsBasedOnDataCompare(results []model.DataResponse) {
	errResultCount := markErrorResultsAndCount(results)

	if len(results) == model.DefaultNodeCount-1 {
		handleTwoResults(results)
	} else if len(results) == model.DefaultNodeCount-2 {
		handleSingleResult(results)
	} else {
		handleFullResults(results, errResultCount)
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

// markErrorResultsAndCount marks the error results and returns the count of error results.
func markErrorResultsAndCount(results []model.DataResponse) (errResultCount int) {
	for i := range results {
		if results[i].Err != nil {
			results[i].InvalidRequest = invalidRequestUnit
			errResultCount++
		}
	}

	return errResultCount
}

// handleTwoResults handles the case when there are two results.
func handleTwoResults(results []model.DataResponse) {
	if results[0].Err == nil && results[1].Err == nil {
		updateRequestBasedOnComparison(results)
	} else {
		markInvalidIfError(results...)
	}
}

// handleSingleResult handles the case when there is only one result.
func handleSingleResult(results []model.DataResponse) {
	if results[0].Err == nil {
		results[0].Request = requestUnit
	} else {
		results[0].InvalidRequest = invalidRequestUnit
	}
}

// handleFullResults handles the case when there are more than two results.
func handleFullResults(results []model.DataResponse, errResultCount int) {
	if errResultCount < len(results) {
		compareAndAssignRequests(results, errResultCount)
	}
}

func updateRequestBasedOnComparison(results []model.DataResponse) {
	if compareData(results[0].Data, results[1].Data) {
		results[0].Request = 2 * requestUnit
		results[1].Request = requestUnit
	} else {
		results[0].Request = requestUnit
	}
}

func markInvalidIfError(results ...model.DataResponse) {
	for i, result := range results {
		if result.Err != nil {
			results[i].InvalidRequest = invalidRequestUnit
		} else {
			results[i].Request = requestUnit
		}
	}
}

func compareAndAssignRequests(results []model.DataResponse, errResultCount int) {
	d0, d1, d2 := results[0].Data, results[1].Data, results[2].Data
	diff01, diff02, diff12 := compareData(d0, d1), compareData(d0, d2), compareData(d1, d2)

	switch errResultCount {
	// results contain 2 errors
	case len(results) - 1:
		results[0].Request = requestUnit
	// results contain 1 error
	case len(results) - 2:
		if diff01 {
			results[0].Request = 2 * requestUnit
			results[1].Request = requestUnit
		} else {
			results[0].Request = requestUnit
		}
	// results contain no error
	case len(results) - 3:
		if diff01 && diff02 && diff12 {
			results[0].Request = 2 * requestUnit
			results[1].Request = requestUnit
			results[2].Request = requestUnit
		} else if diff01 && !diff02 {
			results[0].Request = 2 * requestUnit
			results[1].Request = requestUnit
			results[2].InvalidRequest = invalidRequestUnit
		} else if diff02 && !diff01 {
			results[0].Request = 2 * requestUnit
			results[1].InvalidRequest = invalidRequestUnit
			results[2].Request = requestUnit
		} else if diff12 {
			// if the second result is non-null
			if results[1].Valid {
				results[0].InvalidRequest = invalidRequestUnit
				results[1].Request = requestUnit
				results[2].Request = requestUnit
			} else {
				// it means the last two results are the null data
				results[0].Request = requestUnit
			}
		} else if !diff01 {
			results[0].Request = requestUnit
		}
	}
}
