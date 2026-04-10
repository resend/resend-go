package main

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

func main() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Create and publish a template to use in the automation
	template, err := client.Templates.CreateWithContext(ctx, &resend.CreateTemplateRequest{
		Name:    "welcome-email",
		Subject: "Welcome!",
		Html:    "<strong>Welcome to our service!</strong>",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created template id: %s\n", template.Id)

	_, err = client.Templates.PublishWithContext(ctx, template.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Published template id: %s\n", template.Id)

	// Create an automation
	automation, err := client.Automations.CreateWithContext(ctx, &resend.CreateAutomationRequest{
		Name:   "Welcome Flow",
		Status: resend.AutomationStatusDisabled,
		Steps: []resend.AutomationStep{
			{
				Key:  "trigger_1",
				Type: resend.AutomationStepTypeTrigger,
				Config: map[string]any{
					"event_name": "user.created",
				},
			},
			{
				Key:  "send_1",
				Type: resend.AutomationStepTypeSendEmail,
				Config: map[string]any{
					"template": map[string]any{
						"id": template.Id,
					},
				},
			},
		},
		Connections: []resend.AutomationConnection{
			{From: "trigger_1", To: "send_1"},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created automation id: %s\n", automation.Id)

	// Get an automation
	retrieved, err := client.Automations.GetWithContext(ctx, automation.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Automation name: %s, status: %s\n", retrieved.Name, retrieved.Status)
	for _, step := range retrieved.Steps {
		fmt.Printf("  Step key: %s, type: %s\n", step.Key, step.Type)
	}
	for _, conn := range retrieved.Connections {
		fmt.Printf("  Connection: %s -> %s\n", conn.From, conn.To)
	}

	// List automations
	automations, err := client.Automations.ListWithContext(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total automations: %d\n", len(automations.Data))

	// List automations filtered by status
	status := resend.AutomationStatusEnabled
	limit := 10
	filtered, err := client.Automations.ListWithOptions(ctx, &resend.ListAutomationsOptions{
		Status: &status,
		Limit:  &limit,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Enabled automations: %d\n", len(filtered.Data))

	// Update an automation (enable it)
	updated, err := client.Automations.UpdateWithContext(ctx, automation.Id, &resend.UpdateAutomationRequest{
		Status: resend.AutomationStatusEnabled,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Updated automation id: %s\n", updated.Id)

	// Stop an automation
	stopped, err := client.Automations.StopWithContext(ctx, automation.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Stopped automation id: %s, status: %s\n", stopped.Id, stopped.Status)

	// List automation runs
	runs, err := client.Automations.ListRuns(automation.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total runs: %d\n", len(runs.Data))

	// List runs filtered by status
	runLimit := 5
	filteredRuns, err := client.Automations.ListRunsWithContext(ctx, automation.Id, &resend.ListAutomationRunsOptions{
		Status: []resend.AutomationRunStatus{
			resend.AutomationRunStatusCompleted,
			resend.AutomationRunStatusFailed,
		},
		Limit: &runLimit,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Completed/failed runs: %d\n", len(filteredRuns.Data))

	// Get a single run
	if len(runs.Data) > 0 {
		run, err := client.Automations.GetRunWithContext(ctx, automation.Id, runs.Data[0].Id)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Run id: %s, status: %s, steps: %d\n", run.Id, run.Status, len(run.Steps))
		for _, step := range run.Steps {
			fmt.Printf("  Step key: %s, type: %s, status: %s\n", step.Key, step.Type, step.Status)
		}
	}

	// Delete the automation
	deleted, err := client.Automations.RemoveWithContext(ctx, automation.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted automation id: %s, deleted: %v\n", deleted.Id, deleted.Deleted)

	// Create an event for the wait_for_event step
	event, err := client.Events.CreateWithContext(ctx, &resend.CreateEventRequest{
		Name: "user.verified",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created event id: %s\n", event.Id)

	// Create a multi-step automation using delay and wait_for_event
	// delay config: use "duration" (human-readable string) or "seconds" (number) — not both
	// wait_for_event config: use "timeout" (human-readable string) — timeout_seconds is not supported
	multiStep, err := client.Automations.CreateWithContext(ctx, &resend.CreateAutomationRequest{
		Name:   "Onboarding Flow",
		Status: resend.AutomationStatusDisabled,
		Steps: []resend.AutomationStep{
			{
				Key:  "trigger_1",
				Type: resend.AutomationStepTypeTrigger,
				Config: map[string]any{
					"event_name": "user.created",
				},
			},
			{
				Key:  "delay_1",
				Type: resend.AutomationStepTypeDelay,
				Config: map[string]any{
					"duration": "30 minutes",
				},
			},
			{
				Key:  "wait_1",
				Type: resend.AutomationStepTypeWaitForEvent,
				Config: map[string]any{
					"event_name": "user.verified",
					"timeout":    "1 hour",
				},
			},
			{
				Key:  "send_1",
				Type: resend.AutomationStepTypeSendEmail,
				Config: map[string]any{
					"template": map[string]any{
						"id": template.Id,
					},
				},
			},
		},
		Connections: []resend.AutomationConnection{
			{From: "trigger_1", To: "delay_1"},
			{From: "delay_1", To: "wait_1"},
			{From: "wait_1", To: "send_1", Type: resend.AutomationConnectionTypeEventReceived},
			{From: "wait_1", To: "send_1", Type: resend.AutomationConnectionTypeTimeout},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created multi-step automation id: %s\n", multiStep.Id)

	// Get to verify step keys and config shapes come back correctly
	retrieved2, err := client.Automations.GetWithContext(ctx, multiStep.Id)
	if err != nil {
		panic(err)
	}
	for _, step := range retrieved2.Steps {
		fmt.Printf("  Step key: %s, type: %s, config: %v\n", step.Key, step.Type, step.Config)
	}

	// Clean up
	_, err = client.Automations.RemoveWithContext(ctx, multiStep.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted multi-step automation\n")

	_, err = client.Events.RemoveWithContext(ctx, event.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted event id: %s\n", event.Id)

	// Clean up: delete the template
	removedTemplate, err := client.Templates.RemoveWithContext(ctx, template.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted template id: %s, deleted: %v\n", removedTemplate.Id, removedTemplate.Deleted)
}
