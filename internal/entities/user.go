package entities

type UserJob = string

const Worker = UserJob("worker")
const Admin = UserJob("admin")

type User struct {
	Id          int     `db:"id"`
	Name        string  `db:"name"`
	PhoneNumber string  `db:"phone_humber"`
	JobPosition UserJob `db:"job_position"`
}
