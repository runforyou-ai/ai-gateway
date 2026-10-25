package settings

import (
	"context"
	"errors"
	"strings"

	"github.com/runforyou-ai/ai-gateway/internal/sysdoc"
	"github.com/uptrace/bun"
)

var (
	// ErrInvalidName 表示站点名称无效。
	ErrInvalidName = errors.New("invalid name")
)

// Site 是站点设置视图。
type Site struct {
	Name             string
	RegistrationOpen bool
}

// GetSite 读取站点设置；缺失时返回默认值。
func GetSite(ctx context.Context, db bun.IDB) (Site, error) {
	site, ok, err := sysdoc.GetSettings[sysdoc.Site](ctx, db, sysdoc.SettingsSite, true)
	if err != nil {
		return Site{}, err
	}
	if !ok {
		d := sysdoc.DefaultSite()
		return Site{Name: d.Name, RegistrationOpen: d.RegistrationOpen}, nil
	}
	return Site{Name: site.Name, RegistrationOpen: site.RegistrationOpen}, nil
}

// UpdateSite 写入站点设置。
func UpdateSite(ctx context.Context, db bun.IDB, in Site) (Site, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || len([]rune(name)) > 64 {
		return Site{}, ErrInvalidName
	}
	site := sysdoc.Site{
		Name:             name,
		RegistrationOpen: in.RegistrationOpen,
	}
	if err := sysdoc.PutSettings(ctx, db, sysdoc.SettingsSite, site, nil); err != nil {
		return Site{}, err
	}
	return Site{Name: site.Name, RegistrationOpen: site.RegistrationOpen}, nil
}
