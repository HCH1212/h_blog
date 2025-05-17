package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

// 加载所有文章
func loadPosts() ([]Post, error) {
	var posts []Post

	// 确保content目录存在
	if _, err := os.Stat("content"); os.IsNotExist(err) {
		return nil, fmt.Errorf("content directory does not exist")
	}

	files, err := ioutil.ReadDir("content")
	if err != nil {
		return nil, fmt.Errorf("error reading content directory: %v", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no markdown files found in content directory")
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {
			slug := strings.TrimSuffix(file.Name(), ".md")
			post, err := loadPost(slug)
			if err != nil {
				// 记录错误但继续处理其他文件
				fmt.Printf("Error loading post %s: %v\n", slug, err)
				continue
			}
			posts = append(posts, post)
		}
	}

	if len(posts) == 0 {
		return nil, fmt.Errorf("no valid markdown files found")
	}

	// 使用 sort.Slice 反转顺序
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date > posts[j].Date
	})

	return posts, nil
}

// 加载单个文章
func loadPost(slug string) (Post, error) {
	filePath := filepath.Join("content", slug+".md")
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Post{}, err
	}

	frontMatter, markdownContent, err := parseFrontMatter(content)
	if err != nil {
		return Post{}, err
	}

	// 设置默认值
	post := Post{
		Title: slug, // 默认使用slug作为标题
		Slug:  slug,
	}

	// 解析Front Matter
	if frontMatter != nil {
		if title, ok := frontMatter["title"].(string); ok {
			post.Title = title
		}
		if date, ok := frontMatter["date"].(string); ok {
			post.Date = date
		}
		if categories, ok := frontMatter["categories"].([]interface{}); ok {
			for _, c := range categories {
				if cat, ok := c.(string); ok {
					post.Categories = append(post.Categories, cat)
				}
			}
		}
	}

	// 如果没有在Front Matter中找到标题，尝试从第一行获取
	if post.Title == slug {
		lines := strings.Split(string(markdownContent), "\n")
		title := strings.TrimSpace(lines[0])
		if strings.HasPrefix(title, "# ") {
			post.Title = strings.TrimPrefix(title, "# ")
		}
	}

	post.Content = string(markdownContent)

	return post, nil
}

// 通过分类加载文章
func loadPostsByCategory(category string) ([]Post, error) {
	allPosts, err := loadPosts()
	if err != nil {
		return nil, err
	}

	var filteredPosts []Post
	for _, post := range allPosts {
		for _, cat := range post.Categories {
			if cat == category {
				filteredPosts = append(filteredPosts, post)
				break
			}
		}
	}

	return filteredPosts, nil
}

// 获取分类
func getAllCategories(posts []Post) []Category {
	categoryMap := make(map[string]int)

	for _, post := range posts {
		for _, cat := range post.Categories {
			categoryMap[cat]++
		}
	}

	var categories []Category
	for name, count := range categoryMap {
		categories = append(categories, Category{
			Name:  name,
			Count: count,
		})
	}

	// 按名称排序
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Name < categories[j].Name
	})

	return categories
}

// 解析md头部信息
func parseFrontMatter(content []byte) (map[string]interface{}, []byte, error) {
	re := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)$`)
	matches := re.FindSubmatch(content)

	if len(matches) < 3 {
		return nil, content, nil
	}

	var frontMatter map[string]interface{}
	if err := yaml.Unmarshal(matches[1], &frontMatter); err != nil {
		return nil, nil, err
	}

	return frontMatter, matches[2], nil
}
