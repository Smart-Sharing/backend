package entities

type MachineState = int

const MachineFree = MachineState(0)
const MachineInUse = MachineState(1)

type Machine struct {
	Id      string       `db:"id"`
	State   MachineState `db:"state"`
	Voltage int          `db:"voltage"`
}
