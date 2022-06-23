package domain

type IDGenerator interface {
	New() (string, error)
}
