package ArticleView

import (
	"database/sql"
	"note-gin/models"
	"strings"
)

// RubbishArticleInfo 回收文章信息视图
// 用于回收列表展示的文章简要信息
type RubbishArticleInfo struct {
	ID          int64    `json:"id" form:"id"`                     // 文章ID
	Title       string   `json:"title" form:"title"`               // 文章标题
	UpdatedAt   string   `json:"updated_at" form:"updated_at"`     // 更新时间
	Tags        []string `json:"tags" form:"tags"`                 // 文章标签列表
	Deleted     bool     `form:"deleted" json:"deleted"`           // 是否删除
	DeletedTime string   `form:"deleted_time" json:"deleted_time"` // 删除时间
}

// ToRubbishArticleInfos 将文章模型数组转换为文章信息视图数组
// 批量转换Article模型到ArticleInfo视图
func ToRubbishArticleInfos(articles []models.Article) []RubbishArticleInfo {
	ArticleInfos := make([]RubbishArticleInfo, len(articles))

	for index := range articles {
		ArticleInfos[index].ID = articles[index].ID
		ArticleInfos[index].Title = articles[index].Title
		ArticleInfos[index].UpdatedAt = articles[index].UpdatedAt.Format("2006-1-2")
		ArticleInfos[index].Tags = strings.Split(articles[index].Tags, ",")
		ArticleInfos[index].Deleted = articles[index].Deleted
		ArticleInfos[index].DeletedTime = formatNullTime(articles[index].DeletedTime)
	}

	return ArticleInfos
}

func formatNullTime(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format("2006-01-02")
	}
	return ""
}
