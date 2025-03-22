package flags

import (
	"flag"

	"github.com/dmitrovia/passkeeper/internal/server/models/procattrs/serverpa"
)

func InitFlags(attr *serverpa.ServerProcAttr) {
	flag.StringVar(
		attr.GetDBDSN(),
		"db", attr.GetDefDBDSN(),
		"database connection address.",
	)
	flag.StringVar(
		attr.GetServerAddr(),
		"saddr", attr.GetDefServerAddr(),
		"Port to listen on.",
	)
	flag.StringVar(
		attr.GetConfigPath(),
		"cfgpath", attr.GetDefConfigPath(),
		"Port to listen on.",
	)
	flag.Parse()
}
