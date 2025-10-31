# Changelog

## 3.0.0 - 2025-01-XX

This is a major release with breaking changes to the Contacts API and the addition of several new APIs (Topics, Templates, Webhooks, Inbound, Segments, Contact Properties, Contact Topics). The Contacts API has been redesigned to use struct-based arguments for better ergonomics and to support global contacts.

**Note:** The breaking changes are isolated to the Contacts API only. If you only use the SDK for sending emails, you can upgrade to v3.0.0 without any code changes. The Emails API remains unchanged and fully backward compatible. The Audiences API is deprecated but maintained for backward compatibility.

- ⚠️ Change `Contacts.Get` to accept `*GetContactOptions` instead of positional `audienceId` and `id` parameters
- ⚠️ Change `Contacts.GetWithContext` to accept `*GetContactOptions` instead of positional `audienceId` and `id` parameters
- ⚠️ Change `Contacts.List` to accept `*ListContactsOptions` instead of positional `audienceId` parameter
- ⚠️ Change `Contacts.ListWithContext` to accept `*ListContactsOptions` instead of positional `audienceId` parameter
- ⚠️ Change `Contacts.Remove` to accept `*RemoveContactOptions` instead of positional `audienceId` and `id` parameters
- ⚠️ Change `Contacts.RemoveWithContext` to accept `*RemoveContactOptions` instead of positional `audienceId` and `id` parameters
- Add support for new `Topics` service for managing email topics and subscription preferences [#85](https://github.com/resend/resend-go/pull/85)
- Add support for new `Templates` service for managing email templates [#84](https://github.com/resend/resend-go/pull/84)
- Add support for `TemplateId` and `TemplateAlias` on `SendEmailRequest` for sending emails with templates [#84](https://github.com/resend/resend-go/pull/84)
- Add support for new `Webhooks` service for managing webhook endpoints [#86](https://github.com/resend/resend-go/pull/86)
- Add support for new `Receiving` service for managing inbound email handling and attachments [#83](https://github.com/resend/resend-go/pull/83)
- Add support for new `Segments` service for organizing contacts into segments (replaces `Audiences`)
- Add support for `SegmentId` on `CreateBroadcastRequest` and `UpdateBroadcastRequest` for targeting segments in broadcasts
- Deprecate `Audiences` service in favor of `Segments` service (backward compatible wrapper maintained)
- Add support for new `Contacts.Properties` sub-service for managing custom contact properties [#90](https://github.com/resend/resend-go/pull/90)
- Add support for new `Contacts.Topics` sub-service for managing contact topic subscriptions [#90](https://github.com/resend/resend-go/pull/90)
- Add support for global contacts by making `AudienceId` optional on `CreateContactRequest` and `UpdateContactRequest`
- Add support for `Properties` on `CreateContactRequest`, `UpdateContactRequest`, and `Contact` for custom key-value pairs on global contacts
- Add support for new `Contacts.Segments` sub-service for managing contact membership in segments
