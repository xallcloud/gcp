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
/// Notifications
////////////////////////////////////////////////////////////////////////////////////////////////

//NotificationAdd method that
func NotificationAdd(ctx context.Context, client *datastore.Client, not *dst.Notification) (*dst.Notification, error) {

	// Unique ID for Notification ID
	uid := uuid.New()

	// copy to new record
	n := &dst.Notification{
		NtID:          uid.String(),
		AcID:          not.AcID,
		Priority:      not.Priority,
		Category:      not.Category,
		Destination:   not.Destination,
		Message:       not.Message,
		ResponseTitle: not.ResponseTitle,
		Options:       not.Options,
		Created:       time.Now(),
	}

	// do the insert
	key := datastore.IncompleteKey(dst.KindNotifications, nil)

	dskey, err := client.Put(ctx, key, n)

	if err != nil || key == nil {
		return nil, err
	}

	n.ID = dskey.ID

	return n, nil
}

// NotificationsGetByAcID will return the list of notifications with the same acID
func NotificationsGetByAcID(ctx context.Context, client *datastore.Client, acID string) ([]*dst.Notification, error) {
	// Create a query to fetch all Task entities, ordered by "created".

	log.Println("[NotificationsGetByAcID] will filter by acID:", acID)

	var notifications []*dst.Notification
	query := datastore.NewQuery(dst.KindNotifications).
		Filter("acID =", acID)

	log.Println("[NotificationsGetByAcID] will perform query")

	keys, err := client.GetAll(ctx, query, &notifications)
	if err != nil {
		return nil, err
	}

	log.Println("[NotificationsGetByAcID] Total keys returned", len(keys))

	// Set the ID field on each Notification from the corresponding key.
	for i, key := range keys {
		notifications[i].ID = key.ID
	}

	return notifications, nil
}

// NotificationsListAll returns all the notifications in ascending order of creation time.
func NotificationsListAll(ctx context.Context, client *datastore.Client) ([]*dst.Notification, error) {
	var notifications []*dst.Notification

	// Create a query to fetch all Notifications entities, ordered by "created".
	query := datastore.NewQuery(dst.KindNotifications).Order("created")
	keys, err := client.GetAll(ctx, query, &notifications)
	if err != nil {
		return nil, err
	}

	// Set the id field on each Notifications from the corresponding DataStore key.
	for i, key := range keys {
		notifications[i].ID = key.ID
	}

	return notifications, nil
}

// NotificationsToJSON prints the notifications into JSON to the given writer.
func NotificationsToJSON(w io.Writer, notifications []*dst.Notification) {
	const line = `%s
	{
		"ID": %d,
		"ntID": "%s",
		"acID": "%s",
		"priority": "%d",
		"category": "%s",
		"destination": "%s",
		"message": "%s",
		"responseTitle": "%s",
		"options": "%s",
		"created": "%v"
	}`

	// Use a tab writer to help make results pretty.
	tw := tabwriter.NewWriter(w, 4, 4, 1, ' ', 0) // Min cell size of 8.

	var term = ""
	fmt.Fprintf(tw, "[\n")
	for _, n := range notifications {

		fmt.Fprintf(tw, line, term,
			n.ID,
			n.NtID,
			n.AcID,
			n.Priority,
			n.Category,
			n.Destination,
			n.Message,
			n.ResponseTitle,
			n.Options,
			n.Created,
		)
		term = ","
	}
	fmt.Fprintf(tw, "\n]")
	tw.Flush()
}

// NotificationToJSONString prints the callpoints into JSON to the given writer.
func NotificationToJSONString(n *dst.Notification) string {
	const line = `
	{
		"ID": %d,
		"ntID": "%s",
		"acID": "%s",
		"priority": "%d",
		"category": "%s",
		"destination": "%s",
		"message": "%s",
		"responseTitle": "%s",
		"options": "%s",
		"created": "%v"
	}`

	r := fmt.Sprintf(line,
		n.ID,
		n.NtID,
		n.AcID,
		n.Priority,
		n.Category,
		n.Destination,
		n.Message,
		n.ResponseTitle,
		n.Options,
		n.Created,
	)

	return r
}
