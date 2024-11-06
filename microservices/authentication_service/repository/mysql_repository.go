package repository

import (
	"authentication-service/config"
	. "authentication-service/models"
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlAuthRepository struct {
	DB *gorm.DB
}

func FormatMysqlConnection(user, password, address, port, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, address, port, dbName)
}

func NewMysqlAuthRepository(db *gorm.DB) AuthRepository {
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

	if err := db.AutoMigrate(&User{}, &RefreshToken{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	return MysqlAuthRepository{
		DB: db,
	}
}

func (r MysqlAuthRepository) FindUserByID(id uint) (*User, error) {
	var user User
	result := r.DB.First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r MysqlAuthRepository) FindUserByEmail(email string) (*User, error) {
	var user User
	result := r.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r MysqlAuthRepository) UserExistsByEmail(email string) (bool, error) {
	var count int64
	result := r.DB.Model(&User{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r MysqlAuthRepository) FindTokenByUserID(id uint) (*RefreshToken, error) {
	var token RefreshToken
	result := r.DB.Model(&RefreshToken{}).
		Where("user_id = ? AND revoked = ? AND CURRENT_TIMESTAMP < expires_at", id, false).
		First(&token)

	if result.Error != nil {
		return nil, result.Error
	}

	return &token, nil
}

func (r MysqlAuthRepository) CreateNewRefreshToken(t *RefreshToken) error {
	var oldToken RefreshToken

	//see if prexisting token exists
	result := r.DB.Model(&RefreshToken{}).
		Where("user_id = ? AND revoked = ? AND CURRENT_TIMESTAMP < expires_at", t.UserID, false).
		First(&oldToken)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	//start a transaction
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	//revoke old token if it exists
	if result.RowsAffected > 0 {
		oldToken.Revoked = true
		if err := tx.Save(&oldToken).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	//create new token
	if err := tx.Create(t).Error; err != nil {
		tx.Rollback()
		return err
	}

	//commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r MysqlAuthRepository) RevokeAllTokensByUserID(userId uint) error {
	// Update all refresh tokens for the given userID to be revoked
	result := r.DB.Model(&RefreshToken{}).
		Where("user_id = ?", userId).
		Update("revoked", true)

	// Check for any errors during the update
	if result.Error != nil {
		return result.Error
	}

	return nil
}
