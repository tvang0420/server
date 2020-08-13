// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/worker"
)

// WorkerHandlers is a function that extends the provided base router group
// with the API handlers proxying worker endpoints.
//
// GET       /api/v1/workers/:worker/health
// DELETE    /api/v1/workers/:worker/:org/:repo/:buildNum
func WorkerHandlers(base *gin.RouterGroup) {

	// Worker endpoints
	_workers := base.Group("/workers/:worker")
	{
		_workers.GET("/health", worker.Health)
	} // end of Worker endpoints
}
