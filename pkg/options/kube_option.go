package options

import "time"

type KubeOption struct {
	KubeConfigPath          string
	InformersReSyncInterval time.Duration
}
