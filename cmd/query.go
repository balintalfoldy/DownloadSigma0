package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

const baseURL = "https://productapi.fir.gov.hu/api/v1/products"

func getProductQuery(productId string, filterRules []Rules, startPage int) (*QueryData, error) {
	queries := Queries{}
	// Read the JSON data from the file.
	queryBytes, err := ioutil.ReadFile("products_query.json")
	if err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}

	err = json.Unmarshal(queryBytes, &queries)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling query json: %s", err)
	}

	for _, query := range queries.Queries {
		if query.ProductTypeCodes[0] == productId {
			query.Filter.Rules = filterRules
			query.PagingInfo.StartPage = startPage
			return &query, nil
		}
	}

	return nil, fmt.Errorf("product id %s not found", productId)
}

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query products.",
	Long:  "Query products in the FIR database based on filtering parameters.",
	Run: func(cmd *cobra.Command, args []string) {

		apiKey := os.Getenv("FIR_PROD_API_KEY")
		if apiKey == "" {
			fmt.Println("API key not provided and not found in environment")
			return
		}

		client := NewClient(baseURL, apiKey, 10)

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

		var products []Product
		for startPage := 0; startPage < 100; startPage++ {
			queryData, err := getProductQuery(productId, []Rules{beginPositionFilterStart, beginPositionFilterEnd, relativeOrbitNumFilter}, startPage)
			if err != nil {
				fmt.Println(err)
				return
			}

			queryJSON, err := json.Marshal(queryData)
			if err != nil {
				fmt.Println(err)
				return
			}

			res := QueryResponse{}
			if err := client.PostRequest("query", bytes.NewReader(queryJSON), &res); err != nil {
				fmt.Println(err)
				return
			}
			products = append(products, res.Products...)

			if res.ItemCount >= len(products) {
				break
			}
		}

		fmt.Println(products)

		if downloadFolder != "" {
			client := NewClient(baseURL, apiKey, 600)
			for _, product := range products {
				outPath := fmt.Sprintf("%s/%s.zip", downloadFolder, product.ID)
				url := fmt.Sprintf("%s/zip", product.ID)
				size, err := strconv.ParseInt(product.Metadata.Size, 10, 64)
				if err != nil {
					fmt.Println(err)
					size = 1
				}
				fmt.Printf("Downloading to %s\n", outPath)
				err = client.DownloadFile(outPath, url, size)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringP("id", "i", "", "Product ID")
	queryCmd.Flags().StringSliceP("begin_position", "b", nil, "Begin position")
	queryCmd.Flags().StringP("relative_orbit_number", "r", "", "Relative orbit number")
	queryCmd.Flags().StringP("download_folder", "d", "", "Download folder")
}
