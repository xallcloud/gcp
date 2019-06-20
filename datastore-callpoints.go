package gcp

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"text/tabwriter"
	"time"

	"cloud.google.com/go/datastore"

	dst "github.com/xallcloud/api/datastore"
)

////////////////////////////////////////////////////////////////////////////////////////////////
/// Callpoints
////////////////////////////////////////////////////////////////////////////////////////////////

//CallpointAdd method that
func CallpointAdd(ctx context.Context, client *datastore.Client, cp *dst.Callpoint) (*datastore.Key, error) {

	// first check if there already exists this Callpoint ID:
	cps, err := CallpointGetByCpID(ctx, client, cp.CpID)
	if err != nil {
		return nil, err
	}

	// if has already the value, return key and error
	if len(cps) > 0 {
		return &datastore.Key{ID: cps[0].ID, Kind: dst.KindCallpoints}, fmt.Errorf("cpID allready exists. %d", cps[0].ID)
	}

	// copy to new record
	n := &dst.Callpoint{
		CpID:        cp.CpID,
		Created:     time.Now(),
		AbsAddress:  cp.AbsAddress,
		Label:       cp.Label,
		Description: cp.Description,
		Type:        cp.Type,
		Priority:    cp.Priority,
		Icon:        cp.Icon,
		RawRequest:  cp.RawRequest,
	}

	// do the insert
	key := datastore.IncompleteKey(dst.KindCallpoints, nil)
	return client.Put(ctx, key, n)
}

// CallpointGetByCpID will return the list of callpoints with the same cpID
func CallpointGetByCpID(ctx context.Context, client *datastore.Client, cpID string) ([]*dst.Callpoint, error) {
	// Create a query to fetch all Task entities, ordered by "created".

	log.Println("[CallpointGetByCpID] will filter by cpID:", cpID)

	var callpoints []*dst.Callpoint
	query := datastore.NewQuery(dst.KindCallpoints).
		Filter("cpID =", cpID)

	log.Println("[CallpointGetByCpID] will perform query")

	keys, err := client.GetAll(ctx, query, &callpoints)
	if err != nil {
		return nil, err
	}

	log.Println("[CallpointGetByCpID] Total keys returned", len(keys))

	// Set the ID field on each Callpoint from the corresponding key.
	for i, key := range keys {
		callpoints[i].ID = key.ID
	}

	return callpoints, nil
}

// CallpointsListAll returns all the tasks in ascending order of creation time.
func CallpointsListAll(ctx context.Context, client *datastore.Client) ([]*dst.Callpoint, error) {
	var cps []*dst.Callpoint

	// Create a query to fetch all Task entities, ordered by "created".
	query := datastore.NewQuery(dst.KindCallpoints).Order("created")
	keys, err := client.GetAll(ctx, query, &cps)
	if err != nil {
		return nil, err
	}

	// Set the id field on each Callpoint from the corresponding DataStore key.
	for i, key := range keys {
		cps[i].ID = key.ID
	}

	return cps, nil
}

// CallpointDelete will delete a callpoint from the datastore
func CallpointDelete(ctx context.Context, client *datastore.Client, cpKeyID int64) error {
	return client.Delete(ctx, datastore.IDKey(dst.KindCallpoints, cpKeyID, nil))
}

// CallpointsToJSON prints the callpoints into JSON to the given writer.
func CallpointsToJSON(w io.Writer, cps []*dst.Callpoint) {
	const line = `%s
	{
		"ID": %d,
		"cpID": "%s",
		"label": "%s",
		"type": %d,
		"icon": "%s",				
		"description": "%s",		
		"priority": %d,
		"absAddress": "%s",		
		"created": "%v",
		"changed": "%v",
		"rawRequest": %s
	}`

	// Use a tab writer to help make results pretty.
	tw := tabwriter.NewWriter(w, 4, 4, 1, ' ', 0) // Min cell size of 8.

	var term = ""
	var rawRequest string
	fmt.Fprintf(tw, "[\n")
	for _, c := range cps {
		rawRequest = strings.TrimSpace(c.RawRequest)

		//log.Println("[JSONNotifications] parameter raw:", raw)

		if rawRequest == "" {
			rawRequest = "null"
		}

		fmt.Fprintf(tw, line, term,
			c.ID,
			c.CpID,
			c.Label,
			c.Type,
			c.Icon,
			c.Description,
			c.Priority,
			c.AbsAddress,
			c.Created,
			c.Changed,
			rawRequest,
		)
		term = ","
	}
	fmt.Fprintf(tw, "\n]")
	tw.Flush()
}

// CallpointToJSONString prints the callpoints into JSON to the given writer.
func CallpointToJSONString(c *dst.Callpoint) string {
	const line = `
	{
		"ID": %d,
		"cpID": "%s",
		"label": "%s",
		"type": %d,
		"icon": "%s",				
		"description": "%s",		
		"priority": %d,
		"absAddress": "%s",		
		"created": "%v",
		"changed": "%v",
		"rawRequest": %s
	}`

	rawRequest := strings.TrimSpace(c.RawRequest)

	if rawRequest == "" {
		rawRequest = "null"
	}

	r := fmt.Sprintf(line,
		c.ID,
		c.CpID,
		c.Label,
		c.Type,
		c.Icon,
		c.Description,
		c.Priority,
		c.AbsAddress,
		c.Created,
		c.Changed,
		rawRequest,
	)

	return r
}
