package constants

type LoggerColors struct {
	Reset	  string
	Red		  string
	Green	  string
	Yellow	  string
	Blue	  string
	Magenta	  string
	Cyan	  string
	White	  string
	Black	  string
}

func NewLoggerColors() *LoggerColors {
	return &LoggerColors{
		Reset:	 	"\033[0m",
		Red:	  	"\033[31m",
		Green:	 	"\033[32m",
		Yellow:	  	"\033[33m",
		Blue:	  	"\033[34m",
		Magenta:	"\033[35m",
		Cyan:		"\033[36m",
		White:		"\033[37m",
		Black:		"\033[30m",
	}
}