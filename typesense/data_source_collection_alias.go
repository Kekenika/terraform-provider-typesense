package typesense

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/typesense/typesense-go/typesense"
)

func dataSourceTypesenseCollectionAlias() *schema.Resource {
	return &schema.Resource{
		Description: "Virtual collection name that points to a real collection.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the collection alias",
				Required:    true,
			},
			"collection_name": {
				Type:        schema.TypeString,
				Description: "Target collection names",
				Computed:    true,
			},
		},
		ReadContext: resourceTypesenseCollectionAliasRead,
	}
}

func dataSourceTypesenseCollectionAliasRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	id := d.Id()

	alias, err := client.Alias(id).Retrieve()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	if err := d.Set("name", alias.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("collection_name", alias.CollectionName); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("id")
	return diags
}
