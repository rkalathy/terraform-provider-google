// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/bigtableadmin/v2"
)

func ResourceBigtableAppProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigtableAppProfileCreate,
		Read:   resourceBigtableAppProfileRead,
		Update: resourceBigtableAppProfileUpdate,
		Delete: resourceBigtableAppProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceBigtableAppProfileImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"app_profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The unique name of the app profile in the form '[_a-zA-Z0-9][-_.a-zA-Z0-9]*'.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Long form description of the use case for this app profile.`,
			},
			"ignore_warnings": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `If true, ignore safety checks when deleting/updating the app profile.`,
				Default:     false,
			},
			"instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareResourceNames,
				Description:      `The name of the instance to create the app profile within.`,
			},
			"multi_cluster_routing_use_any": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: `If true, read/write requests are routed to the nearest cluster in the instance, and will fail over to the nearest cluster that is available
in the event of transient errors or delays. Clusters in a region are considered equidistant. Choosing this option sacrifices read-your-writes
consistency to improve availability.`,
				ExactlyOneOf: []string{"single_cluster_routing", "multi_cluster_routing_use_any"},
			},
			"single_cluster_routing": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Use a single-cluster routing policy.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The cluster to which read/write requests should be routed.`,
						},
						"allow_transactional_writes": {
							Type:     schema.TypeBool,
							Optional: true,
							Description: `If true, CheckAndMutateRow and ReadModifyWriteRow requests are allowed by this app profile.
It is unsafe to send these requests to the same table/row/column in multiple clusters.`,
						},
					},
				},
				ExactlyOneOf: []string{"single_cluster_routing", "multi_cluster_routing_use_any"},
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique name of the requested app profile. Values are of the form 'projects/<project>/instances/<instance>/appProfiles/<appProfileId>'.`,
			},
			"multi_cluster_routing_cluster_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The set of clusters to route to. The order is ignored; clusters will be tried in order of distance. If left empty, all clusters are eligible.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"single_cluster_routing"},
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceBigtableAppProfileCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	descriptionProp, err := expandBigtableAppProfileDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	multiClusterRoutingUseAnyProp, err := expandBigtableAppProfileMultiClusterRoutingUseAny(d.Get("multi_cluster_routing_use_any"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("multi_cluster_routing_use_any"); !isEmptyValue(reflect.ValueOf(multiClusterRoutingUseAnyProp)) && (ok || !reflect.DeepEqual(v, multiClusterRoutingUseAnyProp)) {
		obj["multiClusterRoutingUseAny"] = multiClusterRoutingUseAnyProp
	}
	singleClusterRoutingProp, err := expandBigtableAppProfileSingleClusterRouting(d.Get("single_cluster_routing"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("single_cluster_routing"); !isEmptyValue(reflect.ValueOf(singleClusterRoutingProp)) && (ok || !reflect.DeepEqual(v, singleClusterRoutingProp)) {
		obj["singleClusterRouting"] = singleClusterRoutingProp
	}

	obj, err = resourceBigtableAppProfileEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := ReplaceVars(d, config, "{{BigtableBasePath}}projects/{{project}}/instances/{{instance}}/appProfiles?appProfileId={{app_profile_id}}&ignoreWarnings={{ignore_warnings}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new AppProfile: %#v", obj)
	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for AppProfile: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating AppProfile: %s", err)
	}
	if err := d.Set("name", flattenBigtableAppProfileName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}

	// Store the ID now
	id, err := ReplaceVars(d, config, "projects/{{project}}/instances/{{instance}}/appProfiles/{{app_profile_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating AppProfile %q: %#v", d.Id(), res)

	return resourceBigtableAppProfileRead(d, meta)
}

func resourceBigtableAppProfileRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := ReplaceVars(d, config, "{{BigtableBasePath}}projects/{{project}}/instances/{{instance}}/appProfiles/{{app_profile_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for AppProfile: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("BigtableAppProfile %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading AppProfile: %s", err)
	}

	if err := d.Set("name", flattenBigtableAppProfileName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading AppProfile: %s", err)
	}
	if err := d.Set("description", flattenBigtableAppProfileDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading AppProfile: %s", err)
	}
	if err := d.Set("multi_cluster_routing_use_any", flattenBigtableAppProfileMultiClusterRoutingUseAny(res["multiClusterRoutingUseAny"], d, config)); err != nil {
		return fmt.Errorf("Error reading AppProfile: %s", err)
	}
	if err := d.Set("single_cluster_routing", flattenBigtableAppProfileSingleClusterRouting(res["singleClusterRouting"], d, config)); err != nil {
		return fmt.Errorf("Error reading AppProfile: %s", err)
	}

	return nil
}

func resourceBigtableAppProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for AppProfile: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	descriptionProp, err := expandBigtableAppProfileDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	multiClusterRoutingUseAnyProp, err := expandBigtableAppProfileMultiClusterRoutingUseAny(d.Get("multi_cluster_routing_use_any"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("multi_cluster_routing_use_any"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, multiClusterRoutingUseAnyProp)) {
		obj["multiClusterRoutingUseAny"] = multiClusterRoutingUseAnyProp
	}
	singleClusterRoutingProp, err := expandBigtableAppProfileSingleClusterRouting(d.Get("single_cluster_routing"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("single_cluster_routing"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, singleClusterRoutingProp)) {
		obj["singleClusterRouting"] = singleClusterRoutingProp
	}

	obj, err = resourceBigtableAppProfileEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := ReplaceVars(d, config, "{{BigtableBasePath}}projects/{{project}}/instances/{{instance}}/appProfiles/{{app_profile_id}}?ignoreWarnings={{ignore_warnings}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating AppProfile %q: %#v", d.Id(), obj)
	updateMask := []string{}

	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}

	if d.HasChange("multi_cluster_routing_use_any") {
		updateMask = append(updateMask, "multiClusterRoutingUseAny")
	}

	if d.HasChange("single_cluster_routing") {
		updateMask = append(updateMask, "singleClusterRouting")
	}
	// updateMask is a URL parameter but not present in the schema, so ReplaceVars
	// won't set it
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	if d.HasChange("multi_cluster_routing_cluster_ids") && !stringInSlice(updateMask, "multiClusterRoutingUseAny") {
		updateMask = append(updateMask, "multiClusterRoutingUseAny")
	}

	// this api requires the body to define something for all values passed into
	// the update mask, however, multi-cluster routing and single-cluster routing
	// are conflicting, so we can't have them both in the update mask, despite
	// both of them registering as changing. thus, we need to remove whichever
	// one is not defined.
	newRouting, oldRouting := d.GetChange("multi_cluster_routing_use_any")
	if newRouting != oldRouting {
		for i, val := range updateMask {
			if val == "multiClusterRoutingUseAny" && newRouting.(bool) ||
				val == "singleClusterRouting" && oldRouting.(bool) {
				updateMask = append(updateMask[0:i], updateMask[i+1:]...)
				break
			}
		}
	}
	// updateMask is a URL parameter but not present in the schema, so ReplaceVars
	// won't set it
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequestWithTimeout(config, "PATCH", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating AppProfile %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating AppProfile %q: %#v", d.Id(), res)
	}

	return resourceBigtableAppProfileRead(d, meta)
}

func resourceBigtableAppProfileDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for AppProfile: %s", err)
	}
	billingProject = project

	url, err := ReplaceVars(d, config, "{{BigtableBasePath}}projects/{{project}}/instances/{{instance}}/appProfiles/{{app_profile_id}}?ignoreWarnings={{ignore_warnings}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting AppProfile %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequestWithTimeout(config, "DELETE", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "AppProfile")
	}

	log.Printf("[DEBUG] Finished deleting AppProfile %q: %#v", d.Id(), res)
	return nil
}

func resourceBigtableAppProfileImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := ParseImportId([]string{
		"projects/(?P<project>[^/]+)/instances/(?P<instance>[^/]+)/appProfiles/(?P<app_profile_id>[^/]+)",
		"(?P<project>[^/]+)/(?P<instance>[^/]+)/(?P<app_profile_id>[^/]+)",
		"(?P<instance>[^/]+)/(?P<app_profile_id>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := ReplaceVars(d, config, "projects/{{project}}/instances/{{instance}}/appProfiles/{{app_profile_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenBigtableAppProfileName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenBigtableAppProfileDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenBigtableAppProfileMultiClusterRoutingUseAny(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return false
	}

	if v.(map[string]interface{})["clusterIds"] == nil {
		return true
	}

	if len(v.(map[string]interface{})["clusterIds"].([]interface{})) > 0 {
		if err := d.Set("multi_cluster_routing_cluster_ids", v.(map[string]interface{})["clusterIds"]); err != nil {
			return true
		}
	}

	return true
}

func flattenBigtableAppProfileSingleClusterRouting(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["cluster_id"] =
		flattenBigtableAppProfileSingleClusterRoutingClusterId(original["clusterId"], d, config)
	transformed["allow_transactional_writes"] =
		flattenBigtableAppProfileSingleClusterRoutingAllowTransactionalWrites(original["allowTransactionalWrites"], d, config)
	return []interface{}{transformed}
}
func flattenBigtableAppProfileSingleClusterRoutingClusterId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenBigtableAppProfileSingleClusterRoutingAllowTransactionalWrites(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandBigtableAppProfileDescription(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigtableAppProfileMultiClusterRoutingUseAny(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil || !v.(bool) {
		return nil, nil
	}

	obj := bigtableadmin.MultiClusterRoutingUseAny{}

	clusterIds := d.Get("multi_cluster_routing_cluster_ids").([]interface{})

	for _, id := range clusterIds {
		obj.ClusterIds = append(obj.ClusterIds, id.(string))
	}

	return obj, nil
}

func expandBigtableAppProfileSingleClusterRouting(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedClusterId, err := expandBigtableAppProfileSingleClusterRoutingClusterId(original["cluster_id"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedClusterId); val.IsValid() && !isEmptyValue(val) {
		transformed["clusterId"] = transformedClusterId
	}

	transformedAllowTransactionalWrites, err := expandBigtableAppProfileSingleClusterRoutingAllowTransactionalWrites(original["allow_transactional_writes"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAllowTransactionalWrites); val.IsValid() && !isEmptyValue(val) {
		transformed["allowTransactionalWrites"] = transformedAllowTransactionalWrites
	}

	return transformed, nil
}

func expandBigtableAppProfileSingleClusterRoutingClusterId(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigtableAppProfileSingleClusterRoutingAllowTransactionalWrites(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func resourceBigtableAppProfileEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	// Instance is a URL parameter only, so replace self-link/path with resource name only.
	if err := d.Set("instance", GetResourceNameFromSelfLink(d.Get("instance").(string))); err != nil {
		return nil, fmt.Errorf("Error setting instance: %s", err)
	}
	return obj, nil
}
