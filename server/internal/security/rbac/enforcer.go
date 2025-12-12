package rbac

import (
	"sync"
	"time"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/config"
)

type Enforcer struct {
	engine   *casbin.Enforcer
	db       *gorm.DB
	cfg      config.RBACConfig
	mu       sync.RWMutex
	stopChan chan struct{}
}

func NewEnforcer(cfg config.RBACConfig, db *gorm.DB) (*Enforcer, error) {
	engine, err := casbin.NewEnforcer(cfg.ModelPath)
	if err != nil {
		return nil, err
	}
	e := &Enforcer{
		engine:   engine,
		db:       db,
		cfg:      cfg,
		stopChan: make(chan struct{}),
	}
	if err := e.ReloadPolicies(); err != nil {
		return nil, err
	}
	if cfg.AutoReload && cfg.AutoReloadSeconds > 0 {
		go e.autoReload()
	}
	return e, nil
}

func (e *Enforcer) Close() {
	if e.stopChan != nil {
		close(e.stopChan)
	}
}

func (e *Enforcer) autoReload() {
	ticker := time.NewTicker(time.Duration(e.cfg.AutoReloadSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_ = e.ReloadPolicies()
		case <-e.stopChan:
			return
		}
	}
}

func (e *Enforcer) ReloadPolicies() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.engine.ClearPolicy()

	var rows []struct {
		Role       string
		Permission string
	}
	if err := e.db.Table("role_permission").
		Select("role.role_name AS role, permission.permission_name AS permission").
		Joins("JOIN role ON role.id = role_permission.role_id").
		Joins("JOIN permission ON permission.id = role_permission.permission_id").
		Scan(&rows).Error; err != nil {
		return err
	}

	for _, row := range rows {
		_, _ = e.engine.AddPolicy(row.Role, row.Permission, "allow")
	}
	return nil
}

func (e *Enforcer) HasPermission(roles []string, permission string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, role := range roles {
		ok, err := e.engine.Enforce(role, permission, "allow")
		if err == nil && ok {
			return true
		}
	}
	return false
}
