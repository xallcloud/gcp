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

//DeviceAdd will add a new device to the datastore database
func DeviceAdd(ctx context.Context, client *datastore.Client, dv *dst.Device) (*datastore.Key, error) {
	// first check if there already exists this Device by dvID:
	devices, err := DeviceGetByDvID(ctx, client, dv.DvID)
	if err != nil {
		return nil, err
	}
	// if it has already the value, return key and error
	if len(devices) > 0 {
		return &datastore.Key{ID: devices[0].ID, Kind: dst.KindCallpoints}, fmt.Errorf("dvID allready exists. %d", devices[0].ID)
	}
	// copy information into the datastore format
	n := &dst.Device{
		DvID:        dv.DvID,
		Created:     time.Now(),
		Label:       dv.Label,
		Description: dv.Description,
		Type:        dv.Type,
		Priority:    dv.Priority,
		Icon:        dv.Icon,
		IsTwoWay:    dv.IsTwoWay,
		Category:    dv.Category,
		Destination: dv.Destination,
		Settings:    dv.Settings,
		RawRequest:  dv.RawRequest,
	}
	//do the insert into the database
	key := datastore.IncompleteKey(dst.KindDevices, nil)
	return client.Put(ctx, key, n)
}

// DeviceGetByDvID will return the list of devices with the same dvID
func DeviceGetByDvID(ctx context.Context, client *datastore.Client, dvID string) ([]*dst.Device, error) {
	log.Println("[DeviceGetByDvID] will filter by cpID:", dvID)
	var devices []*dst.Device
	// Create a query to fetch all Devices filtered by dvID
	query := datastore.NewQuery(dst.KindDevices).Filter("dvID =", dvID)
	log.Println("[DeviceGetByDvID] will perform query")
	keys, err := client.GetAll(ctx, query, &devices)
	if err != nil {
		return nil, err
	}
	log.Println("[DeviceGetByDvID] Total keys returned", len(keys))
	// Set the ID field on each Callpoint from the corresponding key.
	for i, key := range keys {
		devices[i].ID = key.ID
	}
	return devices, nil
}

// DevicesListAll returns all the devices in ascending order of creation time.
func DevicesListAll(ctx context.Context, client *datastore.Client) ([]*dst.Device, error) {
	log.Println("[DevicesListAll] Get all devices records")
	var devices []*dst.Device
	// Create a query to fetch all Devices entities, ordered by "created".
	query := datastore.NewQuery(dst.KindDevices).Order("created")
	keys, err := client.GetAll(ctx, query, &devices)
	if err != nil {
		return nil, err
	}
	log.Println("[DevicesListAll] Total keys returned", len(keys))
	// Set the id field on each Devices from the corresponding DataStore key.
	for i, key := range keys {
		devices[i].ID = key.ID
	}
	return devices, nil
}

// DeviceDelete will delete a device from the datastore
func DeviceDelete(ctx context.Context, client *datastore.Client, dvKeyID int64) error {
	return client.Delete(ctx, datastore.IDKey(dst.KindDevices, dvKeyID, nil))
}

// DevicesToJSON prints the devices into JSON to the given writer.
func DevicesToJSON(w io.Writer, devices []*dst.Device) {
	const line = `%s
	{
		"ID": %d,
		"dvID": "%s",
		"label": "%s",
		"type": %d,
		"icon": "%s",
		"description": "%s",
		"isTwoWay": %s,
		"category": "%s",
		"destination": "%s",
		"priority": %d,
		"created": "%v",
		"changed": "%v",		
		"settings": %s,
		"rawRequest": %s
	}`
	// Use a tab writer to help make results pretty.
	tw := tabwriter.NewWriter(w, 4, 4, 1, ' ', 0) // Min cell size of 8.
	var term = ""
	var rawRequest, isTwoWayString string
	fmt.Fprintf(tw, "[\n")
	for _, d := range devices {
		rawRequest = strings.TrimSpace(d.RawRequest)
		if rawRequest == "" {
			rawRequest = "null"
		}
		if d.IsTwoWay {
			isTwoWayString = "true"
		} else {
			isTwoWayString = "false"
		}
		fmt.Fprintf(tw, line, term,
			d.ID,
			d.DvID,
			d.Label,
			d.Type,
			d.Icon,
			d.Description,
			isTwoWayString,
			d.Category,
			d.Destination,
			d.Priority,
			d.Created,
			d.Changed,
			d.Settings,
			rawRequest,
		)
		term = ","
	}
	fmt.Fprintf(tw, "\n]")
	tw.Flush()
}

// DeviceToJSONString prints a single device into JSON to the given writer.
func DeviceToJSONString(d *dst.Device) string {
	const line = `
	{
		"ID": %d,
		"dvID": "%s",
		"label": "%s",
		"type": %d,
		"icon": "%s",
		"description": "%s",
		"isTwoWay": %s,
		"category": "%s",
		"destination": "%s",
		"priority": %d,
		"created": "%v",
		"changed": "%v",		
		"settings": %s,
		"rawRequest": %s
	}`
	rawRequest := strings.TrimSpace(d.RawRequest)
	if rawRequest == "" {
		rawRequest = "null"
	}
	isTwoWayString := ""
	if d.IsTwoWay {
		isTwoWayString = "true"
	} else {
		isTwoWayString = "false"
	}
	r := fmt.Sprintf(line,
		d.ID,
		d.DvID,
		d.Label,
		d.Type,
		d.Icon,
		d.Description,
		isTwoWayString,
		d.Category,
		d.Destination,
		d.Priority,
		d.Created,
		d.Changed,
		d.Settings,
		rawRequest,
	)
	return r
}
