package typesense

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTypesenseSynonyms() *schema.Resource {
	return &schema.Resource{
		Description: "Search terms that should be considered equivalent",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the synonyms",
				Computed:    true,
			},
			"collection_name": {
				Type:        schema.TypeString,
				Description: "Name of the collection",
				Computed:    true,
			},
			"synonyms": {
				Type:        schema.TypeList,
				Description: "Target collection names",
				Computed:    true,
				Elem:        &schema.Resource{},
			},
			"root": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Root for one-way synonym",
			},
		},
		ReadContext: resourceTypesenseSynonymsRead,
	}
}
