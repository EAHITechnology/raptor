package emysql

type Account struct {
	Ip       string
	Username string
	Password string
}

type MConfigInfo struct {
	Name            string
	Master          Account
	Slaves          []Account
	Database        string
	Charset         string
	ParseTime       string
	Loc             string
	ReadTimeout     string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxIdleTime int
	ConnMaxLifetime int
	LogMode         bool
}

type mInfo struct {
	Name            string
	Username        string
	Ip              string
	Password        string
	Database        string
	Charset         string
	ParseTime       string
	Loc             string
	ReadTimeout     string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxIdleTime int
	ConnMaxLifetime int
	LogMode         bool
}
