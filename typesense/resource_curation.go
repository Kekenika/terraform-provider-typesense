package typesense

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

func resourceTypesenseCuration() *schema.Resource {
	return &schema.Resource{
		Description: "Promote or exclude certain documents from a query result",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the curation",
				Required:    true,
				ForceNew:    true,
			},
			"collection_name": {
				Type:        schema.TypeString,
				Description: "Name of the collection",
				Required:    true,
			},
			"rule": {
				Type:        schema.TypeList,
				Description: "Rule of this curation",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query": {
							Type:     schema.TypeString,
							Required: true,
						},
						"match": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"exact", "contains"}, false),
						},
					},
				},
			},
			"includes": {
				Type:        schema.TypeList,
				Description: "Documents to include",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Document id to include",
							Required:    true,
						},
						"position": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "Document position",
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},
			"excludes": {
				Type:        schema.TypeList,
				Description: "Documents to exclude",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Document id to exclude",
							Required:    true,
						},
					},
				},
			},
		},
		ReadContext:   resourceTypesenseCurationRead,
		CreateContext: resourceTypesenseCurationUpsert,
		UpdateContext: resourceTypesenseCurationUpsert,
		DeleteContext: resourceTypesenseCurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTypesenseCurationState,
		},
	}
}

func resourceTypesenseCurationUpsert(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	collectionName := d.Get("collection_name").(string)
	overwriteSchema := &api.SearchOverrideSchema{}

	if vs := d.Get("rule").([]interface{}); len(vs) > 0 {
		rule := vs[0].(map[string]interface{})

		overwriteSchema.Rule = api.SearchOverrideRule{
			Match: rule["match"].(string),
			Query: rule["query"].(string),
		}
	}

	if vs := d.Get("includes").([]interface{}); len(vs) > 0 {
		includes := make([]api.SearchOverrideInclude, len(vs))

		for i, v := range vs {
			r := v.(map[string]interface{})

			include := api.SearchOverrideInclude{
				Id: r["id"].(string),
			}

			if v, ok := r["position"].(int); ok {
				include.Position = v
			}

			includes[i] = include
		}

		overwriteSchema.Includes = includes
	}

	if vs := d.Get("excludes").([]interface{}); len(vs) > 0 {
		excludes := make([]api.SearchOverrideExclude, len(vs))

		for i, v := range vs {
			r := v.(map[string]interface{})
			excludes[i] = api.SearchOverrideExclude{
				Id: r["id"].(string),
			}
		}

		overwriteSchema.Excludes = excludes
	}

	override, err := client.Collection(collectionName).Overrides().Upsert(name, overwriteSchema)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s.%s", collectionName, override.Id))
	return diags
}

func resourceTypesenseCurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	collectionName, id, err := splitCollectionRelatedId(d.Id(), "curation")
	if err != nil {
		return diag.FromErr(err)
	}

	override, err := client.Collection(collectionName).Override(id).Retrieve()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	if err := d.Set("name", override.Id); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("rule", flattenCurationRule(override.Rule)); err != nil {
		return diag.FromErr(err)
	}

	if len(override.Includes) > 0 {
		if err := d.Set("includes", flattenCurationIncludes(override.Includes)); err != nil {
			return diag.FromErr(err)
		}
	}

	if len(override.Excludes) > 0 {
		if err := d.Set("excludes", flattenCurationExcludes(override.Excludes)); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceTypesenseCurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*typesense.Client)

	var diags diag.Diagnostics

	collectionName, id, err := splitCollectionRelatedId(d.Id(), "curation")
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Collection(collectionName).Override(id).Delete()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}

func resourceTypesenseCurationState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*typesense.Client)

	collectionName, id, err := splitCollectionRelatedId(d.Id(), "alias")
	if err != nil {
		return nil, err
	}

	override, err := client.Collection(collectionName).Override(id).Retrieve()
	if err != nil {
		return nil, err
	}

	d.SetId(fmt.Sprintf("%s.%s", collectionName, override.Id))
	return []*schema.ResourceData{d}, nil
}

func flattenCurationRule(rule api.SearchOverrideRule) []interface{} {
	res := []interface{}{}
	res[0] = map[string]interface{}{
		"match": rule.Match,
		"query": rule.Query,
	}
	return res
}

func flattenCurationIncludes(includes []api.SearchOverrideInclude) []interface{} {
	ins := make([]interface{}, len(includes))

	for i, include := range includes {
		in := make(map[string]interface{})
		in["id"] = include.Id

		if include.Position > 0 {
			in["position"] = include.Position
		}

		ins[i] = in
	}

	return ins
}

func flattenCurationExcludes(excludes []api.SearchOverrideExclude) []interface{} {
	exs := make([]interface{}, len(excludes))

	for i, exclude := range excludes {
		ex := make(map[string]interface{})
		ex["id"] = exclude.Id
		exs[i] = ex
	}

	return exs
}
