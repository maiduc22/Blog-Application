package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Post struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"text;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Post) Prepare() {
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}
func (p *Post) Validate() error {
	if p.Title == "" {
		return errors.New("Post title can not empty")
	}
	if p.Content == "" {
		return errors.New("Post content can not empty")
	}
	if p.AuthorID < 1 {
		return errors.New("Require author")
	}
	return nil
}
func (p *Post) SavePost(db *gorm.DB) (*Post, error) {
	err := db.Model(&Post{}).Create(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

func (p *Post) FindAllPosts(db *gorm.DB) (*[]Post, error) {
	posts := []Post{}
	err := db.Model(&Post{}).Limit(100).Find(&Post{}).Error
	if err != nil {
		return &[]Post{}, err
	}
	if len(posts) > 0 {
		for _, post := range posts {
			err = db.Model(&User{}).Where("id = ?", post.AuthorID).Take(&post.Author).Error
			if err != nil {
				return &[]Post{}, err
			}
		}
	}
	return &posts, nil
}

func (p *Post) FindPostByID(db *gorm.DB, postid uint64) (*Post, error) {
	err := db.Model(&Post{}).Where("id = ?", postid).Take(&Post{}).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

func (p *Post) UpdatePost(db *gorm.DB) (*Post, error) {
	err := db.Model(&Post{}).Where("id = ?", p.ID).Take(&Post{}).UpdateColumns(
		map[string]interface{}{
			"title":     p.Title,
			"content":   p.Content,
			"update_at": time.Now(),
		},
	).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

func (p *Post) DeletePost(db *gorm.DB) (int64, error) {
	result := db.Model(&Post{}).Where("id = ?", p.ID).Take(&Post{}).Delete(&Post{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (p *Post) GetUsersPosts(db *gorm.DB, uid uint32) (*[]Post, error) {
	posts := []Post{}
	err := db.Model(&Post{}).Where("author_id = ?", uid).Limit(100).Order("created_at").Find(&Post{}).Error
	if err != nil {
		return &[]Post{}, err
	}
	if len(posts) > 0 {
		for _, post := range posts {
			err = db.Model(&User{}).Where("id = ?", post.AuthorID).Take(&User{}).Error
			if err != nil {
				return &[]Post{}, err
			}
		}
	}
	return &posts, nil
}

// When delete an user, need to delete all user's posts
func (p *Post) DeleteAllUserPosts(db *gorm.DB, uid int32) (int64, error) {
	result := db.Model(&Post{}).Where("author_id = ?", uid).Delete(&Post{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
