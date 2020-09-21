package onepassword

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVaultMember() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceVaultMemberRead,
		CreateContext: resourceVaultMemberCreate,
		DeleteContext: resourceVaultMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"vault": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"user": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceVaultMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vaultID, userID, err := resourceVaultMemberExtractID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	m := meta.(*Meta)
	v, err := m.onePassClient.ListVaultMembers(vaultID)
	if err != nil {
		return diag.FromErr(err)
	}

	var found string
	for _, member := range v {
		if member.UUID == userID {
			found = member.UUID
		}
	}

	if found == "" {
		d.SetId("")
		return nil
	}

	d.SetId(resourceVaultMemberBuildID(vaultID, found))
	d.Set("vault", vaultID)
	d.Set("user", found)
	return nil
}

func resourceVaultMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)
	err := m.onePassClient.CreateVaultMember(
		d.Get("vault").(string),
		d.Get("user").(string),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceVaultMemberBuildID(d.Get("vault").(string), d.Get("user").(string)))
	return resourceVaultMemberRead(ctx, d, meta)
}

func resourceVaultMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vaultID, userID, err := resourceVaultMemberExtractID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	m := meta.(*Meta)
	err = m.onePassClient.DeleteVaultMember(
		vaultID,
		userID,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// resourceVaultMemberBuildID will conjoin the vault ID and user ID into a single string
// This is used as the resource ID.
//
// Note that user ID is being lowercased. Some operations require this user ID to be uppercased.
// Use the resourceVaultMemberExtractID function to correctly reverse this encoding.
func resourceVaultMemberBuildID(vaultID, userID string) string {
	return strings.ToLower(vaultID + "-" + strings.ToLower(userID))
}

// resourceVaultMemberExtractID will split the vault ID and user ID from a given resource ID
//
// Note that user ID is being uppercased. Some operations require this user ID to be uppercased.
func resourceVaultMemberExtractID(id string) (vaultID, userID string, err error) {
	spl := strings.Split(id, "-")
	if len(spl) != 2 {
		return "", "", fmt.Errorf("Improperly formatted vault member string. The format \"vaultid-userid\" is expected")
	}
	return spl[0], strings.ToUpper(spl[1]), nil
}
