package models

import (
	"database/sql"
	"errors"
	"time"
)

type Article struct {
	BaseModel
	Title    string `form:"title" json:"title"`
	FolderID int64  `form:"folder_id" json:"folder_id"`
	MkValue  string `form:"mk_value" json:"mk_value"`
	Tags     string `form:"tags" json:"tags"`
}

// Find
func (this Article) Count() (total int) {
	db.Table("article").Where("deleted=?", 0).Count(&total)
	return total
}

func (this Article) GetMany(page int) (articles []Article) {

	db.Table("article").Where("deleted=?", 0).Order("updated_at desc").Offset((page - 1) * 10).Limit(10).
		Select([]string{"id", "updated_at", "title"}).Find(&articles)
	return
}

func (this *Article) GetArticleInfo() {
	db.Where("id=?", this.ID).First(&this)
}
func (this *Article) GetArticleInfoByTitle() {
	db.Where("title=?", this.Title).First(&this)
}

func (this Article) GetDeletedArticle() (articles []Article) {
	db.Find(&articles, "deleted=?", 1)
	return
}

// Create
func (this *Article) Add() {
	db.Create(this)
}

// Update Or Create
func (this *Article) Update() {
	if this.ID != 0 {
		db.Save(&this)
	} else {
		db.Create(this)
	}
}
func (this *Article) SetTag() {
	db.Model(&this).Update("tags", this.Tags)
}

// 回收到垃圾箱
func (this *Article) Delete() {
	db.Model(&this).Update("deleted", true).Update("deleted_time", sql.NullTime{Time: time.Now(), Valid: true})
}

// 真实批量删除【后台】
func (this Article) DeleteMany(ids []string) {
	db.Table("article").Where("id in (?)", ids).Delete(&this)
}

// 清空垃圾箱
func (this Article) ClearRubbish() {
	db.Where("deleted <> 0").Delete(&Article{})
}

// PermanentDeleteOne 永久删除回收站中的单篇文章
func (this *Article) PermanentDeleteOne() {
	db.Where("id = ? AND deleted = ?", this.ID, 1).Delete(&this)
}

// 垃圾箱恢复
func (this Article) Recover() error {
	hasFolder := 0
	db.First(&this)
	db.Table("folder").Where("id=?", this.FolderID).Count(&hasFolder)

	if hasFolder != 0 || this.FolderID == 0 {
		db.Table("article").Where("id=?", this.ID).Update("deleted", 0).Update("deleted_time", nil)
		return nil
	} else {
		return errors.New("父目录不存在！恢复失败")
	}
}

func (this Article) IsExist() bool {
	c := 0
	db.Table("article").Where("title=?", this.Title).Count(&c)
	return c > 0
}
