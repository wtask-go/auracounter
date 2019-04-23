package model

import "time"

// Counter - counter model
type Counter struct {
	CounterID int       `gorm:"primary_key;auto_increment:false;column:counter_id"`
	CreatedAt time.Time `gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"not null;default:current_timestamp on update current_timestamp"`
	Value     int       `gorm:"not null;default:'0';column:value"`
	Increment int       `gorm:"not null;default:'1';column:increment"`
	Lower     int       `gorm:"not null;default:'0';column:lower"`
	Upper     int       `gorm:"not null;default:'1';column:upper"`
}
