package enforcer

import (
	"encoding/json"

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
	srcActivity := &distributor.ActivityResponse{}
	desActivity := &distributor.ActivityResponse{}

	// check if the data is activity response
	if validateData(src, srcActivity) && validateData(des, desActivity) {
		if srcActivity.Data == nil && desActivity.Data == nil {
			return true
		} else if srcActivity.Data != nil && desActivity.Data != nil {
			if _, exist := distributor.MutablePlatformMap[srcActivity.Data.Platform]; !exist {
				// TODO: if false, save the record to the database
				return compareActivity(srcActivity.Data, desActivity.Data)
			}

			return true
		}
	}

	srcActivities := &distributor.ActivitiesResponse{}
	desActivities := &distributor.ActivitiesResponse{}
	// check if the data is activities	response
	if validateData(src, srcActivities) && validateData(des, desActivities) {
		if srcActivities.Data == nil && desActivities.Data == nil {
			return true
		} else if srcActivities.Data != nil && desActivities.Data != nil {
			// exclude the mutable platforms
			srcFeeds, desFeeds := excludeMutableActivity(srcActivities.Data), excludeMutableActivity(desActivities.Data)

			for i := range srcFeeds {
				if !compareActivity(srcFeeds[i], desFeeds[i]) {
					// TODO: if false, save the record to the database
					return false
				}
			}

			return true
		}
	}

	return false
}

// excludeMutableActivity excludes the mutable platforms from the activities.
func excludeMutableActivity(activities []*distributor.Feed) []*distributor.Feed {
	var newActivities []*distributor.Feed

	for i := range activities {
		if _, exist := distributor.MutablePlatformMap[activities[i].Platform]; !exist {
			newActivities = append(newActivities, activities[i])
		}
	}

	return newActivities
}

// compareActivity compares two activities and returns true if they are equal.
func compareActivity(src, des *distributor.Feed) bool {
	if src.ID != des.ID ||
		src.Network != des.Network ||
		src.Index != des.Index ||
		src.From != des.From ||
		src.To != des.To ||
		src.Owner != des.Owner ||
		src.Tag != des.Tag ||
		src.Type != des.Type ||
		src.Platform != des.Platform ||
		len(src.Actions) != len(des.Actions) {
		return false
	}

	if len(src.Actions) > 0 {
		for i := range des.Actions {
			if src.Actions[i].From != des.Actions[i].From ||
				src.Actions[i].To != des.Actions[i].To ||
				src.Actions[i].Tag != des.Actions[i].Tag ||
				src.Actions[i].Type != des.Actions[i].Type {
				return false
			}
		}
	}

	return true
}

// validateData validates the data and returns true if the data is valid.
func validateData(data []byte, target any) bool {
	if err := json.Unmarshal(data, target); err == nil {
		return true
	}

	return false
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

// updateRequestBasedOnComparison updates the requests based on the comparison of the data.
func updateRequestBasedOnComparison(responses []distributor.DataResponse) {
	if compareData(responses[0].Data, responses[1].Data) {
		responses[0].Request = 2 * requestUnit
		responses[1].Request = requestUnit
	} else {
		responses[0].Request = requestUnit
	}
}

// markErrorResponse marks the error responses.
func markErrorResponse(responses ...distributor.DataResponse) {
	for i, result := range responses {
		if result.Err != nil {
			responses[i].InvalidRequest = invalidRequestUnit
		} else {
			responses[i].Request = requestUnit
		}
	}
}

// compareAndAssignRequests compares the data and assigns the requests.
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
