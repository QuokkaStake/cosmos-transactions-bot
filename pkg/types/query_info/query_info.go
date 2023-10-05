package query_info

import "time"

type QueryInfo struct {
	Success bool
	Time    time.Duration
	URL     string
	Node    string
}
