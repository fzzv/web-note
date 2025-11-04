package article_service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/fzzv/go-gin-example/models"
	"github.com/fzzv/go-gin-example/pkg/export"
	"github.com/fzzv/go-gin-example/pkg/gredis"
	"github.com/fzzv/go-gin-example/pkg/logging"
	"github.com/fzzv/go-gin-example/service/cache_service"
	"github.com/unknwon/com"
	"github.com/xuri/excelize/v2"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

func (a *Article) Add() error {
	article := map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"created_by":      a.CreatedBy,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
	}

	if err := models.AddArticle(article); err != nil {
		return err
	}

	return nil
}

func (a *Article) Edit() error {
	return models.EditArticle(a.ID, map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"modified_by":     a.ModifiedBy,
	})
}

func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}

	article, err := models.GetArticle(a.ID)
	log.Print(article)
	if err != nil {
		return nil, err
	}

	gredis.Set(key, article, 3600)
	return article, nil
}

func (a *Article) GetAll() ([]*models.Article, error) {
	var (
		articles, cacheArticles []*models.Article
	)

	cache := cache_service.Article{
		TagID: a.TagID,
		State: a.State,

		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}
	key := cache.GetArticlesKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(key, articles, 3600)
	return articles, nil
}

func (a *Article) Delete() error {
	return models.DeleteArticle(a.ID)
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

func (a *Article) Count() (int, error) {
	return models.GetArticleTotal(a.getMaps())
}

func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}

	return maps
}

func (a *Article) Export() (string, error) {
	if a.PageSize == 0 {
		a.PageSize = 10000
	}
	if a.TagID == 0 {
		a.TagID = -1
	}
	a.State = 1
	articles, err := a.GetAll()
	fmt.Println("articles:", articles)
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "文章信息"
	f.SetSheetName(f.GetSheetName(0), sheetName)

	titles := []string{"ID", "标题", "描述", "内容", "封面图片", "状态", "创建人", "创建时间", "修改人", "修改时间", "标签ID"}
	if err := f.SetSheetRow(sheetName, "A1", &titles); err != nil {
		return "", fmt.Errorf("设置表头失败: %w", err)
	}

	for i, v := range articles {
		createdTime := time.Unix(int64(v.CreatedOn), 0).Format("2006-01-02 15:04:05")
		modifiedTime := time.Unix(int64(v.ModifiedOn), 0).Format("2006-01-02 15:04:05")

		row := i + 2
		values := []interface{}{
			v.ID,
			v.Title,
			v.Desc,
			v.Content,
			v.CoverImageUrl,
			v.State,
			v.CreatedBy,
			createdTime,
			v.ModifiedBy,
			modifiedTime,
			v.TagID,
		}

		cell := fmt.Sprintf("A%d", row)
		if err := f.SetSheetRow(sheetName, cell, &values); err != nil {
			return "", fmt.Errorf("写入第 %d 行数据失败: %w", row, err)
		}
	}

	filename := fmt.Sprintf("articles-%d.xlsx", time.Now().Unix())
	fullPath := export.GetExcelFullPath() + filename

	if err := os.MkdirAll(export.GetExcelFullPath(), 0755); err != nil {
		return "", fmt.Errorf("创建导出目录失败: %w", err)
	}
	if err := f.SaveAs(fullPath); err != nil {
		return "", fmt.Errorf("保存 Excel 文件失败: %w", err)
	}

	return filename, nil
}

func (a *Article) Import(r io.Reader) error {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}
	defer xlsx.Close()

	rows, err := xlsx.GetRows("文章信息")
	if err != nil {
		return err
	}

	for irow, row := range rows {
		if irow == 0 {
			continue
		}
		if len(row) < 3 {
			continue
		}
		title := row[1]
		desc := row[2]
		content := row[3]
		coverImageUrl := row[4]
		state := row[5]
		createdBy := row[6]
		createdTime := row[7]
		modifiedBy := row[8]
		modifiedTime := row[9]
		tagId := row[10]
		models.AddArticle(map[string]interface{}{
			"tag_id":          com.StrTo(tagId).MustInt(),
			"title":           title,
			"desc":            desc,
			"content":         content,
			"cover_image_url": coverImageUrl,
			"state":           com.StrTo(state).MustInt(),
			"created_by":      createdBy,
			"created_time":    createdTime,
			"modified_by":     modifiedBy,
			"modified_time":   modifiedTime,
		})
	}
	return nil
}
