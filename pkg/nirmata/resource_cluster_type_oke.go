package nirmata

import (
	"fmt"
	"log"
	"strings"

	"regexp"
	"time"

	guuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	client "github.com/nirmata/go-client/pkg/client"
)

func resourceOkeClusterType() *schema.Resource {
	return &schema.Resource{
		Create: resourceOkeClusterTypeCreate,
		Read:   resourceOkeClusterTypeRead,
		Update: resourceOkeClusterTypeUpdate,
		Delete: resourceOkeClusterTypeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) > 64 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 64 characters", k))
					}
					if !regexp.MustCompile(`^[\w+=,.@-]*$`).MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must match [\\w+=,.@-]", k))
					}
					return
				},
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"credentials": {
				Type:     schema.TypeString,
				Required: true,
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vm_shape": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if !regexp.MustCompile(`^[\w+=,.@-]*$`).MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must match [\\w+=,.@-]", k))
					}
					return
				},
			},
		},
	}
}

func resourceOkeClusterTypeCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	clouduuid := guuid.New()
	nodepooluuid := guuid.New()

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	credentials := d.Get("credentials").(string)
	region := d.Get("region").(string)
	vmshape := d.Get("vm_shape").(string)

	cloudCredID, err := apiClient.QueryByName(client.ServiceClusters, "CloudCredentials", credentials)
	if err != nil {
		log.Printf("[ERROR] - %v", err)
		return err
	}

	var otherAddons []map[string]interface{}

	otherAddons = append(otherAddons, map[string]interface{}{
		"modelIndex":    "AddOnSpec",
		"name":          "kyverno",
		"addOnSelector": "kyverno",
		"catalog":       "default-addon-catalog",
	},
	)
	credential := map[string]interface{}{
		"id":         cloudCredID.UUID(),
		"service":    "Cluster",
		"modelIndex": "CloudCredentials",
	}

	clustertype := map[string]interface{}{
		"name":        name,
		"description": "",
		"modelIndex":  "ClusterType",
		"spec": map[string]interface{}{
			"clusterMode": "providerManaged",
			"modelIndex":  "ClusterSpec",
			"version":     version,
			"cloud":       "oraclecloudservices",
			"addons": map[string]interface{}{
				"dns":        false,
				"modelIndex": "AddOns",
				"other":      otherAddons,
			},
			"cloudConfigSpec": map[string]interface{}{
				"credentials":   credential,
				"id":            clouduuid,
				"modelIndex":    "CloudConfigSpec",
				"nodePoolTypes": nodepooluuid,
				"okeConfig": map[string]interface{}{
					"region":     region,
					"modelIndex": "OkeClusterConfig",
				},
			},
		},
	}

	nodepoolobj := map[string]interface{}{
		"id":              nodepooluuid,
		"modelIndex":      "NodePoolType",
		"name":            name + "-default-node-pool-type",
		"cloudConfigSpec": clouduuid,
		"spec": map[string]interface{}{
			"modelIndex": "NodePoolSpec",
			"okeConfig": map[string]interface{}{
				"vmshape":    vmshape,
				"modelIndex": "OkeNodePoolConfig",
			},
		},
	}

	txn := make(map[string]interface{})
	var objArr = make([]interface{}, 0)
	objArr = append(objArr, clustertype, nodepoolobj)
	txn["create"] = objArr

	data, err := apiClient.PostFromJSON(client.ServiceClusters, "txn", txn, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to create cluster type  with data : %v", err)
		return err
	}

	changeID := data["changeId"].(string)
	d.SetId(changeID)

	return nil
}

func resourceOkeClusterTypeRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceOkeClusterTypeUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceOkeClusterTypeDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	name := d.Get("name").(string)

	id, err := apiClient.QueryByName(client.ServiceClusters, "clustertypes", name)
	if err != nil {
		log.Printf("[ERROR] - %v", err.Error())
		return err
	}

	params := map[string]string{
		"action": "delete",
	}

	if err := apiClient.Delete(id, params); err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}

		log.Printf("[ERROR] - %v", err.Error())
		return err
	}

	log.Printf("Deleted cluster type %s", name)
	return nil
}
