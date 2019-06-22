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

//ActionAdd will add a new action to the datastore database
func ActionAdd(ctx context.Context, client *datastore.Client, ac *dst.Action) (*datastore.Key, error) {
	// first check if there already exists this Action by acID:
	actions, err := ActionGetByAcID(ctx, client, ac.AcID)
	if err != nil {
		return nil, err
	}
	// if it has already the value, return key and error
	if len(actions) > 0 {
		return &datastore.Key{ID: actions[0].ID, Kind: dst.KindActions}, fmt.Errorf("acID allready exists. %d", actions[0].ID)
	}
	// copy information into the datastore format
	n := &dst.Action{
		AcID:        ac.AcID,
		CpID:        ac.CpID,
		Action:      ac.Action,
		Description: ac.Description,
		Created:     time.Now(),
		RawRequest:  ac.RawRequest,
	}
	//do the insert into the database
	key := datastore.IncompleteKey(dst.KindActions, nil)
	return client.Put(ctx, key, n)
}

// ActionGetByAcID will return the list of actions with the same acID
func ActionGetByAcID(ctx context.Context, client *datastore.Client, acID string) ([]*dst.Action, error) {
	log.Println("[ActionGetByAcID] will filter by cpID:", acID)
	var actions []*dst.Action
	// Create a query to fetch all Actions filtered by acID
	query := datastore.NewQuery(dst.KindActions).Filter("acID =", acID)
	log.Println("[ActionGetByAcID] will perform query")
	keys, err := client.GetAll(ctx, query, &actions)
	if err != nil {
		return nil, err
	}
	log.Println("[ActionGetByAcID] Total keys returned", len(keys))
	// Set the ID field on each Action from the corresponding key.
	for i, key := range keys {
		actions[i].ID = key.ID
	}
	return actions, nil
}

// ActionsListAll returns all the actions in ascending order of creation time.
func ActionsListAll(ctx context.Context, client *datastore.Client) ([]*dst.Action, error) {
	log.Println("[ActionsListAll] Get all action records")
	var actions []*dst.Action
	// Create a query to fetch all Actions entities, ordered by "created".
	query := datastore.NewQuery(dst.KindActions).Order("created")
	keys, err := client.GetAll(ctx, query, &actions)
	if err != nil {
		return nil, err
	}
	log.Println("[ActionsListAll] Total keys returned", len(keys))
	// Set the id field on each Actions from the corresponding DataStore key.
	for i, key := range keys {
		actions[i].ID = key.ID
	}
	return actions, nil
}

// ActionsToJSON prints the actions into JSON to the given writer.
func ActionsToJSON(w io.Writer, actions []*dst.Action) {
	const line = `%s
	{
		"ID": %d,
		"acID": "%s",
		"cpID": "%s",
		"action": "%s",
		"description": "%s",
		"created": "%v",
		"rawRequest": %s
	}`
	// Use a tab writer to help make results pretty.
	tw := tabwriter.NewWriter(w, 4, 4, 1, ' ', 0) // Min cell size of 8.
	var term = ""
	var rawRequest string
	fmt.Fprintf(tw, "[\n")
	for _, d := range actions {
		rawRequest = strings.TrimSpace(d.RawRequest)
		if rawRequest == "" {
			rawRequest = "null"
		}
		fmt.Fprintf(tw, line, term,
			d.ID,
			d.AcID,
			d.CpID,
			d.Action,
			d.Description,
			d.Created,
			rawRequest,
		)
		term = ","
	}
	fmt.Fprintf(tw, "\n]")
	tw.Flush()
}

// ActionToJSONString prints a single action into JSON to the given writer.
func ActionToJSONString(d *dst.Action) string {
	const line = `
	{
		"ID": %d,
		"acID": "%s",
		"cpID": "%s",
		"action": "%s",
		"description": "%s",
		"created": "%v",
		"rawRequest": %s
	}`
	rawRequest := strings.TrimSpace(d.RawRequest)
	if rawRequest == "" {
		rawRequest = "null"
	}
	r := fmt.Sprintf(line,
		d.ID,
		d.AcID,
		d.CpID,
		d.Action,
		d.Description,
		d.Created,
		rawRequest,
	)
	return r
}
