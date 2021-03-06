// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resolve

import (
	"context"
	"fmt"

	"github.com/google/gapid/gapis/database"
	"github.com/google/gapid/gapis/gfxapi"
	"github.com/google/gapid/gapis/service/path"
)

// ResourceMeta returns the metadata for the specified resource.
func ResourceMeta(ctx context.Context, id *path.ID, after *path.Command) (*gfxapi.ResourceMeta, error) {
	obj, err := database.Build(ctx, &ResourceMetaResolvable{id, after})
	if err != nil {
		return nil, err
	}
	return obj.(*gfxapi.ResourceMeta), nil
}

// Resolve implements the database.Resolver interface.
func (r *ResourceMetaResolvable) Resolve(ctx context.Context) (interface{}, error) {
	resources, err := database.Build(ctx, &AllResourceDataResolvable{r.After})
	if err != nil {
		return nil, err
	}
	res, ok := resources.(*ResolvedResources)
	if !ok {
		return nil, fmt.Errorf("Cannot resolve resources at command: %v", r.After)
	}
	id := r.Id.ID()
	val, ok := res.resources[id]
	if !ok {
		return nil, fmt.Errorf("Could not find resource %v", id)
	}
	result := &gfxapi.ResourceMeta{
		IDMap:    res.resourceMap,
		Resource: val,
	}
	return result, nil
}
