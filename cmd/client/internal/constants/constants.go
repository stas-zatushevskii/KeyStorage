package constants

type Mode int

const (
	ModeCreate Mode = iota
	ModeList
)

type ObjType int

const (
	Text ObjType = iota
	Account
	File
	Bank
)

func (m Mode) String() string {
	switch m {
	case ModeCreate:
		return "CREATE"
	case ModeList:
		return "LIST"
	default:
		return "UNKNOWN"
	}
}

func (t ObjType) String() string {
	switch t {
	case Text:
		return "text"
	case Account:
		return "account"
	case File:
		return "file"
	case Bank:
		return "bank"
	default:
		return "unknown"
	}
}
