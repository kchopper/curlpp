package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/alecthomas/chroma/quick"
	"github.com/fatih/color"
	httpclient "github.com/kchopper/curlpp/internal/client"
	"github.com/kchopper/curlpp/internal/config"
	"github.com/tidwall/pretty"
)

func main() {
	// Basic flags
	url := flag.String("url", "", "URL to request")
	method := flag.String("method", "GET", "HTTP method")
	prettyPrint := flag.Bool("pretty", true, "Pretty print response")
	parallel := flag.Int("parallel", 1, "Number of parallel requests")
	retries := flag.Int("retries", 3, "Number of retries for failed requests")
	profile := flag.String("profile", "", "Config profile to use")
	header := flag.String("H", "", "Headers in format 'key:value'")
	selector := flag.String("selector", "", "CSS selector to extract specific elements")

	flag.Parse()

	if *url == "" {
		fmt.Println("URL is required")
		flag.Usage()
		os.Exit(1)
	}

	// Create a basic config
	cfg := &config.Config{
		Current: *profile,
		Profiles: map[string]config.Profile{
			"default": {
				Auth: config.AuthConfig{
					Type:  "bearer",
					Token: os.Getenv("API_TOKEN"), // You can set token via environment variable
				},
			},
		},
	}

	// Create headers map from header flag
	headers := make(map[string]string)
	if *header != "" {
		// TODO: Parse header string into map
	}

	// Create client and make request
	client := httpclient.NewClient(cfg)
	resp, err := client.Do(&httpclient.Request{
		URL:      *url,
		Method:   *method,
		Pretty:   *prettyPrint,
		Parallel: *parallel,
		Retries:  *retries,
		Headers:  headers,
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Print response details
	c := color.New(color.FgGreen)
	c.Printf("Status Code: %d\n", resp.StatusCode)
	c.Printf("Time taken: %v\n", resp.Timing.TotalDuration)

	// Print headers
	color.Blue("\nHeaders:")
	for k, v := range resp.Headers {
		color.Cyan("%s: ", k)
		fmt.Printf("%s\n", v[0])
	}

	// Pretty print response
	fmt.Println("\nResponse Body:")
	if *prettyPrint {
		if isJSON(resp.Body) {
			colored := pretty.Color(pretty.Pretty(resp.Body), nil)
			fmt.Printf("%s\n", colored)
		} else if isHTML(resp.Headers) {
			if err := handleHTMLResponse(resp.Body, *prettyPrint, *selector); err != nil {
				color.Red("Error processing HTML: %v\n", err)
				color.White("%s\n", string(resp.Body))
			}
		} else {
			color.White("%s\n", string(resp.Body))
		}
	} else {
		color.White("%s\n", string(resp.Body))
	}
}

func isJSON(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}

func isHTML(headers http.Header) bool {
	contentType := headers.Get("Content-Type")
	return strings.Contains(contentType, "text/html") || strings.Contains(contentType, "application/html")
}

func handleHTMLResponse(body []byte, prettyPrint bool, selector string) error {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return err
	}

	// Extract title by default
	title := doc.Find("title").Text()
	if title != "" {
		color.Green("Title: %s\n", strings.TrimSpace(title))
	}

	// Extract meta description
	desc, exists := doc.Find("meta[name=description]").Attr("content")
	if exists {
		color.Green("Description: %s\n", desc)
	}

	// If selector is provided, extract matching elements
	if selector != "" {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			color.Yellow("Match %d: %s\n", i+1, strings.TrimSpace(s.Text()))
		})
		return nil
	}

	// Pretty print with syntax highlighting if requested
	if prettyPrint {
		var buf bytes.Buffer
		err := quick.Highlight(&buf, string(body), "html", "terminal", "monokai")
		if err != nil {
			return err
		}
		fmt.Println(buf.String())
	} else {
		fmt.Println(string(body))
	}

	return nil
}
