// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"
)

func TestPostgres_NewMap(t *testing.T) {
	// setup types
	want := &Map{
		BuildService: &Service{
			List: map[string]string{
				"all":  ListBuilds,
				"repo": ListRepoBuilds,
			},
			Select: map[string]string{
				"repo":          SelectRepoBuild,
				"last":          SelectLastRepoBuild,
				"count":         SelectBuildsCount,
				"countByStatus": SelectBuildsCountByStatus,
				"countByRepo":   SelectRepoBuildCount,
			},
			Delete: DeleteBuild,
		},
		LogService: &Service{
			List: map[string]string{
				"all":   ListLogs,
				"build": ListBuildLogs,
			},
			Select: map[string]string{
				"step":    SelectStepLog,
				"service": SelectServiceLog,
			},
			Delete: DeleteLog,
		},
		RepoService: &Service{
			List: map[string]string{
				"all":  ListRepos,
				"user": ListUserRepos,
			},
			Select: map[string]string{
				"repo":        SelectRepo,
				"count":       SelectReposCount,
				"countByUser": SelectUserReposCount,
			},
			Delete: DeleteRepo,
		},
		SecretService: &Service{
			List: map[string]string{
				"all":    ListSecrets,
				"org":    ListOrgSecrets,
				"repo":   ListRepoSecrets,
				"shared": ListSharedSecrets,
			},
			Select: map[string]string{
				"org":         SelectOrgSecret,
				"repo":        SelectRepoSecret,
				"shared":      SelectSharedSecret,
				"countOrg":    SelectOrgSecretsCount,
				"countRepo":   SelectRepoSecretsCount,
				"countShared": SelectSharedSecretsCount,
			},
			Delete: DeleteSecret,
		},
		ServiceService: &Service{
			List: map[string]string{
				"all":   ListServices,
				"build": ListBuildServices,
			},
			Select: map[string]string{
				"build": SelectBuildService,
				"count": SelectBuildServicesCount,
			},
			Delete: DeleteService,
		},
		StepService: &Service{
			List: map[string]string{
				"all":   ListSteps,
				"build": ListBuildSteps,
			},
			Select: map[string]string{
				"build": SelectBuildStep,
				"count": SelectBuildStepsCount,
			},
			Delete: DeleteStep,
		},
		UserService: &Service{
			List: map[string]string{
				"all":  ListUsers,
				"lite": ListLiteUsers,
			},
			Select: map[string]string{
				"user":  SelectUser,
				"name":  SelectUserName,
				"count": SelectUsersCount,
			},
			Delete: DeleteUser,
		},
	}

	// run test
	got := NewMap()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewMap is %v, want %v", got, want)
	}
}
