package contract

// FederationFriendLinkRequestReq 联合友链申请请求。
type FederationFriendLinkRequestReq struct {
	RequesterURL string `json:"requester_url"`
	Message      string `json:"message,omitempty"`
	RSSURL       string `json:"rss_url,omitempty"`
}

// FederationCitationSourcePost 引用来源文章信息。
type FederationCitationSourcePost struct {
	ID    string `json:"id,omitempty"`
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

// FederationCitationRequestReq 跨站引用请求。
type FederationCitationRequestReq struct {
	SourceInstanceURL string                       `json:"source_instance_url"`
	SourcePost        FederationCitationSourcePost `json:"source_post"`
	TargetPostID      string                       `json:"target_post_id"`
	CitationContext   string                       `json:"citation_context,omitempty"`
	CitationType      string                       `json:"citation_type,omitempty"`
}

// FederationCitationDecisionReq 跨站引用批准/拒绝。
type FederationCitationDecisionReq struct {
	CitationID        int64  `json:"citation_id"`
	SourceInstanceURL string `json:"source_instance_url,omitempty"`
	Reason            string `json:"reason,omitempty"`
}

// FederationMentionSourcePost 提及来源文章。
type FederationMentionSourcePost struct {
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

// FederationMentionNotifyReq 跨站提及通知。
type FederationMentionNotifyReq struct {
	SourceInstanceURL string                      `json:"source_instance_url"`
	SourcePost        FederationMentionSourcePost `json:"source_post"`
	MentionedUser     string                      `json:"mentioned_user"`
	MentionContext    string                      `json:"mention_context"`
	MentionType       string                      `json:"mention_type,omitempty"`
}
