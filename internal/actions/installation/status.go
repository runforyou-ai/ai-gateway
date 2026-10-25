package installation

import (
	"context"

	"github.com/runforyou-ai/ai-gateway/internal/model"
	"github.com/runforyou-ai/ai-gateway/internal/sysdoc"
	"github.com/uptrace/bun"
)

// Status 返回是否已有账号及站点名。
// 首个注册用户即为管理员，无单独安装向导 / create-admin CLI。
func Status(ctx context.Context, db bun.IDB) (installed bool, siteName string, err error) {
	n, err := db.NewSelect().Model((*model.User)(nil)).Count(ctx)
	if err != nil {
		return false, "", err
	}
	installed = n > 0

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
