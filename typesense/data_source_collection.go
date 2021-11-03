package typesense

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/typesense/typesense-go/typesense"
)

func dataSourceTypesenseCollection() *schema.Resource {
	return &schema.Resource{
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
		ReadContext: dataSourceTypesenseCollectionRead,
	}
}

func dataSourceTypesenseCollectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	collection, err := client.Collection(d.Get("name").(string)).Retrieve()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Got collection name:%s\n", collection.Name)

	if err := d.Set("name", collection.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("fields", flattenCollectionFields(collection.Fields)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("default_sorting_field", collection.DefaultSortingField); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("num_documents", collection.NumDocuments); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
