package contract

type CreateNavMenuReq struct {
	Name     string  `json:"name"`
	URL      string  `json:"url"`
	ParentID *int64  `json:"parentId,omitempty"`
	Icon     *string `json:"icon,omitempty"`
}

type UpdateNavMenuReq struct {
	Name     string  `json:"name"`
	URL      string  `json:"url"`
	ParentID *int64  `json:"parentId,omitempty"`
	Icon     *string `json:"icon,omitempty"`
	Sort     *int    `json:"sort,omitempty"`
}

type NavMenuOrderItem struct {
	ID       int64  `json:"id"`
	ParentID *int64 `json:"parentId,omitempty"`
	Sort     int    `json:"sort"`
}

type ReorderNavMenuReq struct {
	Items []NavMenuOrderItem `json:"items"`
}
