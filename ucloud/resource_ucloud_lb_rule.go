package ucloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func resourceUCloudLBRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudLBRuleCreate,
		Update: resourceUCloudLBRuleUpdate,
		Read:   resourceUCloudLBRuleRead,
		Delete: resourceUCloudLBRuleDelete,

		Schema: map[string]*schema.Schema{
			"load_balancer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"listener_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"backend_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				ForceNew: true,
			},

			"domain": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"path"},
			},

			"path": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"domain"},
			},
		},
	}
}

func resourceUCloudLBRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*UCloudClient).ulbconn

	req := conn.NewCreatePolicyRequest()
	req.ULBId = ucloud.String(d.Get("load_balancer_id").(string))
	req.VServerId = ucloud.String(d.Get("listener_id").(string))
	req.BackendId = ifaceToStringSlice(d.Get("backend_ids"))

	if val, ok := d.GetOk("domain"); ok {
		req.Type = ucloud.String("Domain")
		req.Match = ucloud.String(val.(string))
	} else if val, ok := d.GetOk("path"); ok {
		req.Type = ucloud.String("Path")
		req.Match = ucloud.String(val.(string))
	} else {
		return fmt.Errorf("error in create lb rule, shoule set one of domain and path")
	}

	resp, err := conn.CreatePolicy(req)

	if err != nil {
		return fmt.Errorf("error in create lb rule, %s", err)
	}

	d.SetId(resp.PolicyId)

	time.Sleep(5 * time.Second)

	return resourceUCloudLBRuleUpdate(d, meta)
}

func resourceUCloudLBRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)

	isChanged := false
	conn := meta.(*UCloudClient).ulbconn
	req := conn.NewUpdatePolicyRequest()
	req.ULBId = ucloud.String(d.Get("load_balancer_id").(string))
	req.VServerId = ucloud.String(d.Get("listener_id").(string))
	req.BackendId = ifaceToStringSlice(d.Get("backend_ids"))
	req.PolicyId = ucloud.String(d.Id())

	if d.HasChange("domain") && !d.IsNewResource() {
		isChanged = true
		req.Type = ucloud.String("Domain")
		req.Match = ucloud.String(d.Get("domain").(string))
		d.SetPartial("domain")
	}

	if d.HasChange("path") && !d.IsNewResource() {
		isChanged = true
		req.Type = ucloud.String("Path")
		req.Match = ucloud.String(d.Get("path").(string))
		d.SetPartial("path")
	}

	if isChanged {
		_, err := conn.UpdatePolicy(req)

		if err != nil {
			return fmt.Errorf("do %s failed in update lb rule %s, %s", "UpdatePolicy", d.Id(), err)
		}

		time.Sleep(5 * time.Second)
	}

	d.Partial(false)

	return resourceUCloudLBRuleRead(d, meta)
}

func resourceUCloudLBRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)

	lbId := d.Get("load_balancer_id").(string)
	listenerId := d.Get("listener_id").(string)

	policySet, err := client.describePolicyById(lbId, listenerId, d.Id())

	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("do %s failed in read lb rule %s, %s", "DescribeVServer", d.Id(), err)
	}

	if policySet.Type == "Path" {
		d.Set("path", policySet.Match)
	}

	if policySet.Type == "Domain" {
		d.Set("domain", policySet.Match)
	}

	return nil
}

func resourceUCloudLBRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.ulbconn

	lbId := d.Get("load_balancer_id").(string)
	listenerId := d.Get("listener_id").(string)

	req := conn.NewDeletePolicyRequest()
	req.VServerId = ucloud.String(listenerId)
	req.PolicyId = ucloud.String(d.Id())

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.DeletePolicy(req); err != nil {
			return resource.NonRetryableError(fmt.Errorf("error in delete lb rule %s, %s", d.Id(), err))
		}

		_, err := client.describePolicyById(lbId, listenerId, d.Id())

		if err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("do %s failed in delete lb rule %s, %s", "DescribeVServer", d.Id(), err))
		}

		return resource.RetryableError(fmt.Errorf("delete lb rule but still exists"))
	})
}
