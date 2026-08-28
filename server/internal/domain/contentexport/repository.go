package contentexport

import (
	"context"
	"errors"
)

var (
	ErrNotFound      = errors.New("export not found")
	ErrExportRunning = errors.New("another export is already running")
	ErrInvalidTicket = errors.New("invalid or expired download ticket")
)

type Repository interface {
	Create(context.Context, *Record) error
	Update(context.Context, *Record) error
	Get(context.Context, string) (*Record, error)
	List(context.Context) ([]Record, error)
	Delete(context.Context, string) error
	MarkInterrupted(context.Context) error
	CreateTicket(context.Context, DownloadTicket) error
	ResolveTicket(context.Context, string) (*Record, error)
	DeleteExpiredTickets(context.Context) error
}
