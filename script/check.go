package script

type CheckPointer interface {
	Check()
	Update()
}
