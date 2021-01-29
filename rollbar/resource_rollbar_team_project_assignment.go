package rollbar

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"strings"
)

func getId(d *schema.ResourceData) string {
	ids := []string{
		strconv.Itoa(d.Get("team_id").(int)),
		strconv.Itoa(d.Get("project_id").(int)),
	}

	return strings.Join(ids, "-")
}

func resourceRollbarTeamProjectAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRollbarTeamProjectAssignmentCreate,
		ReadContext:   resourceRollbarTeamProjectAssignmentRead,
		DeleteContext: resourceRollbarTeamProjectAssignmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceRollbarTeamProjectAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*Config).API
	teamID := d.Get("team_id").(int)
	projectID := d.Get("project_id").(int)

	_, _, err := client.Teams.AssignProject(teamID, projectID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(getId(d))

	return diags
}

func resourceRollbarTeamProjectAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*Config).API

	hasProject, _, err := client.Teams.HasProject(d.Get("team_id").(int), d.Get("project_id").(int))
	if err != nil {
		return diag.FromErr(err)
	}

	if hasProject != true {
		return diag.Errorf("%s", "not found")
	}

	d.SetId(getId(d))

	return diags
}

func resourceRollbarTeamProjectAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*Config).API

	_, err := client.Teams.RemoveProject(d.Get("team_id").(int), d.Get("project_id").(int))
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
