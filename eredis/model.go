package eredis

type RedisInfo struct {
	RedisName      string
	Addr           string
	MaxIdle        int
	MaxActive      int
	IdleTimeout    int64
	ReadTimeout    int64
	WriteTimeout   int64
	ConnectTimeout int64
	Password       string
	Wait           bool
	Database       int
}
