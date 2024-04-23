package enforcer

import (
	"encoding/json"
	"sort"

	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
)

const (
	// a valid response gives 1 point
	validPointUnit = 1
	// an invalid response gives 1 point (in a bad way)
	invalidPointUnit = 1
)

// sortResponseByValidity sorts the responses based on the validity.
func sortResponseByValidity(responses []*model.DataResponse) {
	sort.SliceStable(responses, func(i, j int) bool {
		return (responses[i].Err == nil && responses[j].Err != nil) ||
			(responses[i].Err == nil && responses[j].Err == nil && responses[i].Valid && !responses[j].Valid)
	})
}

// updatePointsBasedOnIdentity updates both  based on responses identity.
func updatePointsBasedOnIdentity(responses []*model.DataResponse) {
	errResponseCount := countAndMarkErrorResponse(responses)

	if len(responses) == model.RequiredQualifiedNodeCount-1 {
		handleTwoResponses(responses)
	} else if len(responses) == model.RequiredQualifiedNodeCount-2 {
		handleSingleResponse(responses)
	} else {
		handleFullResponses(responses, errResponseCount)
	}
}

// isResponseIdentical returns true if two byte slices (responses) are identical.
func isResponseIdentical(src, des []byte) bool {
	srcActivity := &model.ActivityResponse{}
	desActivity := &model.ActivityResponse{}

	// check if the data is activity response
	if isDataValid(src, srcActivity) && isDataValid(des, desActivity) {
		if srcActivity.Data == nil && desActivity.Data == nil {
			return true
		} else if srcActivity.Data != nil && desActivity.Data != nil {
			if _, exist := model.MutablePlatformMap[srcActivity.Data.Platform]; !exist {
				// TODO: if false, save the record to the database
				return isActivityIdentical(srcActivity.Data, desActivity.Data)
			}

			return true
		}
	}

	srcActivities := &model.ActivitiesResponse{}
	desActivities := &model.ActivitiesResponse{}
	// check if the data is activities response
	if isDataValid(src, srcActivities) && isDataValid(des, desActivities) {
		if srcActivities.Data == nil && desActivities.Data == nil {
			return true
		} else if srcActivities.Data != nil && desActivities.Data != nil {
			// exclude the mutable platforms
			srcActivity, desActivity := excludeMutableActivity(srcActivities.Data), excludeMutableActivity(desActivities.Data)

			for i := range srcActivity {
				if !isActivityIdentical(srcActivity[i], desActivity[i]) {
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
func excludeMutableActivity(activities []*model.Activity) []*model.Activity {
	var newActivities []*model.Activity

	for i := range activities {
		if _, exist := model.MutablePlatformMap[activities[i].Platform]; !exist {
			newActivities = append(newActivities, activities[i])
		}
	}

	return newActivities
}

// isActivityIdentical returns true if two Activity are identical.
func isActivityIdentical(src, des *model.Activity) bool {
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

	// check if the inner actions are identical
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

// isDataValid returns true if the data is valid.
func isDataValid(data []byte, target any) bool {
	if err := json.Unmarshal(data, target); err == nil {
		return true
	}

	return false
}

// countAndMarkErrorResponse marks the error results and returns the count.
func countAndMarkErrorResponse(responses []*model.DataResponse) (errResponseCount int) {
	for i := range responses {
		if responses[i].Err != nil {
			responses[i].InvalidPoint = invalidPointUnit
			errResponseCount++
		}
	}

	return errResponseCount
}

// handleTwoResponses handles the case when there are two responses.
func handleTwoResponses(responses []*model.DataResponse) {
	if responses[0].Err == nil && responses[1].Err == nil {
		updateRequestBasedOnComparison(responses)
	} else {
		markErrorResponse(responses...)
	}
}

// handleSingleResponse handles the case when there is only one response.
func handleSingleResponse(responses []*model.DataResponse) {
	if responses[0].Err == nil {
		responses[0].ValidPoint = validPointUnit
	} else {
		responses[0].InvalidPoint = invalidPointUnit
	}
}

// handleFullResponses handles the case when there are more than two results.
func handleFullResponses(responses []*model.DataResponse, errResponseCount int) {
	if errResponseCount < len(responses) {
		compareAndAssignPoints(responses, errResponseCount)
	}
}

// updateRequestBasedOnComparison updates the requests based on the comparison of the data.
func updateRequestBasedOnComparison(responses []*model.DataResponse) {
	if isResponseIdentical(responses[0].Data, responses[1].Data) {
		responses[0].ValidPoint = 2 * validPointUnit
		responses[1].ValidPoint = validPointUnit
	} else {
		responses[0].ValidPoint = validPointUnit
	}
}

// markErrorResponse marks the error responses.
func markErrorResponse(responses ...*model.DataResponse) {
	for i, result := range responses {
		if result.Err != nil {
			responses[i].InvalidPoint = invalidPointUnit
		} else {
			responses[i].ValidPoint = validPointUnit
		}
	}
}

// compareAndAssignPoints compares the data for identity and assigns corresponding points.
func compareAndAssignPoints(responses []*model.DataResponse, errResponseCount int) {
	d0, d1, d2 := responses[0].Data, responses[1].Data, responses[2].Data
	diff01, diff02, diff12 := isResponseIdentical(d0, d1), isResponseIdentical(d0, d2), isResponseIdentical(d1, d2)

	switch errResponseCount {
	// responses contain 2 errors
	case len(responses) - 1:
		responses[0].ValidPoint = validPointUnit
	// responses contain 1 error
	case len(responses) - 2:
		if diff01 {
			responses[0].ValidPoint = 2 * validPointUnit
			responses[1].ValidPoint = validPointUnit
		} else {
			responses[0].ValidPoint = validPointUnit
		}
	// responses contain no errors
	case len(responses) - 3:
		if diff01 && diff02 && diff12 {
			responses[0].ValidPoint = 2 * validPointUnit
			responses[1].ValidPoint = validPointUnit
			responses[2].ValidPoint = validPointUnit
		} else if diff01 && !diff02 {
			responses[0].ValidPoint = 2 * validPointUnit
			responses[1].ValidPoint = validPointUnit
			responses[2].InvalidPoint = invalidPointUnit
		} else if diff02 && !diff01 {
			responses[0].ValidPoint = 2 * validPointUnit
			responses[1].InvalidPoint = invalidPointUnit
			responses[2].ValidPoint = validPointUnit
		} else if diff12 {
			// if the second response is non-null
			if responses[1].Valid {
				responses[0].InvalidPoint = invalidPointUnit
				responses[1].ValidPoint = validPointUnit
				responses[2].ValidPoint = validPointUnit
			} else {
				// the last two responses must include null data
				responses[0].ValidPoint = validPointUnit
			}
		} else if !diff01 {
			responses[0].ValidPoint = validPointUnit
		}
	}
}
