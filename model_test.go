package curd

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"ma_system/pkg/orm"
	"testing"
	"time"
)

// ProductBrand 品牌表
type ProductBrand struct {
	ID                  int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name                string         `gorm:"column:name;not null;comment:品牌名称" json:"name"`                                     // 品牌名称
	FirstLetter         string         `gorm:"column:first_letter;not null;comment:首字母" json:"first_letter"`                      // 首字母
	Sequence            int32          `gorm:"column:sequence;not null;comment:排序" json:"sequence"`                               // 排序
	Status              int32          `gorm:"column:status;not null;comment:品牌状态;1=正常;0=隐藏" json:"status"`                       // 品牌状态;1=正常;0=隐藏
	ProductCount        int32          `gorm:"column:product_count;not null;comment:产品数量" json:"product_count"`                   // 产品数量
	ProductCommentCount int32          `gorm:"column:product_comment_count;not null;comment:产品评论数量" json:"product_comment_count"` // 产品评论数量
	BrandStory          string         `gorm:"column:brand_story;not null;comment:品牌故事" json:"brand_story"`                       // 品牌故事
	Cover               int64          `gorm:"column:cover;comment:品牌logo;对应uploads" json:"cover"`                                // 品牌logo;对应uploads
	BannerCover         int64          `gorm:"column:banner_cover;comment:品牌横幅照片l;对应uploads" json:"banner_cover"`                 // 品牌横幅照片l;对应uploads
	CreatedAt           time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at;not null" json:"deleted_at"`
}

func (p ProductBrand) TableName() string {
	return "product_brands"
}

var dsn = "mathias:123456@tcp(localhost:3306)/mm_system?parseTime=true&loc=Local"

func TestNewModel(t *testing.T) {

	db := orm.MustNewMysql(&orm.Config{
		DSN: dsn,
	})

	productBrandModel := NewModel[ProductBrand](db.DB)
	var conditions []map[string]interface{}
	conditions = append(conditions, map[string]interface{}{
		"name":         "三星",
		"letter_first": "S",
	})
	conditions = append(conditions, map[string]interface{}{
		"name":         "小米",
		"letter_first": "M",
	})
	productBrands, total, err := productBrandModel.Select(context.Background(), &ModelBO{
		PageNo:     1,
		PageSize:   10,
		Conditions: conditions,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(productBrands, total)
	t.Log(productBrands[0].Name)
}

func TestModel_FindOne(t *testing.T) {
	db := orm.MustNewMysql(&orm.Config{
		DSN: dsn,
	})

	productBrandModel := NewModel[ProductBrand](db.DB)
	productBrand, err := productBrandModel.IsMaster(true).FindOne(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(productBrand.ID)
}

func TestModel_Insert(t *testing.T) {
	db := orm.MustNewMysql(&orm.Config{
		DSN: dsn,
	})

	productBrandModel := NewModel[ProductBrand](db.DB)
	insert, err := productBrandModel.Insert(context.Background(), map[string]interface{}{
		"name":                  "小米ccccc",
		"first_letter":          "M",
		"sequence":              1,
		"status":                1,
		"product_count":         1,
		"product_comment_count": 1,
		"brand_story":           "小米的故事",
		"cover":                 1,
		"banner_cover":          1,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(insert.ID)
}

func TestModel_Update(t *testing.T) {
	db := orm.MustNewMysql(&orm.Config{
		DSN: dsn,
	})

	productBrandModel := NewModel[ProductBrand](db.DB)
	err := productBrandModel.Update(context.Background(), 68, map[string]interface{}{
		"name":                  "苹果bbbbbbbb",
		"first_letter":          "A",
		"sequence":              1,
		"status":                1,
		"product_count":         1,
		"product_comment_count": 1,
		"brand_story":           "小米的故事",
		"cover":                 1,
		"banner_cover":          1,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestModel_Delete(t *testing.T) {
	db := orm.MustNewMysql(&orm.Config{
		DSN: dsn,
	})

	productBrandModel := NewModel[ProductBrand](db.DB)
	productBrand, err := productBrandModel.IsMaster(true).FindOne(context.Background(), 68)

	err = productBrandModel.IsMaster(true).Delete(context.Background(), productBrand)
	if err != nil {
		t.Fatal(err)
	}
}
