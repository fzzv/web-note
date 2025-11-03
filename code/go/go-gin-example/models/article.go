package models

import "github.com/jinzhu/gorm"

type Article struct {
	Model

	// gorm:index，用于声明这个字段为索引，如果使用了自动迁移功能则会有所影响，不使用则无影响
	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	CoverImageUrl string `json:"cover_image_url"`
	State         int    `json:"state"`
}

func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ? AND deleted_on = ? ", id, 0).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if article.ID > 0 {
		return true, nil
	}

	return false, nil
}

func GetArticleTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&Article{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

/*
Preload是什么东西，为什么查询可以得出每一项的关联Tag？
Preload就是一个预加载器，它会执行两条SQL，
分别是SELECT * FROM blog_articles;和SELECT * FROM blog_tag WHERE id IN (1,2,3,4);
那么在查询出结构后，gorm内部处理对应的映射逻辑，将其填充到Article的Tag中，会特别方便，并且避免了循环查询
*/
func GetArticles(pageNum int, pageSize int, maps interface{}) ([]*Article, error) {
	var articles []*Article
	err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return articles, nil
}

/*
Article是如何关联到Tag？
能够达到关联，首先是gorm本身做了大量的约定俗成
Article有一个结构体成员是TagID，就是外键。gorm会通过类名+ID的方式去找到这两个类之间的关联关系
Article有一个结构体成员是Tag，就是我们嵌套在Article里的Tag结构体，我们可以通过Related进行关联查询
*/
func GetArticle(id int) (*Article, error) {
	var article Article
	err := db.Where("id = ? AND deleted_on = ? ", id, 0).First(&article).Related(&article.Tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &article, nil
}

func EditArticle(id int, data interface{}) error {
	if err := db.Model(&Article{}).Where("id = ? AND deleted_on = ? ", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func AddArticle(data map[string]interface{}) error {
	article := Article{
		TagID:         data["tag_id"].(int),
		Title:         data["title"].(string),
		Desc:          data["desc"].(string),
		Content:       data["content"].(string),
		CreatedBy:     data["created_by"].(string),
		State:         data["state"].(int),
		CoverImageUrl: data["cover_image_url"].(string),
	}
	if err := db.Create(&article).Error; err != nil {
		return err
	}

	return nil
}

func DeleteArticle(id int) error {
	if err := db.Where("id = ?", id).Delete(Article{}).Error; err != nil {
		return err
	}

	return nil
}

// func (article *Article) BeforeCreate(scope *gorm.Scope) error {
// 	// time.Now().Unix() 返回当前的时间戳
// 	scope.SetColumn("CreatedOn", time.Now().Unix())
// 	return nil
// }

// func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
// 	scope.SetColumn("ModifiedOn", time.Now().Unix())
// 	return nil
// }

func CleanAllArticle() bool {
	db.Unscoped().Where("deleted_on != ?", 0).Delete(&Article{})
	return true
}
