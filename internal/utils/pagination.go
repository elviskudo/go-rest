package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func Search(c *gin.Context, fields []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := c.Query("search")
		if query == "" {
			return db
		}

		sql := ""
		args := []interface{}{}
		for i, field := range fields {
			if i > 0 {
				sql += " OR "
			}
			sql += field + " LIKE ?"
			args = append(args, "%"+query+"%")
		}
		return db.Where(sql, args...)
	}
}

func Sort(c *gin.Context, allowedFields map[string]bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		sort := c.Query("sort")
		order := c.DefaultQuery("order", "asc")

		if sort == "" {
			return db
		}

		if allowedFields[sort] {
			if order != "asc" && order != "desc" {
				order = "asc"
			}
			return db.Order(sort + " " + order)
		}

		return db
	}
}
