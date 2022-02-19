package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LogicAppActionCustomResource struct{}

func TestAccLogicAppActionCustom_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_custom", "test")
	r := LogicAppActionCustomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppActionCustom_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_custom", "test")
	r := LogicAppActionCustomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_logic_app_action_custom"),
		},
	})
}

func TestAccLogicAppActionCustom_stepRemoval(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_custom", "add_two")
	r := LogicAppActionCustomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.threestep(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.twostep(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (LogicAppActionCustomResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	return actionExists(ctx, clients, state)
}

func (r LogicAppActionCustomResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_custom" "test" {
  name         = "action%d"
  logic_app_id = azurerm_logic_app_workflow.test.id

  body = <<BODY
{
    "description": "A variable to configure the auto expiration age in days. Configured in negative number. Default is -30 (30 days old).",
    "inputs": {
        "variables": [
            {
                "name": "ExpirationAgeInDays",
                "type": "Integer",
                "value": -30
            }
        ]
    },
    "runAfter": {},
    "type": "InitializeVariable"
}
BODY

}
`, r.template(data), data.RandomInteger)
}

func (r LogicAppActionCustomResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_custom" "import" {
  name         = azurerm_logic_app_action_custom.test.name
  logic_app_id = azurerm_logic_app_action_custom.test.logic_app_id
  body         = azurerm_logic_app_action_custom.test.body
}
`, r.basic(data))
}

func (r LogicAppActionCustomResource) threestep(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_custom" "init" {
  name         = "init"
  logic_app_id = azurerm_logic_app_workflow.test.id

  body = <<BODY
{
    "inputs": {
        "variables": [
            {
                "name": "var1",
                "type": "Integer",
                "value": 1
            }
        ]
    },
    "runAfter": {},
    "type": "InitializeVariable"
}
BODY
}

resource "azurerm_logic_app_action_custom" "add_one" {
  name         = "add_one"
  logic_app_id = azurerm_logic_app_workflow.test.id

  body = <<BODY
{
    "inputs": {
		"name": "var1",
		"value": 1
    },
    "runAfter": {
		"${azurerm_logic_app_action_custom.init.name}": ["Succeeded"]
	},
    "type": "IncrementVariable"
}
BODY
}

resource "azurerm_logic_app_action_custom" "add_two" {
  name         = "add_two"
  logic_app_id = azurerm_logic_app_workflow.test.id

  body = <<BODY
{
    "inputs": {
		"name": "var1",
		"value": 2
    },
    "runAfter": {
		"${azurerm_logic_app_action_custom.add_one.name}": ["Succeeded"]
	},
    "type": "IncrementVariable"
}
BODY
}
`, r.template(data))
}

func (r LogicAppActionCustomResource) twostep(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_custom" "init" {
  name         = "init"
  logic_app_id = azurerm_logic_app_workflow.test.id

  body = <<BODY
{
    "inputs": {
        "variables": [
            {
                "name": "var1",
                "type": "Integer",
                "value": 1
            }
        ]
    },
    "runAfter": {},
    "type": "InitializeVariable"
}
BODY
}

resource "azurerm_logic_app_action_custom" "add_two" {
  name         = "add_two"
  logic_app_id = azurerm_logic_app_workflow.test.id

  body = <<BODY
{
    "inputs": {
		"name": "var1",
		"value": 2
    },
    "runAfter": {
		"${azurerm_logic_app_action_custom.init.name}": ["Succeeded"]
	},
    "type": "IncrementVariable"
}
BODY
}
`, r.template(data))
}

func (LogicAppActionCustomResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
