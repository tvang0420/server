// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"
)

func TestPostgres_createRepoService(t *testing.T) {
	// setup types
	want := &Service{
		List: map[string]string{
			"all":  ListRepos,
			"user": ListUserRepos,
			"org":  ListOrgRepos,
		},
		Select: map[string]string{
			"repo":        SelectRepo,
			"count":       SelectReposCount,
			"countByUser": SelectUserReposCount,
		},
		Delete: DeleteRepo,
	}

	// run test
	got := createRepoService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createRepoService is %v, want %v", got, want)
	}
}
