package typesense

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

func resourceTypesenseCollection() *schema.Resource {
	return &schema.Resource{
		Description: "Group of related documents which are roughly equivalent to a table in a relational database.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Name of the collection.",
				Required:    true,
			},
			"fields": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"facet": {
							Type:        schema.TypeBool,
							ForceNew:    true,
							Optional:    true,
							Computed:    true,
							Description: "Facetable field",
						},
						"index": {
							Type:        schema.TypeBool,
							ForceNew:    true,
							Optional:    true,
							Default:     true,
							Description: "Index field",
						},
						"optional": {
							Type:        schema.TypeBool,
							ForceNew:    true,
							Optional:    true,
							Computed:    true,
							Description: "Optional field",
						},
						"type": {
							Type:        schema.TypeString,
							ForceNew:    true,
							Required:    true,
							Description: "Field type",
							ValidateFunc: validation.StringInSlice([]string{
								"string",
								"int32",
								"int64",
								"float",
								"bool",
								"string[]",
								"int32[]",
								"int64[]",
								"float[]",
								"bool[]",
								"geopoint",
								"auto",
							}, false),
						},
					},
				},
			},
			"default_sorting_field": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"num_documents": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		ReadContext:   resourceTypesenseCollectionRead,
		CreateContext: resourceTypesenseCollectionCreate,
		UpdateContext: resourceTypesenseCollectionUpdate,
		DeleteContext: resourceTypesenseCollectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTypesenseCollectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	schema := &api.CollectionSchema{}

	if v := d.Get("name"); v != "" {
		schema.Name = v.(string)
	}

	if v := d.Get("default_sorting_field"); v != "" {
		schema.DefaultSortingField = v.(string)
	}

	fields := []api.Field{}
	for _, vs := range d.Get("fields").([]interface{}) {
		v := vs.(map[string]interface{})

		field := api.Field{
			Name: v["name"].(string),
			Type: v["type"].(string),
		}

		if value := v["facet"]; value != "" {
			field.Facet = value.(bool)
		}

		if value := v["optional"]; value != "" {
			field.Optional = value.(bool)
		}

		if value := v["index"]; value != "" {
			field.Index = boolPointer(value.(bool))
		}

		fields = append(fields, field)
	}

	schema.Fields = fields

	collection, err := client.Collections().Create(schema)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(collection.Name)
	return resourceTypesenseCollectionRead(ctx, d, meta)
}

func resourceTypesenseCollectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	id := d.Id()

	collection, err := client.Collection(id).Retrieve()
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

// Typesense doesn't offer any update API for collections. The team has plans to offer it, see the following issue for details.
// https://github.com/typesense/typesense/issues/96
func resourceTypesenseCollectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	id := d.Id()

	_, err := client.Collection(id).Delete()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	schema := &api.CollectionSchema{}

	if v := d.Get("name"); v != "" {
		schema.Name = v.(string)
	}

	if v := d.Get("default_sorting_field"); v != "" {
		schema.DefaultSortingField = v.(string)
	}

	fields := []api.Field{}
	for _, vs := range d.Get("fields").([]interface{}) {
		v := vs.(map[string]interface{})

		field := api.Field{
			Name: v["name"].(string),
			Type: v["type"].(string),
		}

		if value := v["facet"]; value != "" {
			field.Facet = value.(bool)
		}

		if value := v["optional"]; value != "" {
			field.Optional = value.(bool)
		}

		if value := v["index"]; value != "" {
			field.Index = boolPointer(value.(bool))
		}

		fields = append(fields, field)
	}

	schema.Fields = fields

	_, err = client.Collections().Create(schema)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTypesenseCurationRead(ctx, d, meta)
}

func resourceTypesenseCollectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	id := d.Id()

	_, err := client.Collection(id).Delete()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}

func flattenCollectionFields(fields []api.Field) []interface{} {
	if fields != nil {
		fis := make([]interface{}, len(fields))

		for i, field := range fields {
			fi := make(map[string]interface{})
			fi["name"] = field.Name
			fi["facet"] = field.Facet
			fi["index"] = field.Index
			fi["optional"] = field.Optional
			fi["type"] = field.Type
			fis[i] = fi
		}

		return fis
	}

	return make([]interface{}, 0)
}
