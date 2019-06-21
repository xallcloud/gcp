package gcp

//////////////////////////////////////////////////////////
// events: Visibility flags
//////////////////////////////////////////////////////////

// VisibilityServer is internal on events
const VisibilityServer = "server"

// VisibilityAll is set to all on events
const VisibilityAll = "all"

//////////////////////////////////////////////////////////
// events: Event Types flags
//////////////////////////////////////////////////////////

// EvTypeStart indicates the begining of a notification
const EvTypeStart = "start"

// EvTypeEnded indicates the ending of a notification
const EvTypeEnded = "ended"

// EvSubTypeStartStep1 indicates the begining of a notification
const EvSubTypeStartStep1 = "1"

// EvTypeServices indicates the begining of a notification
const EvTypeServices = "services"

// EvTypeDevices indicates the begining of a notification
const EvTypeDevices = "devices"

// EvSubTypeReaching that the device is being reached
const EvSubTypeReaching = "reaching"

// EvSubTypeDelivered that the device is being reached
const EvSubTypeDelivered = "delivered"

// EvSubTypeFailed that the device is being reached
const EvSubTypeFailed = "failed"

// EvSubTypeTimeout that the device is being reached
const EvSubTypeTimeout = "timeout"

// EvSubTypeReply that the device is being reached
const EvSubTypeReply = "reply"
