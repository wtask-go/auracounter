package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/wtask-go/auracounter/internal/api"
	"github.com/wtask-go/auracounter/internal/counter/datastore/mysql/model"
)

// GetNumber - return current counter value
func (s *storage) GetNumber() (int, error) {
	c, err := s.getCounter()
	// If settings was not set, getCounter() return empty model
	// If err == nil, errors.Wrapf() return nil
	// So method for counter without settings set will always return (0, nil)
	num := 0
	if c != nil {
		num = c.CurrentValue
	}
	return num, errors.Wrapf(err, "mysql.Repository: failed to get current counter (#%d) value", s.cid)
}

func (s *storage) IncrementNumber() (int, error) {
	c := &model.Counter{}
	table := s.db.NewScope(c).TableName()
	tx := s.db.Begin()
	if tx.Error != nil {
		return 0, errors.Wrap(tx.Error, "mysql.Repository: failed to lock before increment")
	}
	err := tx.Exec(
		"INSERT INTO "+table+" (counter_id, current_value, delta, max) VALUES (?, ?, ?, ?) "+
			"ON DUPLICATE KEY UPDATE current_value=IF(current_value+delta>max,0,current_value+delta)",
		s.cid,
		1, // not 0 !!! if insert, hence previous counter value was 0, but new is 1 for default
		1,
		api.MaxInt,
	).Error
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "mysql.Repository: failed to increment counter")
	}
	if err := tx.First(c, s.cid).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return 0, errors.New("mysql.Repository: failed to find incremented counter")
		}
		return 0, errors.Wrap(err, "mysql.Repository: failed to complete increment")
	}
	tx.Commit()

	return c.CurrentValue, nil
}

func (s *storage) SetSettings(delta, max int) error {
	// transaction is unnecessary
	table := s.db.NewScope(&model.Counter{}).TableName()
	err := s.db.Exec(
		"INSERT INTO "+table+" (counter_id, current_value, delta, max) VALUES (?, ?, ?, ?) "+
			"ON DUPLICATE KEY UPDATE delta=?, max=?;",
		// table,
		s.cid,
		0, // now 0 !!! if insert, hence it is primary initialization
		delta,
		max, // max int
		delta,
		max,
	).Error

	return errors.Wrap(err, "mysql.Repository: failed to set settings")
}

// getCounter - loads complete counter model.
// If record was not found, return empty model (CounterID==0)
func (s *storage) getCounter() (*model.Counter, error) {
	c := &model.Counter{}
	if err := s.db.First(c, s.cid).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrapf(err, "mysql.Repository: failed to get counter (#%d)", s.cid)
	}
	return c, nil
}
