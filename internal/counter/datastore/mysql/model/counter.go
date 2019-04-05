package model

import "time"

// Counter - counter model
type Counter struct {
	CounterID    int       `gorm:"primary_key;column:counter_id"`
	CreatedAt    time.Time `gorm:"not null;default:current_timestamp"`
	UpdatedAt    time.Time `gorm:"not null;default:current_timestamp on update current_timestamp"`
	CurrentValue int       `gorm:"not null;default:'0';column:current_value"`
	Delta        int       `gorm:"not null;default:'1';column:delta"`
	Max          int       `gorm:"not null;default:'1000';column:max"`
}
