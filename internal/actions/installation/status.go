package installation

import (
	"context"

	"github.com/runforyou-ai/lite-ai-gateway/internal/sysdoc"
	"github.com/uptrace/bun"
)

// Status 返回是否已安装及站点名。
func Status(ctx context.Context, db bun.IDB) (installed bool, siteName string, err error) {
	inst, ok, err := sysdoc.GetState[sysdoc.Installation](ctx, db, sysdoc.StateInstallation, true)
	if err != nil {
		return false, "", err
	}
	if ok && inst.Installed {
		installed = true
	}
	site, ok, err := sysdoc.GetSettings[sysdoc.Site](ctx, db, sysdoc.SettingsSite, true)
	if err != nil {
		return false, "", err
	}
	if ok {
		siteName = site.Name
	} else {
		siteName = sysdoc.DefaultSite().Name
	}
	return installed, siteName, nil
}
