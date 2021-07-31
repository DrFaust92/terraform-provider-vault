package vault

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceConsulAccessCredentials_basic(t *testing.T) {
	backend := acctest.RandomWithPrefix("tf-test-backend")
	name := acctest.RandomWithPrefix("tf-test-name")
	token := "026a0c16-87cd-4c2d-b3f3-fb539f592b7e"

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConsulAccessCredentialsConfigBase(backend, name, token),
				// Check: resource.ComposeTestCheckFunc(
				// 	resource.TestCheckResourceAttr("data.vault_consul_access_credentials.test", "security_token", ""),
				// 	resource.TestCheckResourceAttr("data.vault_consul_access_credentials.test", "type", "creds"),
				// 	resource.TestCheckResourceAttrSet("data.vault_consul_access_credentials.test", "lease_id"),
				// ),
			},
			{
				Config: testAccDataSourceConsulAccessCredentialsConfig_basic(backend, name, token),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vault_consul_access_credentials.test", "security_token", ""),
					resource.TestCheckResourceAttr("data.vault_consul_access_credentials.test", "type", "creds"),
					resource.TestCheckResourceAttrSet("data.vault_consul_access_credentials.test", "lease_id"),
				),
			},
		},
	})
}

func testAccDataSourceConsulAccessCredentialsConfigBase(backend, name, token string) string {
	return fmt.Sprintf(`
resource "vault_consul_secret_backend" "test" {
  path                      = "%s"
  description               = "test description"
  default_lease_ttl_seconds = 3600
  max_lease_ttl_seconds     = 86400
  address                   = "127.0.0.1:8500"
  token                     = "%s"
}

resource "vault_consul_secret_backend_role" "test" {
  backend = vault_consul_secret_backend.test.path
  name    = "%s"

  policies = [
    "foo"
  ]
}
`, backend, token, name)
}

func testAccDataSourceConsulAccessCredentialsConfig_basic(backend, name, token string) string {
	return testAccDataSourceConsulAccessCredentialsConfigBase(backend, name, token) + `
	data "vault_consul_access_credentials" "test" {
		backend = vault_consul_secret_backend_role.test.path
		role    = vault_consul_secret_backend_role.test.name
	  }
	`
}
