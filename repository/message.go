package repository

import (
	"danmu/model"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

type MessageRepository struct {
	db *gorm.DB
}

var messageRepository MessageRepository

func GetMessageRepository() MessageRepository {
	return messageRepository
}

func SetUpMessageRepository() {
	var err error
	user := viper.GetString("mysql.user")
	password := viper.GetString("mysql.password")
	dbName := viper.GetString("mysql.dbName")
	dsn := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v?charset=utf8mb4&parseTime=True&loc=Local",user,password,dbName)
	messageRepository.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy:schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		panic("failed to connect database")
	}
}

func (m MessageRepository) SaveRoomMessage(msg *model.Message) error {
	return m.db.Model(&model.Message{}).Save(msg).Error
}

func (m MessageRepository) GetRoomMessageList(room string) ([]model.Message, error) {
	var msgList []model.Message
	err := m.db.Model(&model.Message{}).Where("room=?", room).Order("id").Find(&msgList).Limit(20).Error
	if err != nil {
		log.Println("get msgList failed room=%+v err=%v", room, err)
		return nil, err
	}
	return msgList, nil
}
