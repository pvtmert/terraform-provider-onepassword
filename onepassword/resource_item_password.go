package onepassword

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceItemPassword() *schema.Resource {
	return &schema.Resource{
		Read:   resourceItemPasswordRead,
		Create: resourceItemPasswordCreate,
		Delete: resourceItemDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				err := resourceItemPasswordRead(d, meta)
				return []*schema.ResourceData{d}, err
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vault": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"section": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     sectionSchema,
			},
			"url": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: urlValidate,
			},
		},
	}
}

func resourceItemPasswordRead(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	vaultId := d.Get("vault").(string)
	err, v := m.onePassClient.ReadItem(getId(d), vaultId)
	log.Printf("[DEBUG] %v", v)
	if err != nil {
		return err
	}
	if v.Template != Category2Template(PasswordCategory) {
		return errors.New("Item is not from " + string(PasswordCategory))
	}

	d.SetId(v.Uuid)
	d.Set("name", v.Overview.Title)
	if err := d.Set("url", v.Overview.Url); err != nil {
		return err
	}
	d.Set("notes", v.Details.Notes)
	d.Set("tags", v.Overview.Tags)
	d.Set("vault", v.Vault)
	d.Set("password", v.Details.Password)
	return d.Set("section", v.ProcessSections())
}

func resourceItemPasswordCreate(d *schema.ResourceData, meta interface{}) error {
	item := &Item{
		Vault:    d.Get("vault").(string),
		Template: Category2Template(PasswordCategory),
		Overview: Overview{
			Title: d.Get("name").(string),
			Url:   d.Get("url").(string),
			Tags:  ParseTags(d),
		},
		Details: Details{
			Notes:    d.Get("notes").(string),
			Password: d.Get("password").(string),
			Sections: ParseSections(d),
		},
	}
	m := meta.(*Meta)
	err := m.onePassClient.CreateItem(item)
	if err != nil {
		return err
	}
	d.SetId(item.Uuid)
	return nil
}