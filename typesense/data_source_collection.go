package typesense

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTypesenseCollection() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceTypesenseCollectionRead,
		Description: "Group of related documents which are roughly equivalent to a table in a relational database.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the collection.",
			},
			"fields": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"facet": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Facetable field",
						},
						"index": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Index field",
						},
						"optional": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Optional field",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
					},
				},
			},
			"default_sorting_field": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"num_documents": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}
