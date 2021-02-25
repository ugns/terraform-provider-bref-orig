package bref

import (
    "context"
    "encoding/json"
    "net/http"
    "fmt"
    "strconv"
    "time"

    cleanhttp "github.com/hashicorp/go-cleanhttp"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBrefLambdaLayer() *schema.Resource {
    return &schema.Resource{
        Description: "Bref PHP Lambda Layer.",
        
        ReadContext: dataSourceBrefLambdaLayerRead,

        Schema: map[string]*schema.Schema{
            "layer_name": {
                Type: schema.TypeString,
                Required: true,
            },
            "version": {
                Type: schema.TypeInt,
                Computed: true,
            },
            "arn": {
                Type: schema.TypeString,
                Computed: true,
            },
            "layer_arn": {
                Type: schema.TypeString,
                Computed: true,
            },
        },
    }
}

func dataSourceBrefLambdaLayerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    cli := cleanhttp.DefaultPooledClient()

    layerName := d.Get("layer_name").(string)

    brefVersion := m.(*Config).Version
    accountId := m.(*Config).AccountId
    region := m.(*Config).Region

    // Warning or errors can be collected in a slice type
    var diags diag.Diagnostics

    req, err := http.NewRequest("GET", fmt.Sprintf("https://raw.githubusercontent.com/brefphp/bref/%s/layers.json", brefVersion), nil)
    if err != nil {
        return diag.Errorf("Unable to build request for %s version of Bref runtime layers", brefVersion)
    }

    r, err := cli.Do(req)
    if err != nil {
        return diag.Errorf("Error retrieving Bref runtime layers: %s", err.Error())
    }
    defer r.Body.Close()

    var layers map[string]interface{}
    err = json.NewDecoder(r.Body).Decode(&layers)
    if err != nil {
        return diag.Errorf("Error parsing Bref runtime layers: %s", err.Error())
    }

    regions := layers[layerName].(map[string]interface{})
    version, err := strconv.Atoi(regions[region].(string))
    if err != nil {
        return diag.Errorf("Unable to locate a Bref v%s lambda layer version for %s in %s region", brefVersion, layerName, region)
    }
    arn := fmt.Sprintf("arn:aws:lambda:%s:%s:layer:%s:%d", region, accountId, layerName, version)
    
    if err := d.Set("version", version); err != nil {
        return diag.FromErr(err)
    }

    if err := d.Set("arn", arn); err != nil {
        return diag.FromErr(err)
    }

    if err := d.Set("layer_arn", arn); err != nil {
        return diag.FromErr(err)
    }

    // always run
    d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

    return diags
}