/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/5 15:26
 **/
package model

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	User struct {
		ID          uint   `gorm:"primarykey"`
		Name        string `gorm:"not null;unique;"`
		Alias       string `gorm:"not null"`
		Tel         string `gorm:"type:varchar(11);unique;"`
		Email       string `gorm:"type:varchar(100);unique;not null"`
		Password    string `gorm:"not null"`
		Status      uint32 `gorm:"DEFAULT:1;not null"`
		EmailStatus uint32 `gorm:"DEFAULT:2;not null"`
		CreatedAt   time.Time
		UpdatedAt   time.Time
		DeletedAt   gorm.DeletedAt `gorm:"index"`
	}
)

func (m *User) TableName() string {
	return "user"
}

func (m *User) Add(ctx context.Context) error {
	return MainDB().Table(m.TableName()).WithContext(ctx).Create(m).Error
}

func (m *User) Adds(ctx context.Context, data []User) (count int64, err error) {
	tx := MainDB().Table(m.TableName()).WithContext(ctx).CreateInBatches(data, 200)
	err = tx.Error
	count = tx.RowsAffected
	return
}

func (m *User) Del(ctx context.Context, wheres map[string][]interface{}) (count int64, err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	tx := db.Delete(m)
	err = tx.Error
	count = tx.RowsAffected
	return
}
func (m *User) GetAll(ctx context.Context, data *[]User, wheres map[string][]interface{}) (err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	err = db.Find(&data).Error
	return
}
func (m *User) Get(ctx context.Context, start int, size int, data *[]User, wheres map[string][]interface{}, isDelete bool) (total int64, err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	if isDelete {
		db = db.Unscoped().Where("deleted_at is not null")
	} else {
		db = db.Where(map[string]interface{}{"deleted_at": nil})
	}
	tx := db.Limit(size).Offset(start).Find(data)
	err = tx.Error
	total, err = m.Count(ctx, wheres, isDelete)
	return
}

func (m *User) GetById(ctx context.Context, IgnoreDel bool) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	if !IgnoreDel {
		db = db.Unscoped()
	}
	return db.First(m).Error
}

func (m *User) GetByWhere(ctx context.Context, wheres map[string][]interface{}) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.First(m).Error
}

func (m *User) ExistWhere(ctx context.Context, wheres map[string][]interface{}) bool {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	first := db.First(m)
	return first.RowsAffected != 0
}

func (m *User) UpdatesWhere(ctx context.Context, wheres map[string][]interface{}) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Updates(m).Error
}

func (m *User) UpdateWhere(ctx context.Context, wheres map[string][]interface{}, column string, value interface{}) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Update(column, value).Error
}

func (m *User) UpdateStatus(ctx context.Context, status uint32) error {
	return MainDB().Table(m.TableName()).WithContext(ctx).Where("id=?", m.ID).Update("status", status).Error
}

func (m *User) UpdateEmailStatus(ctx context.Context, status uint32) error {
	return MainDB().Table(m.TableName()).WithContext(ctx).Where("id=?", m.ID).Update("email_status", status).Error
}

func (m *User) DelRes(ctx context.Context, wheres map[string][]interface{}) (count int64, err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	tx := db.Update("deleted_at", nil)
	err = tx.Error
	count = tx.RowsAffected
	return
}

func (m *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password))
	return err == nil
}

func (m *User) SetPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	m.Password = string(bytes)
	return nil
}

func (m *User) Count(ctx context.Context, wheres map[string][]interface{}, isDelete bool) (count int64, err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	if isDelete {
		db = db.Unscoped().Where("deleted_at is not null")
	} else {
		db = db.Where(map[string]interface{}{"deleted_at": nil})
	}
	tx := db.Count(&count)
	return count, tx.Error
}
