package gcp

import (
	"context"
	"fmt"
	"io"
	"log"
	"text/tabwriter"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"

	dst "github.com/xallcloud/api/datastore"
)

//EventAdd will add a new Event to the datastore database
func EventAdd(ctx context.Context, client *datastore.Client, ev *dst.Event) (*datastore.Key, error) {
	// Generate a new Unique ID for the event
	uid := uuid.New()
	// copy information into the datastore format
	e := &dst.Event{
		EvID:          uid.String(),
		NtID:          ev.NtID,
		CpID:          ev.CpID,
		DvID:          ev.DvID,
		Visibility:    ev.Visibility,
		EvType:        ev.EvType,
		EvSubType:     ev.EvSubType,
		EvDescription: ev.EvDescription,
		Created:       time.Now(),
	}
	//do the insert into the database
	key := datastore.IncompleteKey(dst.KindEvents, nil)
	return client.Put(ctx, key, e)
}

// EventsGetByCpID will return the list of events with the same cpID
func EventsGetByCpID(ctx context.Context, client *datastore.Client, cpID string) ([]*dst.Event, error) {
	log.Println("[EventsGetByCpID] will filter by cpID:", cpID)
	var events []*dst.Event
	// Create a query to fetch all Events filtered by acID
	query := datastore.NewQuery(dst.KindEvents).Filter("cpID =", cpID)
	log.Println("[EventsGetByCpID] will perform query")
	keys, err := client.GetAll(ctx, query, &events)
	if err != nil {
		return nil, err
	}
	log.Println("[EventsGetByCpID] Total keys returned", len(keys))
	// Set the ID field on each Event from the corresponding key.
	for i, key := range keys {
		events[i].ID = key.ID
	}
	return events, nil
}

// EventsGetByAcID will return the list of events with the same acID
func EventsGetByAcID(ctx context.Context, client *datastore.Client, acID string) ([]*dst.Event, error) {
	log.Println("[EventsGetByAcID] will filter by acID:", acID)
	log.Println("[EventsGetByAcID] first get matching notification based on acID:", acID)
	notifications, err := NotificationsGetByAcID(ctx, client, acID)
	if err != nil {
		return nil, err
	}
	// it has no matching notification. Return empty
	if len(notifications) <= 0 {
		return nil, nil
	}
	log.Println("[EventsGetByNtID] will filter by NtID:", notifications[0].NtID)

	// will contain all events with the same Action ID
	var allEvents []*dst.Event
	//will contain the events of each notification
	var events []*dst.Event
	// for each notification, get all the events
	for _, not := range notifications {
		events, err = EventsGetByNtID(ctx, client, not.NtID)
		if err != nil {
			return nil, err
		}
		//append to allevents array each individual result
		if len(events) > 0 {
			for _, e := range events {
				allEvents = append(allEvents, e)
			}
		}
	}
	return allEvents, nil
}

// EventsGetByNtID will return the list of events with the same ntID
func EventsGetByNtID(ctx context.Context, client *datastore.Client, ntID string) ([]*dst.Event, error) {
	log.Println("[EventsGetByNtID] will filter by ntID:", ntID)
	var events []*dst.Event
	// Create a query to fetch all Events filtered by acID
	query := datastore.NewQuery(dst.KindEvents).Filter("ntID =", ntID).Order("created")
	log.Println("[EventsGetByNtID] will perform query")
	keys, err := client.GetAll(ctx, query, &events)
	if err != nil {
		return nil, err
	}
	log.Println("[EventsGetByNtID] Total keys returned", len(keys))
	// Set the ID field on each Event from the corresponding key.
	for i, key := range keys {
		events[i].ID = key.ID
	}
	return events, nil
}

// EventsListAll returns all the events in ascending order of creation time.
func EventsListAll(ctx context.Context, client *datastore.Client) ([]*dst.Event, error) {
	log.Println("[EventsListAll] Get all events records")
	var events []*dst.Event
	// Create a query to fetch all Events entities, ordered by "created".
	query := datastore.NewQuery(dst.KindEvents).Order("created")
	keys, err := client.GetAll(ctx, query, &events)
	if err != nil {
		return nil, err
	}
	log.Println("[EventsListAll] Total keys returned", len(keys))
	// Set the id field on each Events from the corresponding DataStore key.
	for i, key := range keys {
		events[i].ID = key.ID
	}
	return events, nil
}

// EventsToJSON prints the events into JSON to the given writer.
func EventsToJSON(w io.Writer, events []*dst.Event) {
	const line = `%s
	{
		"ID": %d,
		"evID": "%s",
		"ntID": "%s",
		"cpID": "%s",
		"dvID": "%s",
		"visibility": "%s",
		"evType": "%s",
		"evSubType": "%s",
		"evDescription": "%s",
		"created": "%v"
	}`
	// Use a tab writer to help make results pretty.
	tw := tabwriter.NewWriter(w, 4, 4, 1, ' ', 0) // Min cell size of 8.
	var term = ""
	fmt.Fprintf(tw, "[\n")
	for _, d := range events {
		fmt.Fprintf(tw, line, term,
			d.ID,
			d.EvID,
			d.NtID,
			d.CpID,
			d.DvID,
			d.Visibility,
			d.EvType,
			d.EvSubType,
			d.EvDescription,
			d.Created,
		)
		term = ","
	}
	fmt.Fprintf(tw, "\n]")
	tw.Flush()
}

// EventToJSONString prints the callpoints into JSON to the given writer.
func EventToJSONString(d *dst.Event) string {
	const line = `
	{
		"ID": %d,
		"evID": "%s",
		"ntID": "%s",
		"cpID": "%s",
		"dvID": "%s",
		"visibility": "%s",
		"evType": "%s",
		"evSubType": "%s",
		"evDescription": "%s",
		"created": "%v"
	}`
	r := fmt.Sprintf(line,
		d.ID,
		d.EvID,
		d.NtID,
		d.CpID,
		d.DvID,
		d.Visibility,
		d.EvType,
		d.EvSubType,
		d.EvDescription,
		d.Created,
	)
	return r
}
