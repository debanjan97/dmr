package types

type Notification struct {
	WorkerId string
	Status   bool
}

type Instruction struct {
	Operation string
	Operand   []interface{}
	Status    string
}
