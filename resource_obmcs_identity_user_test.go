// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package main

import (
	"testing"
	"time"

	"github.com/MustWin/baremetal-sdk-go"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/stretchr/testify/suite"
)

type ResourceIdentityUserTestSuite struct {
	suite.Suite
	Client       *baremetal.Client
	Provider     terraform.ResourceProvider
	Providers    map[string]terraform.ResourceProvider
	TimeCreated  time.Time
	Config       string
	ResourceName string
	Res          *baremetal.User
}

func (s *ResourceIdentityUserTestSuite) SetupTest() {
	s.Client = GetTestProvider()

	s.Provider = Provider(
		func(d *schema.ResourceData) (interface{}, error) {
			return s.Client, nil
		},
	)

	s.Providers = map[string]terraform.ResourceProvider{
		"oci": s.Provider,
	}
	s.TimeCreated, _ = time.Parse("2006-Jan-02", "2006-Jan-02")
	s.Config = `
		resource "oci_identity_user" "t" {
			name = "-tf-user"
			description = "automated test user"
		}
	`
	s.Config += testProviderConfig()

	s.ResourceName = "oci_identity_user.t"
	s.Res = &baremetal.User{
		ID:            "id!",
		Name:          "-tf-user",
		Description:   "automated test user",
		CompartmentID: "cid!",
		State:         baremetal.ResourceActive,
		TimeCreated:   s.TimeCreated,
	}

}

func (s *ResourceIdentityUserTestSuite) TestCreateResourceIdentityUser() {

	resource.UnitTest(s.T(), resource.TestCase{
		Providers: s.Providers,
		Steps: []resource.TestStep{
			{
				ImportState:       true,
				ImportStateVerify: true,
				Config:            s.Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(s.ResourceName, "name", s.Res.Name),
					resource.TestCheckResourceAttr(s.ResourceName, "description", s.Res.Description),

					resource.TestCheckResourceAttr(s.ResourceName, "state", s.Res.State),
					resource.TestCheckResourceAttrSet(s.ResourceName, "time_created"),
				),
			},
		},
	})
}

func (s *ResourceIdentityUserTestSuite) TestCreateResourceIdentityUserPolling() {
	s.Res.State = baremetal.ResourceCreating

	u := *s.Res
	u.State = baremetal.ResourceActive

	resource.UnitTest(s.T(), resource.TestCase{
		Providers: s.Providers,
		Steps: []resource.TestStep{
			{
				ImportState:       true,
				ImportStateVerify: true,
				Config:            s.Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(s.ResourceName, "state", baremetal.ResourceActive),
				),
			},
		},
	})
}

func (s *ResourceIdentityUserTestSuite) TestUpdateResourceIdentityUserDescription() {

	c := `
		resource "oci_identity_user" "t" {
			name = "-tf-user"
			description = "automated test user updated"
		}
	`
	c += testProviderConfig()

	resource.UnitTest(s.T(), resource.TestCase{
		Providers: s.Providers,
		Steps: []resource.TestStep{
			{
				ImportState:       true,
				ImportStateVerify: true,
				Config:            s.Config,
			},
			{
				Config: c,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(s.ResourceName, "description", "automated test user updated"),
				),
			},
		},
	})
}

func (s *ResourceIdentityUserTestSuite) TestUpdateResourceIdentityUserNameShouldCreateNew() {

	c := `
		resource "oci_identity_user" "t" {
			name = "-tf-user2"
			description = "automated test user"
		}
	`

	c += testProviderConfig()

	resource.UnitTest(s.T(), resource.TestCase{
		Providers: s.Providers,
		Steps: []resource.TestStep{
			{
				ImportState:       true,
				ImportStateVerify: true,
				Config:            s.Config,
			},
			{
				Config: c,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(s.ResourceName, "name", "-tf-user2"),
				),
			},
		},
	})
}

func (s *ResourceIdentityUserTestSuite) TestDeleteResourceIdentityUser() {

	resource.UnitTest(s.T(), resource.TestCase{
		Providers: s.Providers,
		Steps: []resource.TestStep{
			{
				ImportState:       true,
				ImportStateVerify: true,
				Config:            s.Config,
			},
			{
				Config:  s.Config,
				Destroy: true,
			},
		},
	})

}

func TestResourceIdentityUserTestSuite(t *testing.T) {
	suite.Run(t, new(ResourceIdentityUserTestSuite))
}
