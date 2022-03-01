package emysql

import (
	"errors"
)

var (
	ErrNameInvalid     error = errors.New("name invalid")
	ErrUserNameInvalid error = errors.New("username invalid")
	ErrPassWordInvalid error = errors.New("password invalid")
	ErrMasterInvalid   error = errors.New("master invalid")
	ErrIPInvalid       error = errors.New("IP invalid")
	ErrDBNameInvalid   error = errors.New("DB name invalid")
	ErrDBLoggerNil     error = errors.New("DB logger nil")
)

const (
	MAX_IDLE_CONNS     = 2
	MAX_OPEN_CONNS     = 4
	CONN_MAX_LiFE_TIME = 300
	CONN_MAX_IDLE_TIME = 300
)
