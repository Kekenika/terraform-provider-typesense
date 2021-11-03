package typesense

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/typesense/typesense-go/typesense"
)

func resourceTypesenseDocument() *schema.Resource {
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
				Required:    true,
				MaxItems:    1,
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
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTypesenseDocumentUpsert(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	var collectionName string

	if v, ok := d.GetOk("collection_name"); ok {
		collectionName = v.(string)
	}

	var document map[string]interface{}

	if v, ok := d.GetOk("document"); ok {
		document = v.(map[string]interface{})
	}

	var id string
	if v, ok := d.GetOk("id"); ok {
		id = v.(string)
	}

	document["id"] = id

	_, err := client.Collection(collectionName).Documents().Create(document)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return diags
}

func resourceTypesenseDocumentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	id := d.Id()
	collectionName := d.Get("collection_name").(string)

	doc, err := client.Collection(collectionName).Document(id).Retrieve()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	if err := d.Set("document", doc); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceTypesenseDocumentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	id := d.Id()
	collectionName := d.Get("collection_name").(string)

	_, err := client.Collection(collectionName).Document(id).Delete()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}
