package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicateEmail = errors.New("邮箱已存在")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAOer interface {
	Insert(ctx context.Context, u User) error
	QueryByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, u User) error
	QueryById(ctx context.Context, userId int64) (User, error)
	QueryByPhone(ctx context.Context, phone string) (User, error)
}

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAOer {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if res, ok := err.(*mysql.MySQLError); ok {
		const duplicateErrNum uint16 = 1062
		if res.Number == duplicateErrNum {
			// 邮箱冲突
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) QueryByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) Update(ctx context.Context, u User) error {
	return dao.db.WithContext(ctx).Model(&User{}).Where("id = ?", u.Id).Updates(u).Error
}

func (dao *UserDAO) QueryById(ctx context.Context, userId int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", userId).First(&u).Error
	return u, err
}

func (dao *UserDAO) QueryByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 代表这是一个可以为 NULL 的列
	Email    sql.NullString `gorm:"unique"`
	Password string
	Birthday int64
	Nickname string `gorm:"type:varchar(100)"`
	AboutMe  string `gorm:"type=varchar(4096)"`
	// 代表这是一个可以为 NULL 的列
	Phone sql.NullString `gorm:"unique"`

	// 时区，UTC 0 的毫秒数
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64

	// json 存储
	//Addr string
}
