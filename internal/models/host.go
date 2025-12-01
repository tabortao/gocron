package models

import (
	"gorm.io/gorm"
)

// 主机
type Host struct {
	Id        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string `json:"name" gorm:"type:varchar(64);not null"`
	Alias     string `json:"alias" gorm:"type:varchar(32);not null;default:''"`
	Port      int    `json:"port" gorm:"not null;default:5921"`
	Remark    string `json:"remark" gorm:"type:varchar(100);not null;default:''"`
	BaseModel `json:"-" gorm:"-"`
	Selected  bool `json:"-" gorm:"-"`
}

// 新增
func (host *Host) Create() (insertId int, err error) {
	result := Db.Create(host)
	if result.Error == nil {
		insertId = host.Id
	}

	return insertId, result.Error
}

func (host *Host) UpdateBean(id int) (int64, error) {
	result := Db.Model(&Host{}).Where("id = ?", id).
		Select("name", "alias", "port", "remark").
		Updates(host)
	return result.RowsAffected, result.Error
}

// 更新
func (host *Host) Update(id int, data CommonMap) (int64, error) {
	updateData := make(map[string]interface{})
	for k, v := range data {
		updateData[k] = v
	}
	result := Db.Model(&Host{}).Where("id = ?", id).UpdateColumns(updateData)
	return result.RowsAffected, result.Error
}

// 删除
func (host *Host) Delete(id int) (int64, error) {
	result := Db.Delete(&Host{}, id)
	return result.RowsAffected, result.Error
}

func (host *Host) Find(id int) error {
	return Db.First(host, id).Error
}

func (host *Host) NameExists(name string, id int) (bool, error) {
	var count int64
	query := Db.Model(&Host{}).Where("name = ?", name)
	if id != 0 {
		query = query.Where("id != ?", id)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (host *Host) List(params CommonMap) ([]Host, error) {
	host.parsePageAndPageSize(params)
	list := make([]Host, 0)
	query := Db.Order("id DESC")
	host.parseWhere(query, params)
	err := query.Limit(host.PageSize).Offset(host.pageLimitOffset()).Find(&list).Error

	return list, err
}

func (host *Host) AllList() ([]Host, error) {
	list := make([]Host, 0)
	err := Db.Select("name", "port").Order("id DESC").Find(&list).Error

	return list, err
}

func (host *Host) Total(params CommonMap) (int64, error) {
	var count int64
	query := Db.Model(&Host{})
	host.parseWhere(query, params)
	err := query.Count(&count).Error
	return count, err
}

// 解析where
func (host *Host) parseWhere(query *gorm.DB, params CommonMap) {
	if len(params) == 0 {
		return
	}
	id, ok := params["Id"]
	if ok && id.(int) > 0 {
		query.Where("id = ?", id)
	}
	name, ok := params["Name"]
	if ok && name.(string) != "" {
		query.Where("name = ?", name)
	}
}
