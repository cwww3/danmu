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

type MySQLRepository struct {
	db *gorm.DB
}

var mysqlRepository MySQLRepository

func GetMySQLRepository() MySQLRepository {
	return mysqlRepository
}

func SetUpMySQLRepository() {
	var err error
	user := viper.GetString("mysql.user")
	password := viper.GetString("mysql.password")
	dbName := viper.GetString("mysql.dbName")
	dsn := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v?charset=utf8mb4&parseTime=True&loc=Local", user, password, dbName)
	mysqlRepository.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	mysqlRepository.db.AutoMigrate(new(model.Room))
	if err != nil {
		panic("failed to connect database")
	}
}

func (m MySQLRepository) SaveRoomMessage(msg *model.Message) error {
	return m.db.Model(&model.Message{}).Save(msg).Error
}

func (m MySQLRepository) GetRoomMessageList(room string) ([]model.Message, error) {
	var msgList []model.Message
	err := m.db.Model(&model.Message{}).Where("room=?", room).Order("id").Find(&msgList).Limit(20).Error
	if err != nil {
		log.Printf("get msgList failed room=%+v err=%v\n", room, err)
		return nil, err
	}
	return msgList, nil
}

func (m MySQLRepository) SaveRoom(room *model.Room) (uint, error) {
	err := m.db.Model(&model.Room{}).Save(room).Error
	return room.ID, err
}

func (m MySQLRepository) GetRoomByUser(user string) (*model.Room, error) {
	room := new(model.Room)
	err := m.db.Model(&model.Room{}).First(room, "user = ?", user).Error
	return room, err
}
