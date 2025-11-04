package tag_service

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fzzv/go-gin-example/models"
	"github.com/fzzv/go-gin-example/pkg/export"
	"github.com/fzzv/go-gin-example/pkg/gredis"
	"github.com/fzzv/go-gin-example/pkg/logging"
	"github.com/fzzv/go-gin-example/service/cache_service"
	"github.com/xuri/excelize/v2"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	if t.State >= 0 {
		data["state"] = t.State
	}

	return models.EditTag(t.ID, data)
}

func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)

	cache := cache_service.Tag{
		State: t.State,

		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	key := cache.GetTagsKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(key, tags, 3600)
	return tags, nil
}

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}

func (t *Tag) Export() (string, error) {
	// 导出时需要获取所有数据，设置一个很大的 PageSize
	if t.PageSize == 0 {
		t.PageSize = 10000 // 设置一个足够大的值，或者改为 -1 表示不限制
	}
	tags, err := t.GetAll()
	fmt.Println("tags:", tags)
	if err != nil {
		return "", err
	}

	// 创建 Excel 文件
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "标签信息"
	// 默认第一个 sheet，重命名为“标签信息”
	f.SetSheetName(f.GetSheetName(0), sheetName)

	// 表头
	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	if err := f.SetSheetRow(sheetName, "A1", &titles); err != nil {
		return "", fmt.Errorf("设置表头失败: %w", err)
	}

	// 填充数据
	for i, v := range tags {
		// 格式化时间（假设 CreatedOn 和 ModifiedOn 是 int64 或 int 的 Unix 时间戳）
		createdTime := time.Unix(int64(v.CreatedOn), 0).Format("2006-01-02 15:04:05")
		modifiedTime := time.Unix(int64(v.ModifiedOn), 0).Format("2006-01-02 15:04:05")

		row := i + 2 // 从第2行开始（A1 是表头）
		values := []interface{}{
			v.ID,
			v.Name,
			v.CreatedBy,
			createdTime,
			v.ModifiedBy,
			modifiedTime,
		}

		cell := fmt.Sprintf("A%d", row)
		if err := f.SetSheetRow(sheetName, cell, &values); err != nil {
			return "", fmt.Errorf("写入第 %d 行数据失败: %w", row, err)
		}
	}

	// 生成文件名
	filename := fmt.Sprintf("tags-%d.xlsx", time.Now().Unix())
	fullPath := export.GetExcelFullPath() + filename

	// 创建导出目录
	if err := os.MkdirAll(export.GetExcelFullPath(), 0755); err != nil {
		return "", fmt.Errorf("创建导出目录失败: %w", err)
	}
	// 保存文件
	if err := f.SaveAs(fullPath); err != nil {
		return "", fmt.Errorf("保存 Excel 文件失败: %w", err)
	}

	return filename, nil
}
