package mysql

import metrics "github.com/rcrowley/go-metrics"

type State struct {
	Sessions          metrics.Counter
	LoginAttempts     metrics.Counter
	LoginFailures     metrics.Counter
	DecoderErrors     metrics.Counter
	ProtocolErrors    metrics.Counter
	UpgradedToSsl     metrics.Counter
	AuthSwitchRequest metrics.Counter
	QueriesParsed     metrics.Counter
	QueriesParseError metrics.Counter
}

func newState() *State {
	return &State{
		Sessions:          metrics.NewCounter(),
		LoginAttempts:     metrics.NewCounter(),
		LoginFailures:     metrics.NewCounter(),
		ProtocolErrors:    metrics.NewCounter(),
		UpgradedToSsl:     metrics.NewCounter(),
		AuthSwitchRequest: metrics.NewCounter(),
		QueriesParsed:     metrics.NewCounter(),
		QueriesParseError: metrics.NewCounter(),
	}
}
