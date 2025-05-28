package ArticleView

import (
	"note-gin/models"
	"strings"
	"time"
)

// ArticleInfo 文章信息视图
// 用于列表展示的文章简要信息
type ArticleInfo struct {
	ID        int64    `json:"id" form:"id"`                 // 文章ID
	Title     string   `json:"title" form:"title"`           // 文章标题
	UpdatedAt string   `json:"updated_at" form:"updated_at"` // 更新时间
	Tags      []string `json:"tags" form:"tags"`             // 文章标签列表
}

// ToArticleInfos 将文章模型数组转换为文章信息视图数组
// 批量转换Article模型到ArticleInfo视图
func ToArticleInfos(articles []models.Article) []ArticleInfo {
	ArticleInfos := make([]ArticleInfo, len(articles))

	for index := range articles {
		ArticleInfos[index].ID = articles[index].ID
		ArticleInfos[index].Title = articles[index].Title
		ArticleInfos[index].UpdatedAt = articles[index].UpdatedAt.Format("2006-01-02")
		ArticleInfos[index].Tags = strings.Split(articles[index].Tags, ",")
	}
	return ArticleInfos
}

// ToArticle 将文章信息视图转换为文章模型
// 从ArticleInfo视图创建Article模型对象
func ToArticle(articleInfo ArticleInfo) models.Article {
	article := models.Article{}
	article.ID = articleInfo.ID
	article.Title = articleInfo.Title
	article.Tags = strings.Join(articleInfo.Tags, ",")
	article.UpdatedAt, _ = time.Parse("2006-01-02", articleInfo.UpdatedAt)
	return article
}
