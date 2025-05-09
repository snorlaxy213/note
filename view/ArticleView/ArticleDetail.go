package ArticleView

import "note-gin/models"

// ArticleDetail 文章详情视图
// 包含文章的基本信息和Markdown内容
type ArticleDetail struct {
	ID      int64  `form:"id" json:"id"`       // 文章ID
	Title   string `form:"title" json:"title"` // 文章标题
	MkValue string `form:"mkValue" json:"mkValue"` // Markdown格式的文章内容
}

// ToArticleDetail 将数据模型转换为文章详情视图
// 从Article模型创建ArticleDetail视图对象
func ToArticleDetail(article models.Article) ArticleDetail {
	articleDetail := ArticleDetail{
		ID:      article.ID,
		Title:   article.Title,
		MkValue: article.MkValue,
	}
	return articleDetail
}
