package FolderService

import (
	"note-gin/models"
	"note-gin/pkg/RedisClient"
	"note-gin/pkg/utils"
	"note-gin/view/ArticleView"
	"note-gin/view/FolderView"
	"time"
)

// GetFolderPath 调用 models 层的方法，递归获取指定 FolderID 的完整目录路径
// FolderID: 当前要查找其父路径的文件夹ID
// DirPath: 指向一个 int64 切片的指针，用于存储路径上的文件夹ID
func GetFolderPath(FolderID int64, DirPath *[]int64) {
	models.Folder{}.GetFolderPath(FolderID, DirPath)
}

// GetFolderByTitle 根据文件夹标题获取文件夹信息
// folder_title: 要查询的文件夹标题
// 返回值: 转换后的 FolderView.FolderInfo 对象
func GetFolderByTitle(folder_title string) FolderView.FolderInfo {
	folderInfo := FolderView.ToFolderInfo(models.Folder{Title: folder_title})
	return folderInfo
}

// GetSubFile 获取指定文件夹下的子文件和子文件夹 (分页)
// folder_title: 父文件夹的标题
// page: 请求的页码
// 返回值: folderInfos (子文件夹信息列表), articleInfos (子文章信息列表), total (子文件和子文件夹总数)
func GetSubFile(folderTitle string, page int) ([]FolderView.FolderInfo, []ArticleView.ArticleInfo, int) {
	folder := models.Folder{}

	folder.Title = folderTitle
	folder.GetFolderByTitle() // 根据标题获取文件夹的详细信息 (主要是ID)

	// 根据页码查找这个目录下的全部文件和文件夹，并获取总数
	folders, articles, total := folder.GetSubFile(page)
	articleInfos := ArticleView.ToArticleInfos(articles) // 将文章模型列表转换为视图模型列表
	folderInfos := FolderView.ToFolderInfos(folders)     // 将文件夹模型列表转换为视图模型列表

	return folderInfos, articleInfos, total

}

// ChangeNav 根据当前页码和文件夹标题更新导航路径
// page: 当前请求的页码字符串
// folder_title: 当前文件夹的标题
// 返回值: 更新后的导航路径字符串切片
// 注意: 只有当 page 为 "1" 时 (通常意味着切换到新的文件夹)，才会更新 Redis 中的导航缓存
func ChangeNav(page string, folder_title string) []string {
	var nav []string // 如果是访问新文件夹 (AccessFolder)，则需要加载导航；如果是页码跳转，则不需要加载，前端保留之前的导航
	if page == "1" { // page=1 才可能是切换到其他目录
		nav = RedisClient.ChangeFolderNav(folder_title) // 改变 Redis 中存储的当前目录路径缓存
		nav = append(nav, "Home")                       // 在导航路径末尾添加 "Home"
	}
	return nav
}

// GetSubFolders 获取指定父文件夹ID下的所有直接子文件夹，用于前端选择列表
// id: 父文件夹的ID字符串
// 返回值: FolderView.FolderSelectView 切片，用于前端 el-cascader 或类似组件
func GetSubFolders(id string) []FolderView.FolderSelectView {
	folder := models.Folder{}
	folder.ID = int64(utils.StrToInt(id)) // 将字符串ID转换为 int64
	folders := folder.GetSubFolders()     // 获取所有子文件夹

	// 创建指定长度的切片，可以直接通过索引赋值
	folderSelectList := make([]FolderView.FolderSelectView, len(folders))
	for i := range folders {
		folderSelectList[i] = FolderView.FolderSelectView{
			Value: folders[i].ID,                    // 文件夹ID作为值
			Label: folders[i].Title,                 // 文件夹标题作为标签
			Leaf:  folders[i].CountSubFolder() <= 0, // 如果没有子文件夹，则标记为叶子节点
		}
	}
	return folderSelectList
}

// Update 更新文件夹信息
// folderInfo: 包含待更新文件夹信息的 FolderView.FolderInfo 对象
func Update(folderInfo FolderView.FolderInfo) {
	folder := FolderView.ToFolder(folderInfo) // 将视图模型转换为数据模型
	folder.Update()                           // 调用数据模型的更新方法
}

// Add 添加一个新的文件夹
// title: 新文件夹的标题
// fatherTitle: 父文件夹的标题 (用于确定新文件夹的 FolderID)
func Add(title string, fatherTitle string) {
	folder := models.Folder{}
	folder.Title = title // 设置新文件夹的标题

	father := models.Folder{}
	father.Title = fatherTitle
	father.GetFolderByTitle() // 根据父文件夹标题获取其ID

	folder.FolderID = father.ID // 设置新文件夹的父ID

	// 显式设置时间字段为当前时间
	now := time.Now()
	folder.CreatedAt = now
	folder.UpdatedAt = now

	folder.Add() // 调用数据模型的添加方法
}

// Delete 删除指定ID的文件夹及其内容 (逻辑删除)
// id: 要删除的文件夹的ID字符串
// 返回值: 被删除文件夹的ID (int64)
func Delete(id string) int64 {
	folder := models.Folder{}
	folder.ID = int64(utils.StrToInt(id)) // 将字符串ID转换为 int64
	folder.Delete()                       // 调用数据模型的删除方法 (通常是递归删除)
	return folder.ID
}
