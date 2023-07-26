package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/IceflowRE/go-multiprogressbar"
	"github.com/nathan-fiscaletti/consolesize-go"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

const baseURL = "https://productapi.fir.gov.hu/api/v1/products"
const concurrentDownloads = 4 // number of Wait groups for parallel processing

func getProductQuery(productId string, filterRules []Rules, startPage int) (*QueryData, error) {
	queries := Queries{}
	// Read the JSON data from the file.
	queryBytes, err := os.ReadFile("products_query.json")
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

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download products.",
	Long:  "Download products from the FIR database based on filtering parameters.",
	Run: func(cmd *cobra.Command, args []string) {

		apiKey := os.Getenv("FIR_PROD_API_KEY")
		if apiKey == "" {
			fmt.Println("API key not provided and not found in environment")
			return
		}

		client := NewClient(baseURL, apiKey, 20)

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
		if downloadFolder == "" {
			fmt.Println("Please provide a download folder")
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

		_, rows := consolesize.GetConsoleSize()
		fmt.Printf("\033[%d;0H", rows)
		fmt.Printf("A total of %d products will be downloaded\n", len(products))

		client = NewClient(baseURL, apiKey, 600)

		mpb := multiprogressbar.New()

		for _, p := range products {
			size, err := strconv.ParseInt(p.Metadata.Size, 10, 64)
			if err != nil {
				fmt.Println(err)
				size = 1
			}
			pBar := GetProgressbar(int(size))
			mpb.Add(pBar)
		}

		var wg sync.WaitGroup

		for i := 0; i < len(products); i += concurrentDownloads {
			end := i + concurrentDownloads
			if end > len(products) {
				end = len(products)
			}
			group := products[i:end]

			for i, p := range group {

				outPath := fmt.Sprintf("%s/%s.zip", downloadFolder, p.ID)
				url := fmt.Sprintf("%s/zip", p.ID)
				size, err := strconv.ParseInt(p.Metadata.Size, 10, 64)
				if err != nil {
					fmt.Println(err)
					size = 1
				}
				bar := mpb.Get(i)

				wg.Add(1)
				go func(o string, u string, s int64, b *progressbar.ProgressBar) {

					defer func() {
						wg.Done()
					}()

					if err := client.DownloadFile(o, u, s, b); err != nil {
						b.Exit()
						fmt.Printf("Error downloading %s: %s\n", u, err)
					} else {
						b.Finish()
					}
				}(outPath, url, size, bar)
			}

			wg.Wait()
			cols, rows := consolesize.GetConsoleSize()
			fmt.Printf("\033[%d;%dH", rows, cols)
		}

		fmt.Println("\nAll files downloaded")

	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringP("id", "i", "", "Product ID")
	downloadCmd.Flags().StringSliceP("begin_position", "b", nil, "Begin position")
	downloadCmd.Flags().StringP("relative_orbit_number", "r", "", "Relative orbit number")
	downloadCmd.Flags().StringP("download_folder", "d", "", "Download folder")
}
