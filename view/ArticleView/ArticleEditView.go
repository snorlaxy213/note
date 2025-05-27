package ArticleView

import "note-gin/models"

type ArticleEditView struct {
	ID          int64   `form:"id" json:"id"`
	CreatedAt   string  `json:"created_at" form:"created_at"`
	UpdatedAt   string  `json:"updated_at" form:"updated_at"`
	Title       string  `json:"title" form:"title"`
	DirPath     []int64 `json:"dir_path" form:"dir_path"`
	FolderID    int64   `json:"folder_id" form:"folder_id"`
	FolderTitle string  `json:"folder_title" form:"folder_title"`
	MkValue     string  `form:"mkValue" json:"mkValue"`
}

// ToEditArticleDetail 将数据模型转换为文章详情视图(编辑)
// 从Article模型创建ArticleDetail视图对象
func ToEditArticleDetail(article models.Article) ArticleEditView {
	articleDetail := ArticleEditView{
		ID:       article.ID,
		Title:    article.Title,
		MkValue:  article.MkValue,
		FolderID: article.FolderID,
	}
	return articleDetail
}
