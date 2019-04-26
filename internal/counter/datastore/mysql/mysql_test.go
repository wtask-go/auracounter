// +build integration

// This integration test uses connection with MySQL test database.
// As expected, all connection params are stored in environment variables and must exist before test will run.
// Check `deployments/config.test.env`
package mysql

import (
	"fmt"
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

// clearDB - drops all known tables in the database
func clearDB(checker *gorm.DB) error {
	return checker.DropTable(
		&model.Counter{},
	).Error
}

func TestDatastore(t *testing.T) {
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

	clearDB(checker)

	storage, err := NewStorage(cfg.CounterDB.DSN(), WithTablePrefix(cfg.CounterDB.TablePrefix))
	if err != nil {
		t.Errorf("Unable to create mysql storage: %v", err)
	}
	for i, test := range DatastoreSuite(checker, storage) {
		t.Run(fmt.Sprintf("test #%d", i+1), test)
	}
}

func DatastoreSuite(checker *gorm.DB, storage counter.Storage) []test {
	return []test{
		// run StorageEnsureLatest first,
		// when the test was successful it should guarantee appropriate db structure
		StorageEnsureLatest(checker, storage),
		RepositoryEnsureSettings(checker, storage.Repository()),
		RepositoryGetValue(checker, storage.Repository()),
		RepositoryIncrease(checker, storage.Repository()),
		RepositorySetSettings(checker, storage.Repository()),
	}
}

func StorageEnsureLatest(checker *gorm.DB, storage counter.Storage) test {
	return func(t *testing.T) {
		t.Log("TEST: Storage.(mysql).EnsureLatest()")
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

func RepositoryEnsureSettings(checker *gorm.DB, repository counter.Repository) test {
	return func(t *testing.T) {
		t.Log("TEST: Repository.(mysql).EnsureSettings()")
		// empty database
		// if err := storage.EnsureLatest(); err != nil {
		// 	t.Errorf("Unable to ensure the database has latest version: %v", err)
		// }

		checker.Delete(&model.Counter{}) // should delete all records
		settings := counter.DefaultSettings()
		t.Logf("Case: empty database and default counter.Settings %+v", settings)

		if err := repository.EnsureSettings(1, settings); err != nil {
			t.Errorf("Method failed (empty database): %v", err)
		}

		c := &model.Counter{}
		if err := checker.First(c, 1).Error; err != nil {
			t.Errorf("Failed to load counter.Settings: %v", err)
		}
		loaded := &counter.Settings{
			StartFrom: c.Lower,
			Increment: c.Increment,
			Lower:     c.Lower,
			Upper:     c.Upper,
		}
		t.Logf("Loaded saved counter.Settings: %+v", loaded)

		if *settings != *loaded {
			t.Errorf("Loaded unexpected counter.Settings: %v", loaded)
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
		t.Logf("Case: non-empty database and custom counter.Settings %+v", settings)

		if err := repository.EnsureSettings(1, settings); err != nil {
			t.Errorf("Method failed (non-empty database): %v", err)
		}

		if err := checker.First(c, 1).Error; err != nil {
			t.Errorf("Failed to load counter.Settings: %v", err)
		}
		loaded = &counter.Settings{
			StartFrom: c.Lower,
			Increment: c.Increment,
			Lower:     c.Lower,
			Upper:     c.Upper,
		}
		t.Logf("Loaded saved counter.Settings: %+v", loaded)

		if *settings == *loaded {
			t.Errorf("Loaded unexpected counter.Settings: %v", loaded)
		}

		if c.Value != 100 || loaded.Increment != 10 || loaded.Upper != 100 {
			t.Errorf("Method violates counter.Settings integrity: %v", loaded)
		}
	}
}

func RepositoryGetValue(checker *gorm.DB, repository counter.Repository) test {
	return func(t *testing.T) {
		t.Log("TEST: Repository.(mysql).GetValue()")

		checker.Delete(&model.Counter{}) // should delete all records
		t.Logf("Case: empty database")

		_, err := repository.GetValue(1)
		if err == nil {
			t.Error("Expected error for non-existed counter, got nothing")
		}
		t.Logf("Got expected error: %v", err)

		c := &model.Counter{
			CounterID: 1,
			Value:     100,
			Increment: 10,
			Lower:     0,
			Upper:     1000,
		}
		checker.Save(c)
		t.Logf("Case: non-empty database")

		v, err := repository.GetValue(1)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if v != c.Value {
			t.Errorf("Expected %d, got %d", c.Value, v)
		}
	}
}

func RepositoryIncrease(checker *gorm.DB, repository counter.Repository) test {
	return func(t *testing.T) {
		t.Log("TEST: Repository.(mysql).Increase()")

		checker.Delete(&model.Counter{}) // should delete all records
		t.Logf("Case: empty database")

		_, err := repository.Increase(1)
		if err == nil {
			t.Error("Expected error for non-existed counter, got nothing")
		}
		t.Logf("Got expected error: %v", err)

		c := &model.Counter{
			CounterID: 1,
			Value:     990,
			Increment: 10,
			Lower:     0,
			Upper:     1000,
		}
		checker.Save(c)

		t.Logf("Case: non-empty database")
		v, err := repository.Increase(1)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if v != c.Value+c.Increment {
			t.Errorf("Expected %d, got %d", c.Value+c.Increment, v)
		}

		t.Logf("Case: reaching the upper limit")
		v, err = repository.Increase(1)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if v != c.Lower {
			t.Errorf("Expected %d, got %d", c.Lower, v)
		}
	}
}

func RepositorySetSettings(checker *gorm.DB, repository counter.Repository) test {
	return func(t *testing.T) {
		t.Log("TEST: Repository.(mysql).SetSettings()")

		checker.Delete(&model.Counter{}) // should delete all records
		t.Logf("Case: empty database")

		initial := &counter.Settings{
			StartFrom: 100,
			Increment: 10,
			Lower:     0,
			Upper:     1000,
		}
		if err := repository.SetSettings(1, initial); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		c := &model.Counter{}
		if err := checker.First(c, 1).Error; err != nil {
			t.Errorf("Failed to load counter.Settings: %v", err)
		}
		loaded := &counter.Settings{
			StartFrom: c.Value,
			Increment: c.Increment,
			Lower:     c.Lower,
			Upper:     c.Upper,
		}
		t.Logf("Loaded saved counter.Settings: %+v", loaded)

		if *initial != *loaded {
			t.Errorf("Loaded unexpected counter.Settings: %v", loaded)
		}

		t.Logf("Case: empty database")

		final := &counter.Settings{
			// StartFrom must not affect existing counter
			StartFrom: 500,
			Increment: 100,
			Lower:     100,
			Upper:     10000,
		}
		if err := repository.SetSettings(1, final); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if err := checker.First(c, 1).Error; err != nil {
			t.Errorf("Failed to load counter.Settings: %v", err)
		}
		if c.Value != initial.StartFrom {
			t.Errorf("Unexpected model.Counter.Value (%d) after settings were set", c.Value)
		}

		loaded = &counter.Settings{
			StartFrom: final.StartFrom,
			Increment: c.Increment,
			Lower:     c.Lower,
			Upper:     c.Upper,
		}
		if *final != *loaded {
			t.Errorf("Loaded unexpected counter.Settings: %v", loaded)
		}
	}
}
