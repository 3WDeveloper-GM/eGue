package adapter

type DBImplementation interface {
	GetDBAdmin() string
	GetDBPassword() string
	GetDBURL() string
}
