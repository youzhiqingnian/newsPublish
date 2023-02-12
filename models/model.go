package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// 定义一个结构体
type Itcast struct {
	Id       int
	Name     string
	PassWord string
}

// 定义一个结构体
type User struct {
	Id       int
	Name     string
	PassWord string
	//Pass_Word
	Articles []*Article `orm:"reverse(many)"`
}

type Article struct {
	Id       int       `orm:"pk;auto"`
	ArtiName string    `orm:"size(20)"`
	Atime    time.Time `orm:"auto_now"`
	Acount   int       `orm:"default(0);null"`
	Acontent string    `orm:"size(500)"`
	Aimg     string    `orm:"size(100)"`

	ArticleType *ArticleType `orm:"rel(fk)"`
	Users       []*User      `orm:"rel(m2m)"`
}

// 类型表
type ArticleType struct {
	Id       int
	TypeName string `orm:"size(20)"`

	Articles []*Article `orm:"reverse(many)"`
}

func init() {
	// 操作数据库代码
	// 连接数据库字符串
	// "root:123456@tcp(127.0.0.1:3306")
	// 用户名：密码@tcp(127.0.0.1:3306)/数据库名称?charset=utf8
	/*	conn, err := sql.Open("mysql", "root:cjx8279742@tcp(127.0.0.1:3306)/test?charset=utf8")
		if err != nil {
			Println("连接错误", err)
			Println("链接错误", err)
			return
		}*/

	// 关闭数据库
	//defer conn.Close()
	// 创建表
	/*_, err = conn.Exec("create table user(name VARCHAR(40) , password VARCHAR(40));")
	if err != nil {
		Println("创建表失败", err)
		Println("创建表失败", err)
		return
	}*/

	// 插入
	//conn.Exec("insert INTO user(name, password) values(?,?)", "chuanzhi", "heima")

	// 查询
	/*res, err := conn.Query("select name from user")
	var name string
	for res.Next() {
		res.Scan(&name)
		Println(name)
	}*/
	//

	// 增删改查

	// ORM操作数据库

	// 获取连接对象
	orm.RegisterDataBase("default", "mysql", "root:cjx8279742@tcp(127.0.0.1:3306)/test?charset=utf8")

	// 创建表
	// 注册表
	orm.RegisterModel(new(User), new(Article), new(ArticleType))
	// 生成表
	// 第一个参数是数据库别名，第二个参数是是否强制更新（每次启动项目的时候先drop表，会造成数据丢失）,第三个参数是否可见，生成表的时候的sql语句执行过程要不要看见
	orm.RunSyncdb("default", false, true)

	// 在ORM里面__（双下划线是由特殊含义的），避免定义变量的时候 Pass_Word，在数据库里面变量名就会成为pass__word，结构体中的PassWord变量在数据库表中会变为pass_word
	// 在ORM中，如果表中没有主键，会将名称为Id，类型为int的变量设置为主键
	// 操作表

}
