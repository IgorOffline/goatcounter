package goatcounter

import (
	"context"
	"fmt"
	"time"

	"zgo.at/errors"
	"zgo.at/zcache/v2"
	"zgo.at/zdb"
	"zgo.at/zstd/zjson"
	"zgo.at/zstd/zruntime"
)

type BosmangStat struct {
	ID        int64     `db:"site_id"`
	Codes     string    `db:"codes"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	LastMonth int       `db:"last_month"`
	Total     int       `db:"total"`
	Avg       int       `db:"avg"`
}

type BosmangStats []BosmangStat

// List stats for all sites, for all time.
func (a *BosmangStats) List(ctx context.Context) error {
	err := zdb.Select(ctx, a, "load:bosmang.List")
	if err != nil {
		return errors.Wrap(err, "BosmangStats.List")
	}
	return nil
}

func ListCache(ctx context.Context) map[string]struct {
	Size  int64
	Items map[string]string
} {
	c := make(map[string]struct {
		Size  int64
		Items map[string]string
	})

	caches := map[string]interface {
		ItemsAny() map[any]zcache.Item[any]
	}{
		"sites":          cacheSites(ctx),
		"ua":             cacheUA(ctx),
		"browsers":       cacheBrowsers(ctx),
		"systems":        cacheSystems(ctx),
		"paths":          cachePaths(ctx),
		"loc":            cacheLoc(ctx),
		"changed_titles": cacheChangedTitles(ctx),
		//"loader":         handler.loader.conns,
	}

	for name, f := range caches {
		var (
			content = f.ItemsAny()
			s       = zruntime.SizeOf(content)
			items   = make(map[string]string)
		)
		for k, v := range content {
			items[fmt.Sprintf("%v", k)] = fmt.Sprintf("%s\n", zjson.MustMarshalIndent(v.Object, "", "  "))
			s += c[name].Size + zruntime.SizeOf(v.Object)
		}
		c[name] = struct {
			Size  int64
			Items map[string]string
		}{s / 1024, items}
	}

	{
		var (
			name    = "sites_host"
			content = cacheSitesHost(ctx).Items()
			s       = zruntime.SizeOf(content)
			items   = make(map[string]string)
		)
		for k, v := range content {
			items[k] = fmt.Sprintf("%d", v)
			s += c[name].Size + zruntime.SizeOf(v)
		}
		c[name] = struct {
			Size  int64
			Items map[string]string
		}{s / 1024, items}
	}
	return c
}
