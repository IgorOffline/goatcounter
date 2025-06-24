package goatcounter

import (
	"context"
	"strconv"
	"strings"

	"zgo.at/zdb"
)

func pgIn(ctx context.Context) zdb.SQL {
	if zdb.SQLDialect(ctx) == zdb.DialectPostgreSQL {
		return "= any"
	}
	return "in"
}

func pgArray[T ~int8 | ~int16 | ~int32 | ~int64](ctx context.Context, p []T) any {
	if zdb.SQLDialect(ctx) == zdb.DialectSQLite {
		return p
	}

	if len(p) == 0 {
		var zero T
		return zero
	}

	var b strings.Builder
	b.Grow(len(p) * 3)
	b.WriteString("'{")
	for i, pp := range p {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.FormatInt(int64(pp), 10))
	}
	b.WriteString("}'")
	return zdb.SQL(b.String())
}

var sqlArrayEscaper = strings.NewReplacer(`'`, `''`, `"`, `\"`)

func pgArrayString(ctx context.Context, p []string) any {
	if zdb.SQLDialect(ctx) == zdb.DialectSQLite {
		return p
	}
	if len(p) == 0 {
		return ""
	}

	b := new(strings.Builder)
	b.Grow(len(p) * 3)
	b.WriteString("'{")
	for i, pp := range p {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteByte('"')
		sqlArrayEscaper.WriteString(b, pp)
		b.WriteByte('"')
	}
	b.WriteString("}'")
	return zdb.SQL(b.String())
}
