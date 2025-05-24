package models

import (
	"database/sql"
	"time"
)

type BaseModel struct {
	ID          int64        `gorm:"primary_key" form:"title" json:"title"`
	CreatedAt   time.Time    `form:"created_at" json:"created_at"`
	UpdatedAt   time.Time    `form:"updated_at" json:"updated_at"`
	Deleted     bool         `form:"deleted" json:"deleted"`
	DeletedTime sql.NullTime `form:"deleted_time" json:"deleted_time"`
}

// 增
func AddOne(obj interface{}) {

}
func AddMany(obj []interface{}) {

}

// 删
func DeleteOne(obj interface{}) {

}
func DeleteMany(attr map[string]string) {

}

// 查
func FindOne(obj interface{}) {

}
func FindByPage(obj interface{}, page int) {

}

// 改
func Update(attr map[string]string) {

}
func UpdateMany(obj interface{}, attr map[string]string) {

}
