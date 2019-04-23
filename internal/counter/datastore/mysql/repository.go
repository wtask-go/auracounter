package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/wtask-go/auracounter/internal/counter"
	"github.com/wtask-go/auracounter/internal/counter/datastore/mysql/model"
)

func (s *storage) EnsureSettings(counterID int, defaults *counter.Settings) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.Wrapf(tx.Error, "mysql.EnsureSettings(#%d): failed to begin transaction.", counterID)
	}
	err := tx.Where(&model.Counter{CounterID: counterID}).
		Attrs(&model.Counter{
			Value:     defaults.StartFrom,
			Increment: defaults.Increment,
			Lower:     defaults.Lower,
			Upper:     defaults.Upper,
		}).FirstOrCreate(&model.Counter{}). // ignore result, but will check error
		Error
	if err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "mysql.EnsureSettings(#%d): failed", counterID)
	}

	return errors.Wrapf(tx.Commit().Error, "mysql.EnsureSettings(#%d): commit failed", counterID)
}

// Get - return current counter value
func (s *storage) GetValue(counterID int) (int, error) {
	c := &model.Counter{}
	if err := s.db.First(c, counterID).Error; err != nil {
		// same here if record not found
		return 0, errors.Wrapf(err, "mysql.GetValue(#%d): failed", counterID)
	}
	return c.Value, nil
}

// Increase - increase counter using previously stored settings without validating its consistency.
// Method returns highly likely calculated counter value.
// If counter/counter settings were not prepared before calling `mysql.Increase`, method will fail.
// See `mysql.EnsureSettings`.
func (s *storage) Increase(counterID int) (int, error) {
	// table := s.db.NewScope(c).TableName()
	tx := s.db.Begin()
	if tx.Error != nil {
		return 0, errors.Wrapf(tx.Error, "mysql.Increase(#%d): failed to begin transaction", counterID)
	}
	c := &model.Counter{}
	if err := tx.First(c, counterID).Error; err != nil {
		tx.Rollback()
		// same here if record not found
		return 0, errors.Wrapf(err, "mysql.Increase(#%d): failed to get counter", counterID)
	}
	result := c.Value + c.Increment
	if result > c.Upper {
		result = c.Lower
	}
	// TODO test if also model updated
	err := tx.Model(c).
		Update("value", gorm.Expr("IF(value+increment>upper,lower,value+increment)")).
		Error
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrapf(err, "mysql.Increase(#%d): failed", counterID)
	}
	err = errors.Wrapf(tx.Commit().Error, "mysql.Increase(#%d): commit failed", counterID)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (s *storage) SetSettings(counterID int, settings *counter.Settings) error {
	// we need transaction due to sequential select, insert/update queries
	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.Wrapf(tx.Error, "mysql.SetSettings(#%d): failed to begin transaction", counterID)
	}
	var (
		original = &model.Counter{}
		err      error
	)
	switch err = tx.First(original, counterID).Error; {
	default:
		tx.Rollback()
		return errors.Wrapf(err, "mysql.SetSettings(#%d): failed to get counter", counterID)
	case err == nil:
		// update
		err = tx.Model(original).
			Updates(&model.Counter{
				Increment: settings.Increment,
				Lower:     settings.Lower,
				Upper:     settings.Upper,
			}).Error
	case err == gorm.ErrRecordNotFound:
		// insert
		err = tx.Create(&model.Counter{
			CounterID: counterID,
			Value:     settings.StartFrom,
			Increment: settings.Increment,
			Lower:     settings.Lower,
			Upper:     settings.Upper,
		}).Error
	}

	if err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "mysql.SetSettings(#%d): failed to set %v", counterID, *settings)
	}

	return errors.Wrapf(tx.Commit().Error, "mysql.SetSettings(#%d): failed to commit changes", counterID)
}
