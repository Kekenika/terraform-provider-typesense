package typesense

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTypesenseDocument() *schema.Resource {
	return &schema.Resource{
		Description: "Item in a collection",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "ID of the document",
				Required:    true,
			},
			"collection_name": {
				Type:        schema.TypeString,
				Description: "Name of the collection",
				Required:    true,
			},
			"document": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Document's body",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		ReadContext: resourceTypesenseDocumentRead,
	}
}
