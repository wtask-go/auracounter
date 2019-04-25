// +build integration

// This integration test uses connection with MySQL test database.
// As expected, all connection params are stored in environment variables and must exist before test will run.
// Check `deployments/config.test.env`
package mysql

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/wtask-go/auracounter/internal/config/env"
	"github.com/wtask-go/auracounter/internal/counter"
	"github.com/wtask-go/auracounter/internal/counter/datastore/mysql/model"
)

type test func(*testing.T)

// connectDB - creates independent connection to database used for tests.
func connectDB(dsn, tablePrefix string) *gorm.DB {
	iface, err := NewStorage(dsn, WithTablePrefix(tablePrefix))
	if err != nil {
		panic(errors.Wrap(err, "Failed to build storage interface"))
	}
	impl, ok := iface.(*storage)
	if !ok {
		panic(errors.New("Typecast failed, expected mysql.storage"))
	}
	return impl.db
}

func clearDB(checker *gorm.DB) error {
	return checker.DropTable(
		&model.Counter{},
	).Error
}

func TestMySQLDatastore(t *testing.T) {
	// package-level setup
	cfg, err := env.NewApplicationConfig("TEST_")
	if err != nil {
		panic(errors.Wrap(err, "Poor testing environment"))
	}

	checker := connectDB(cfg.CounterDB.DSN(), cfg.CounterDB.TablePrefix)
	defer func() {
		clearDB(checker)
		checker.Close()
	}()

	storage, err := NewStorage(cfg.CounterDB.DSN(), WithTablePrefix(cfg.CounterDB.TablePrefix))
	if err != nil {
		t.Errorf("Unable to create mysql storage: %v", err)
	}
	suite := []test{
		StorageEnsureLatest(checker, storage),
		RepositoryEnsureSettings(checker, storage),
	}
	for _, test := range suite {
		clearDB(checker)
		test(t)
	}
}

func StorageEnsureLatest(checker *gorm.DB, storage counter.Storage) test {
	return func(t *testing.T) {
		t.Log("Storage.EnsureLatest()")
		if checker.HasTable(&model.Counter{}) {
			t.Error("Unable to start test, model.Counter table exists in the database")
		}
		if err := storage.EnsureLatest(); err != nil {
			t.Errorf("EnsureLatest() for mysql failed: %v", err)
		}
		if !checker.HasTable(&model.Counter{}) {
			// expected previous err == nil, so table must exits
			t.Error("Unexpected EnsureLatest() behaviour, model.Counter does not exist in the database")
		}
	}
}

func RepositoryEnsureSettings(checker *gorm.DB, storage counter.Storage) test {
	return func(t *testing.T) {
		t.Log("Repository.EnsureSettings()")
		// empty database
		if err := storage.EnsureLatest(); err != nil {
			t.Errorf("Unable to ensure database has latest version: %v", err)
		}

		checker.Delete(&model.Counter{}) // should delete all records

		settings := counter.DefaultSettings()
		if err := storage.Repository().EnsureSettings(1, settings); err != nil {
			t.Errorf("Repository.EnsureSettings(): failed (empty database): %v", err)
		}

		c := &model.Counter{}
		if err := checker.First(c, 1).Error; err != nil {
			t.Errorf("Repository.EnsureSettings(): failed to load counter (1): %v", err)
		}

		loaded := &counter.Settings{
			StartFrom: c.Lower,
			Increment: c.Increment,
			Lower:     c.Lower,
			Upper:     c.Upper,
		}
		if *settings != *loaded {
			t.Errorf("[1] Loaded unexpected counter.Settings: %v", loaded)
		}

		// again with non empty database
		c.Value = 100
		c.Increment = 10
		c.Upper = 100
		checker.Save(c)

		settings = &counter.Settings{
			StartFrom: 1000,
			Increment: 100,
			Lower:     33,
			Upper:     10000,
		}

		if err := storage.Repository().EnsureSettings(1, settings); err != nil {
			t.Errorf("Repository.EnsureSettings(): method failed: %v", err)
		}

		if err := checker.First(c, 1).Error; err != nil {
			t.Errorf("Repository.EnsureSettings(): failed to load counter (2): %v", err)
		}
		loaded = &counter.Settings{
			StartFrom: c.Lower,
			Increment: c.Increment,
			Lower:     c.Lower,
			Upper:     c.Upper,
		}

		if *settings == *loaded {
			t.Errorf("[2] Loaded unexpected counter.Settings: %v", loaded)
		}

		if c.Value != 100 || loaded.Increment != 10 || loaded.Upper != 100 {
			t.Errorf("[3] Loaded unexpected counter.Settings: %v", loaded)
		}
	}
}
