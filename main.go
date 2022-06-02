package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name  string `gorm:"type:varchar(20);not null"`
	Phone string `gorm:"type:varchar(20);not null;unique"`
	Pwd   string `gorm:"size:255;not null"`
}

func main() {
	db := InitDB()
	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		phone := ctx.PostForm("phone")
		pwd := ctx.PostForm("pwd")

		if len(phone) != 11 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
			return
		}
		if len(pwd) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
			return
		}
		if len(name) == 0 {
			name = RandomString(10)
		}

		if isPhoneExist(db, phone) {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户存在"})
			return
		}

		newUser := User{
			Name:  name,
			Phone: phone,
			Pwd:   pwd,
		}
		fmt.Println("newUser-----------------------", newUser.Name, newUser.Phone)
		db.Create(&newUser)
		ctx.JSON(200, gin.H{
			"message": "注册成功",
		})
	})
	print(r.Run())
}
func isPhoneExist(db *gorm.DB, phone string) bool {
	var user User
	db.Where("phone = ?", phone).First(&user)

	if user.ID != 0 {
		return true
	}
	return false
}

func RandomString(length int) (str string) {
	var letters = []byte("zxcvbnmasdfghjklqwertyuiopZXCVBNMASDFGHJKLQWERTYUIOP")
	res := make([]byte, length)
	rand.Seed(time.Now().Unix())
	for i := range res {
		res[i] = letters[rand.Intn(len(letters))]
	}
	str = string(res)
	return
}

func InitDB() *gorm.DB {
	username := "root"   //账号
	password := "123456" //密码
	host := "127.0.0.1"  //数据库地址，可以是Ip或者域名
	port := 3306         //数据库端口
	Dbname := "aiguigu"  //数据库名
	timeout := "10s"     //连接超时，10秒

	//拼接下dsn参数, dsn格式可以参考上面的语法，这里使用Sprintf动态拼接dsn参数，因为一般数据库连接参数，我们都是保存在配置文件里面，需要从配置文件加载参数，然后拼接dsn。
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	db.AutoMigrate(&User{})
	if err != nil {
		fmt.Println("数据库连接出错")
	}
	return db
}

// 任意一个类型都可以赋值给空接口
type Ha map[string]interface{}
