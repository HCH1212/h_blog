package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

type Post struct {
	Title      string
	Content    template.HTML
	Slug       string
	Date       string
	Categories []string
}

type Category struct {
	Name  string
	Count int
}

func main() {
	r := gin.Default()

	// 加载模板
	r.LoadHTMLGlob("template/*")

	// 静态文件
	r.Static("/static", "./static")

	// 首页路由
	r.GET("/", func(c *gin.Context) {
		posts, err := loadPosts()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"posts":         posts,
			"allCategories": getAllCategories(posts),
		})
	})

	// 文章详情路由
	r.GET("/post/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		post, err := loadPost(slug)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.HTML(http.StatusOK, "post.html", gin.H{
			"post": post,
		})
	})

	// 分类路由
	r.GET("/category/:name", func(c *gin.Context) {
		name := c.Param("name")
		posts, err := loadPostsByCategory(name)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"posts":         posts,
			"category":      name,
			"allCategories": getAllCategories(posts),
		})
	})

	// 启动服务器
	r.Run(":8888")
}
