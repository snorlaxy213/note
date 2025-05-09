package RedisClient

import (
	"encoding/json"
	"note-gin/pkg/utils"
	"note-gin/view/ArticleView"
	"time"
)

// GetTempEdit 获取临时编辑内容
// 从Redis中获取临时保存的文章编辑内容
func GetTempEdit(article_view *ArticleView.ArticleEditView) bool {
	isExist := RedisClient.Exists("temp_edit").Val()
	if isExist == 1 {
		s := RedisClient.Get("temp_edit").Val()

		err := json.Unmarshal([]byte(s), article_view)
		utils.ErrReport(err)
		return true
	} else {
		return false
	}
}

// SaveTempEdit 保存临时编辑内容
// 将文章编辑内容序列化后保存到Redis中，有效期为1天
func SaveTempEdit(temp ArticleView.ArticleEditView) string {
	s, _ := json.Marshal(temp)                                 //直接序列化存储了 因为还需要考虑没有ID的临时编辑
	return RedisClient.Set("temp_edit", s, time.Hour*24).Val() //1天
}

// DeleteTempEdit 删除临时编辑内容
// 从Redis中删除临时保存的文章编辑内容
func DeleteTempEdit() int64 {
	return RedisClient.Del("temp_edit").Val()
}
