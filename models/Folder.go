package models

import (
	"note-gin/config"
	"time"
)

// PageSize 全局变量，用于分页查询时每页的项目数量，从应用配置中读取
var PageSize = config.Conf.AppConfig.PageSize

// Folder 代表一个文件夹的结构体
type Folder struct {
	BaseModel        // 嵌入 BaseModel，包含ID, CreatedAt, UpdatedAt, Deleted, DeletedTime等通用字段
	Title     string // 文件夹标题
	FolderID  int64  // 父文件夹ID，如果为0，则表示为根文件夹
}

// GetRootFolder 获取所有根文件夹 (folder_id 为 0 的文件夹)
func (this Folder) GetRootFolder() (roots []Folder) {
	db.Find(&roots, "folder_id=?", 0)
	return
}

// GetFolderPath 递归获取指定 FolderID 的完整目录路径 (从根目录到当前目录的ID列表)
// FolderID: 当前要查找其父路径的文件夹ID
// DirPath: 指向一个 int64 切片的指针，用于存储路径上的文件夹ID
func (this Folder) GetFolderPath(FolderID int64, DirPath *[]int64) {
	if FolderID == 0 { // 如果 FolderID 为0，表示已到达根目录或初始调用ID为0，则返回
		return
	}
	folder := Folder{}
	db.Where("id=?", FolderID).First(&folder) // 根据 FolderID 查找文件夹

	if folder.FolderID != 0 { // 如果当前文件夹有父文件夹
		*DirPath = append([]int64{folder.FolderID}, *DirPath...) // 将父文件夹ID添加到路径的前面
		this.GetFolderPath(folder.FolderID, DirPath)             // 递归查找父文件夹的路径
	} else { // 如果没有父文件夹 (即当前文件夹是路径中的某个根节点下的直接子节点)，则返回
		return
	}
}

// GetFolderByID 根据文件夹的 ID 获取文件夹信息，并更新接收者 this
func (this Folder) GetFolderByID() {
	db.Where("id=?", this.ID).First(&this)
	return
}

// GetSubFile 获取指定文件夹下的子文件和子文件夹 (分页)
// page: 请求的页码
// 返回值: fds (子文件夹列表), articles (子文章列表), total (子文件和子文件夹总数)
func (this Folder) GetSubFile(page int) (fds []Folder, articles []Article, total int) {
	if PageSize <= 0 { // 如果全局 PageSize 未正确配置或为0，则设置默认值为10
		PageSize = 10
	}
	fds = this.GetSubFolderOnPage(page, PageSize) // 获取当前页的子文件夹
	total = this.CountSubFile()                   // 获取子文件和子文件夹的总数
	fdsCount := len(fds)                          // 当前页获取到的子文件夹数量

	if fdsCount < PageSize && fdsCount > 0 { // 如果子文件夹数量小于每页大小但大于0 (即子文件夹未填满一页，但有子文件夹)
		// 计算剩余空间应填充的文章数量
		articles = this.GetSubArticle(PageSize-fdsCount, 0) // 获取文章以填补页面剩余空间，从第一页的文章开始取
	} else if fdsCount == 0 { // 如果当前页没有子文件夹 (可能意味着所有子文件夹已显示完毕，或者此页开始显示文章)
		SubFolderCount := this.CountSubFolder()                           // 子文件夹总数
		offset := PageSize - (SubFolderCount % PageSize)                  // 计算文章开始的偏移量，确保文章列表接在文件夹列表之后
		page = page - ((SubFolderCount / PageSize) + 1)                   // 调整页码以正确获取文章分页
		articles = this.GetSubArticle(PageSize, offset+(page-1)*PageSize) // 获取对应页码的文章
	}
	return
}

// GetSubFolders 获取当前文件夹下的所有子文件夹
func (this Folder) GetSubFolders() (folders []Folder) {
	db.Table("folder").Where("folder_id=?", this.ID).Find(&folders)
	return
}

// GetSubFolderOnPage 获取当前文件夹下指定页码和页面大小的子文件夹
// page: 请求的页码
// PageSize: 每页的项目数量
func (this Folder) GetSubFolderOnPage(page, PageSize int) (fds []Folder) {
	db.Limit(PageSize).Offset((page-1)*PageSize).Find(&fds, "folder_id=?", this.ID)
	return
}

// GetSubArticle 获取当前文件夹下指定数量和偏移量的子文章
// limit: 需要获取的文章数量
// offset: 查询结果的偏移量
func (this Folder) GetSubArticle(limit, offset int) (articles []Article) {
	db.Limit(limit).Offset(offset).Where("deleted=?", 0).Select([]string{"id", "title", "updated_at", "tags"}).Find(&articles, "folder_id=?", this.ID)
	return
}

// GetFolderInfo 根据接收者 this 的当前字段值 (通常是ID) 获取文件夹的完整信息
func (this Folder) GetFolderInfo() {
	db.Where(this).First(&this)
}

// GetFolderByTitle 根据文件夹的标题获取文件夹信息，并更新接收者 this
func (this Folder) GetFolderByTitle() {
	db.Where("title=?", this.Title).First(&this)
}

// CountSubFile 计算当前文件夹下子文件和子文件夹的总数
func (this Folder) CountSubFile() int {
	sum := this.CountSubFolder() + this.CountSubArticle() // 总数 = 子文件夹数量 + 子文章数量
	return sum
}

// CountSubFolder 计算当前文件夹下的子文件夹数量
func (this Folder) CountSubFolder() (count int) {
	db.Table("folder").Where("folder_id=?", this.ID).Count(&count)
	return
}

// CountSubArticle 计算当前文件夹下未被删除的子文章数量
func (this Folder) CountSubArticle() (count int) {
	db.Model(&Article{}).Where("folder_id=? and deleted=?", this.ID, 0).Count(&count)
	return
}

// Add 创建一个新的文件夹记录到数据库
func (this *Folder) Add() {
	db.Create(this)
}

// Update 更新文件夹信息，主要是标题和更新时间
func (this *Folder) Update() {
	db.Model(this).Where("id=?", this.ID).Updates(map[string]interface{}{"title": this.Title, "updated_at": time.Now()})
}

// Delete 递归删除文件夹及其下的所有子文件夹和文章 (逻辑删除)
func (this *Folder) Delete() {
	db.Delete(this)    // 删除当前文件夹 (gorm的软删除)
	deleteDFS(this.ID) // 递归删除子项
}

// deleteDFS 深度优先搜索并递归删除指定文件夹ID下的所有文章和子文件夹 (逻辑删除)
// FolderID: 要删除其内容的父文件夹ID
func deleteDFS(FolderID int64) {
	// 将该文件夹下的所有文章标记为已删除
	db.Table("article").Where("folder_id=?", FolderID).Update("deleted", true)
	sub_folder := []Folder{}
	// 查找该文件夹下的所有直接子文件夹
	db.Find(&sub_folder, "folder_id=?", FolderID)
	for index := range sub_folder { // 遍历所有子文件夹
		sub_folder[index].Delete() // 递归调用 Delete 方法删除每个子文件夹
	}
}
