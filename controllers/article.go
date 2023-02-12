package controllers

import (
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"math"
	"path"
	"shanghai/models"
	"time"
)

type ArticleController struct {
	beego.Controller
}

// 展示文章列表页
func (this *ArticleController) ShowArticleList() {

	// session判断
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}

	// 获取数据
	// 高级查询
	// 指定表
	o := orm.NewOrm()
	qs := o.QueryTable("Article") // queryseter 查询集
	// 查询所有
	var articles []models.Article
	/*_, err := qs.All(&articles)
	if err != nil {
		println("查询数据错误")
	}*/
	//查询总记录数
	// 根据响应的类型查询响应的文章
	typeName := this.GetString("select")
	var count int64
	// 每页有几条数据
	pageSize := 2
	//math.Floor(1.1) // 地板函数，向下取整   1.9=1

	// 获取页码
	pageIndex, err := this.GetInt("pageIndex")
	println(pageIndex)
	if err != nil {
		pageIndex = 1
	}
	//println(pageIndex)

	// 获取数据
	// 作用就是获取数据库部分数据，第一个参数，获取几条数据;第二个参数，从哪条数据开始获取,返回值还是querySeter
	// 起始位置计算
	start := (pageIndex - 1) * pageSize
	if typeName == "" {
		count, _ = qs.Count()
	} else {
		count, _ = qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()
	}
	// 获取总页数
	pageCount := math.Ceil(float64(count) / float64(pageSize)) // 天花板函数向上取整，1.1=2

	//println("start=", start)
	//qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)

	// 获取文章类型
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types

	if typeName == "" {
		qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)
	} else {
		qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles)
	}

	// 传递数据
	this.Data["title"] = "文章列表"
	this.Data["userName"] = userName
	this.Data["typeName"] = typeName
	this.Data["pageIndex"] = pageIndex
	this.Data["pageCount"] = int(pageCount)
	this.Data["count"] = count
	this.Data["articles"] = articles

	// 指定视图布局
	this.Layout = "layout.html"
	this.TplName = "index.html"
}

// 展示添加文章页面
func (this *ArticleController) ShowAddArticle() {
	//查询所有类型数据，并展示
	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)

	// 传递数据
	userName := this.GetSession("userName")
	this.Data["userName"] = userName.(string)
	this.Data["types"] = types
	this.Data["title"] = "添加文章"
	this.Layout = "layout.html"
	this.TplName = "add.html"
}

// 获取添加文章数据
func (this *ArticleController) HandleAddArticle() {
	// 1.获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")

	// 2.校验数据
	if articleName == "" || content == "" {
		this.Data["errmsg"] = "添加数据不完整"
		this.TplName = "add.html"
		return
	}

	// 处理文件上传
	filePath := UploadFile(&this.Controller, "uploadname")

	// 3.处理数据
	// 插入操作
	o := orm.NewOrm()

	var article models.Article

	article.ArtiName = articleName
	article.Acontent = content
	article.Aimg = filePath
	// 给文章添加类型
	// 获取类型数据
	typeName := this.GetString("select")
	//根据名称查询对象
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType, "TypeName")

	article.ArticleType = &articleType

	o.Insert(&article)

	// 4.返回页面
	this.Redirect("/articl/showArticleList", 302)

}

// 展示文章详情页面
func (this *ArticleController) ShowArticleDetail() {
	// 获取数据
	id, err := this.GetInt("articleId")
	// 校验数据
	if err != nil {
		println("传递的链接错误")
	}

	println("id===", id)

	// 操作数据
	o := orm.NewOrm()
	var article models.Article
	article.Id = id

	//o.Read(&article)
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("Id", id).One(&article)

	// 修改阅读量
	article.Acount += 1
	o.Update(&article)

	// 多对多插入浏览记录,把users插入article中,m2m是个多对多操作对象
	m2m := o.QueryM2M(&article, "Users")
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}

	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")

	// 插入操作
	m2m.Add(user)

	// 查询 没法去重
	//o.LoadRelated(&article, "Users")
	var users []models.User
	// Distinct是去重
	o.QueryTable("User").Filter("Articles__Article__Id", id).Distinct().All(&users)

	// 返回视图页面
	this.Data["title"] = "文章详情"
	this.Data["users"] = users
	this.Data["article"] = article
	this.Data["userName"] = userName.(string)
	this.Layout = "layout.html"
	this.TplName = "content.html"
}

// 展示编辑文章页面
func (this *ArticleController) ShowUpdateArticle() {
	// 获取数据
	id, err := this.GetInt("articleId")
	// 校验数据
	if err != nil {
		println("传递的链接错误")
	}

	println("id===", id)

	// 操作数据
	// 查询相应文章
	o := orm.NewOrm()
	var article models.Article
	article.Id = id

	o.Read(&article)

	// 返回视图页面
	userName := this.GetSession("userName")
	this.Data["title"] = "编辑文章"
	this.Data["userName"] = userName.(string)
	this.Data["article"] = article
	this.Layout = "layout.html"
	this.TplName = "update.html"
}

func UploadFile(this *beego.Controller, filePath string) string {
	file, head, err := this.GetFile(filePath)
	if head.Filename == "" {
		return "NoImg"
	}
	if err != nil {
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = "add.html"
		return ""
	}
	defer file.Close()
	// 1.文件大小
	if head.Size > 5000000 {
		this.Data["errmsg"] = "文件太大，请重新上传"
		this.TplName = "add.html"
		return ""
	}

	// 2.文件格式
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		this.Data["errmsg"] = "文件格式错误，请重新上传"
		this.TplName = "add.html"
		return ""
	}

	// 3.防止重名
	fileName := time.Now().Format("2006-01-02-15:04:05") + ext

	this.SaveToFile(filePath, "./static/img/"+fileName)

	return "/static/img/" + fileName
}

// 处理编辑页面数据
func (this *ArticleController) HandleUpdateArticle() {
	// 获取数据
	id, err := this.GetInt("articleId")
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	// 处理文件上传
	filePath := UploadFile(&this.Controller, "uploadname")

	// 数据校验
	if err != nil || articleName == "" || content == "" || filePath == "" {
		println("请求错误")
		return
	}
	// 数据处理
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article)
	if err != nil {
		println("更新的文章不存在")
		return
	}
	article.ArtiName = articleName
	article.Acontent = content
	if filePath != "NoImg" {
		article.Aimg = filePath
	}
	o.Update(&article)

	// 返回视图
	this.Redirect("/article/showArticleList", 302)

}

// 删除文章
func (this *ArticleController) DeleteArticle() {
	// 获取数据
	id, err := this.GetInt("articleId")

	// 校验数据
	if err != nil {
		println("删除文章请求路径错误")
		return
	}

	// 操作数据
	// 删除操作
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Delete(&article)

	// 返回视图
	this.Redirect("/articl/showArticleList", 302)

}

// 展示添加分类页面
func (this *ArticleController) ShowAddType() {
	// 获取数据
	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)

	//传递数据
	userName := this.GetSession("userName")
	this.Data["title"] = "添加分类"
	this.Data["userName"] = userName.(string)
	this.Data["types"] = types
	this.Layout = "layout.html"
	this.TplName = "addType.html"
}

// 处理添加类型数据
func (this *ArticleController) HandleAddType() {

	// 获取数据
	typeName := this.GetString("typeName")
	// 校验数据
	if typeName == "" {
		println("信息不完整，请重新输入")
		return
	}
	// 处理数据
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Insert(&articleType)

	// 返回视图
	this.Redirect("/articl/addType", 302)

}
