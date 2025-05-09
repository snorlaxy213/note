package FolderController

import (
	"github.com/gin-gonic/gin"
	"note-gin/pkg/RedisClient"
	"note-gin/pkg/utils"
	"note-gin/service/FolderService"
	"note-gin/view/common"
)

// GetCurrentNav 处理获取当前导航路径的请求
// @Summary 获取当前导航路径
// @Description 从 Redis 获取并返回当前导航路径，默认追加 "Home"
// @Tags Folder
// @Produce json
// @Success 200 {object} common.Response "成功响应，数据为导航路径列表"
// @Router /nav/current [get]
func GetCurrentNav(c *gin.Context) {
	nav := RedisClient.GetCurrentNav() // 从 Redis 获取导航信息
	nav = append(nav, "Home")          // 在导航末尾添加 "Home"
	c.JSON(200, common.OkWithData("", nav))
}

// GetSubFile 处理获取指定文件夹下子文件和子文件夹列表的请求 (分页)
// @Summary 获取子文件和文件夹列表
// @Description 根据文件夹标题和页码，获取其下的子文件和子文件夹，并返回导航信息
// @Tags Folder
// @Produce json
// @Param page path string true "页码"
// @Param title query string false "文件夹标题"
// @Success 200 {object} common.FileList "成功响应，包含文件列表、文件夹列表、导航和总数"
// @Router /folder/subfile/{page} [get]
func GetSubFile(c *gin.Context) {
	page := c.Param("page")                 // 从路径参数中获取页码
	folder_title := c.Query("title")        // 从查询参数中获取文件夹标题
	// 调用服务层获取子文件、子文件夹和总数
	folderInfos, articleInfos, total := FolderService.GetSubFile(folder_title, utils.StrToInt(page))
	//导航
	nav := FolderService.ChangeNav(page, folder_title) // 获取并更新导航信息
	resp := common.FileList{                           // 构建响应体
		Folders:  folderInfos,
		Articles: articleInfos,
		Nav:      nav,
		Total:    total,
	}
	c.JSON(200, resp)
}

// GetSubFolders 处理编辑器中懒加载获取子文件夹列表的请求
// @Summary 获取子文件夹列表 (懒加载)
// @Description 根据父文件夹ID，获取其直接子文件夹列表，用于编辑器目录树的懒加载
// @Tags Folder
// @Produce json
// @Param id path string true "父文件夹ID"
// @Success 200 {object} common.Response "成功响应，数据为子文件夹选择列表"
// @Router /folder/subfolders/{id} [get]
func GetSubFolders(c *gin.Context) {
	id := c.Param("id") // 从路径参数中获取父文件夹ID
	// 调用服务层获取子文件夹列表
	folderSelectList := FolderService.GetSubFolders(id)
	c.JSON(200, common.OkWithData("", folderSelectList))
}
