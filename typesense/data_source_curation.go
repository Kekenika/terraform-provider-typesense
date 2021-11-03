package typesense

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/typesense/typesense-go/typesense"
)

func dataSourceTypesenseCuration() *schema.Resource {
	return &schema.Resource{
		Description: "Promote or exclude certain documents from a query result",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the curation",
				Required:    true,
			},
			"collection_name": {
				Type:        schema.TypeString,
				Description: "Name of the collection",
				Required:    true,
			},
			"rule": {
				Type:        schema.TypeList,
				Description: "Rule of this curation",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"match": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"includes": {
				Type:        schema.TypeList,
				Description: "Documents to include",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Document id to include",
							Computed:    true,
						},
						"position": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Document position",
						},
					},
				},
			},
			"excludes": {
				Type:        schema.TypeList,
				Description: "Documents to exclude",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Document id to exclude",
							Computed:    true,
						},
					},
				},
			},
		},
		ReadContext: dataSourceTypesenseCurationRead,
	}
}

func dataSourceTypesenseCurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	collectionName := d.Get("collection_name").(string)
	id := fmt.Sprintf("%s.%s", collectionName, name)

	override, err := client.Collection(collectionName).Override(name).Retrieve()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	if err := d.Set("name", override.Id); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("collection_name", collectionName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("rule", flattenCurationRule(override.Rule)); err != nil {
		return diag.FromErr(err)
	}

	if len(override.Includes) > 0 {
		if err := d.Set("includes", flattenCurationIncludes(override.Includes)); err != nil {
			return diag.FromErr(err)
		}
	}

	if len(override.Excludes) > 0 {
		if err := d.Set("excludes", flattenCurationExcludes(override.Excludes)); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(id)
	return diags
}
