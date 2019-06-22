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

//AssignmentAdd will add a new assignments to the datastore database
func AssignmentAdd(ctx context.Context, client *datastore.Client, asgn *dst.Assignment) (*datastore.Key, error) {

	// first check if there already exists this Assignment by asID:
	assignmnets, err := AssignmentGetByAsID(ctx, client, asgn.AsID)
	if err != nil {
		return nil, err
	}
	// if it has already the value, return key and error
	if len(assignmnets) > 0 {
		return &datastore.Key{ID: assignmnets[0].ID, Kind: dst.KindAssignments}, fmt.Errorf("asID allready exists. %d", assignmnets[0].ID)
	}
	// copy information into the datastore format
	n := &dst.Assignment{
		AsID:        asgn.AsID,
		Created:     time.Now(),
		Changed:     time.Now(),
		Description: asgn.Description,
		CpID:        asgn.CpID,
		DvID:        asgn.DvID,
		Level:       asgn.Level,
		Settings:    asgn.Settings,
		RawRequest:  asgn.RawRequest,
	}
	//do the insert into the database
	key := datastore.IncompleteKey(dst.KindAssignments, nil)
	return client.Put(ctx, key, n)
}

// AssignmentGetByAsID will return the list of assignments with the same asID
func AssignmentGetByAsID(ctx context.Context, client *datastore.Client, asID string) ([]*dst.Assignment, error) {
	log.Println("[AssignmentGetByAsID] will filter by asID:", asID)
	var assignments []*dst.Assignment
	// Create a query to fetch all Assignments filtered by asID
	query := datastore.NewQuery(dst.KindAssignments).Filter("asID =", asID)
	keys, err := client.GetAll(ctx, query, &assignments)
	if err != nil {
		return nil, err
	}
	log.Println("[AssignmentGetByAsID] Total keys returned", len(keys))
	// Set the ID field on each Assignment from the corresponding key.
	for i, key := range keys {
		assignments[i].ID = key.ID
	}
	return assignments, nil
}

// AssignmentsByCpID will return the list of assignments and its information with the same cpID
func AssignmentsByCpID(ctx context.Context, client *datastore.Client, cpID string) ([]*dst.Assignment, error) {
	log.Println("[AssignmentsByCpID] will filter by cpID:", cpID)
	var assignments []*dst.Assignment
	// Create a query to fetch all Actions entities, ordered by "created".
	query := datastore.NewQuery(dst.KindAssignments).Filter("cpID =", cpID)
	keys, err := client.GetAll(ctx, query, &assignments)
	if err != nil {
		return nil, err
	}
	log.Println("[AssignmentsByCpID] Total keys returned", len(keys))
	// Set the ID field on each assignments from the corresponding key.
	for i, key := range keys {
		assignments[i].ID = key.ID
	}
	// Get device and callpoint associated to each assignment
	for i, a := range assignments {
		// Load Callpoint information from database
		cps, err := CallpointGetByCpID(ctx, client, a.CpID)
		if err == nil && len(cps) > 0 {
			assignments[i].CallpointObj.CpID = cps[0].CpID
			assignments[i].CallpointObj.Label = cps[0].Label
			assignments[i].CallpointObj.Priority = cps[0].Priority
			assignments[i].CallpointObj.AbsAddress = cps[0].AbsAddress
			assignments[i].CallpointObj.Type = cps[0].Type
			assignments[i].CallpointObj.Icon = cps[0].Icon
			assignments[i].CallpointObj.Description = cps[0].Description
		}
		// Load Device information from database
		dvs, err := DeviceGetByDvID(ctx, client, a.DvID)
		if err == nil && len(dvs) > 0 {
			assignments[i].DeviceObj.DvID = dvs[0].DvID
			assignments[i].DeviceObj.Label = dvs[0].Label
			assignments[i].DeviceObj.Priority = dvs[0].Priority
			assignments[i].DeviceObj.Type = dvs[0].Type
			assignments[i].DeviceObj.Icon = dvs[0].Icon
			assignments[i].DeviceObj.Description = dvs[0].Description
			assignments[i].DeviceObj.IsTwoWay = dvs[0].IsTwoWay
			assignments[i].DeviceObj.Category = dvs[0].Category
			assignments[i].DeviceObj.Settings = dvs[0].Settings
			assignments[i].DeviceObj.RawRequest = dvs[0].RawRequest
		}
	}
	return assignments, nil
}

// AssignmentsToJSON prints the assignments into JSON to the given writer.
func AssignmentsToJSON(w io.Writer, asgs []*dst.Assignment) {
	const line = `%s
	{
		"ID": %d,
		"asID": "%s",				
		"cpID": "%s",
		"dvID": "%s",
		"description": "%s",
		"level": %d,
		"created": "%v",
		"changed": "%v",
		"settings": %s,
		"rawRequest": %s,
		"callpoint": %s,
		"device": %s
	}`
	// Use a tab writer to help make results pretty.
	tw := tabwriter.NewWriter(w, 4, 4, 1, ' ', 0) // Min cell size of 8.
	var term = ""
	var rawRequest, rawCallpoint, rawDevice string
	fmt.Fprintf(tw, "[\n")
	for _, a := range asgs {
		rawRequest = strings.TrimSpace(a.RawRequest)
		if rawRequest == "" {
			rawRequest = "null"
		}
		rawCallpoint = CallpointToJSONString(&a.CallpointObj)
		rawDevice = DeviceToJSONString(&a.DeviceObj)
		fmt.Fprintf(tw, line, term,
			a.ID,
			a.AsID,
			a.CpID,
			a.DvID,
			a.Description,
			a.Level,
			a.Created,
			a.Changed,
			a.Settings,
			rawRequest,
			rawCallpoint,
			rawDevice,
		)
		term = ","
	}
	fmt.Fprintf(tw, "\n]")
	tw.Flush()
}
