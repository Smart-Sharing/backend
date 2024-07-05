package entities

import "time"

type SessionState = int

const (
	SessionActive = SessionState(0)
	SessionPause  = SessionState(1)
	SessionStoped = SessionState(2)
)

type Session struct {
	Id             int          `db:"id"`
	State          SessionState `db:"state"`
	MachineId      string       `db:"machine_id"`
	WorkerId       int          `db:"worker_id"`
	DatetimeStart  time.Time    `db:"datetime_start"`
	DatetimeFinish time.Time    `db:"datetime_finish"`
}
