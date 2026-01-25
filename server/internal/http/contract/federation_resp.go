package contract

import "time"

// FederationFriendLinkResponseResp 联合友链申请响应。
type FederationFriendLinkResponseResp struct {
	ApplicationID int64  `json:"applicationId"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

// FederationCitationResponseResp 引用请求响应。
type FederationCitationResponseResp struct {
	CitationID int64  `json:"citation_id"`
	Status     string `json:"status"`
}

// FederationCitationDecisionResp 引用审批响应。
type FederationCitationDecisionResp struct {
	CitationID int64  `json:"citation_id"`
	Status     string `json:"status"`
}

// FederationMentionNotifyResp 提及通知响应。
type FederationMentionNotifyResp struct {
	MentionID int64 `json:"mention_id"`
	Delivered bool  `json:"delivered"`
}

// FederationPostAuthorResp 联合时间线作者信息。
type FederationPostAuthorResp struct {
	Name   string  `json:"name"`
	URL    *string `json:"url,omitempty"`
	Avatar *string `json:"avatar,omitempty"`
}

// FederationPostResp 联合时间线文章条目。
type FederationPostResp struct {
	ID             string                   `json:"id"`
	URL            string                   `json:"url"`
	Title          string                   `json:"title"`
	Summary        string                   `json:"summary"`
	ContentPreview *string                  `json:"content_preview,omitempty"`
	Author         FederationPostAuthorResp `json:"author"`
	PublishedAt    time.Time                `json:"published_at"`
	UpdatedAt      *time.Time               `json:"updated_at,omitempty"`
	CoverImage     *string                  `json:"cover_image,omitempty"`
	Language       *string                  `json:"language,omitempty"`
	AllowCitation  bool                     `json:"allow_citation"`
	AllowComment   bool                     `json:"allow_comment"`
}

// FederationTimelineResp 联合时间线响应。
type FederationTimelineResp struct {
	Items []FederationPostResp `json:"items"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
}

// FederationPostDetailResp 文章详情响应。
type FederationPostDetailResp struct {
	Post         FederationPostResp   `json:"post"`
	RelatedPosts []FederationPostResp `json:"related_posts,omitempty"`
}
