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

type (
	storage struct {
		db     *gorm.DB
		dsn    string
		prefix string
	}

	storageOption func() (func(*storage), error)
)

// failedOption - helper to expose error from option builder
func failedOption(err error) storageOption {
	return func() (func(*storage), error) {
		return nil, err
	}
}

// properOption - helper to expose setter from option builder
func properOption(setter func(*storage)) storageOption {
	return func() (func(*storage), error) {
		return setter, nil
	}
}

// setup - set storage options
func (s *storage) setup(options ...storageOption) error {
	if s == nil {
		return nil
	}
	for _, option := range options {
		if option == nil {
			continue
		}
		setter, err := option()
		if err != nil {
			return err
		}
		if setter != nil {
			setter(s)
		}
	}
	return nil
}

// WithTablePrefix - common custom prefix for underlying database table(s).
// By default, storage does not use prefix for  table names.
func WithTablePrefix(prefix string) storageOption {
	return properOption(func(s *storage) {
		s.prefix = prefix
	})
}

// NewStorage - implements counter.Storage interface to store cyclic incremental counter with mysql.
// If storage was created without errors, you may use it after has ensured it has latest version
// and is up-to-date, see `EnsureLatest()` method.
func NewStorage(dsn string, options ...storageOption) (counter.Storage, error) {
	s := (&storage{
		dsn: strings.TrimPrefix(dsn, "mysql://"),
	})

	if s.dsn == "" {
		return nil, errors.New("mysql.NewStorage: required DSN is missed")
	}

	if err := s.setup(options...); err != nil {
		return nil, errors.Wrap(err, "mysql.NewStorage: option error")
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

	s.db.LogMode(false).
		SingularTable(true)

	return s, nil
}

// EnsureLatest - make sure underlying database has latest version and is up-to-date to store counter.
func (s *storage) EnsureLatest() error {
	err := s.db.
		Set("gorm:table_options", "COLLATE='utf8_general_ci' ENGINE=InnoDB").
		AutoMigrate(&model.Counter{}).
		Error
	return errors.Wrap(err, "mysql.EnsureLatest: failed")
}

// Close - close and free all used connections and resources.
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
