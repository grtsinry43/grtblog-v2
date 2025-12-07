package media

import "time"

type UploadFile struct {
	ID        int64
	Name      string
	Path      string
	Type      string
	Size      int64
	CreatedAt time.Time
}

type Photo struct {
	ID          int64
	URL         string
	Location    *string
	Device      *string
	Shade       *string
	Description *string
	CreatedAt   time.Time
}
