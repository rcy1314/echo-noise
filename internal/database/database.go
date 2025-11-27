package database

import (
    "errors"
    "fmt"
    "os"
    "time"

    "github.com/lin-snow/ech0/config"
    "github.com/lin-snow/ech0/internal/models"
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "github.com/glebarez/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func getEnvOrConfig(envKey, configValue string) string {
    if value := os.Getenv(envKey); value != "" {
        return value
    }
    return configValue
}

func InitDB() error {
    if DB != nil {
        return nil
    }

    dbType := getEnvOrConfig("DB_TYPE", config.Config.Database.Type)
    var err error

    // 配置 GORM
    gormConfig := &gorm.Config{
        PrepareStmt: true,
        Logger: logger.Default.LogMode(logger.Silent),
    }

    switch dbType {
    case "sqlite":
        dbPath := getEnvOrConfig("DB_PATH", config.Config.Database.Path)
        dir := fmt.Sprintf("%s", dbPath[:len(dbPath)-len("/noise.db")])
        if err := os.MkdirAll(dir, os.ModePerm); err != nil {
            return fmt.Errorf("创建数据库目录失败: %v", err)
        }
        DB, err = gorm.Open(sqlite.Open(dbPath), gormConfig)

    case "postgres":
        host := getEnvOrConfig("DB_HOST", config.Config.Database.Host)
        port := getEnvOrConfig("DB_PORT", config.Config.Database.Port)
        user := getEnvOrConfig("DB_USER", config.Config.Database.User)
        password := getEnvOrConfig("DB_PASSWORD", config.Config.Database.Password)
        dbname := getEnvOrConfig("DB_NAME", config.Config.Database.DBName)
        sslmode := getEnvOrConfig("DB_SSL_MODE", "disable")
        timezone := getEnvOrConfig("DB_TIMEZONE", "Asia/Shanghai")
        
        dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s pool_max_conns=20",
            host, user, password, dbname, port, sslmode, timezone)
        DB, err = gorm.Open(postgres.Open(dsn), gormConfig)

    case "mysql":
        host := getEnvOrConfig("DB_HOST", config.Config.Database.Host)
        port := getEnvOrConfig("DB_PORT", config.Config.Database.Port)
        user := getEnvOrConfig("DB_USER", config.Config.Database.User)
        password := getEnvOrConfig("DB_PASSWORD", config.Config.Database.Password)
        dbname := getEnvOrConfig("DB_NAME", config.Config.Database.DBName)
        charset := getEnvOrConfig("DB_CHARSET", "utf8mb4")
        
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=5s&readTimeout=5s&writeTimeout=5s&maxAllowedPacket=0&interpolateParams=true",
            user, password, host, port, dbname, charset)
        DB, err = gorm.Open(mysql.New(mysql.Config{
            DSN: dsn,
            DefaultStringSize: 191,
        }), gormConfig)

    default:
        return fmt.Errorf("不支持的数据库类型: %s", dbType)
    }

    if err != nil {
        return fmt.Errorf("数据库连接失败: %v", err)
    }

    // 配置连接池
    sqlDB, err := DB.DB()
    if err != nil {
        return fmt.Errorf("获取数据库实例失败: %v", err)
    }

    // 设置连接池参数
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    sqlDB.SetConnMaxIdleTime(time.Minute * 10)

    models.SetDB(DB)

    if err = models.MigrateDB(DB); err != nil {
        return fmt.Errorf("数据库迁移失败: %v", err)
    }

    return nil
}

func GetSetting() (*models.Setting, error) {
    if DB == nil {
        return nil, errors.New("数据库未初始化")
    }

    var setting models.Setting
    if err := DB.First(&setting).Error; err != nil {
        return nil, err
    }
    return &setting, nil
}

func UpdateSetting(allowRegistration bool) error {
    if DB == nil {
        return errors.New("数据库未初始化")
    }

    return DB.Model(&models.Setting{}).Where("1 = 1").Update("allow_registration", allowRegistration).Error
}

func GetDB() (*gorm.DB, error) {
    if DB == nil {
        return nil, errors.New("数据库未初始化")
    }
    return DB, nil
}

func GetSystemStatus() (map[string]interface{}, error) {
    if DB == nil {
        return nil, errors.New("数据库未初始化")
    }

    var setting models.Setting
    if err := DB.First(&setting).Error; err != nil {
        return nil, err
    }

    var adminUser models.User
    if err := DB.Where("is_admin = ?", true).First(&adminUser).Error; err != nil {
        return nil, err
    }

    var totalMessages int64
    if err := DB.Model(&models.Message{}).Count(&totalMessages).Error; err != nil {
        return nil, err
    }

    return map[string]interface{}{
        "username":          adminUser.Username,
        "isAdmin":          adminUser.IsAdmin,
        "total_messages":    totalMessages,
        "allowRegistration": setting.AllowRegistration,
    }, nil
}

func ReconnectDB() error {
    if DB != nil {
        sqlDB, err := DB.DB()
        if err == nil {
            if err := sqlDB.Close(); err != nil {
                return fmt.Errorf("关闭数据库连接失败: %v", err)
            }
        }
        DB = nil
    }
    return InitDB()
}
