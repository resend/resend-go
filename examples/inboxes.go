package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v4"
)

func inboxesExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Create an inbox on one of your domains.
	created, err := client.Inboxes.CreateWithContext(ctx, &resend.CreateInboxRequest{
		EmailAddress: "support@acme.com",
		Name:         "Customer Support",
		FromName:     "Ada from Support",
	})
	if err != nil {
		panic(err)
	}
	inboxId := created.Id
	fmt.Println("Created inbox id: " + inboxId)

	// Get accepts either the inbox ID or the inbox email address.
	inbox, err := client.Inboxes.GetWithContext(ctx, "support@acme.com")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inbox %s has %d unread threads\n", inbox.EmailAddress, inbox.Unread)

	inboxes, err := client.Inboxes.ListWithContext(ctx, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d inboxes\n", len(inboxes.Data))

	_, err = client.Inboxes.UpdateWithContext(ctx, inboxId, &resend.UpdateInboxRequest{
		Name: "Support",
	})
	if err != nil {
		panic(err)
	}

	// Configure the inbox agent. Omitted fields keep their current value.
	tone := resend.InboxAgentToneConciseCasual
	_, err = client.Inboxes.Agent.UpdateWithContext(ctx, inboxId, &resend.UpdateInboxAgentRequest{
		Tone:    &tone,
		Actions: []resend.InboxAgentAction{resend.InboxAgentActionDraftReply},
	})
	if err != nil {
		panic(err)
	}
	agent, err := client.Inboxes.Agent.GetWithContext(ctx, inboxId)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Agent actions: %v\n", agent.Actions)

	// Labels are named, colored tags that belong to one inbox.
	label, err := client.Inboxes.Labels.CreateWithContext(ctx, inboxId, &resend.CreateInboxLabelRequest{
		Name:  "Urgent",
		Color: "#E93D82",
	})
	if err != nil {
		panic(err)
	}

	// List the unread threads in the Inbox folder.
	read := false
	threads, err := client.Inboxes.Threads.ListWithContext(ctx, inboxId, &resend.ListInboxThreadsOptions{
		Folders: []resend.InboxFolder{resend.InboxFolderInbox},
		Read:    &read,
	})
	if err != nil {
		panic(err)
	}

	// Search the threads of the inbox.
	results, err := client.Inboxes.Threads.SearchWithContext(ctx, inboxId, &resend.SearchInboxThreadsOptions{
		Query: "invoice -draft",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Search found %d threads\n", len(results.Data))

	if len(threads.Data) > 0 {
		threadId := threads.Data[0].Id

		// Label the thread and archive it.
		threadLabels, err := client.Inboxes.Threads.Labels.AddWithContext(ctx, inboxId, threadId, &resend.AddInboxThreadLabelsRequest{
			Labels: []string{label.Id},
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Thread has %d labels\n", len(threadLabels.Data))

		_, err = client.Inboxes.Threads.UpdateWithContext(ctx, inboxId, threadId, &resend.UpdateInboxThreadRequest{
			Folder: resend.InboxMoveThreadFolderArchive,
		})
		if err != nil {
			panic(err)
		}

		emails, err := client.Inboxes.Threads.Emails.ListWithContext(ctx, inboxId, threadId, nil)
		if err != nil {
			panic(err)
		}
		emailId := emails.Data[len(emails.Data)-1].Id

		// Reply right away. For a received message, the reply goes to its Reply-To addresses, or to
		// its sender. For a sent message, it goes to that message's `to` recipients.
		reply, err := client.Inboxes.Threads.Emails.ReplyWithContext(ctx, inboxId, threadId, emailId, &resend.ReplyInboxThreadEmailRequest{
			Text: "Thanks, we are on it.",
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("Sent reply email id: " + reply.Id)

		// Forward the message with an extra attachment.
		forward, err := client.Inboxes.Threads.Emails.ForwardWithContext(ctx, inboxId, threadId, emailId, &resend.ForwardInboxThreadEmailRequest{
			To:   []string{"billing@acme.com"},
			Text: "Flagging this for you.",
			Attachments: []*resend.Attachment{{
				Content:  []byte("Order 1041"),
				Filename: "order.txt",
			}},
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("Sent forward email id: " + forward.Id)

		// Or save a reply draft to review before you send it.
		draft, err := client.Inboxes.Drafts.CreateWithContext(ctx, inboxId, &resend.CreateInboxDraftRequest{
			Text:           "Draft reply.",
			ThreadId:       threadId,
			ReplyToEmailId: emailId,
		})
		if err != nil {
			panic(err)
		}
		// Send with options: IdempotencyKey. A retry with the same key within 24 hours returns the
		// first response.
		sent, err := client.Inboxes.Drafts.SendWithOptions(ctx, inboxId, draft.Id, &resend.SendInboxDraftRequest{
			Headers: map[string]string{"X-Ticket-Id": "4821"},
		}, &resend.SendInboxDraftOptions{
			IdempotencyKey: "send-draft-" + draft.Id,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("Sent draft email id: " + sent.Id)

		// Remove moves the thread to trash and schedules it for deletion in 30 days.
		_, err = client.Inboxes.Threads.RemoveWithContext(ctx, inboxId, threadId)
		if err != nil {
			panic(err)
		}
	}

	// A draft without a thread starts a new thread when you send it.
	draft, err := client.Inboxes.Drafts.CreateWithContext(ctx, inboxId, &resend.CreateInboxDraftRequest{
		To:      []string{"ada@example.org"},
		Subject: "Welcome",
		Text:    "Hello Ada.",
		// Draft attachments need Content, Filename and ContentType.
		Attachments: []*resend.InboxDraftAttachmentInput{{
			Content:     []byte("Welcome guide"),
			Filename:    "guide.txt",
			ContentType: "text/plain",
		}},
	})
	if err != nil {
		panic(err)
	}
	subject := "Welcome aboard"
	_, err = client.Inboxes.Drafts.UpdateWithContext(ctx, inboxId, draft.Id, &resend.UpdateInboxDraftRequest{
		Subject: &subject,
	})
	if err != nil {
		panic(err)
	}
	drafts, err := client.Inboxes.Drafts.ListWithContext(ctx, inboxId, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d drafts\n", len(drafts.Data))
	_, err = client.Inboxes.Drafts.RemoveWithContext(ctx, inboxId, draft.Id)
	if err != nil {
		panic(err)
	}

	_, err = client.Inboxes.Labels.RemoveWithContext(ctx, inboxId, label.Id)
	if err != nil {
		panic(err)
	}

	removed, err := client.Inboxes.RemoveWithContext(ctx, inboxId)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Removed inbox %s: %v\n", removed.Id, removed.Deleted)
}
