package loggerEntity

// Definición del tipo LogLevel
type LogLevel uint8

const (
	INFO LogLevel = iota + 1
	WARNING
	ERROR
	SUCCESS
	PANIC
)

func Name(logLevel LogLevel) string {
	switch logLevel {
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case SUCCESS:
		return "SUCCESS"
	case PANIC:
		return "PANIC"
	default:
		return "UNKNOWN"
	}
}

// Devuelve el log como un string
func (l LogLevel) String() string {
	return Name(l)
}
