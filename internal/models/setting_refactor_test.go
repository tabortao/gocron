package models

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestSettingRefactorBackwardCompatibility 测试重构后的向后兼容性
func TestSettingRefactorBackwardCompatibility(t *testing.T) {
	// 创建内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&Setting{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	// 保存原始数据库连接
	oldDb := Db
	Db = db
	defer func() { Db = oldDb }()

	setting := &Setting{}

	// 测试 LogRetentionDays
	t.Run("LogRetentionDays", func(t *testing.T) {
		// 测试获取不存在的配置（应返回默认值0）
		days := setting.GetLogRetentionDays()
		if days != 0 {
			t.Errorf("expected 0, got %d", days)
		}

		// 测试创建配置
		if err := setting.UpdateLogRetentionDays(30); err != nil {
			t.Errorf("failed to update: %v", err)
		}

		// 测试获取已存在的配置
		days = setting.GetLogRetentionDays()
		if days != 30 {
			t.Errorf("expected 30, got %d", days)
		}

		// 测试更新已存在的配置
		if err := setting.UpdateLogRetentionDays(60); err != nil {
			t.Errorf("failed to update: %v", err)
		}

		days = setting.GetLogRetentionDays()
		if days != 60 {
			t.Errorf("expected 60, got %d", days)
		}
	})

	// 测试 LogCleanupTime
	t.Run("LogCleanupTime", func(t *testing.T) {
		// 测试获取不存在的配置（应返回默认值"03:00"）
		time := setting.GetLogCleanupTime()
		if time != "03:00" {
			t.Errorf("expected '03:00', got '%s'", time)
		}

		// 测试创建配置
		if err := setting.UpdateLogCleanupTime("02:00"); err != nil {
			t.Errorf("failed to update: %v", err)
		}

		// 测试获取已存在的配置
		time = setting.GetLogCleanupTime()
		if time != "02:00" {
			t.Errorf("expected '02:00', got '%s'", time)
		}

		// 测试更新已存在的配置
		if err := setting.UpdateLogCleanupTime("04:00"); err != nil {
			t.Errorf("failed to update: %v", err)
		}

		time = setting.GetLogCleanupTime()
		if time != "04:00" {
			t.Errorf("expected '04:00', got '%s'", time)
		}
	})

	// 测试 LogFileSizeLimit
	t.Run("LogFileSizeLimit", func(t *testing.T) {
		// 测试获取不存在的配置（应返回默认值0）
		size := setting.GetLogFileSizeLimit()
		if size != 0 {
			t.Errorf("expected 0, got %d", size)
		}

		// 测试创建配置
		if err := setting.UpdateLogFileSizeLimit(100); err != nil {
			t.Errorf("failed to update: %v", err)
		}

		// 测试获取已存在的配置
		size = setting.GetLogFileSizeLimit()
		if size != 100 {
			t.Errorf("expected 100, got %d", size)
		}

		// 测试更新已存在的配置
		if err := setting.UpdateLogFileSizeLimit(200); err != nil {
			t.Errorf("failed to update: %v", err)
		}

		size = setting.GetLogFileSizeLimit()
		if size != 200 {
			t.Errorf("expected 200, got %d", size)
		}
	})

	// 测试数据库中的实际记录
	t.Run("DatabaseRecords", func(t *testing.T) {
		var count int64
		db.Model(&Setting{}).Where("code = ?", SystemCode).Count(&count)
		if count != 3 {
			t.Errorf("expected 3 system settings, got %d", count)
		}

		// 验证每个配置的值
		var settings []Setting
		db.Where("code = ?", SystemCode).Find(&settings)

		valueMap := make(map[string]string)
		for _, s := range settings {
			valueMap[s.Key] = s.Value
		}

		if valueMap[LogRetentionDaysKey] != "60" {
			t.Errorf("expected '60', got '%s'", valueMap[LogRetentionDaysKey])
		}

		if valueMap[LogCleanupTimeKey] != "04:00" {
			t.Errorf("expected '04:00', got '%s'", valueMap[LogCleanupTimeKey])
		}

		if valueMap[LogFileSizeLimitKey] != "200" {
			t.Errorf("expected '200', got '%s'", valueMap[LogFileSizeLimitKey])
		}
	})
}

// TestSettingHelperMethods 测试辅助方法
func TestSettingHelperMethods(t *testing.T) {
	// 创建内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&Setting{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	// 保存原始数据库连接
	oldDb := Db
	Db = db
	defer func() { Db = oldDb }()

	setting := &Setting{}

	t.Run("getSettingValue", func(t *testing.T) {
		// 测试获取不存在的值
		value, err := setting.getSettingValue("test", "key1")
		if err == nil {
			t.Error("expected error for non-existent setting")
		}
		if value != "" {
			t.Errorf("expected empty string, got '%s'", value)
		}

		// 创建一个配置
		db.Create(&Setting{Code: "test", Key: "key1", Value: "value1"})

		// 测试获取存在的值
		value, err = setting.getSettingValue("test", "key1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if value != "value1" {
			t.Errorf("expected 'value1', got '%s'", value)
		}
	})

	t.Run("updateOrCreateSetting", func(t *testing.T) {
		// 测试创建新配置
		err := setting.updateOrCreateSetting("test", "key2", "value2")
		if err != nil {
			t.Errorf("failed to create: %v", err)
		}

		var s Setting
		db.Where("code = ? AND `key` = ?", "test", "key2").First(&s)
		if s.Value != "value2" {
			t.Errorf("expected 'value2', got '%s'", s.Value)
		}

		// 测试更新已存在的配置
		err = setting.updateOrCreateSetting("test", "key2", "value2_updated")
		if err != nil {
			t.Errorf("failed to update: %v", err)
		}

		db.Where("code = ? AND `key` = ?", "test", "key2").First(&s)
		if s.Value != "value2_updated" {
			t.Errorf("expected 'value2_updated', got '%s'", s.Value)
		}

		// 验证只有一条记录
		var count int64
		db.Model(&Setting{}).Where("code = ? AND `key` = ?", "test", "key2").Count(&count)
		if count != 1 {
			t.Errorf("expected 1 record, got %d", count)
		}
	})
}
