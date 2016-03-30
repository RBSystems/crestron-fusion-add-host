package main

//RoomInfo is Struct to be used when unmarshalling items from a csv document
type RoomInfo struct {
	IPAddress   string
	Hostname    string
	RoomName    string
	Coordinates string
}

//Signal is a struct built to reflect as signal item in JSON
type Signal struct {
	AttributeName   string
	AttributeID     string
	AttributeType   int
	JoinNumber      int
	LogicalOperator int
}

//Symbol is the struct built to reflect the Symbols item in JSON
type Symbol struct {
	ConnectInfo   string
	IPID          int
	SymbolName    string
	ProcessorID   string
	ProcessorName string
	Port          int
	SecurePort    int
	Version       string
	Signals       []Signal
}

//Room represents a room object in Fusion
type Room struct {
	RoomName              string
	Description           string
	ParentNodeID          string
	TimeZoneID            string
	Symbols               []Symbol
	GroupwareUsername     string
	GroupwarePassword     string
	GroupwareURL          string
	GroupwareProviderType string
}

//ElasticSearchResponse Object to reflect an elastic search response
type ElasticSearchResponse struct {
	Hits ElasticSearchHitWrapper
}

//ElasticSearchHitWrapper wrapper for hit object
type ElasticSearchHitWrapper struct {
	Total int
	Hits  []ElasitcSearchHit
}

//ElasitcSearchHit is one 'hit' in the system - basically each different item returned by the query.
type ElasitcSearchHit struct {
	Index  string
	Type   string
	ID     string
	Score  string
	Source ElasticSearchConfigSource `json:"_source"`
}

//ElasticSearchConfigSource is the source subdirectory of the hit - basically what we're actually
//putting into the config index.
type ElasticSearchConfigSource struct {
	MacAddress  string
	Description string
	Serial      string
	Hostname    string
	IPAddress   string
	Room        ElasticSearchRoomInfo
}

//ElasticSearchRoomInfo is the roomInfor in the ElasticSearchConfigSource
type ElasticSearchRoomInfo struct {
	Building     string
	NameOrNumber string
	Coordinates  string
	Floor        string
}

//Config represents the unmarshalled items in the config file
type Config struct {
	ElasticSearchConfigInfoAddress string
	SignalDefinitionFile           string
	GroupWarePassword              string
	GroupWareUsername              string
	GroupwareURL                   string
	GroupwareProviderType          string
	ParentNodeID                   string
	TimeZoneID                     string
	Version                        string
	FusionAddress                  string
	IPID                           int
	Port                           int
	SecurePort                     int
}
