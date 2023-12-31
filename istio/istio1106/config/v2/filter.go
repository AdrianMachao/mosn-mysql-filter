/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v2

import (
	envoy_extensions_filters_http_rbac_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/rbac/v3"
	"mosn.io/mosn/istio/istio1106/mixer/v1/config/client"
)

type RBACConfig struct {
	envoy_extensions_filters_http_rbac_v3.RBAC
	Version string `json:"version"`
}

const (
	IstioStats = "istio.stats"
	MIXER      = "mixer"
	JwtAuthn   = "jwt_authn"
	RBAC       = "rbac"
)

type Mixer struct {
	client.HttpClientConfig
}
