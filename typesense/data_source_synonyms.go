package typesense

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/typesense/typesense-go/typesense"
)

func dataSourceTypesenseSynonyms() *schema.Resource {
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
				Type:        schema.TypeList,
				Description: "Target collection names",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"root": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Root for one-way synonym",
			},
		},
		ReadContext: dataSourceTypesenseSynonymsRead,
	}
}

func dataSourceTypesenseSynonymsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	collectionName := d.Get("collection_name").(string)

	id := fmt.Sprintf("%s.%s", collectionName, name)

	synonym, err := client.Collection(collectionName).Synonym(name).Retrieve()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	if err := d.Set("collection_name", collectionName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", synonym.Id); err != nil {
		return diag.FromErr(err)
	}

	log.Println(synonym.Synonyms)

	if err := d.Set("synonyms", synonym.Synonyms); err != nil {
		return diag.FromErr(err)
	}

	if synonym.Root != nil && *synonym.Root != "" {
		if err := d.Set("root", synonym.Root); err != nil {
			if err := d.Set("root", synonym.Root); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(id)

	return diags
}
