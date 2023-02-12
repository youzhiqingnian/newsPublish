package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.Data["data"] = "china"
	c.TplName = "test.html"
}

func (c *MainController) Post() {
	c.Data["data"] = "上海一期最棒"
	c.TplName = "test.html"

}

func (c *MainController) ShowGet() {
	// 获取ORM对象
	/*o := orm.NewOrm()
	// 执行某个操作函数 增删改查
	// 插入操作
	// 获取插入对象
	var user models.User
	user.Name = "heima"
	user.PassWord = "chuanzhi"
	// 插入操作
	count, err := o.Insert(&user)
	if err != nil {
		Println("插入失败")
		return
	}
	Println(count)*/

	// 查询操作
	// 查询对象
	/*var user models.Itcast

	user.Id = 1

	//err := o.Read(&user, "Id")
	err := o.Read(&user) // 如果查询指定的字段是Id，在Read方法中第二个字段参数可以不传，默认就能查，效果和指定了id是一样的
	if err != nil {
		Println("查询失败")
	}

	Println(user)*/

	// 更新操作
	/*var user models.Itcast
	user.Id = 1
	// 查询一下Id=1的数据存在不存在，不存在的话就不更新
	err := o.Read(&user)
	if err != nil {
		println("要更新的数据不存在")
	}

	user.Name = "上海一期"
	count, err := o.Update(&user)
	if err != nil {
		println("更新失败")
	}

	println(count)*/

	// 删除操作
	/*var user models.Itcast
	user.Id = 1
	// 如果不查询直接删除，删除对象的主键要有值
	count, err := o.Delete(&user)
	if err != nil {
		println("删除失败")
	}
	println(count)*/

	c.Data["data"] = "上海"
	c.TplName = "index.html"

}
