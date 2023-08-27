package mysql

type Session struct {
	state       State
	expectedSeq uint8
}

type State uint8

const (
	Init State = iota
	ChallengeReq
	ChallengeResp41
	ChallengeResp320
	SslPt
	AuthSwitchReq
	AuthSwitchReqOld
	AuthSwitchResp
	AuthSwitchMore
	ReqResp
	Req
	Resync
	NotHandled
	Error
)

func (s Session) setState(state State) {

}

func (s Session) getState() State {
	return s.state
}

func (s Session) getExpectedSeq() uint8 {
	return s.expectedSeq
}

func (s Session) setExpectedSeq(seq uint8) {
}
