package nodes_manager

import (
	"main/pkg/types"
	"sync"
)

type ReportQueue struct {
	Data  []types.Report
	Size  int
	Mutes sync.Mutex
}

func NewReportQueue(size int) ReportQueue {
	return ReportQueue{Data: make([]types.Report, 0), Size: size}
}

func (q *ReportQueue) Add(report types.Report) {
	q.Mutes.Lock()

	if len(q.Data) >= q.Size {
		_, q.Data = q.Data[0], q.Data[1:]
	}

	q.Data = append(q.Data, report)
	q.Mutes.Unlock()
}

func (q *ReportQueue) Has(msg types.Report) bool {
	for _, elem := range q.Data {
		if elem.Reportable.GetHash() == msg.Reportable.GetHash() {
			return true
		}
	}

	return false
}
