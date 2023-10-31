package models

import "gorm.io/gorm"

type Regra struct {
    gorm.Model
    Data string `gorm:"column:data"`
}