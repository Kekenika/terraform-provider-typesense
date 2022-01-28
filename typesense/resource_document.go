package typesense

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/typesense/typesense-go/typesense"
)

func resourceTypesenseDocument() *schema.Resource {
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
				ForceNew:    true,
				Description: "Document's body",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		ReadContext:   resourceTypesenseDocumentRead,
		CreateContext: resourceTypesenseDocumentUpsert,
		UpdateContext: resourceTypesenseDocumentUpsert,
		DeleteContext: resourceTypesenseDocumentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTypesenseDocumentState,
		},
	}
}

func resourceTypesenseDocumentUpsert(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var collectionName string

	if v, ok := d.GetOk("collection_name"); ok {
		collectionName = v.(string)
	}

	var document map[string]interface{}

	if v, ok := d.GetOk("document"); ok {
		document = v.(map[string]interface{})
	}

	var id string
	if v, ok := document["id"]; ok {
		id = v.(string)
	} else {
		return diag.Errorf("id required for document")
	}

	_, err := client.Collection(collectionName).Documents().Create(document)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s.%s", collectionName, id))
	return resourceTypesenseCurationRead(ctx, d, meta)
}

func resourceTypesenseDocumentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	collectionName, id, err := splitCollectionRelatedId(d.Id(), "document")
	if err != nil {
		return diag.FromErr(err)
	}

	doc, err := client.Collection(collectionName).Document(id).Retrieve()
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

	return diags
}

func resourceTypesenseDocumentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	collectionName, id, err := splitCollectionRelatedId(d.Id(), "document")
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Collection(collectionName).Document(id).Delete()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}

func resourceTypesenseDocumentState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*typesense.Client)

	collectionName, id, err := splitCollectionRelatedId(d.Id(), "document")
	if err != nil {
		return nil, err
	}

	doc, err := client.Collection(collectionName).Document(id).Retrieve()
	if err != nil {
		return nil, err
	}

	d.SetId(fmt.Sprintf("%s.%s", collectionName, doc["id"]))
	return []*schema.ResourceData{d}, nil
}
