package bref

import (
    "fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"AWS_REGION",
					"AWS_DEFAULT_REGION",
				}, nil),
				Description: "AWS Region of Bref PHP Layers. Can be specified with the `AWS_REGION` " +
					"or `AWS_DEFAULT_REGION` environment variable.",
				InputDefault: "us-east-1",
			},
			"bref_version": {
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("BREF_VERSION", "1.2.0"),
				Description: "The Bref PHP Runtime Version to work with. Can be specified with the " +
					"`BREF_VERSION` environment variable.",
			},
			"account_id": {
				Type: schema.TypeString,
				Optional: true,
				Default: "209497400698",
				Description: "The Bref PHP Lambda Layer AWS Account ID.",
			},
		},

		ResourcesMap:   map[string]*schema.Resource{},

		DataSourcesMap: map[string]*schema.Resource{
			"bref_lambda_layer": dataSourceBrefLambdaLayer(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func init() {
	schema.DescriptionKind = schema.StringMarkdown

	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		if s.Deprecated != "" {
			desc += " " + s.Deprecated
		}
		return strings.TrimSpace(desc)
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Region:    d.Get("region").(string),
		Version:   d.Get("bref_version").(string),
		AccountId: d.Get("account_id").(string),
	}

	return &config, nil
}

type Config struct {
	Region    string
	Version   string
	AccountId string
}