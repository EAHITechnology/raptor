package skip_list

type Iter interface {
	Seek(key []byte) error
	Get() ([]byte, error)
	Next() ([]byte, error)
	HasNext() bool
	Close()
}

type SkipList interface {
	Get(key []byte) ([]byte, bool, error)
	Put(key, value []byte) error
	Del(key []byte) error
	Len() int64
	GetIter() Iter
}

type SkipListConf struct {
	Typ string
}

func NewSkipList(config SkipListConf) (SkipList, error) {
	switch config.Typ {
	case "", "default":
		return NewDefaultSkipListImpl()
	default:
		return nil, nil
	}
}
