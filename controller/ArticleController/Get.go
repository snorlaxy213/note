package ArticleController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"note-gin/pkg/HttpCode"
	"note-gin/pkg/utils"
	"note-gin/service/ArticleService"
	"note-gin/view/ArticleView"
	"note-gin/view/common"
	"strings"
)

// DeleteMany 接收文章ID数组并执行批量删除操作
func DeleteMany(c *gin.Context) {
	ArticleService.DeleteMany(c.QueryArray("items[]"))
	c.JSON(200, common.OkWithMsg("删除成功!"))
}

// GetArticleByPage 分页获取文章列表
// 根据页码返回对应页的文章列表
func GetArticleByPage(c *gin.Context) {
	page := utils.StrToInt(c.Param("page"))
	articleInfos, total := ArticleService.GetArticleByPage(page)
	c.JSON(200, common.DataList{
		Items: articleInfos,
		Total: int64(total),
	})
}

// GetArticleDetail 获取文章详情
// 根据文章ID返回文章的详细内容
func GetArticleDetail(c *gin.Context) {
	articleDetail := ArticleService.GetArticleDetail(c.Param("id"))
	c.JSON(HttpCode.SUCCESS, articleDetail)
}

// GetRubbishArticles 获取回收站中的文章
// 返回已删除但未彻底清除的文章列表
func GetRubbishArticles(c *gin.Context) {
	respDataList := ArticleService.GetRubbishArticles()
	c.JSON(HttpCode.SUCCESS, respDataList)
}

// ArticleRecover 恢复回收站中的文章
// 将已删除的文章恢复到正常状态
func ArticleRecover(c *gin.Context) {
	err := ArticleService.ArticleRecover(c.Query("id"))
	if err != nil {
		c.JSON(HttpCode.ERROR_RECOVER, common.ErrorWithMsg(HttpCode.HttpMsg[HttpCode.ERROR_RECOVER]))
	} else {
		c.JSON(200, common.OkWithMsg("恢复成功！"))
	}
}

// TempArticleEditSave 保存文章编辑临时草稿
// 将当前编辑状态的文章保存到临时存储中
func TempArticleEditSave(c *gin.Context) {
	articleEditView := ArticleView.ArticleEditView{}
	err := c.ShouldBind(&articleEditView)
	utils.ErrReport(err)
	flag := ArticleService.TempArticleEditSave(articleEditView)
	if flag {
		c.JSON(HttpCode.SUCCESS, common.OkWithMsg("文章暂存草稿箱,1天后失效！"))
	} else {
		c.JSON(HttpCode.ERROR_TEMP_SAVE, common.OkWithMsg(HttpCode.HttpMsg[HttpCode.ERROR_TEMP_SAVE]))
	}
}

// TempArticleEditGet 获取临时保存的文章草稿
// 从临时存储中获取上次编辑的文章内容
func TempArticleEditGet(c *gin.Context) {
	if articleEditView, ok := ArticleService.TempArticleEditGet(); ok {
		c.JSON(200, common.OkWithData("", articleEditView))
	} else {
		c.JSON(200, common.OkWithData("获取失败", articleEditView))
	}
}

// TempArticleEditDelete 删除临时保存的文章草稿
// 清除临时存储中的文章编辑内容
func TempArticleEditDelete(c *gin.Context) {
	flag := ArticleService.TempArticleEditDelete()
	if flag == 1 {
		c.JSON(200, common.OkWithMsg("清除成功!"))
	} else {
		c.JSON(200, common.OkWithMsg("清除失败:"+string(flag)))
	}
}

// ArticleDownLoad 下载文章
// 将文章内容作为文件提供下载
func ArticleDownLoad(c *gin.Context) {
	filename, MkValue := ArticleService.ArticleDownLoad(c.Param("id"))
	//文件命名
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	io.Copy(c.Writer, strings.NewReader(MkValue))
}

// EditArticleDetail 获取文章编辑信息
// 根据文章ID获取文章的编辑视图
func EditArticleDetail(c *gin.Context) {
	articleEditView := ArticleView.ArticleEditView{}
	articleEditView.ID = int64(utils.StrToInt(c.Param("id")))
	ArticleService.Edit(&articleEditView)
	c.JSON(HttpCode.SUCCESS, common.OkWithData("", articleEditView))
}
