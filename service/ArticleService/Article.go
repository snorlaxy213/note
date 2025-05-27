package ArticleService

import (
	"errors"
	"mime/multipart"
	"note-gin/models"
	"note-gin/pkg/HttpCode"
	"note-gin/pkg/RedisClient"
	"note-gin/pkg/utils"
	"note-gin/service/FolderService"
	"note-gin/view/ArticleView"
	"note-gin/view/common"
	"strings"
	"time"
)

// ArticleDownLoad 下载文章
// 获取文章标题和内容用于下载
func ArticleDownLoad(ID string) (string, string) {
	article := GetArticleDetail(ID)
	return article.Title, article.MkValue
}

// GetArticleByPage 分页获取文章
// 根据页码获取对应页的文章列表和总数
func GetArticleByPage(page int) ([]ArticleView.ArticleInfo, int) {
	articles := models.Article{}.GetMany(page)
	total := models.Article{}.Count()
	ArticleInfos := ArticleView.ToArticleInfos(articles)
	return ArticleInfos, total
}

// GetArticleDetail 获取文章详情
// 根据ID获取文章的详细信息
func GetArticleDetail(ID string) ArticleView.ArticleDetail {
	article := models.Article{}
	article.ID = int64(utils.StrToInt(ID))
	article.GetArticleInfo()
	articleDetail := ArticleView.ToArticleDetail(article)
	return articleDetail
}

// GetEditArticleDetail 获取文章详情(编辑)
// 根据ID获取文章的详细信息
func GetEditArticleDetail(ID string) ArticleView.ArticleEditView {
	article := models.Article{}
	article.ID = int64(utils.StrToInt(ID))
	article.GetArticleInfo()
	editArticleDetail := ArticleView.ToEditArticleDetail(article)
	return editArticleDetail
}

// ClearRubbish 清空回收站
// 永久删除所有在回收站中的文章
func ClearRubbish() {
	models.Article{}.ClearRubbish()
}

// Delete 删除文章
// 将文章移动到回收站
func Delete(ID string) int64 {
	article := models.Article{}
	article.ID = int64(utils.StrToInt(ID))
	article.Delete()
	return article.ID
}

// DeleteMany 批量删除文章
// 将多篇文章移动到回收站
func DeleteMany(IDs []string) {
	models.Article{}.DeleteMany(IDs)
}

// GetRubbishArticles 获取回收站文章
// 返回所有在回收站中的文章列表
func GetRubbishArticles() common.DataList {
	articles := models.Article{}.GetDeletedArticle()
	articleInfos := ArticleView.ToRubbishArticleInfos(articles) // 将回收文章模型列表转换为视图模型列表
	respDataList := common.DataList{
		Items: articleInfos,
		Total: int64(len(articles)),
	}
	return respDataList
}

// ArticleRecover 恢复文章
// 将回收站中的文章恢复到正常状态
func ArticleRecover(ID string) error {
	article := models.Article{}
	article.ID = int64(utils.StrToInt(ID))
	return article.Recover()
}

// Add 添加文章
// 创建新文章并设置相关属性
func Add(articleEditView *ArticleView.ArticleEditView) {
	article := models.Article{}
	article.Title = articleEditView.Title
	if articleEditView.FolderTitle != "Home" {
		article.FolderID = FolderService.GetFolderByTitle(articleEditView.FolderTitle).ID
	}
	article.Add() //这里调用的方法必须是指针类型
	articleEditView.FolderID = article.FolderID
	articleEditView.DirPath = append(articleEditView.DirPath, articleEditView.FolderID) //先添加自己的根目录
	models.Folder{}.GetFolderPath(articleEditView.FolderID, &articleEditView.DirPath)   //查找路径
}

// Update 更新文章
// 保存文章的修改内容
func Update(articleEditView *ArticleView.ArticleEditView) {
	article := models.Article{}
	article.ID = articleEditView.ID
	article.UpdatedAt = time.Now()
	if len(articleEditView.DirPath) != 0 {
		article.FolderID = articleEditView.DirPath[len(articleEditView.DirPath)-1]
	}

	article.MkValue = articleEditView.MkValue
	article.Title = articleEditView.Title
	article.Update()

	articleEditView.UpdatedAt = article.UpdatedAt.Format("2006-01-02")
	articleEditView.CreatedAt = article.UpdatedAt.Format("2006-01-02")
	articleEditView.ID = article.ID
}

// Edit 编辑文章,获取文章的目录路径信息
func Edit(articleEditView *ArticleView.ArticleEditView) {
	//目录路径回溯
	articleEditView.DirPath = append(articleEditView.DirPath, articleEditView.FolderID) //先添加自己的根目录
	FolderService.GetFolderPath(articleEditView.FolderID, &articleEditView.DirPath)     //查找路径
}

// SetTag 设置文章标签
// 更新文章的标签信息
func SetTag(articleInfo ArticleView.ArticleInfo) {
	article := ArticleView.ToArticle(articleInfo)
	article.SetTag()
}

// TempArticleEditGet 获取临时编辑内容
// 从Redis中获取临时保存的文章编辑内容
func TempArticleEditGet() (ArticleView.ArticleEditView, bool) {
	articleEditView := ArticleView.ArticleEditView{}
	ok := RedisClient.GetTempEdit(&articleEditView)
	return articleEditView, ok
}

// TempArticleEditDelete 删除临时编辑内容
// 从Redis中删除临时保存的文章编辑内容
func TempArticleEditDelete() int64 {
	return RedisClient.DeleteTempEdit()
}

// TempArticleEditSave 保存临时编辑内容
// 将文章编辑内容临时保存到Redis中
func TempArticleEditSave(articleEditView ArticleView.ArticleEditView) bool {
	flag := RedisClient.SaveTempEdit(articleEditView)
	if strings.ToLower(flag) == "ok" {
		return true
	} else {
		return false
	}
}

// UploadArticle 上传文章
// 处理上传的Markdown文件并创建或更新文章
func UploadArticle(files map[string][]*multipart.FileHeader, folder_title string, file_name *string) (bool, error) {
	folder_id := FolderService.GetFolderByTitle(folder_title).ID
	for name, file := range files {
		names := strings.Split(name, ".")
		typeName := names[1]
		if typeName != "md" {
			return false, errors.New(HttpCode.HttpMsg[HttpCode.ERROR_FILE_TYPE])
		}

		fp, _ := file[0].Open()
		b := make([]byte, file[0].Size)
		fp.Read(b)

		article := models.Article{}
		article.Title = names[0]
		*file_name = article.Title
		isExist := article.IsExist()
		if isExist != true {
			article.FolderID = folder_id
			article.MkValue = string(b)
			article.Add()
			return true, nil
		} else { //存在同名文件则更新 不管是否是在同一个目录下  【整个系统不允许出现同名文件】
			article.GetArticleInfoByTitle()
			article.FolderID = folder_id
			article.MkValue = string(b)
			article.Update()
			return false, errors.New(HttpCode.HttpMsg[HttpCode.FILE_IS_EXIST_AND_UPDATE])
		}

	}
	return false, nil
}
