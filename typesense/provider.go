package typesense

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/typesense/typesense-go/typesense"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("TYPESENSE_API_KEY", nil),
				Description: "API Key to access the Typesense server. This can also be set via the `TYPESENSE_API_KEY` environment variable.",
			},
			"api_address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TYPESENSE_API_ADDRESS", nil),
				Description: "URL of the Typesense server. This can also be set via the `TYPESENSE_API_ADDRESS` environment variable.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"typesense_collection":       dataSourceTypesenseCollection(),
			"typesense_collection_alias": dataSourceTypesenseCollectionAlias(),
			"typesense_curation":         dataSourceTypesenseCuration(),
			"typesense_document":         dataSourceTypesenseDocument(),
			"typesense_synonyms":         dataSourceTypesenseSynonyms(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"typesense_collection":       resourceTypesenseCollection(),
			"typesense_collection_alias": resourceTypesenseCollectionAlias(),
			"typesense_document":         resourceTypesenseDocument(),
			"typesense_curation":         resourceTypesenseCuration(),
			"typesense_synonyms":         resourceTypesenseSynonyms(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	opts := []typesense.ClientOption{}

	if v := d.Get("api_key").(string); v != "" {
		opts = append(opts, typesense.WithAPIKey(v))
	}

	if v := d.Get("api_address").(string); v != "" {
		opts = append(opts, typesense.WithServer(v))
	}

	return typesense.NewClient(opts...), nil
}
