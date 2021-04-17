package scs

type CheckPointer interface {
	Check()
	Update()
}
