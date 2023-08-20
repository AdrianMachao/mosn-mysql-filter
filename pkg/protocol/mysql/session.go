package mysql

type Session struct {
	state State
}

type State int

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
	return Init
}

func (s Session) getExpectedSeq() uint8 {
	return 0
}

func (s Session) setExpectedSeq(seq uint8) {
}
