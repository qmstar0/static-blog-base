package main

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/log"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

const (
	TemplateIndex      = "index.html"
	TemplatePostItem   = "post-item.html"
	TemplatePostDetail = "post-detail.html"
	TemplatePagination = "pagination.html"
)

type Data struct {
	Posts      []Post
	Categories []Category
	Tags       []Tag
	Recent     []Post
}

type BuilderConfig struct {
	ContentDir  string `yaml:"content_dir" toml:"content_dir" mapstructure:"content_dir"`
	OutputDir   string `yaml:"output_dir" toml:"output_dir" mapstructure:"output_dir"`
	TemplateDir string `yaml:"template_dir" toml:"template_dir" mapstructure:"template_dir"`
	AssetsDir   string `yaml:"assets_dir" toml:"assets_dir" mapstructure:"assets_dir"`
	AutoClear   bool   `yaml:"auto_clear" toml:"auto_clear" mapstructure:"auto_clear"`
}

type Builder struct {
	Config BuilderConfig

	Data Data

	template *template.Template
}

func NewBuilder(config BuilderConfig) *Builder {
	matches, err := filepath.Glob(filepath.Join(config.ContentDir, "*.md"))
	if err != nil {
		log.Fatal("在读取ContentDir时发生错误", "err", err)
	}
	var posts = make(Posts, 0, len(matches))
	for _, match := range matches {

		doc := MarkdownToHTML(match)

		posts = append(posts, Post{
			Title:     doc.Title,
			Uri:       doc.Uri,
			Category:  doc.Category,
			Tags:      doc.Tags,
			CreatedAt: doc.CreatedAt,
			Desc:      template.HTML(doc.Desc),
			Content:   template.HTML(doc.Content),
		})
	}

	glob, err := template.New("builder").
		Funcs(template.FuncMap{
			"timestamp": TimeStamp,
		}).
		ParseGlob(filepath.Join(config.TemplateDir, "*.html"))
	if err != nil {
		log.Fatal("加载template错误", "err", err)
	}

	Sort(posts)

	return &Builder{
		Config: config,
		Data: Data{
			Posts:      posts,
			Categories: posts.Categories(),
			Tags:       posts.Tags(),
			Recent:     posts.Recent(),
		},
		template: glob,
	}
}

func (b *Builder) Build() {
	if b.Config.AutoClear {
		ClearDir(b.Config.OutputDir)
	}
	if b.Config.AssetsDir != "" {
		b.CopyAssets(b.Config.AssetsDir, b.Config.OutputDir)
	}

	b.BuildIndex(".")
	b.BuildPages(".", b.Data.Posts, 13)

	var wg sync.WaitGroup
	for _, category := range b.Data.Categories {
		wg.Add(1)
		go func(c Category) {
			defer wg.Done()
			path := filepath.Join("categories", c.Name)
			log.Info("正在构建：", "target", path)
			b.BuildIndex(path)
			b.BuildPages(path, c.Posts, 13)
		}(category)
	}

	for _, tag := range b.Data.Tags {
		wg.Add(1)
		go func(t Tag) {
			defer wg.Done()
			path := filepath.Join("tags", t.Name)
			log.Info("正在构建：", "target", path)
			b.BuildIndex(path)
			b.BuildPages(path, t.Posts, 13)
		}(tag)
	}

	b.BuildPosts("post")
	wg.Wait()
}

func (b *Builder) BuildPosts(dir string) {

	var wg sync.WaitGroup
	for _, post := range b.Data.Posts {
		wg.Add(1)
		go func(p Post) {
			defer wg.Done()
			path := filepath.Join(b.Config.OutputDir, dir, p.Uri, "index.html")
			CreateDir(path)

			file, err := os.Create(path)
			if err != nil {
				log.Fatal("创建文件时", "err", err)
			}
			defer file.Close()

			err = b.template.ExecuteTemplate(file, TemplatePostDetail, map[string]any{
				"Post":  p,
				"Aside": b.Data,
			})
			if err != nil {
				log.Fatal("创建文件时", "err", err)
			}
		}(post)
	}
	wg.Wait()
}

func (b *Builder) BuildIndex(dir string) {
	path := filepath.Join(b.Config.OutputDir, dir, "index.html")
	CreateDir(path)
	file, err := os.Create(path)
	if err != nil {
		log.Fatal("创建文件时", "err", err)
	}
	defer file.Close()
	err = b.template.ExecuteTemplate(file, TemplateIndex, b.Data)
	if err != nil {
		log.Fatal("渲染html时", "err", err)
	}
}

func (b *Builder) BuildPages(dir string, posts []Post, perPageNum int) {
	totalPages := (len(posts) + perPageNum - 1) / perPageNum
	var (
		buffers = make([]bytes.Buffer, totalPages)
	)

	for i, post := range posts {
		page := i / perPageNum
		err := b.template.ExecuteTemplate(&buffers[page], TemplatePostItem, post)
		if err != nil {
			log.Fatal("渲染html时", "err", err)
		}
	}

	for i, buffer := range buffers {
		path := filepath.Join(b.Config.OutputDir, dir, fmt.Sprintf("page_%d.html", i+1))
		CreateDir(path)

		var data = make(map[string]any)
		SetupPaginationData(data, totalPages, i+1)
		err := b.template.ExecuteTemplate(&buffer, TemplatePagination, data)
		if err != nil {
			log.Fatal("渲染html时", "err", err)
		}

		if err := os.WriteFile(
			path,
			buffer.Bytes(),
			0755,
		); err != nil {
			log.Fatal("向文件写入数据时", "err", err)
		}
	}
}

func (b *Builder) CopyAssets(from, to string) {
	var err error
	var output []byte
	switch goos := runtime.GOOS; goos {
	case "windows":
		output, err = exec.Command("xcopy", from, filepath.Join(to, from), "/E", "/I", "/H").Output()
	case "linux", "darwin":
		command := exec.Command("cp", "-r", "-v", from, filepath.Join(to, from))
		log.Info(command.String())
		output, err = command.Output()
	}
	log.Info(string(output))
	if err != nil {
		log.Fatal("拷贝静态资源时", "err", err)
	}
}
