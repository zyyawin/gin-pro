package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

//获取标签
func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

func GetTagTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)

	return
}

//增加标签
func AddTag(name string, state int, createBy string) bool {
	db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createBy,
	})
	return true
}

//检查是否存在name
func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name=?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}
	return false
}

//id
func ExistTagByID(id int) bool {
	var tag Tag
	err := db.Select("id").Where("id = ? AND deleted_on = ? ", id, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	if tag.ID > 0 {
		return true
	}

	return false
}

//修改标签
//UPDATE `blog_tag` SET `modified_by` = 'eidt1', `modified_on` = 1571120144, `name` = 'eidt1', `state` = 1  WHERE (id=2)
func EditTag(id int, data interface{}) bool {

	//第一版
	db.Model(&Tag{}).Where("id=?", id).Update(data)
	return true
	/*if err := db.Model(&Tag{}).Where("id=? AND deleted_on=?", id, 0).Update(data).Error; err != nil {
		return err
	}
	return nil
	*/

}

//删除标签
// DELETE FROM `blog_tag`  WHERE (id=1)
func DeleteTag(id int) bool {
	db.Where("id=?", id).Delete(&Tag{})
	return true
}

/*
这属于gorm的Callbacks，可以将回调方法定义为模型结构的指针，在创建、更新、查询、删除时将被调用，如果任何回调返回错误，gorm将停止未来操作并回滚所有更改。


gorm所支持的回调方法：

创建：BeforeSave、BeforeCreate、AfterCreate、AfterSave
更新：BeforeSave、BeforeUpdate、AfterUpdate、AfterSave
删除：BeforeDelete、AfterDelete
查询：AfterFind
*/
func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())
	return nil
}

//编写硬删除代码
func CleanAllTag() bool {
	db.Unscoped().Where("deleted_on!=?", 0).Delete(&Tag{})
	return true

}
