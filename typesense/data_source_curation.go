package typesense

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		ReadContext: resourceTypesenseCurationRead,
	}
}
