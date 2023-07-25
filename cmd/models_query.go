package cmd

import "time"

type Queries struct {
	Queries []QueryData `json:"queries"`
}

type QueryData struct {
	ProductTypeCodes []string `json:"productTypeCodes"`
	Filter           struct {
		Discriminator string  `json:"discriminator"`
		Condition     string  `json:"condition"`
		Rules         []Rules `json:"rules"`
	} `json:"filter"`
	PagingInfo struct {
		StartPage int `json:"startPage"`
		PageSize  int `json:"pageSize"`
	} `json:"pagingInfo"`
	SortingInfo []struct {
		SortColumn    string `json:"sortColumn"`
		SortDirection string `json:"sortDirection"`
		AttributeType string `json:"attributeType"`
	} `json:"sortingInfo"`
}

type Rules struct {
	ID       string `json:"id"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
	Type     string `json:"type"`
}

type QueryResponse struct {
	Products    []Product     `json:"products"`
	SortingInfo []SortingInfo `json:"sortingInfo"`
	PagingInfo  PagingInfo    `json:"pagingInfo"`
	ItemCount   int           `json:"itemCount"`
}

type Product struct {
	Geom          string `json:"geom"`
	GeomExtent    string `json:"geomExtent"`
	GeomFootprint string `json:"geomFootprint"`
	Metadata      struct {
		Size                     string    `json:"size"`
		Format                   string    `json:"format"`
		Status                   string    `json:"status"`
		Filename                 string    `json:"filename"`
		Footprint                string    `json:"footprint"`
		SensorType               string    `json:"sensorType"`
		Timeliness               string    `json:"timeliness"`
		CycleNumber              string    `json:"cycleNumber"`
		EndPosition              string    `json:"endPosition"`
		OrbitNumber              string    `json:"orbitNumber"`
		ProductType              string    `json:"productType"`
		SliceNumber              string    `json:"sliceNumber"`
		GmlFootprint             string    `json:"gmlFootprint"`
		PlatformName             string    `json:"platformName"`
		ProductClass             string    `json:"productClass"`
		ProductLevel             string    `json:"productLevel"`
		BeginPosition            string    `json:"beginPosition"`
		InstrumentName           string    `json:"instrumentName"`
		OrbitDirection           string    `json:"orbitDirection"`
		ProcessingDate           time.Time `json:"processingDate"`
		AcquisitionType          string    `json:"acquisitionType"`
		LastOrbitNumber          string    `json:"lastOrbitNumber"`
		PhaseIdentifier          string    `json:"phaseIdentifier"`
		PlatformNssdcid          string    `json:"platformNssdcid"`
		ProcessingLevel          string    `json:"processingLevel"`
		SwathIdentifier          string    `json:"swathIdentifier"`
		PlatformShortName        string    `json:"platformShortName"`
		DataTakeIdentifier       string    `json:"dataTakeIdentifier"`
		ProductComposition       string    `json:"productComposition"`
		InstrumentShortName      string    `json:"instrumentShortName"`
		RelativeOrbitNumber      string    `json:"relativeOrbitNumber"`
		PolarisationChannels     string    `json:"polarisationChannels"`
		SensorOperationalMode    string    `json:"sensorOperationalMode"`
		LastRelativeOrbitNumber  string    `json:"lastRelativeOrbitNumber"`
		ProductClassDescription  string    `json:"productClassDescription"`
		PlatformSerialIdentifier string    `json:"platformSerialIdentifier"`
	} `json:"metadata"`
	Invalid            bool   `json:"invalid"`
	InvalidationDate   any    `json:"invalidationDate"`
	InvalidationReason any    `json:"invalidationReason"`
	ID                 string `json:"id"`
	Created            string `json:"created"`
	ProductTypeName    string `json:"productTypeName"`
	ProductTypeCode    string `json:"productTypeCode"`
	FolderPath         string `json:"folderPath"`
}

type SortingInfo struct {
	SortColumn    string `json:"sortColumn"`
	SortDirection string `json:"sortDirection"`
	AttributeType string `json:"attributeType"`
}

type PagingInfo struct {
	StartPage int `json:"startPage"`
	Offset    any `json:"offset"`
	PageSize  int `json:"pageSize"`
}
