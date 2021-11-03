package typesense

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

func resourceTypesenseSynonyms() *schema.Resource {
	return &schema.Resource{
		Description: "Search terms that should be considered equivalent",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the synonyms",
				Required:    true,
			},
			"collection_name": {
				Type:        schema.TypeString,
				Description: "Name of the collection",
				Required:    true,
			},
			"synonyms": {
				Type:     schema.TypeList,
				Default:  "Target collection names",
				Required: true,
				Elem:     &schema.Resource{},
			},
			"root": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Root for one-way synonym",
			},
		},
		ReadContext:   resourceTypesenseSynonymsRead,
		CreateContext: resourceTypesenseSynonymsUpsert,
		UpdateContext: resourceTypesenseSynonymsUpsert,
		DeleteContext: resourceTypesenseSynonymsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTypesenseSynonymsUpsert(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	collectionName := d.Get("collection_name").(string)
	synonymSchema := &api.SearchSynonymSchema{
		Synonyms: interfaceArrayToStringArray(d.Get("synonyms").([]interface{})),
	}

	if v := d.Get("root"); v != nil {
		synonymSchema.Root = v.(string)
	}

	synonym, err := client.Collection(collectionName).Synonyms().Upsert(name, synonymSchema)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(synonym.Id)
	return diags
}

func resourceTypesenseSynonymsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	id := d.Id()
	collectionName := d.Get("collection_name").(string)

	synonym, err := client.Collection(collectionName).Synonym(id).Retrieve()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	if err := d.Set("name", synonym.Id); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("synonyms", synonym.Synonyms); err != nil {
		return diag.FromErr(err)
	}

	if synonym.Root != "" {
		if err := d.Set("root", synonym.Root); err != nil {
			if err := d.Set("root", synonym.Root); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return diags
}

func resourceTypesenseSynonymsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	id := d.Id()
	collectionName := d.Get("collection_name").(string)

	_, err := client.Collection(collectionName).Synonym(id).Delete()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}
