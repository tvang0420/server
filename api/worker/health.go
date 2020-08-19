// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/workers/{worker} workers WorkerHealth
//
// Check if the worker's API is available
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: worker
//   description: worker hostname to get the health of
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully 'ping'-ed the Vela worker's API
//     schema:
//       type: string

// Health represents the API handler to
// report the health status for a Vela worker.
func Health(c *gin.Context) {
	workerURL := "http://worker_%s:8080/health"
	response, err := http.Get(fmt.Sprintf(workerURL, c.Param("worker")))
	if err != nil {
		retErr := fmt.Errorf("unable to check health of %s: %w", workerURL, err)
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		retErr := fmt.Errorf("unable to check health of %s: %w", workerURL, err)
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}
	c.JSON(response.StatusCode, strings.Trim(string(respBody), "\""))
}
