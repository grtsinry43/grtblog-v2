package contract

// FriendLinkApplicationReq 友链申请请求。
type FriendLinkApplicationReq struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	Message     string `json:"message"`
	RSSURL      string `json:"rssUrl"`
}
