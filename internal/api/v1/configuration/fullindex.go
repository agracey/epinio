// Copyright © 2021 - 2023 SUSE LLC
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configuration

import (
	"github.com/epinio/epinio/helpers/kubernetes"
	"github.com/epinio/epinio/internal/api/v1/response"
	"github.com/epinio/epinio/internal/application"
	"github.com/epinio/epinio/internal/auth"
	"github.com/epinio/epinio/internal/cli/server/requestctx"
	"github.com/epinio/epinio/internal/configurations"
	apierror "github.com/epinio/epinio/pkg/api/core/v1/errors"

	"github.com/gin-gonic/gin"
)

// FullIndex handles the API endpoint GET /configurations
// It lists all the known applications in all namespaces, with and without workload.
func (hc Controller) FullIndex(c *gin.Context) apierror.APIErrors {
	ctx := c.Request.Context()
	user := requestctx.User(ctx)

	cluster, err := kubernetes.GetCluster(ctx)
	if err != nil {
		return apierror.InternalError(err)
	}

	allConfigurations, err := configurations.List(ctx, cluster, "")
	if err != nil {
		return apierror.InternalError(err)
	}
	filteredConfigurations := auth.FilterResources(user, allConfigurations)

	appsOf, err := application.BoundAppsNames(ctx, cluster, "")
	if err != nil {
		return apierror.InternalError(err)
	}

	responseData, err := makeResponse(ctx, appsOf, filteredConfigurations)
	if err != nil {
		return apierror.InternalError(err)
	}

	response.OKReturn(c, responseData)
	return nil
}
