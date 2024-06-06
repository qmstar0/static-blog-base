package main

import "html/template"

type Post struct {
	Title     string
	Uri       string
	Category  string
	Tags      []string
	CreatedAt int64
	Desc      template.HTML
	Content   template.HTML
}

type Posts []Post

func (c Posts) Len() int {
	return len(c)
}

func (c Posts) Less(i, j int) bool {
	return c[i].CreatedAt > c[j].CreatedAt
}

func (c Posts) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type Category struct {
	Name  string `json:"name"`
	Posts Posts  `json:"contents"`
}

type Tag struct {
	Name  string `json:"name"`
	Posts Posts  `json:"contents"`
}

func (c Posts) Categories() []Category {
	var categoryFilter = make(map[string]Posts)
	for _, content := range c {
		if content.Category == "" {
			break
		}
		categoryFilter[content.Category] = append(categoryFilter[content.Category], content)
	}

	var result = make([]Category, 0, len(categoryFilter))
	for s, contents := range categoryFilter {
		result = append(result, Category{
			Name:  s,
			Posts: contents,
		})
	}
	return result
}

func (c Posts) Tags() []Tag {
	var tagFilter = make(map[string]Posts)
	for _, content := range c {
		for _, t := range content.Tags {
			if t == "" {
				continue
			}
			tagFilter[t] = append(tagFilter[t], content)
		}
	}
	var result = make([]Tag, 0, len(tagFilter))
	for s, contents := range tagFilter {
		result = append(result, Tag{
			Name:  s,
			Posts: contents,
		})
	}
	return result
}

func (c Posts) Recent() Posts {
	var result = make(Posts, 0, 5)

	for i, post := range c {
		if i >= 5 {
			break
		}
		result = append(result, post)
	}

	return result
}
