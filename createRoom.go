package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getRoomCSV(filePath string) []RoomInfo {
	fmt.Printf("Importing the room information from %v \n", filePath)
	f, err := os.Open(filePath)

	check(err)

	reader := csv.NewReader(f)

	values, err := reader.ReadAll()
	check(err)

	//-1 to account for the headers
	toReturn := make([]RoomInfo, len(values)-1)
	count := 0

	for k := range values {
		//Bypass the header
		if k == 0 {
			continue
		}

		//For now assume order is going to hostname,ipaddress.
		//TODO: make this dynamic
		info := RoomInfo{Hostname: values[k][0], IPAddress: values[k][1], RoomName: values[k][2], Coordinates: values[k][3]}

		//Uncomment to see everything that's out
		//fmt.Println("IPAddress: ", info.IPAddress, " HOSTNAME: ", info.Hostname)
		toReturn[k-1] = info
		count++
	}

	//fmt.Printf("%+v \n", toReturn)

	fmt.Printf("Done importing. %v room(s) found.\n", count)
	return toReturn
}

//Import the configuration information from a JSON file
func importConfig(configPath string) Config {
	fmt.Printf("Importing the configuration information from %v\n", configPath)

	f, err := ioutil.ReadFile(configPath)
	check(err)

	var configurationData Config
	json.Unmarshal(f, &configurationData)

	//TODO: this needs to be more robust.
	if configurationData.ElasticSearchConfigInfoAddress == "" {
		panic("Invalid Configuration File.")
	}

	fmt.Printf("Done. Configuration data: \n %+v \n", configurationData)
	return configurationData
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

//Build the room object
func buildRoom(info RoomInfo, configuration Config) Room {

	//Get the signals from the configuration file
	f, err := ioutil.ReadFile(configuration.SignalDefinitionFile)
	check(err)
	var signals []Signal

	json.Unmarshal(f, &signals)

	fmt.Printf("Signals: %+v\n", signals)

	//Let's start building you
	var room Room
	var symbol Symbol

	symbol.ConnectInfo = info.IPAddress
	symbol.IPID = configuration.IPID
	symbol.SymbolName = info.Hostname + "Fusion"
	symbol.ProcessorID = info.IPAddress
	symbol.ProcessorName = info.Hostname
	symbol.Port = configuration.Port
	symbol.SecurePort = configuration.SecurePort
	symbol.Version = configuration.Version
	symbol.Signals = signals

	room.GroupwarePassword = configuration.GroupWarePassword
	room.GroupwareProviderType = configuration.GroupwareProviderType
	room.GroupwareURL = configuration.GroupwareURL
	room.GroupwareUsername = configuration.GroupWareUsername
	room.ParentNodeID = configuration.ParentNodeID
	room.RoomName = info.RoomName
	room.TimeZoneID = configuration.TimeZoneID

	//TODO: Where do we want to get the description from?
	room.Description = ""
	room.Symbols = []Symbol{symbol}

	return room
}

//get the room info for reporting from elastic search
func getRoominfoElasticSearch(address string) []RoomInfo {
	fmt.Printf("Importing the room information from Elastic Search at: %v \n", address)

	resp, err := http.Get(address)
	check(err)

	var response = ElasticSearchResponse{}
	bits, err := ioutil.ReadAll(resp.Body)
	check(err)

	json.Unmarshal(bits, &response)

	//.Printf("Received %+v\n", response)

	var toReturn []RoomInfo

	for index := range response.Hits.Hits {
		current := response.Hits.Hits[index].Source
		IPAddress := current.IPAddress
		Coordinates := current.Room.Coordinates
		Hostname := current.Hostname
		RoomName := current.Room.Building + " " + current.Room.NameOrNumber

		toReturn = append(toReturn, RoomInfo{IPAddress, Hostname, RoomName, Coordinates})
	}

	fmt.Printf("Done. Information %+v \n", toReturn)

	return toReturn
}

func sendRoom(room RoomInfo, config Config) {
	fmt.Printf("Sending room %v\n", room.Hostname)
	roomToSend := buildRoom(room, config)

	b, err := json.Marshal(roomToSend)
	check(err)

	//We should probably check the response here, but if it doesn't succeed err will go bad.
	resp, err := http.Post(config.FusionAddress, "application/json", bytes.NewBuffer(b))

	fmt.Printf("Stuff being sent \n %s \n\n", b)

	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
	}

	if resp.StatusCode == 200 {
		fmt.Printf("Success!")
	} else {
		fmt.Println("ERROR. Status: ", resp.StatusCode)
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%s\n", b)
	}

}

func addAllRooms(config Config, rooms []RoomInfo) {
	for k := range rooms {
		curRoom := rooms[k]
		sendRoom(curRoom, config)
	}
}

func main() {
	var ConfigFileLocation = flag.String("config", "./config.json", "The locaton of the config file.")
	var roomSource = flag.Int("source", 0, "The source of the room info to import into Fusion. 0 for elastic search. 1 for CSV. Default 0.")
	var roomInfo []RoomInfo

	flag.Parse()

	config := importConfig(*ConfigFileLocation)

	fmt.Println("RoomSource", *roomSource)

	if *roomSource == 0 {
		roomInfo = getRoominfoElasticSearch(config.ElasticSearchConfigInfoAddress)
	} else if *roomSource == 1 {
		roomInfo = getRoomCSV(config.CSVRoomInfoLocation)
	} else {
		roomInfo = getRoominfoElasticSearch(config.ElasticSearchConfigInfoAddress)
	}

	addAllRooms(config, roomInfo)
}
