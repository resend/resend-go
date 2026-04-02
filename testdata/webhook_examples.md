# Resend webhook example JSON payloads

Complete webhook event examples extracted from `llms-full.txt` and 4 missing pulled from the HTML docs. Each example retains the source JSON values and fields.

## Event index

- [`contact.created`](#contactcreated)
- [`contact.deleted`](#contactdeleted)
- [`contact.topics.updated`](#contacttopicsupdated)
- [`contact.updated`](#contactupdated)
- [`domain.created`](#domaincreated)
- [`domain.deleted`](#domaindeleted)
- [`domain.updated`](#domainupdated)
- [`email.bounced`](#emailbounced)
- [`email.clicked`](#emailclicked)
- [`email.complained`](#emailcomplained)
- [`email.delivered`](#emaildelivered)
- [`email.delivery_delayed`](#emaildelivery_delayed)
- [`email.failed`](#emailfailed)
- [`email.opened`](#emailopened)
- [`email.received`](#emailreceived)
- [`email.scheduled`](#emailscheduled)
- [`email.sent`](#emailsent)
- [`email.suppressed`](#emailsuppressed)
- [`inbox.created`](#inboxcreated)
- [`inbox.deleted`](#inboxdeleted)
- [`inbox.draft.created`](#inboxdraftcreated)
- [`inbox.draft.deleted`](#inboxdraftdeleted)
- [`inbox.draft.sent`](#inboxdraftsent)
- [`inbox.draft.updated`](#inboxdraftupdated)
- [`inbox.email.received`](#inboxemailreceived)
- [`inbox.email.sent`](#inboxemailsent)
- [`inbox.thread.assigned`](#inboxthreadassigned)
- [`inbox.thread.created`](#inboxthreadcreated)
- [`inbox.thread.folder.updated`](#inboxthreadfolderupdated)
- [`inbox.thread.labels.updated`](#inboxthreadlabelsupdated)
- [`inbox.thread.unassigned`](#inboxthreadunassigned)
- [`inbox.updated`](#inboxupdated)
- [`suppression.added`](#suppressionadded)
- [`suppression.removed`](#suppressionremoved)
- [`topic.created`](#topiccreated)
- [`topic.deleted`](#topicdeleted)
- [`topic.updated`](#topicupdated)

## contact.created

```json
{
  "type": "contact.created",
  "created_at": "2026-11-17T19:32:22.980Z",
  "data": {
    "id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
    "audience_id": "78261eea-8f8b-4381-83c6-79fa7120f1cf",
    "segment_ids": [
      "78261eea-8f8b-4381-83c6-79fa7120f1cf"
    ],
    "created_at": "2026-11-17T19:32:22.980Z",
    "updated_at": "2026-11-17T19:32:22.980Z",
    "email": "steve.wozniak@gmail.com",
    "first_name": null,
    "last_name": null,
    "unsubscribed": false
  }
}
```

## contact.deleted

```json
{
  "type": "contact.deleted",
  "created_at": "2026-11-17T19:32:22.980Z",
  "data": {
    "id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
    "audience_id": "78261eea-8f8b-4381-83c6-79fa7120f1cf",
    "segment_ids": [
      "78261eea-8f8b-4381-83c6-79fa7120f1cf"
    ],
    "created_at": "2026-11-10T15:11:59.110Z",
    "updated_at": "2026-11-17T19:32:22.980Z",
    "email": "steve.wozniak@gmail.com",
    "first_name": "Steve",
    "last_name": "Wozniak",
    "unsubscribed": false
  }
}
```

## contact.topics.updated

```json
{
  "type": "contact.topics.updated",
  "created_at": "2026-02-12T10:00:00.000Z",
  "data": {
    "email": "steve.wozniak@gmail.com",
    "topics": [
      {
        "id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
        "subscription": "opt_in"
      },
      {
        "id": "07d84122-7224-4881-9c31-1c048e204602",
        "subscription": "opt_out"
      }
    ]
  }
}
```

## contact.updated

```json
{
  "type": "contact.updated",
  "created_at": "2026-10-11T23:47:56.678Z",
  "data": {
    "id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
    "audience_id": "78261eea-8f8b-4381-83c6-79fa7120f1cf",
    "segment_ids": [
      "78261eea-8f8b-4381-83c6-79fa7120f1cf"
    ],
    "created_at": "2026-10-10T15:11:54.110Z",
    "updated_at": "2026-10-11T23:47:56.678Z",
    "email": "steve.wozniak@gmail.com",
    "first_name": "Steve",
    "last_name": "Wozniak",
    "unsubscribed": false
  }
}
```

## domain.created

```json
{
  "type": "domain.created",
  "created_at": "2026-11-17T19:32:22.980Z",
  "data": {
    "id": "d91cd9bd-1176-453e-8fc1-35364d380206",
    "name": "example.com",
    "status": "not_started",
    "created_at": "2026-04-26T20:21:26.347Z",
    "region": "us-east-1",
    "capabilities": {
      "sending": "enabled",
      "receiving": "disabled"
    },
    "records": [
      {
        "record": "SPF",
        "name": "send",
        "type": "MX",
        "ttl": "Auto",
        "status": "not_started",
        "value": "feedback-smtp.us-east-1.amazonses.com",
        "priority": 10
      },
      {
        "record": "SPF",
        "name": "send",
        "value": "\"v=spf1 include:amazonses.com ~all\"",
        "type": "TXT",
        "ttl": "Auto",
        "status": "not_started"
      },
      {
        "record": "DKIM",
        "name": "resend._domainkey",
        "value": "p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDsc4Lh8xilsngyKEgN2S84+21gn+x6SEXtjWvPiAAmnmggr5FWG42WnqczpzQ/mNblqHz4CDwUum6LtY6SdoOlDmrhvp5khA3cd661W9FlK3yp7+jVACQElS7d9O6jv8VsBbVg4COess3gyLE5RyxqF1vYsrEXqyM8TBz1n5AGkQIDAQA2",
        "type": "TXT",
        "status": "not_started",
        "ttl": "Auto"
      }
    ]
  }
}
```

## domain.deleted

```json
{
  "type": "domain.deleted",
  "created_at": "2026-11-17T19:32:22.980Z",
  "data": {
    "id": "d91cd9bd-1176-453e-8fc1-35364d380206",
    "name": "example.com",
    "status": "not_started",
    "created_at": "2026-04-26T20:21:26.347Z",
    "region": "us-east-1",
    "capabilities": {
      "sending": "enabled",
      "receiving": "disabled"
    },
    "records": [
      {
        "record": "SPF",
        "name": "send",
        "type": "MX",
        "ttl": "Auto",
        "status": "not_started",
        "value": "feedback-smtp.us-east-1.amazonses.com",
        "priority": 10
      },
      {
        "record": "SPF",
        "name": "send",
        "value": "\"v=spf1 include:amazonses.com ~all\"",
        "type": "TXT",
        "ttl": "Auto",
        "status": "not_started"
      },
      {
        "record": "DKIM",
        "name": "resend._domainkey",
        "value": "p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDsc4Lh8xilsngyKEgN2S84+21gn+x6SEXtjWvPiAAmnmggr5FWG42WnqczpzQ/mNblqHz4CDwUum6LtY6SdoOlDmrhvp5khA3cd661W9FlK3yp7+jVACQElS7d9O6jv8VsBbVg4COess3gyLE5RyxqF1vYsrEXqyM8TBz1n5AGkQIDAQA2",
        "type": "TXT",
        "status": "not_started",
        "ttl": "Auto"
      }
    ]
  }
}
```

## domain.updated

```json
{
  "type": "domain.updated",
  "created_at": "2026-11-17T19:32:22.980Z",
  "data": {
    "id": "d91cd9bd-1176-453e-8fc1-35364d380206",
    "name": "example.com",
    "status": "not_started",
    "created_at": "2026-04-26T20:21:26.347Z",
    "region": "us-east-1",
    "capabilities": {
      "sending": "enabled",
      "receiving": "enabled"
    },
    "records": [
      {
        "record": "SPF",
        "name": "send",
        "type": "MX",
        "ttl": "Auto",
        "status": "not_started",
        "value": "feedback-smtp.us-east-1.amazonses.com",
        "priority": 10
      },
      {
        "record": "SPF",
        "name": "send",
        "value": "\"v=spf1 include:amazonses.com ~all\"",
        "type": "TXT",
        "ttl": "Auto",
        "status": "not_started"
      },
      {
        "record": "DKIM",
        "name": "resend._domainkey",
        "value": "p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDsc4Lh8xilsngyKEgN2S84+21gn+x6SEXtjWvPiAAmnmggr5FWG42WnqczpzQ/mNblqHz4CDwUum6LtY6SdoOlDmrhvp5khA3cd661W9FlK3yp7+jVACQElS7d9O6jv8VsBbVg4COess3gyLE5RyxqF1vYsrEXqyM8TBz1n5AGkQIDAQA2",
        "type": "TXT",
        "status": "not_started",
        "ttl": "Auto"
      },
      {
        "name": "inbound.yourdomain.tld",
        "priority": 10,
        "record": "Receiving MX",
        "status": "pending",
        "ttl": "Auto",
        "type": "MX",
        "value": "inbound-smtp.us-east-1.amazonaws.com"
      }
    ]
  }
}
```

## email.bounced

```json
{
  "type": "email.bounced",
  "created_at": "2026-11-22T23:41:12.126Z",
  "data": {
    "broadcast_id": "8b146471-e88e-4322-86af-016cd36fd216",
    "created_at": "2026-11-22T23:41:11.894Z",
    "email_id": "56761188-7520-42d8-8898-ff6fc54ce618",
    "message_id": "<111-222-333@email.example.com>",
    "from": "Acme <onboarding@resend.dev>",
    "to": [
      "delivered@resend.dev"
    ],
    "subject": "Sending this example",
    "template_id": "43f68331-0622-4e15-8202-246a0388854b",
    "bounce": {
      "message": "The recipient's email address is on the suppression list because it has a recent history of producing hard bounces.",
      "subType": "Suppressed",
      "type": "Permanent"
    },
    "tags": {
      "category": "confirm_email"
    }
  }
}
```

## email.clicked

```json
{
  "type": "email.clicked",
  "created_at": "2026-11-22T23:41:12.126Z",
  "data": {
    "broadcast_id": "8b146471-e88e-4322-86af-016cd36fd216",
    "created_at": "2026-11-22T23:41:11.894Z",
    "email_id": "56761188-7520-42d8-8898-ff6fc54ce618",
    "message_id": "<111-222-333@email.example.com>",
    "from": "Acme <onboarding@resend.dev>",
    "to": [
      "delivered@resend.dev"
    ],
    "click": {
      "ipAddress": "122.115.53.11",
      "link": "https://resend.com",
      "timestamp": "2026-11-24T05:00:57.163Z",
      "userAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15"
    },
    "subject": "Sending this example",
    "template_id": "43f68331-0622-4e15-8202-246a0388854b",
    "tags": {
      "category": "confirm_email"
    }
  }
}
```

## email.complained

```json
{
  "type": "email.complained",
  "created_at": "2026-02-22T23:41:12.126Z",
  "data": {
    "broadcast_id": "8b146471-e88e-4322-86af-016cd36fd216",
    "created_at": "2026-02-22T23:41:11.894Z",
    "email_id": "56761188-7520-42d8-8898-ff6fc54ce618",
    "message_id": "<111-222-333@email.example.com>",
    "from": "Acme <onboarding@resend.dev>",
    "to": [
      "delivered@resend.dev"
    ],
    "subject": "Sending this example",
    "template_id": "43f68331-0622-4e15-8202-246a0388854b",
    "tags": {
      "category": "confirm_email"
    }
  }
}
```

## email.delivered

```json
{
  "type": "email.delivered",
  "created_at": "2026-02-22T23:41:12.126Z",
  "data": {
    "broadcast_id": "8b146471-e88e-4322-86af-016cd36fd216",
    "created_at": "2026-02-22T23:41:11.894Z",
    "email_id": "56761188-7520-42d8-8898-ff6fc54ce618",
    "message_id": "<111-222-333@email.example.com>",
    "from": "Acme <onboarding@resend.dev>",
    "to": [
      "delivered@resend.dev"
    ],
    "subject": "Sending this example",
    "template_id": "43f68331-0622-4e15-8202-246a0388854b",
    "tags": {
      "category": "confirm_email"
    }
  }
}
```

## email.delivery_delayed

```json
{
  "type": "email.delivery_delayed",
  "created_at": "2026-02-22T23:41:12.126Z",
  "data": {
    "broadcast_id": "8b146471-e88e-4322-86af-016cd36fd216",
    "created_at": "2026-02-22T23:41:11.894Z",
    "email_id": "56761188-7520-42d8-8898-ff6fc54ce618",
    "message_id": "<111-222-333@email.example.com>",
    "from": "Acme <onboarding@resend.dev>",
    "to": [
      "delivered@resend.dev"
    ],
    "subject": "Sending this example",
    "template_id": "43f68331-0622-4e15-8202-246a0388854b",
    "tags": {
      "category": "confirm_email"
    }
  }
}
```

## email.failed

```json
{
  "type": "email.failed",
  "created_at": "2026-11-22T23:41:12.126Z",
  "data": {
    "broadcast_id": "8b146471-e88e-4322-86af-016cd36fd216",
    "created_at": "2026-11-22T23:41:11.894Z",
    "email_id": "56761188-7520-42d8-8898-ff6fc54ce618",
    "from": "Acme <onboarding@resend.dev>",
    "to": [
      "delivered@resend.dev"
    ],
    "subject": "Sending this example",
    "template_id": "43f68331-0622-4e15-8202-246a0388854b",
    "failed": {
      "reason": "reached_daily_quota"
    },
    "tags": {
      "category": "confirm_email"
    }
  }
}
```

## email.opened

```json
{
  "type": "email.opened",
  "created_at": "2026-02-22T23:41:12.126Z",
  "data": {
    "broadcast_id": "8b146471-e88e-4322-86af-016cd36fd216",
    "created_at": "2026-02-22T23:41:11.894Z",
    "email_id": "56761188-7520-42d8-8898-ff6fc54ce618",
    "message_id": "<111-222-333@email.example.com>",
    "from": "Acme <onboarding@resend.dev>",
    "to": [
      "delivered@resend.dev"
    ],
    "subject": "Sending this example",
    "template_id": "43f68331-0622-4e15-8202-246a0388854b",
    "tags": {
      "category": "confirm_email"
    }
  }
}
```

## email.received

```json
{
  "type": "email.received",
  "created_at": "2026-02-22T23:41:12.126Z",
  "data": {
    "email_id": "56761188-7520-42d8-8898-ff6fc54ce618",
    "created_at": "2026-02-22T23:41:11.894Z",
    "from": "onboarding@resend.dev",
    "to": [
      "delivered@resend.dev"
    ],
    "bcc": [],
    "cc": [],
    "received_for": [
      "forwarded@example.com"
    ],
    "message_id": "<111-222-333@email.example.com>",
    "subject": "Sending this example",
    "attachments": [
      {
        "id": "2a0c9ce0-3112-4728-976e-47ddcd16a318",
        "filename": "avatar.png",
        "content_type": "image/png",
        "content_disposition": "inline",
        "content_id": "img001"
      }
    ]
  }
}
```

## email.scheduled

```json
{
  "type": "email.scheduled",
  "created_at": "2026-02-22T23:41:12.126Z",
  "data": {
    "broadcast_id": "8b146471-e88e-4322-86af-016cd36fd216",
    "created_at": "2026-02-22T23:41:11.894Z",
    "email_id": "56761188-7520-42d8-8898-ff6fc54ce618",
    "from": "Acme <onboarding@resend.dev>",
    "to": [
      "delivered@resend.dev"
    ],
    "subject": "Sending this example",
    "template_id": "43f68331-0622-4e15-8202-246a0388854b",
    "tags": {
      "category": "confirm_email"
    }
  }
}
```

## email.sent

```json
{
  "type": "email.sent",
  "created_at": "2026-02-22T23:41:12.126Z",
  "data": {
    "broadcast_id": "8b146471-e88e-4322-86af-016cd36fd216",
    "created_at": "2026-02-22T23:41:11.894Z",
    "email_id": "56761188-7520-42d8-8898-ff6fc54ce618",
    "message_id": "<111-222-333@email.example.com>",
    "from": "Acme <onboarding@resend.dev>",
    "to": [
      "delivered@resend.dev"
    ],
    "subject": "Sending this example",
    "template_id": "43f68331-0622-4e15-8202-246a0388854b",
    "tags": {
      "category": "confirm_email"
    }
  }
}
```

## email.suppressed

```json
{
  "type": "email.suppressed",
  "created_at": "2026-11-22T23:41:12.126Z",
  "data": {
    "broadcast_id": "8b146471-e88e-4322-86af-016cd36fd216",
    "created_at": "2026-11-22T23:41:11.894Z",
    "email_id": "56761188-7520-42d8-8898-ff6fc54ce618",
    "from": "Acme <onboarding@resend.dev>",
    "to": [
      "delivered@resend.dev"
    ],
    "subject": "Sending this example",
    "template_id": "43f68331-0622-4e15-8202-246a0388854b",
    "suppressed": {
      "message": "Resend has suppressed sending to this address because it is on the account-level suppression list. This does not count toward your bounce rate metric",
      "type": "OnAccountSuppressionList"
    },
    "tags": {
      "category": "confirm_email"
    }
  }
}
```

## inbox.created

```json
{
  "type": "inbox.created",
  "created_at": "2026-09-29T12:00:00.000Z",
  "data": {
    "source": "api",
    "inbox_id": "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1",
    "inbox": {
      "object": "inbox",
      "id": "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1",
      "name": "Customer Support",
      "email_address": "support@example.com",
      "domain_id": "d91cd9bd-1176-453e-8fc1-35364d380206",
      "receiving_address": null,
      "from_name": "Ada from Support",
      "unread": 0,
      "drafts": 0,
      "last_received": null,
      "created_at": "2026-09-29T12:00:00.000Z",
      "updated_at": "2026-09-29T12:00:00.214Z"
    }
  }
}
```

## inbox.deleted

```json
{
  "type": "inbox.deleted",
  "created_at": "2026-09-29T12:00:00.000Z",
  "data": {
    "source": "api",
    "inbox_id": "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1"
  }
}
```

## inbox.draft.created

```json
{
  "type": "inbox.draft.created",
  "created_at": "2026-09-29T12:00:00.000Z",
  "data": {
    "source": "api",
    "inbox_id": "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1",
    "thread_id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
    "draft_id": "9a2b7c4d-1e3f-4a5b-8c6d-0e1f2a3b4c5d",
    "thread": {
      "object": "inbox_thread",
      "id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
      "subject": "Question about my invoice",
      "folder": "inbox",
      "labels": [],
      "read": false
    },
    "draft": {
      "object": "inbox_draft",
      "id": "9a2b7c4d-1e3f-4a5b-8c6d-0e1f2a3b4c5d",
      "type": "reply",
      "to": [
        "steve.wozniak@gmail.com"
      ],
      "cc": [],
      "bcc": [],
      "subject": "Re: Question about my invoice",
      "thread_id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
      "reply_to_email_id": "4ef9a417-02e9-4d39-ad75-9611e0fcc33c",
      "email_id": null,
      "created_at": "2026-09-29T11:58:00.000Z",
      "updated_at": "2026-09-29T12:00:00.000Z",
      "attachments": []
    }
  }
}
```

## inbox.draft.deleted

```json
{
  "type": "inbox.draft.deleted",
  "created_at": "2026-09-29T12:00:00.000Z",
  "data": {
    "source": "dashboard",
    "inbox_id": "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1",
    "thread_id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
    "draft_id": "9a2b7c4d-1e3f-4a5b-8c6d-0e1f2a3b4c5d",
    "thread": {
      "object": "inbox_thread",
      "id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
      "subject": "Question about my invoice",
      "folder": "inbox",
      "labels": [],
      "read": false
    }
  }
}
```

## inbox.draft.sent

```json
{
  "type": "inbox.draft.sent",
  "created_at": "2026-09-29T12:00:00.000Z",
  "data": {
    "source": "api",
    "inbox_id": "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1",
    "thread_id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
    "email_id": "81ce9b8d-1b0a-4d02-a4d9-2a91f5ab6397",
    "draft_id": "9a2b7c4d-1e3f-4a5b-8c6d-0e1f2a3b4c5d",
    "thread": {
      "object": "inbox_thread",
      "id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
      "subject": "Question about my invoice",
      "folder": "inbox",
      "labels": [],
      "read": true
    },
    "draft": {
      "object": "inbox_draft",
      "id": "9a2b7c4d-1e3f-4a5b-8c6d-0e1f2a3b4c5d",
      "type": "reply",
      "to": [
        "steve.wozniak@gmail.com"
      ],
      "cc": [],
      "bcc": [],
      "subject": "Re: Question about my invoice",
      "thread_id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
      "reply_to_email_id": "4ef9a417-02e9-4d39-ad75-9611e0fcc33c",
      "email_id": "81ce9b8d-1b0a-4d02-a4d9-2a91f5ab6397",
      "created_at": "2026-09-29T11:58:00.000Z",
      "updated_at": "2026-09-29T12:00:00.000Z",
      "attachments": []
    }
  }
}
```

## inbox.draft.updated

```json
{
  "type": "inbox.draft.updated",
  "created_at": "2026-09-29T12:00:00.000Z",
  "data": {
    "source": "agent",
    "inbox_id": "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1",
    "thread_id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
    "draft_id": "9a2b7c4d-1e3f-4a5b-8c6d-0e1f2a3b4c5d",
    "thread": {
      "object": "inbox_thread",
      "id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
      "subject": "Question about my invoice",
      "folder": "inbox",
      "labels": [],
      "read": false
    },
    "draft": {
      "object": "inbox_draft",
      "id": "9a2b7c4d-1e3f-4a5b-8c6d-0e1f2a3b4c5d",
      "type": "reply",
      "to": [
        "steve.wozniak@gmail.com"
      ],
      "cc": [],
      "bcc": [],
      "subject": "Re: Question about my invoice",
      "thread_id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
      "reply_to_email_id": "4ef9a417-02e9-4d39-ad75-9611e0fcc33c",
      "email_id": null,
      "created_at": "2026-09-29T11:58:00.000Z",
      "updated_at": "2026-09-29T12:00:00.000Z",
      "attachments": []
    }
  }
}
```

## inbox.email.received

```json
{
  "type": "inbox.email.received",
  "created_at": "2026-09-29T12:00:00.000Z",
  "data": {
    "source": "system",
    "inbox_id": "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1",
    "thread_id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
    "email_id": "4ef9a417-02e9-4d39-ad75-9611e0fcc33c",
    "thread": {
      "object": "inbox_thread",
      "id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
      "subject": "Question about my invoice",
      "folder": "inbox",
      "labels": [],
      "read": false
    },
    "email": {
      "id": "4ef9a417-02e9-4d39-ad75-9611e0fcc33c",
      "direction": "inbound",
      "from": "Steve Wozniak <steve.wozniak@gmail.com>",
      "to": [
        "support@example.com"
      ],
      "cc": [],
      "bcc": [],
      "reply_to": [],
      "subject": "Question about my invoice",
      "message_id": "<CAF7c1f0a2e5d3b@mail.gmail.com>",
      "attachments": [
        {
          "id": "2a0c9ce0-3112-4728-976e-47ddcd16a318",
          "filename": "invoice.pdf",
          "size": 48213
        }
      ],
      "read": false,
      "received_at": "2026-09-29T12:00:00.000Z"
    }
  }
}
```

## inbox.email.sent

```json
{
  "type": "inbox.email.sent",
  "created_at": "2026-09-29T12:00:00.000Z",
  "data": {
    "source": "system",
    "inbox_id": "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1",
    "thread_id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
    "email_id": "4ef9a417-02e9-4d39-ad75-9611e0fcc33c",
    "thread": {
      "object": "inbox_thread",
      "id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
      "subject": "Question about my invoice",
      "folder": "inbox",
      "labels": [],
      "read": true
    },
    "email": {
      "id": "4ef9a417-02e9-4d39-ad75-9611e0fcc33c",
      "direction": "outbound",
      "from": "Ada from Support <support@example.com>",
      "to": [
        "steve.wozniak@gmail.com"
      ],
      "cc": [],
      "bcc": [],
      "reply_to": [],
      "subject": "Re: Question about my invoice",
      "message_id": "<0100019a4ef9a417@email.amazonses.com>",
      "attachments": [],
      "read": true,
      "received_at": "2026-09-29T12:00:00.000Z"
    }
  }
}
```

## inbox.thread.assigned

```json
{
  "type": "inbox.thread.assigned",
  "created_at": "2026-09-29T12:00:00.000Z",
  "data": {
    "assignee_email": "ada@example.com",
    "assigned_by_email": "grace@example.com",
    "previous_assignee_email": null,
    "source": "dashboard",
    "inbox_id": "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1",
    "thread_id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
    "thread": {
      "object": "inbox_thread",
      "id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
      "subject": "Question about my invoice",
      "folder": "inbox",
      "labels": [],
      "read": false
    }
  }
}
```

## inbox.thread.created

```json
{
  "type": "inbox.thread.created",
  "created_at": "2026-09-29T12:00:00.000Z",
  "data": {
    "source": "system",
    "inbox_id": "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1",
    "thread_id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
    "email_id": "4ef9a417-02e9-4d39-ad75-9611e0fcc33c",
    "thread": {
      "object": "inbox_thread",
      "id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
      "subject": "Question about my invoice",
      "folder": "inbox",
      "labels": [],
      "read": false
    },
    "email": {
      "id": "4ef9a417-02e9-4d39-ad75-9611e0fcc33c",
      "direction": "inbound",
      "from": "Steve Wozniak <steve.wozniak@gmail.com>",
      "to": [
        "support@example.com"
      ],
      "cc": [],
      "bcc": [],
      "reply_to": [],
      "subject": "Question about my invoice",
      "message_id": "<CAF7c1f0a2e5d3b@mail.gmail.com>",
      "attachments": [
        {
          "id": "2a0c9ce0-3112-4728-976e-47ddcd16a318",
          "filename": "invoice.pdf",
          "size": 48213
        }
      ],
      "read": false,
      "received_at": "2026-09-29T12:00:00.000Z"
    }
  }
}
```

## inbox.thread.folder.updated

```json
{
  "type": "inbox.thread.folder.updated",
  "created_at": "2026-09-29T12:00:00.000Z",
  "data": {
    "from": "inbox",
    "to": "archive",
    "source": "dashboard",
    "inbox_id": "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1",
    "thread_id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
    "thread": {
      "object": "inbox_thread",
      "id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
      "subject": "Question about my invoice",
      "folder": "archive",
      "labels": [],
      "read": true
    }
  }
}
```

## inbox.thread.labels.updated

```json
{
  "type": "inbox.thread.labels.updated",
  "created_at": "2026-09-29T12:00:00.000Z",
  "data": {
    "added": [
      {
        "id": "c0a8012e-3b4f-4d7a-9e21-5f6a7b8c9d0e",
        "name": "Billing",
        "color": "#F76B15"
      }
    ],
    "removed": [],
    "source": "dashboard",
    "inbox_id": "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1",
    "thread_id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
    "thread": {
      "object": "inbox_thread",
      "id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
      "subject": "Question about my invoice",
      "folder": "inbox",
      "labels": [
        {
          "id": "c0a8012e-3b4f-4d7a-9e21-5f6a7b8c9d0e",
          "name": "Billing",
          "color": "#F76B15"
        }
      ],
      "read": false
    }
  }
}
```

## inbox.thread.unassigned

```json
{
  "type": "inbox.thread.unassigned",
  "created_at": "2026-09-29T12:00:00.000Z",
  "data": {
    "previous_assignee_email": "ada@example.com",
    "source": "dashboard",
    "inbox_id": "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1",
    "thread_id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
    "thread": {
      "object": "inbox_thread",
      "id": "7c1f0a2e-5d3b-4e8a-9f61-2b8d4c6e1a90",
      "subject": "Question about my invoice",
      "folder": "inbox",
      "labels": [],
      "read": false
    }
  }
}
```

## inbox.updated

```json
{
  "type": "inbox.updated",
  "created_at": "2026-09-29T12:00:00.000Z",
  "data": {
    "source": "dashboard",
    "inbox_id": "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1",
    "inbox": {
      "object": "inbox",
      "id": "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1",
      "name": "Customer Support",
      "email_address": "support@example.com",
      "domain_id": "d91cd9bd-1176-453e-8fc1-35364d380206",
      "receiving_address": null,
      "from_name": "Ada from Support",
      "unread": 3,
      "drafts": 2,
      "last_received": "2026-09-29T11:42:08.000Z",
      "created_at": "2026-09-01T09:30:00.000Z",
      "updated_at": "2026-09-03T09:30:00.000Z"
    }
  }
}
```

## suppression.added

```json
{
  "type": "suppression.added",
  "created_at": "2026-11-17T19:32:22.980Z",
  "data": {
    "id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
    "email": "steve.wozniak@gmail.com",
    "origin": "bounce",
    "source_id": "4ef9a417-02e9-4d39-ad75-9611e0fcc33c",
    "created_at": "2026-11-17T19:32:22.980Z"
  }
}
```

## suppression.removed

```json
{
  "type": "suppression.removed",
  "created_at": "2026-11-17T19:32:22.980Z",
  "data": {
    "id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
    "email": "steve.wozniak@gmail.com",
    "origin": "manual",
    "source_id": null,
    "created_at": "2026-11-15T08:12:45.120Z"
  }
}
```

## topic.created

```json
{
  "type": "topic.created",
  "created_at": "2026-02-12T10:00:00.000Z",
  "data": {
    "id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
    "name": "Product Updates",
    "description": "New features and improvements",
    "default_subscription": "opt_in",
    "deleted": false,
    "created_at": "2026-02-12T10:00:00.000Z",
    "updated_at": "2026-02-12T10:00:00.000Z"
  }
}
```

## topic.deleted

```json
{
  "type": "topic.deleted",
  "created_at": "2026-02-12T10:00:00.000Z",
  "data": {
    "id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
    "name": "Product Updates",
    "description": "New features and improvements",
    "default_subscription": "opt_in",
    "deleted": true,
    "created_at": "2026-02-12T10:00:00.000Z",
    "updated_at": "2026-02-12T10:00:00.000Z"
  }
}
```

## topic.updated

```json
{
  "type": "topic.updated",
  "created_at": "2026-02-12T10:00:00.000Z",
  "data": {
    "id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
    "name": "Product Updates",
    "description": "New features and improvements",
    "default_subscription": "opt_in",
    "deleted": false,
    "created_at": "2026-02-12T10:00:00.000Z",
    "updated_at": "2026-02-12T10:00:00.000Z"
  }
}
```
