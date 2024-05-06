package config

type ConfigService interface {
	ServerAddress() string
	GetResourceURL() string
	GetPathDB() string
	GetPathIndex() string
	GetPathStopwords() string
}
