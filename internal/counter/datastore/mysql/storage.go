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
		db       *gorm.DB
		dsn      string
		prefix   string
		defaults *counter.Settings
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

func NewStorage(options ...storageOption) (counter.Storage, error) {
	s := (&storage{
		defaults: counter.DefaultSettings(),
	}).apply(options...)
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

func WithDSN(dsn string) storageOption {
	dsn = strings.TrimPrefix(dsn, "mysql://")
	if dsn == "" {
		return failedOption(errors.Errorf("invalid data source name %q", dsn))
	}
	return properOption(func(s *storage) {
		s.dsn = dsn
	})
}

func WithTablePrefix(prefix string) storageOption {
	return properOption(func(s *storage) {
		s.prefix = prefix
	})
}

func WithCounterSettings(defaults *counter.Settings) storageOption {
	if defaults == nil {
		return failedOption(errors.New("unable to use nil as default counter settings"))
	}
	return properOption(func(s *storage) {
		s.defaults = defaults
	})
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
