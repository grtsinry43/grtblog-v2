package handler

import (
	"encoding/json"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/federationconfig"
)

type federationPolicy struct {
	AllowCitation                 *bool `json:"allow_citation"`
	AllowMention                  *bool `json:"allow_mention"`
	AutoApproveFriendlink         *bool `json:"auto_approve_friendlink"`
	AutoApproveFriendlinkCitation *bool `json:"auto_approve_friendlink_citation"`
}

func parseFederationPolicy(settings federationconfig.Settings) federationPolicy {
	if len(settings.DefaultPolicies) == 0 {
		return federationPolicy{}
	}
	var policy federationPolicy
	if err := json.Unmarshal(settings.DefaultPolicies, &policy); err != nil {
		return federationPolicy{}
	}
	return policy
}

func policyBool(val *bool, fallback bool) bool {
	if val == nil {
		return fallback
	}
	return *val
}
