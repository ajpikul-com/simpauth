package uwho

import (
	"github.com/google/uuid"
)

type UserStatus int64

const (
	UNKNOWN UserStatus = iota
	KNOWN
	EXPIRED
	AUTHORIZED
	SPOKEN // TODO: not implemented, so one module can hijack whole process
	LOGGEDOUT
)

type userinfo struct {
	Status  UserStatus
	session uuid.NullUUID
	Data    *[]interface{}
}

func (u *userinfo) ReconcileStatus(status UserStatus) {
	if status > u.Status {
		u.Status = status
	}
}

func (u *userinfo) Append(data interface{}) {
	*u.Data = append(*u.Data, data)
}

func (u *userinfo) SetStatus(status UserStatus) {
	u.Status = status
}

func (u *userinfo) IsStatus(status UserStatus) bool {
	return status == u.Status
}

func (u *userinfo) GetStatus() UserStatus {
	return u.Status
}

func (u *userinfo) NewSession() {
	u.session.Valid = true
	u.session.UUID = uuid.New()
}

func (u *userinfo) SetSessionPending(uuid uuid.UUID) {
	u.SetSession(uuid)
}

func (u *userinfo) SetSession(uuid uuid.UUID) {
	u.session.Valid = true
	u.session.UUID = uuid
}

func (u *userinfo) GetSession() uuid.NullUUID {
	return u.session
}

func newUserinfo() *userinfo {
	return &userinfo{
		Status:  UNKNOWN,
		Data:    new([]interface{}),
		session: uuid.NullUUID{Valid: false},
	}
}
