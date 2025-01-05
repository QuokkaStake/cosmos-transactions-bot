package yaml_config

import "gopkg.in/guregu/null.v4"

type MetricsConfig struct {
	Enabled    null.Bool `default:"true"  yaml:"enabled"`
	ListenAddr string    `default:":9580" yaml:"listen-addr"`
}
