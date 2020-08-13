// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executors

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/go-vela/types/library"

	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/util"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Retrieve gets the repo in the given context
func Retrieve(c *gin.Context) []library.Executor {
	return FromContext(c)
}

// Establish sets the repo in the given context
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {

		b := build.Retrieve(c)

		// prepare the request to the worker
		client := &http.Client{}
		endpoint := fmt.Sprintf("%s/api/v1/executors", b.GetHost())
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			retErr := fmt.Errorf("unable to form request to %s: %w", endpoint, err)
			util.HandleError(c, http.StatusBadRequest, retErr)
			return
		}
		// add the token to authenticate to the worker to the headers
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

		e := new([]library.Executor)

		err = json.Unmarshal(respBody, e)
		if err != nil {
			retErr := fmt.Errorf("unable to parse response from %s: %w", endpoint, err)
			util.HandleError(c, http.StatusBadRequest, retErr)
			return
		}

		ToContext(c, *e)
		c.Next()
	}
}
