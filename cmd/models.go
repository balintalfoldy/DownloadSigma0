package cmd

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
