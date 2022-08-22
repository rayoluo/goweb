package routes

import (
	v1 "ginblog/api/v1"
	"ginblog/middleware"
	"ginblog/model"
	"ginblog/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// func createMyRender() multitemplate.Renderer {
// 	p := multitemplate.NewRenderer()
// 	p.AddFromFiles("admin", "web/admin/dist/index.html")
// 	p.AddFromFiles("front", "web/front/dist/index.html")
// 	return p
// }

func InitRouter() {
	gin.SetMode(utils.AppMode)
	r := gin.New()
	// r.HTMLRender = createMyRender()
	r.Use(middleware.Log())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	// r.Static("/static", "./web/front/dist/static")
	// r.Static("/admin", "./web/admin/dist")
	// r.StaticFile("/favicon.ico", "/web/front/dist/favicon.ico")
	//
	// r.GET("/", func(c *gin.Context) {
	// 	c.HTML(200, "front", nil)
	// })
	//
	// r.GET("/admin", func(c *gin.Context) {
	// 	c.HTML(200, "admin", nil)
	// })

	/*
		后台管理路由接口
	*/
	auth := r.Group("api/v1")
	auth.Use(middleware.JwtToken())
	{
		// 用户模块的路由接口 已看
		auth.GET("admin/users", v1.GetUsers)
		auth.PUT("user/:id", v1.EditUser)
		auth.DELETE("user/:id", v1.DeleteUser)
		//修改密码
		auth.PUT("admin/changepw/:id", v1.ChangeUserPassword)
		// 分类模块的路由接口 已看
		auth.GET("admin/category", v1.GetCate)
		auth.POST("category/add", v1.AddCategory)
		auth.PUT("category/:id", v1.EditCate)
		auth.DELETE("category/:id", v1.DeleteCate)
		// 标签模块的路由接口
		auth.GET("admin/tag", v1.GetTagList)
		auth.POST("tag/add", v1.AddTag)
		auth.PUT("tag/:id", v1.EditTag)
		auth.DELETE("tag/:id", v1.DeleteTag)
		auth.GET("admin/tags", v1.GetAllTags)

		// 文章标签模块的路由接口
		// url中的id参数指定文章id
		auth.POST("articleTag/:id", v1.AddArtTags)
		// 根据文章id参数删除该文章的所有标签
		auth.DELETE("articleTag/art/:id", v1.DeleteByArtID)
		// 根据标签id删除该标签对应的所有文章
		auth.DELETE("articleTag/tag/:id", v1.DeleteByTagID)
		// 根据文章ID获取所有的标签
		auth.GET("articleTag/art/:id", v1.GetTagsByArtID)
		// 根据标签ID获取所有的文章
		auth.GET("articleTag/tag/:id", v1.GetArtsByTagID)

		// 文章模块的路由接口
		auth.GET("admin/article/info/:id", v1.GetArtInfoEdit)
		auth.GET("admin/article", v1.GetArtsWithTags)
		auth.POST("article/add", v1.AddArticle)
		auth.PUT("article/:id", v1.EditArt)
		auth.DELETE("article/:id", v1.DeleteArt)
		// 上传文件
		auth.POST("upload", v1.UpLoad)
		// 更新个人设置
		auth.GET("admin/profile/:id", v1.GetProfile)
		auth.PUT("profile/:id", v1.UpdateProfile)
		// 评论模块
		auth.GET("comment/list", v1.GetCommentList)
		auth.DELETE("delcomment/:id", v1.DeleteComment)
		auth.PUT("checkcomment/:id", v1.CheckComment)
		auth.PUT("uncheckcomment/:id", v1.UnCheckcomment)
	}

	/*
		前端展示页面接口
	*/
	router := r.Group("api/v1")

	// test
	store := model.GetSessionStore()
	store.Options(sessions.Options{MaxAge: 3600})
	router.Use(sessions.Sessions("sessionid", store))
	{
		// 用户信息模块 这个模块的接口没有完全用到
		// 生成验证码并发送到邮箱
		router.GET("users/captcha", v1.GenCaptcha)
		// 前端注册用户
		router.POST("user/add", v1.AddUser)
		// 前端修改密码
		router.PUT("users/password", v1.ResetPassword)
		// 前端修改个人信息
		router.PUT("users/update/info/:id", v1.UpdateUserInfo)

		router.GET("user/:id", v1.GetUserInfo)
		router.GET("users", v1.GetUsers)

		// 文章分类信息模块
		router.GET("category", v1.GetCateAll)
		router.GET("category/:id", v1.GetCateInfo)

		// 文章标签模块
		// 获取所有标签
		router.GET("tag", v1.GetAllTags)
		router.GET("tag/:id", v1.GetTagInfo)
		// 根据文章id获取标签
		router.GET("tags/:id", v1.GetTagsByArtID)

		// 文章模块
		router.GET("article", v1.GetArtsWithTags)
		router.GET("article/list/:id", v1.GetCateArtTags)
		router.GET("article/info/:id", v1.GetArtInfo)
		router.GET("article/all", v1.GetAllArt)
		// 根据文章id获取推荐文章
		router.GET("article/recommend/:id", v1.GetRecArts)
		// 根据标签ID获取所有的文章
		router.GET("article/listByTag/:id", v1.GetArtsByTagID)

		// 登录控制模块
		router.POST("loginfront", v1.LoginFront)
		// 前端退出登录
		router.GET("logout", v1.Logout)
		// 后端登录接口
		router.POST("login", v1.Login)

		// 前端头像上传
		router.POST("users/avatar/upload", v1.UpLoad)

		// 获取个人设置信息
		router.GET("profile/:id", v1.GetProfile)

		// 评论模块 这个模块的接口还没有适配
		// 新增评论id是登录用户的id
		router.POST("addcomment/:id", v1.AddComment)
		// 前端获取评论
		router.GET("comments", v1.GetCommentListFront)
		// 前端获取指定ID的父级评论下的回复，分页显示
		router.GET("comments/replies/:commentId", v1.GetCommentReplies)
		router.GET("comment/info/:id", v1.GetComment)
		// router.GET("commentfront/:id", v1.GetCommentListFront)

		router.GET("commentcount/:id", v1.GetCommentCount)

		// 点赞模块
		router.POST("articles/like/:artId", v1.SaveArticleLike)
		router.POST("comments/like/:commentId", v1.SaveCommentLike)
		router.GET("articles/likeCount/:artId", v1.GetArtLikeCount)
	}

	// session控制会话
	// sess := r.Group("api/v1")
	// store := model.GetSessionStore()
	// store.Options(sessions.Options{MaxAge: 3600, HttpOnly: true})
	// sess.Use(sessions.Sessions("mysession", store))
	// {
	// 	// 前端登录接口
	// 	sess.POST("loginfront", v1.LoginFront)
	// 	// 前端头像上传
	// 	sess.POST("users/avatar/upload/:id", v1.UpLoadFront)
	//
	// }

	_ = r.Run(utils.HttpPort)

}
