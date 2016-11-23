package subscription

import (
	"encoding/json"
	"fmt"
	"github.com/griesbacher/check_azure/azureHttp"
	"github.com/griesbacher/check_x"
)

type JsonSubscription struct {
	Value []struct {
		DisplayName          string `json:"displayName"`
		ID                   string `json:"id"`
		State                string `json:"state"`
		SubscriptionID       string `json:"subscriptionId"`
		SubscriptionPolicies struct {
			LocationPlacementID string `json:"locationPlacementId"`
			QuotaID             string `json:"quotaId"`
			SpendingLimit       string `json:"spendingLimit"`
		} `json:"subscriptionPolicies"`
	} `json:"value"`
}

func Display(ac azureHttp.AzureConnector) error {
	body, typ, err := ac.Request("2014-04-01", "", "")
	if err != nil {
		return err
	}
	if typ != azureHttp.ContentTypeJSON {
		return azureHttp.ContentError(azureHttp.ContentTypeJSON, typ)
	}
	var subs JsonSubscription
	err = json.Unmarshal(body, &subs)
	if err != nil {
		return err
	}
	msg := ""
	for i, sub := range subs.Value {
		msg += fmt.Sprintf("%d: %s\n\tID: %s\n\tSpendingLimit: %s\n", i, sub.DisplayName, sub.SubscriptionID, sub.SubscriptionPolicies.SpendingLimit)
	}
	check_x.LongExit(check_x.OK, "see below", msg)
	return nil
}
