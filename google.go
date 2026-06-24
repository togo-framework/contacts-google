// Package google is a Google People API driver for togo contacts. Blank-import
// it and set CONTACTS_DRIVER=google plus a People-API OAuth access token.
//
// The OAuth flow (obtaining/refreshing the access token) is the app's job; this
// driver consumes a token from GOOGLE_CONTACTS_TOKEN (or GOOGLE_ACCESS_TOKEN).
package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/togo-framework/contacts"
	"github.com/togo-framework/togo"
)

const apiBase = "https://people.googleapis.com/v1"
const personFields = "names,emailAddresses,phoneNumbers,photos,organizations"

func init() {
	contacts.RegisterDriver("google", func(k *togo.Kernel) (contacts.ContactsProvider, error) {
		tok := os.Getenv("GOOGLE_CONTACTS_TOKEN")
		if tok == "" {
			tok = os.Getenv("GOOGLE_ACCESS_TOKEN")
		}
		if tok == "" {
			return nil, fmt.Errorf("contacts-google: GOOGLE_CONTACTS_TOKEN (a People-API OAuth access token) not set")
		}
		return &provider{token: tok, client: &http.Client{Timeout: 20 * time.Second}}, nil
	})
}

type provider struct {
	token  string
	client *http.Client
}

type apiPerson struct {
	ResourceName   string `json:"resourceName"`
	Names          []struct{ DisplayName string `json:"displayName"` } `json:"names"`
	EmailAddresses []struct{ Value string `json:"value"` }             `json:"emailAddresses"`
	PhoneNumbers   []struct{ Value string `json:"value"` }             `json:"phoneNumbers"`
	Photos         []struct{ URL string `json:"url"` }                 `json:"photos"`
	Organizations  []struct{ Name string `json:"name"` }               `json:"organizations"`
}

func (p *provider) toContact(a apiPerson) contacts.Contact {
	c := contacts.Contact{ID: a.ResourceName, Source: "google"}
	if len(a.Names) > 0 {
		c.Name = a.Names[0].DisplayName
	}
	for _, e := range a.EmailAddresses {
		c.Emails = append(c.Emails, e.Value)
	}
	for _, ph := range a.PhoneNumbers {
		c.Phones = append(c.Phones, ph.Value)
	}
	if len(a.Photos) > 0 {
		c.Photo = a.Photos[0].URL
	}
	if len(a.Organizations) > 0 {
		c.Org = a.Organizations[0].Name
	}
	return c
}

func (p *provider) List(ctx context.Context, pageToken string) ([]contacts.Contact, string, error) {
	q := url.Values{}
	q.Set("personFields", personFields)
	q.Set("pageSize", "200")
	if pageToken != "" {
		q.Set("pageToken", pageToken)
	}
	var out struct {
		Connections   []apiPerson `json:"connections"`
		NextPageToken string      `json:"nextPageToken"`
	}
	if err := p.get(ctx, apiBase+"/people/me/connections?"+q.Encode(), &out); err != nil {
		return nil, "", err
	}
	cs := make([]contacts.Contact, 0, len(out.Connections))
	for _, a := range out.Connections {
		cs = append(cs, p.toContact(a))
	}
	return cs, out.NextPageToken, nil
}

func (p *provider) Get(ctx context.Context, id string) (*contacts.Contact, error) {
	var a apiPerson
	if err := p.get(ctx, apiBase+"/"+id+"?personFields="+personFields, &a); err != nil {
		return nil, err
	}
	c := p.toContact(a)
	return &c, nil
}

func (p *provider) get(ctx context.Context, u string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+p.token)
	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("contacts-google: People API %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}
