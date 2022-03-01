package elog

type LogConfig struct {
	LogTyp         string
	Prefix         string
	Dir            string
	AutoClearHours int
	Depth          int
	LogMaxByteNum  int64
	LogLevel       LogLevel
	Format         LogFormat
}
