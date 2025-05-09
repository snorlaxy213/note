package ArticleController

import (
	"github.com/gin-gonic/gin"
	"note-gin/pkg/HttpCode"
	"note-gin/service/ArticleService"
	"note-gin/view/common"
)

// Delete 删除文章
// 将文章移动到回收站
func Delete(c *gin.Context) {
	ID := ArticleService.Delete(c.Query("id"))
	c.JSON(HttpCode.SUCCESS, common.OkWithData("成功移动到垃圾箱 定期清除哟！", ID))
}

// ClearRubbish 清空回收站
// 永久删除回收站中的所有文章
func ClearRubbish(c *gin.Context) {
	ArticleService.ClearRubbish()
	c.JSON(HttpCode.SUCCESS, common.OkWithMsg("清空成功！"))
}
