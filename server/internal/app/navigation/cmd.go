package navigation

type CreateNavMenuCmd struct {
	Name     string
	URL      string
	ParentID *int64
	Icon     *string
}

type UpdateNavMenuCmd struct {
	ID       int64
	Name     string
	URL      string
	ParentID *int64
	Icon     *string
	Sort     *int
}

type NavMenuOrderItem struct {
	ID       int64
	ParentID *int64
	Sort     int
}
