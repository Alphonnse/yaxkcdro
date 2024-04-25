package config

type ConfigService interface {
	GetResourceURL() string
	GetPathDB() string
	GetPathIndex() string
	GetPathStopwords() string
}
