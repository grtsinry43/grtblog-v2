package navigation

import "time"

type NavMenu struct {
	ID        int64
	Name      string
	URL       string
	Sort      int
	ParentID  *int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type FooterSection struct {
	ID        int64
	Title     string
	Links     map[string]any
	CreatedAt time.Time
	UpdatedAt time.Time
}
