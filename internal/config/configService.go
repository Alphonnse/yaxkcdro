package config

type ConfigService interface {
	GetResourceURL() string
	GetPathDB() string
	GetPathStopwords() string
}
