package toml_config

import "gopkg.in/guregu/null.v4"

type MetricsConfig struct {
	Enabled    null.Bool `default:"true"  toml:"enabled"`
	ListenAddr string    `default:":9580" toml:"listen-addr"`
}
