// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/go-vela/server/database"
	"github.com/gin-gonic/gin"
)

// Database is a middleware function that initializes the database and
// attaches to the context of every http.Request.
func Database(d database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		database.ToContext(c, d)
		c.Next()
	}
}
