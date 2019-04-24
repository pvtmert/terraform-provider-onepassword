package onepassword

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceVault() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVaultRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Optional:    true,
				Description: "Vault name.",
			},
		},
	}
}

func dataSourceVaultRead(d *schema.ResourceData, meta interface{}) error {
	return resourceVaultRead(d, meta)
}