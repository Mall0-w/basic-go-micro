package repository

import (
	"fmt"
	"user-service/config"
	. "user-service/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlUserRepository struct {
	DB *gorm.DB
}

func FormatMysqlConnection(user, password, address, port, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, address, port, dbName)
}

func NewMysqlUserRepository(db *gorm.DB) UserRepository {
	//go to default db if none is given

	if db == nil {
		conf := config.LoadConfig()
		dsn := FormatMysqlConnection(conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBPort, conf.DBName)
		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	return MysqlUserRepository{
		DB: db,
	}
}

func (r MysqlUserRepository) FindByID(id int64) (*User, error) {
	var user User
	result := r.DB.First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r MysqlUserRepository) ExistsByID(id int64) (bool, error) {
	var count int64
	result := r.DB.Model(&User{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r MysqlUserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	result := r.DB.Model(&User{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r MysqlUserRepository) Create(user *User) (*User, error) {
	result := r.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r MysqlUserRepository) Update(user *User) (*User, error) {
	result := r.DB.Save(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r MysqlUserRepository) Delete(id int64) error {
	result := r.DB.Delete(&User{ID: id})
	return result.Error
}
