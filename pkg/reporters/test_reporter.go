package reporters

import (
	"errors"
	"main/pkg/constants"
	"main/pkg/types"
)

type TestReporter struct {
	FailToSend   bool
	FailToInit   bool
	ReporterName string
}

func (r *TestReporter) Init() error {
	if r.FailToInit {
		return errors.New("fail to init")
	}

	return nil
}

func (r *TestReporter) Start() {

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
