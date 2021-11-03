package typesense

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/typesense/typesense-go/typesense"
)

func dataSourceTypesenseDocument() *schema.Resource {
	return &schema.Resource{
		Description: "Item in a collection",
		Schema: map[string]*schema.Schema{
			"collection_name": {
				Type:        schema.TypeString,
				Description: "Name of the collection",
				Required:    true,
			},
			"document": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "Document's body",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		ReadContext: dataSourceTypesenseDocumentRead,
	}
}

func dataSourceTypesenseDocumentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	collectionName := d.Get("collection_name").(string)

	document := d.Get("document").(map[string]interface{})

	var docId string
	if v, ok := document["id"]; ok {
		docId = v.(string)
	} else {
		return diag.Errorf("id required")
	}

	id := fmt.Sprintf("%s.%s", collectionName, docId)

	doc, err := client.Collection(collectionName).Document(docId).Retrieve()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	if err := d.Set("document", doc); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("collection_name", collectionName); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return diags
}
