package reporters

import (
	"errors"
	"main/pkg/constants"
	"main/pkg/types"
)

type TestReporter struct {
	FailToSend   bool
	ReporterName string
}

func (r *TestReporter) Init() {

}

func (r *TestReporter) Name() string {
	return r.ReporterName
}

func (r *TestReporter) Type() string {
	return constants.ReporterTypeTelegram
}

func (r *TestReporter) Send(report types.Report) error {
	if r.FailToSend {
		return errors.New("send error")
	}

	return nil
}
