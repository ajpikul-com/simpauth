package uwho

import ()

type UserStatus int64

const (
	UNKNOWN UserStatus = iota
	KNOWN
	EXPIRED
	AUTHORIZED
	SPOKEN // TODO: not implemented, so one module can hijack whole process
	LOGGEDOUT
)

func (u *UserStatus) StatusStr() string {
	switch *u {
	case UNKNOWN:
		return "UNKNOWN"
	case KNOWN:
		return "KNOWN"
	case EXPIRED:
		return "EXPIRED"
	case AUTHORIZED:
		return "AUTHORIZED"
	case SPOKEN:
		return "SPOKEN"
	case LOGGEDOUT:
		return "LOGGEDOUT"
	}
	return ""
}

func NewUserStatus() *UserStatus {
	u := UNKNOWN
	return &u
}
func (u *UserStatus) ReconcileStatus(status UserStatus) {
	if status > *u {
		*u = status
	}
}

func (u *UserStatus) SetStatus(status UserStatus) {
	*u = status
}

func (u *UserStatus) IsStatus(status UserStatus) bool {
	return status == *u
}

func (u *UserStatus) GetStatus() UserStatus {
	return *u
}
