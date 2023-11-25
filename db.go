package db

import (
	"fmt"

	"github.com/gopherslab/redbook/pkg/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MigrationConfig struct {
	Dialect string `yaml:"dialect"`
	Verbose bool   `yaml:"verbose"`
	Down    bool   `yaml:"down"`
}

type Config struct {
	Host       string           `yaml:"host"`
	Name       string           `yaml:"name"`
	User       string           `yaml:"user"`
	Password   string           `yaml:"password"`
	Port       int              `yaml:"port"`
	Timeout    int              `yaml:"timeout"`
	Migrations *MigrationConfig `yaml:"migrations"`
}

func NewDB(config *Config, log log.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s "+
		"sslmode=disable  sslmode=disable TimeZone=Asia/Calcutta connect_timeout=%d",
		config.Host, config.Port, config.User, config.Password, config.Name, config.Timeout)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Errorf("unable to open db: %v", err)
		return nil, fmt.Errorf("unable to open db coonnection: %w", err)
	}

	return db, nil
}
