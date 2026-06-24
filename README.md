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
