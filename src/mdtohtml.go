package main

//
import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/charmbracelet/log"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
	"os"
	"path/filepath"
	"regexp"
)

var uriCheck, _ = regexp.Compile("^[a-zA-Z0-9\\-_]+$")

type MarkdownMeta struct {
	Title     string   `json:"title" toml:"title" yaml:"title"`
	Uri       string   `json:"uri" toml:"uri" yaml:"uri"`
	Category  string   `json:"category" toml:"category" yaml:"category"`
	Tags      []string `json:"tags" toml:"tags" yaml:"tags"`
	CreatedAt int64    `json:"createdAt" toml:"createdAt" yaml:"createdAt"`
	Desc      string   `json:"desc" toml:"desc" yaml:"desc"`
}

type MarkdownDoc struct {
	MarkdownMeta
	Content string `json:"content" toml:"content" yaml:"content"`
}

func (p MarkdownMeta) Verify() error {
	if p.Desc == "" {
		return errors.New("desc 为空")
	}
	if p.CreatedAt == 0 {
		return errors.New("createdAt 为空")
	}
	if p.Title == "" {
		return errors.New("title 为空")
	}
	return p.VerifyUri()
}

func (p MarkdownMeta) VerifyUri() error {
	if p.Uri == "" {
		return errors.New("uri 为空")
	}
	if !uriCheck.MatchString(p.Uri) {
		return errors.New("uri 格式错误")
	}
	return nil
}

func (p Post) String() (string, error) {
	marshal, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

var defaultMarkdownToHTMLTool = goldmark.New(
	goldmark.WithExtensions(
		&frontmatter.Extender{},
		extension.GFM,
		extension.NewFootnote(
			extension.WithFootnoteIDPrefix("footnote-"),
		),
		highlighting.NewHighlighting(
			highlighting.WithStyle("github"),
		),
		extension.CJK,
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
		parser.WithAttribute(),
	),

	goldmark.WithRendererOptions(
		html.WithXHTML(),
		html.WithHardWraps(),
	),
)

func MarkdownToHTML(path string) MarkdownDoc {
	var (
		buf  bytes.Buffer
		doc  MarkdownDoc
		meta MarkdownMeta
	)

	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("读取文件错误", "err", err)
	}

	ctx := parser.NewContext()

	if err = defaultMarkdownToHTMLTool.Convert(file, &buf, parser.WithContext(ctx)); err != nil {
		log.Fatal("markdown渲染为html时", "err", err)
	}
	doc.Content = buf.String()

	metadata := frontmatter.Get(ctx)
	if metadata != nil {
		if err = metadata.Decode(&meta); err != nil {
			log.Fatal("解码Markdown元数据时", "err", err)
		}
	} else {
		log.Fatal("解析Markdown元数据时", "err", err)
	}
	if err = meta.Verify(); err != nil {
		log.Fatal("验证数据时", "filepath", path, "filename", filepath.Base(path), "err", err)
	}
	buf.Reset()

	err = defaultMarkdownToHTMLTool.Convert([]byte(meta.Desc), &buf)
	if err != nil {
		log.Fatal("将desc渲染为html时", "err", err)
	}
	doc.MarkdownMeta = meta
	doc.Desc = buf.String()
	return doc
}
