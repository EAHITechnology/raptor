package emysql

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type mClient struct {
	Master *gorm.DB
	Slave  []*gorm.DB
}

type mysqlManager struct {
	lock     sync.RWMutex
	dbclient map[string]*mClient
}

type Logger interface {
	Printf(f string, args ...interface{})
}

func checkMysqlConfig(m MConfigInfo) error {
	if m.Name == "" {
		return ErrNameInvalid
	}
	if m.Master.Ip == "" {
		return ErrMasterInvalid
	}
	if m.Master.Username == "" {
		return ErrUserNameInvalid
	}
	if m.Master.Password == "" {
		return ErrPassWordInvalid
	}
	if m.MaxIdleConns == 0 {
		m.MaxIdleConns = MAX_IDLE_CONNS
	}
	if m.MaxOpenConns == 0 {
		m.MaxOpenConns = MAX_OPEN_CONNS
	}
	if m.ConnMaxLifetime == 0 {
		m.ConnMaxLifetime = CONN_MAX_LiFE_TIME
	}
	if m.ConnMaxIdleTime == 0 {
		m.ConnMaxIdleTime = CONN_MAX_IDLE_TIME
	}
	return nil
}

func newMysql(m mInfo, l Logger) (*gorm.DB, error) {
	connProto := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%s&loc=%s&readTimeout=%s",
		m.Username, m.Password,
		m.Ip, m.Database,
		m.Charset, m.ParseTime,
		m.Loc, m.ReadTimeout,
	)

	gormConfig := &gorm.Config{}
	if m.LogMode {
		gormConfig.Logger = logger.New(l, logger.Config{
			LogLevel: logger.Info,
			Colorful: true,
		})
	}

	db, err := gorm.Open(mysql.Open(connProto), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(m.ConnMaxLifetime))
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(m.ConnMaxLifetime))
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)

	return db, nil
}

/*
NewMysql retuen mClient
*/
func NewMysql(m MConfigInfo, l Logger) (*mClient, error) {
	if err := checkMysqlConfig(m); err != nil {
		return nil, err
	}

	mc := &mClient{}

	mInfos := []mInfo{mInfo{
		Name:            m.Name,
		Username:        m.Master.Username,
		Ip:              m.Master.Ip,
		Password:        m.Master.Password,
		Database:        m.Database,
		Charset:         m.Charset,
		ParseTime:       m.ParseTime,
		Loc:             m.Loc,
		ReadTimeout:     m.ReadTimeout,
		MaxIdleConns:    m.MaxIdleConns,
		MaxOpenConns:    m.MaxOpenConns,
		ConnMaxLifetime: m.ConnMaxLifetime,
		LogMode:         m.LogMode,
	}}

	if len(m.Slaves) != 0 {
		for _, slave := range m.Slaves {
			mInfos = append(mInfos, mInfo{
				Name:            m.Name,
				Username:        slave.Username,
				Ip:              slave.Ip,
				Password:        slave.Password,
				Database:        m.Database,
				Charset:         m.Charset,
				ParseTime:       m.ParseTime,
				Loc:             m.Loc,
				ReadTimeout:     m.ReadTimeout,
				MaxIdleConns:    m.MaxIdleConns,
				MaxOpenConns:    m.MaxOpenConns,
				ConnMaxLifetime: m.ConnMaxLifetime,
				LogMode:         m.LogMode,
			})
		}
	}

	for idx, m := range mInfos {
		db, err := newMysql(m, l)
		if err != nil {
			return nil, err
		}

		if idx == 0 {
			mc.Master = db
			continue
		}
		mc.Slave = append(mc.Slave, db)
	}
	return mc, nil
}

func (m *mClient) GetMaster() *gorm.DB {
	return m.Master
}

func (m *mClient) GetSlave() *gorm.DB {
	return m.Slave[0]
}

func (m *mClient) Close() error {
	sqldb, err := m.Master.DB()
	if err != nil {
		return err
	}
	if err := sqldb.Close(); err != nil {
		return err
	}

	for _, s := range m.Slave {
		sqldb, err := s.DB()
		if err != nil {
			return err
		}
		if err := sqldb.Close(); err != nil {
			return err
		}
	}
	return err
}

// The NewMysqlSingle function provides a way
// to instantiate mysql through configuration.
// If using the NewMysqlSingle function, use
// the GetClient function provided in the package
// to manipulate the mClient instance.
func NewMysqlSingle(ms []MConfigInfo, l Logger) error {
	mManager.lock.Lock()
	defer mManager.lock.Unlock()

	mManager = &mysqlManager{
		dbclient: make(map[string]*mClient),
	}

	for _, m := range ms {
		mclient, err := NewMysql(m, l)
		if err != nil {
			return err
		}
		mManager.dbclient[m.Name] = mclient
	}
	return nil
}

func GetClient(dbname string) (*mClient, error) {
	mManager.lock.RLock()
	defer mManager.lock.RUnlock()

	db, ok := mManager.dbclient[dbname]
	if !ok {
		return nil, fmt.Errorf("db not init")
	}
	return db, nil
}

func CloseMysql() error {
	mManager.lock.Lock()
	defer mManager.lock.Unlock()

	for _, v := range mManager.dbclient {
		v.Close()
	}
	return nil
}
