package usersessioncookie

import "time"

type ReqBySess interface {
	StateToString() (string, time.Duration)
	StringToState(string)
}
