package ArticleController

import (
	"github.com/gin-gonic/gin"
	"note-gin/pkg/HttpCode"
	"note-gin/pkg/logging"
	"note-gin/service/ArticleService"
	"note-gin/view/ArticleView"
	"note-gin/view/common"
)

// SetTag 设置文章标签
// 接收文章信息并更新标签
func SetTag(c *gin.Context) {
	articleInfo := ArticleView.ArticleInfo{}
	_ = c.ShouldBind(&articleInfo)
	ArticleService.SetTag(articleInfo)
	return
}

// Update 更新文章
// 接收编辑后的文章信息并保存
func Update(c *gin.Context) {
	articleEditView := ArticleView.ArticleEditView{}
	err := c.ShouldBind(&articleEditView)
	if err != nil {
		logging.Error(err.Error())
	}
	ArticleService.Update(&articleEditView)
	c.JSON(HttpCode.SUCCESS, common.OkWithData("文章保存成功！", articleEditView))
}
