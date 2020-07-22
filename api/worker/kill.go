// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/worker/executors"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/workers/{worker} workers Kill
//
// Check if the worker's API is available
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// responses:
//   '200':
//     description: Successfully 'ping'-ed Vela worker's API
//     schema:
//       type: string

// Kill represents the API handler to
// report the health status for a Vela worker.
func Kill(c *gin.Context) {

	r := repo.Retrieve(c)
	b := build.Retrieve(c)
	e := executors.Retrieve(c)

	if !strings.EqualFold(b.GetStatus(), "running") {
		retErr := fmt.Errorf("Found build %s/%d but its status was %s", r.GetFullName(), b.GetNumber(), b.GetStatus())
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}

	for _, executor := range e {
		// check each executor to see if it's running the given build
		if strings.EqualFold(executor.Repo.GetFullName(), r.GetFullName()) &&
			*executor.GetBuild().ID == int64(b.GetNumber()) {
			//proceed with killing the build
			client := &http.Client{}
			endpoint := fmt.Sprintf("%s/api/v1/executors/%d/build/kill", b.GetHost(), executor.GetID()+1)
			req, err := http.NewRequest("DELETE", endpoint, nil)
			if err != nil {
				retErr := fmt.Errorf("unable to form request to %s: %w", endpoint, err)
				util.HandleError(c, http.StatusBadRequest, retErr)
				return
			}
			// add the token to authenticate to the worker
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("VELA_SECRET")))

			// make the request to the worker
			resp, err := client.Do(req)
			if err != nil {
				retErr := fmt.Errorf("unable to connect to %s: %w", endpoint, err)
				util.HandleError(c, http.StatusBadRequest, retErr)
				return
			}
			defer resp.Body.Close()

			// Read Response Body
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				retErr := fmt.Errorf("unable to read response from %s: %w", endpoint, err)
				util.HandleError(c, http.StatusBadRequest, retErr)
				return
			}
			c.JSON(resp.StatusCode, strings.Trim(string(respBody), "\""))
			//c.JSON(http.StatusOK, fmt.Sprintf("Deleted build %s/%d", e.Repo.GetFullName(), e.Build.GetID()))
			return
		}
	}

	retErr := fmt.Errorf("Unable to find a build for %s/%d on worker %s", r.GetFullName(), b, b.GetHost())
	util.HandleError(c, http.StatusInternalServerError, retErr)
}

// strings.EqualFold(e.Build.GetStatus(), "running")
