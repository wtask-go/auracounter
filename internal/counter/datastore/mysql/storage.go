package mysql

import (
	"strings"

	"github.com/wtask-go/auracounter/internal/counter/datastore/mysql/model"

	"github.com/pkg/errors"

	"github.com/wtask-go/auracounter/internal/counter"

	"github.com/jinzhu/gorm"

	// "go-sql-driver/mysql" initialization via gorm wrapper
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type storage struct {
	db     *gorm.DB
	dsn    string
	prefix string
	cid    int
}

func NewStorage(options ...storageOption) (counter.Storage, error) { 
	s := (&storage{}).apply(options...)
	if s.dsn == "" {
		return nil, errors.New("mysql.NewStorage: required DSN is missed")
	}
	if s.cid < 1 {
		return nil, errors.Errorf("mysql.NewStorage: invalid counter ID (%d)", s.cid)
	}

	if s.prefix != "" {
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return s.prefix + defaultTableName
		}
	}

	db, err := gorm.Open("mysql", s.dsn)
	if err != nil {
		return nil, errors.Wrap(err, "mysql.NewStorage: failed to open DB connection")
	}
	s.db = db

	s.db.
		LogMode(false).
		SingularTable(true)
	err = s.db.
		Set("gorm:table_options", "COLLATE='utf8_general_ci' ENGINE=InnoDB").
		AutoMigrate(&model.Counter{}).
		Error
	if err != nil {
		return nil, errors.Wrap(err, "mysql.NewStorage: failed to prepare counter storage")
	}
	return s, nil
}

type storageOption func(*storage)

func WithDSN(dsn string) storageOption {
	return func(s *storage) {
		s.dsn = strings.TrimPrefix(dsn, "mysql://")
	}
}

func WithTablePrefix(prefix string) storageOption {
	return func(s *storage) {
		s.prefix = prefix
	}
}

func WithCounterID(cid int) storageOption {
	return func(s *storage) {
		s.cid = cid
	}
}

func (s *storage) apply(options ...storageOption) *storage {
	if s == nil {
		return nil
	}
	for _, o := range options {
		if o != nil {
			o(s)
		}
	}
	return s
}

func (s *storage) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *storage) Repository() counter.Repository {
	if s == nil {
		return nil
	}
	return s
}
