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

////////////////////////////////////////////////////////////////////////////////////////////////
/// Events
////////////////////////////////////////////////////////////////////////////////////////////////

//EventAdd method that
func EventAdd(ctx context.Context, client *datastore.Client, ev *dst.Event) (*datastore.Key, error) {

	// Unique ID for Event ID
	uid := uuid.New()

	// copy to new record
	e := &dst.Event{
		EvID:          uid.String(),
		NtID:          ev.NtID,
		CpID:          ev.CpID,
		DvID:          ev.DvID,
		Visibility:    ev.Visibility,
		EvType:        ev.EvType,
		EvSubType:     ev.EvSubType + appVersion,
		EvDescription: ev.EvDescription,
		Created:       time.Now(),
	}

	// do the insert
	key := datastore.IncompleteKey(dst.KindEvents, nil)
	return client.Put(ctx, key, e)
}

// EventsGetByCpID will return the list of events with the same cpID
func EventsGetByCpID(ctx context.Context, client *datastore.Client, cpID string) ([]*dst.Event, error) {
	// Create a query to fetch all Task entities, ordered by "created".

	log.Println("[EventsGetByCpID] will filter by cpID:", cpID)

	var events []*dst.Event
	query := datastore.NewQuery(dst.KindEvents).
		Filter("cpID =", cpID)

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

// EventsGetByAcID will return the list of events with the same cpID
func EventsGetByAcID(ctx context.Context, client *datastore.Client, acID string) ([]*dst.Event, error) {
	// Create a query to fetch all Task entities, ordered by "created".

	log.Println("[EventsGetByAcID] will filter by acID:", acID)

	log.Println("[EventsGetByAcID] first get matching notification based on acID:", acID)

	// first check if there already exists this Action ID in notifications:
	nts, err := NotificationsGetByAcID(ctx, client, acID)
	if err != nil {
		return nil, err
	}

	// it has no matching action. Return empty
	if len(nts) <= 0 {
		return nil, nil
	}

	log.Println("[EventsGetByNtID] will filter by NtID:", nts[0].NtID)

	var events []*dst.Event

	var eventsFinal []*dst.Event

	// Set the ID field on each Event from the corresponding key.
	for _, not := range nts {

		events, err = EventsGetByNtID(ctx, client, not.NtID)
		if err != nil {
			return nil, err
		}

		//append to final events array
		if len(events) > 0 {
			for _, e := range events {
				eventsFinal = append(eventsFinal, e)
			}
		}

	}

	return eventsFinal, nil
}

// EventsGetByNtID will return the list of events with the same ntID
func EventsGetByNtID(ctx context.Context, client *datastore.Client, ntID string) ([]*dst.Event, error) {
	// Create a query to fetch all Task entities, ordered by "created".

	log.Println("[EventsGetByNtID] will filter by ntID:", ntID)

	var events []*dst.Event
	query := datastore.NewQuery(dst.KindEvents).
		Filter("ntID =", ntID).
		Order("created")

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
	var evs []*dst.Event

	// Create a query to fetch all Events entities, ordered by "created".
	query := datastore.NewQuery(dst.KindEvents).Order("created")
	keys, err := client.GetAll(ctx, query, &evs)
	if err != nil {
		return nil, err
	}

	// Set the id field on each Events from the corresponding DataStore key.
	for i, key := range keys {
		evs[i].ID = key.ID
	}

	return evs, nil
}

// EventsToJSON prints the events into JSON to the given writer.
func EventsToJSON(w io.Writer, evs []*dst.Event) {
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
	for _, d := range evs {

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
