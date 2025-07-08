package entity

type Logger interface {
	WithField(key string, value interface{}) Logger
	Err(err error) Logger
	WithFields(fields map[string]interface{}) Logger

	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
	Fatal(msg string)
}
