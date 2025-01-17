/*
Copyright 2021 The Kubermatic Kubernetes Platform contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aws

import (
	"context"
	"errors"
	"fmt"

	"k8c.io/dashboard/v2/pkg/provider"
	kubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1"
)

const (
	authFailure = "AuthFailure"
)

type AmazonEC2 struct {
	dc                *kubermaticv1.DatacenterSpecAWS
	secretKeySelector provider.SecretKeySelectorValueFunc

	// clientSet is used during tests only
	clientSet *ClientSet
}

// NewCloudProvider returns a new AmazonEC2 provider.
func NewCloudProvider(dc *kubermaticv1.Datacenter, secretKeyGetter provider.SecretKeySelectorValueFunc) (*AmazonEC2, error) {
	if dc.Spec.AWS == nil {
		return nil, errors.New("datacenter is not an AWS datacenter")
	}

	return &AmazonEC2{
		dc:                dc.Spec.AWS,
		secretKeySelector: secretKeyGetter,
	}, nil
}

var _ provider.CloudProvider = &AmazonEC2{}

func (a *AmazonEC2) getClientSet(ctx context.Context, cloud kubermaticv1.CloudSpec) (*ClientSet, error) {
	if a.clientSet != nil {
		return a.clientSet, nil
	}

	accessKeyID, secretAccessKey, assumeRoleARN, assumeRoleExternalID, err := GetCredentialsForCluster(cloud, a.secretKeySelector)
	if err != nil {
		return nil, err
	}

	return GetClientSet(ctx, accessKeyID, secretAccessKey, assumeRoleARN, assumeRoleExternalID, a.dc.Region)
}

func (a *AmazonEC2) DefaultCloudSpec(ctx context.Context, spec *kubermaticv1.CloudSpec) error {
	return nil
}

// ValidateCloudSpec validates the fields that the user can override while creating
// a cluster. We only check those that must pre-exist in the AWS account
// (i.e. the security group and VPC), because the others (like route table)
// will be created if they do not yet exist / are not explicitly specified.
// TL;DR: This validation does not need to be extended to cover more than
// VPC and SG.
func (a *AmazonEC2) ValidateCloudSpec(ctx context.Context, spec kubermaticv1.CloudSpec) error {
	client, err := a.getClientSet(ctx, spec)
	if err != nil {
		return fmt.Errorf("failed to get API client: %w", err)
	}

	// Some settings require the vpc to be set
	if spec.AWS.SecurityGroupID != "" {
		if spec.AWS.VPCID == "" {
			return fmt.Errorf("VPC must be set when specifying a security group")
		}
	}

	if spec.AWS.VPCID != "" {
		vpc, err := getVPCByID(ctx, client.EC2, spec.AWS.VPCID)
		if err != nil {
			return err
		}

		if spec.AWS.SecurityGroupID != "" {
			if _, err = getSecurityGroupByID(ctx, client.EC2, vpc, spec.AWS.SecurityGroupID); err != nil {
				return err
			}
		}
	}

	return nil
}

// ValidateCloudSpecUpdate verifies whether an update of cloud spec is valid and permitted.
func (a *AmazonEC2) ValidateCloudSpecUpdate(_ context.Context, oldSpec kubermaticv1.CloudSpec, newSpec kubermaticv1.CloudSpec) error {
	if oldSpec.AWS == nil || newSpec.AWS == nil {
		return errors.New("'aws' spec is empty")
	}

	if oldSpec.AWS.VPCID != "" && oldSpec.AWS.VPCID != newSpec.AWS.VPCID {
		return fmt.Errorf("updating AWS VPC ID is not supported (was %s, updated to %s)", oldSpec.AWS.VPCID, newSpec.AWS.VPCID)
	}

	if oldSpec.AWS.RouteTableID != "" && oldSpec.AWS.RouteTableID != newSpec.AWS.RouteTableID {
		return fmt.Errorf("updating AWS route table ID is not supported (was %s, updated to %s)", oldSpec.AWS.RouteTableID, newSpec.AWS.RouteTableID)
	}

	if oldSpec.AWS.SecurityGroupID != "" && oldSpec.AWS.SecurityGroupID != newSpec.AWS.SecurityGroupID {
		return fmt.Errorf("updating AWS security group ID is not supported (was %s, updated to %s)", oldSpec.AWS.SecurityGroupID, newSpec.AWS.SecurityGroupID)
	}

	if oldSpec.AWS.ControlPlaneRoleARN != "" && oldSpec.AWS.ControlPlaneRoleARN != newSpec.AWS.ControlPlaneRoleARN {
		return fmt.Errorf("updating AWS control plane ARN is not supported (was %s, updated to %s)", oldSpec.AWS.ControlPlaneRoleARN, newSpec.AWS.ControlPlaneRoleARN)
	}

	if oldSpec.AWS.InstanceProfileName != "" && oldSpec.AWS.InstanceProfileName != newSpec.AWS.InstanceProfileName {
		return fmt.Errorf("updating AWS instance profile name is not supported (was %s, updated to %s)", oldSpec.AWS.InstanceProfileName, newSpec.AWS.InstanceProfileName)
	}

	return nil
}
