package pkg

type Opration int

const (
	Stop Opration = iota
	Start
)

type Sigle struct {
	Pname string
	Name  string
	S     Opration
}
