package dao

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	UpdateById(ctx context.Context, entity User) error
	FindById(ctx *gin.Context, uid int64) (User, error)
	//FindById(ctx context.Context, uid int64) (User, bizerror)
	//FindByPhone(ctx context.Context, phone string) (User, bizerror)
	//FindByWechat(ctx context.Context, openId string) (User, bizerror)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreatedAt = now
	u.UpdatedAt = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 用户冲突，邮箱冲突
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) UpdateById(ctx context.Context, entity User) error {
	return dao.db.WithContext(ctx).Model(&entity).Where("id = ?", entity.ID).
		Updates(map[string]any{
			"UpdatedAt": time.Now().UnixMilli(),
			"nickname":  entity.Nickname,
			"birthday":  entity.Birthday,
			"about_me":  entity.AboutMe,
		}).Error
}

func (dao *GORMUserDAO) FindById(ctx *gin.Context, uid int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", uid).First(&u).Error
	return u, err
}

type User struct {
	ID        int64  `gorm:"primaryKey,autoIncrement"`
	Email     string `gorm:"unique"`
	Password  string
	Nickname  string
	Birthday  string
	AboutMe   string
	CreatedAt int64
	UpdatedAt int64
}
