package utils

import (
	"flag"
)

// KDAuditInitFlags is used to store kdaudit flags
// that are given during the startup
type KDAuditInitFlags struct {
	ClientAuthType string
	LocalKubeCFG   string
}

// ParseFlags for parsing flags into KDAuditInitFlags
func ParseFlags() KDAuditInitFlags {
	clientAuthFlag := flag.String(
		"clientAuthType",
		"inCluster",
		"[localCFG/inCluster] - Default: localCFG",
	)
	localCFG := flag.String(
		"localCFG",
		"~/.kube/config",
		"Provide the path to local kubecfg",
	)

	flag.Parse()

	kdAuditInitFlags := KDAuditInitFlags{
		ClientAuthType: *clientAuthFlag,
		LocalKubeCFG:   *localCFG,
	}

	return kdAuditInitFlags
}
