package typesense

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

func resourceTypesenseCollectionAlias() *schema.Resource {
	return &schema.Resource{
		Description: "Virtual collection name that points to a real collection.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the collection alias",
				Required:    true,
				ForceNew:    true,
			},
			"collection_name": {
				Type:        schema.TypeString,
				Description: "Target collection names",
				Required:    true,
			},
		},
		ReadContext:   resourceTypesenseCollectionAliasRead,
		CreateContext: resourceTypesenseCollectionAliasUpsert,
		UpdateContext: resourceTypesenseCollectionAliasUpsert,
		DeleteContext: resourceTypesenseCollectionAliasDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTypesenseCollectionAliasUpsert(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	name := d.Get("name").(string)
	aliasSchema := &api.CollectionAliasSchema{
		CollectionName: d.Get("collection_name").(string),
	}

	alias, err := client.Aliases().Upsert(name, aliasSchema)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(alias.Name)
	return resourceTypesenseCollectionAliasRead(ctx, d, meta)
}

func resourceTypesenseCollectionAliasRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return diags
}

func resourceTypesenseCollectionAliasDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	id := d.Id()

	_, err := client.Alias(id).Delete()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}
