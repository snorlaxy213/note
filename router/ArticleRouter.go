package router

import "note-gin/controller/ArticleController"

// ArticleRouter 文章相关路由配置
// 设置所有与文章操作相关的API路由
func ArticleRouter(base string) {
	r := Router.Group("/" + base)

	r.GET("/download/:id", ArticleController.ArticleDownLoad)      // 下载文章
	r.GET("/many/:page", ArticleController.GetArticleByPage)       // 分页获取文章
	r.GET("/get/:id", ArticleController.GetArticleDetail)          // 获取文章详情
	r.GET("/clear_rubbish", ArticleController.ClearRubbish)        // 清空回收站
	r.GET("/delete", ArticleController.Delete)                     // 删除文章
	r.GET("/delete/many", ArticleController.DeleteMany)            // 批量删除文章
	r.GET("/rubbish", ArticleController.GetRubbishArticles)        // 获取回收站文章
	r.GET("/recover", ArticleController.ArticleRecover)            // 恢复文章
	r.GET("/temp_get", ArticleController.TempArticleEditGet)       // 获取临时编辑内容
	r.GET("/temp_delete", ArticleController.TempArticleEditDelete) // 删除临时编辑内容
	r.POST("/temp_save", ArticleController.TempArticleEditSave)    // 保存临时编辑内容
	r.POST("/add", ArticleController.Add)                          // 添加文章
	r.GET("/edit/:id", ArticleController.EditArticleDetail)        // 获取编辑文章详情
	r.POST("/update", ArticleController.Update)                    // 更新文章
	r.POST("/set_tag", ArticleController.SetTag)                   // 设置文章标签
	r.POST("/upload_md", ArticleController.UploadArticle)          // 上传Markdown文件

}
