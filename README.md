<!-- togo-header -->
<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/contacts-google</h1>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/contacts-google"><img src="https://pkg.go.dev/badge/github.com/togo-framework/contacts-google.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/contacts-google
```

<!-- /togo-header -->

# contacts-google — Google People API driver for togo contacts

A [togo](https://to-go.dev) **contacts** driver backed by the **Google People API**.
Imports/syncs a user's Google Contacts behind the togo `ContactsProvider` interface.

```bash
togo install togo-framework/contacts        # the base
togo install togo-framework/contacts-google # this driver
```

```env
CONTACTS_DRIVER=google
GOOGLE_CONTACTS_TOKEN=ya29.…   # a People-API OAuth access token
```

The OAuth flow (consent + token refresh, scope `contacts.readonly`) is handled by
your app; this driver consumes the access token and calls
`people/me/connections`. Once selected it powers the base API:

```go
svc, _ := contacts.FromKernel(k)
all, _ := svc.Sync(ctx)   // every Google contact, normalized
```

MIT

<!-- togo-sponsors -->
---

<div align="center">
  <h3>Premium sponsors</h3>
  <p>
    <a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp;
    <a href="https://one-studio.co"><strong>One Studio</strong></a>
  </p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- /togo-sponsors -->
