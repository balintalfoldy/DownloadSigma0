package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

var apiKey string

func getApiKey() error {

	if apiKey == "" {
		apiKey = os.Getenv("MY_APP_API_KEY")
	}
	if apiKey == "" {
		return fmt.Errorf("API key not provided and not found in environment")
	}

	return nil
}

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query products.",
	Long:  "Query products in the FIR database based on filtering parameters.",
	Run: func(cmd *cobra.Command, args []string) {

		const url = "https://productapi.fir.gov.hu/api/v1/products/query"

		apiKey, err := cmd.Flags().GetString("api_key")
		if err != nil {
			fmt.Println("Invalid argument")
		}

		getApiKey()
		if err != nil {
			fmt.Println(err)
			return
		}

		productId, err := cmd.Flags().GetString("id")
		if err != nil {
			fmt.Println("Invalid argument")
			return
		}
		beginPosition, err := cmd.Flags().GetStringSlice("begin_position")
		if err != nil {
			fmt.Println("Invalid argument")
			return
		}
		relOrbit, err := cmd.Flags().GetString("relative_orbit_number")
		if err != nil {
			fmt.Println("Invalid argument")
			return
		}
		downloadFolder, err := cmd.Flags().GetString("download_folder")
		if err != nil {
			fmt.Println("Invalid argument")
			return
		}

		beginPositionFilterStart := Rules{
			ID:       "beginPosition",
			Operator: "greater_or_equal",
			Value:    beginPosition[0],
			Type:     "DateTime",
		}

		beginPositionFilterEnd := Rules{
			ID:       "beginPosition",
			Operator: "less_or_equal",
			Value:    beginPosition[1],
			Type:     "DateTime",
		}

		relativeOrbitNumFilter := Rules{
			ID:       "relativeOrbitNumber",
			Operator: "equal",
			Value:    relOrbit,
			Type:     "Integer",
		}

		queries := Queries{}
		// Read the JSON data from the file.
		queryBytes, err := ioutil.ReadFile("products_query.json")
		if err != nil {
			fmt.Println(err)
			return
		}

		err = json.Unmarshal(queryBytes, &queries)
		if err != nil {
			fmt.Println(err)
			return
		}

		var queryData QueryData

		for _, query := range queries.Queries {
			if query.ProductTypeCodes[0] == productId {
				queryData = query
				break
			}
		}

		if queryData.ProductTypeCodes[0] == "" {
			fmt.Println("Product not found")
			return
		}

		queryData.Filter.Rules = []Rules{beginPositionFilterStart, beginPositionFilterEnd, relativeOrbitNumFilter}

		queryJSON, err := json.Marshal(queryData)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(queryJSON))

		resp, err := PostRequest(url, apiKey, bytes.NewReader(queryJSON))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(resp)
		fmt.Println(downloadFolder)
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringP("id", "i", "", "Product ID")
	queryCmd.Flags().StringSliceP("begin_position", "b", nil, "Begin position")
	queryCmd.Flags().StringP("relative_orbit_number", "r", "", "Relative orbit number")
	queryCmd.Flags().StringP("download_folder", "d", "", "Download folder")
}
