package probe

type CheckPointer interface {
	Check()
	Update(*Probe)
}
