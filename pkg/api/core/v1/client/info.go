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

package client

import (
	"encoding/json"

	api "github.com/epinio/epinio/internal/api/v1"
	"github.com/epinio/epinio/pkg/api/core/v1/models"
)

// Info returns information about Epinio and its components
func (c *Client) Info() (models.InfoResponse, error) {
	var resp models.InfoResponse

	data, err := c.get(api.Routes.Path("Info"))
	if err != nil {
		return resp, err
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, err
	}

	c.log.V(1).Info("response decoded", "response", resp)

	return resp, nil
}
