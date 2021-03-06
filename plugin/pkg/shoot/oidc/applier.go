// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package oidc

import (
	"github.com/gardener/gardener/pkg/apis/garden"
	settingsv1alpha1 "github.com/gardener/gardener/pkg/apis/settings/v1alpha1"
	"github.com/gardener/gardener/pkg/utils"
)

// ApplyOIDCConfiguration applies preset OpenID Connect configuration to the shoot.
func ApplyOIDCConfiguration(shoot *garden.Shoot, preset *settingsv1alpha1.OpenIDConnectPresetSpec) {
	if shoot == nil || preset == nil {
		return
	}
	useRequiredClaims, err := utils.CheckVersionMeetsConstraint(shoot.Spec.Kubernetes.Version, ">= 1.11")
	if err != nil {
		// Don't mutate the resource anymore, because the version is invalid
		// and it'll be caught by validation.
		return
	}

	var client *garden.OpenIDConnectClientAuthentication
	if preset.Client != nil {
		client = &garden.OpenIDConnectClientAuthentication{
			Secret:      preset.Client.Secret,
			ExtraConfig: preset.Client.ExtraConfig,
		}
	}
	oidc := &garden.OIDCConfig{
		CABundle:             preset.Server.CABundle,
		ClientID:             &preset.Server.ClientID,
		GroupsClaim:          preset.Server.GroupsClaim,
		GroupsPrefix:         preset.Server.GroupsPrefix,
		IssuerURL:            &preset.Server.IssuerURL,
		SigningAlgs:          preset.Server.SigningAlgs,
		UsernameClaim:        preset.Server.UsernameClaim,
		UsernamePrefix:       preset.Server.UsernamePrefix,
		ClientAuthentication: client,
	}

	if useRequiredClaims {
		oidc.RequiredClaims = preset.Server.RequiredClaims
	}

	if shoot.Spec.Kubernetes.KubeAPIServer == nil {
		shoot.Spec.Kubernetes.KubeAPIServer = &garden.KubeAPIServerConfig{}
	}
	shoot.Spec.Kubernetes.KubeAPIServer.OIDCConfig = oidc
}
