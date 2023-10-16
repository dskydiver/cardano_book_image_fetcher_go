package pkg

import "cardano-book-image-fetcher/pkg/api"

func VerifyPolicyID(policyID string, collections []api.Collection) bool {
	for _, collection := range collections {
		if collection.CollectionID == policyID {
			return true
		}
	}
	return false
}
