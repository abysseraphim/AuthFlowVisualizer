package input

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func URLHandler(path string) ([]string, error) {
	TmpFiles := []string{}
	SRCs := []string{}
	if !strings.Contains(path, "http://") && !strings.Contains(path, "https://") {
		path = "https://" + path
	}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return TmpFiles, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0 Safari/537.36")
	req.Header.Set("Referer", path)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := (*client).Do(req)
	if err != nil {
		return TmpFiles, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return TmpFiles, errors.New("[e]status not okay")
	}

	body, err := io.ReadAll(resp.Body) // body is a slice of bytes.
	if err != nil {
		return TmpFiles, err
	}

	doc, err := html.Parse(bytes.NewReader(body)) // html.Parse needs an io reader...
	if err != nil {
		return TmpFiles, err
	}

	var getScripts func(*html.Node)

	getScripts = func(node *html.Node) {

		if node == nil {
			return
		}

		if node.Type == html.ElementNode && node.Data == "script" {
			// fmt.Println("Script Found")
			for _, attr := range node.Attr {
				if strings.ToLower(attr.Key) == "src" {
					SRCs = append(SRCs, attr.Val)
					// fmt.Println(attr.Val)
					break
				}
			}

			if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
				// fmt.Println(node.FirstChild.Data)
				tmp, err := CreateTemp(string(node.FirstChild.Data))
				if err != nil {
					fmt.Println("[e]Failed to write Contents of Inline Script to TempFile")
				}
				TmpFiles = append(TmpFiles, tmp)
			}
		}

		// for initialization; condition; update
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			getScripts(child)
		}
		// this form of for loop is equivalent to:
		/*
			child := node.FirstChild

			for child != nil {
				// this part makes the function recursive:
				getScript(child)

				// there is an html tag and children are head, body,.. we are switching between these children.
				child = child.NextSibling
			}
		*/
	}

	getScripts(doc)

	baseURL, err := url.Parse(path) // create a url object
	if err != nil {
		return TmpFiles, err
	}

	resolvedSRCs := []string{}

	for _, src := range SRCs {

		scriptURL, err := url.Parse(src)
		if err != nil {
			continue
		}

		if scriptURL.Scheme != "" &&
			scriptURL.Scheme != "http" &&
			scriptURL.Scheme != "https" {
			continue
		}

		// there are three types of inputs in SRCs: http://..../file.js, /abs/file.js, rel/file.js
		// ResolveReference knows to start from base url if /abs, append if rel
		fullURL := baseURL.ResolveReference(scriptURL)

		resolvedSRCs = append(resolvedSRCs, fullURL.String())
	}

	fmt.Println("[*]Resolved URLs:")
	for _, u := range resolvedSRCs {
		fmt.Println(u)
	}

	// write to tempfile
	for _, url := range resolvedSRCs {

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0 Safari/537.36")
		req.Header.Set("Referer", path)

		response, err := (*client).Do(req)
		if err != nil {
			fmt.Println("URL:", url, "Didn't Resolve.")
			continue
		}

		ScriptText, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			fmt.Println("URL:", url, "Didn't Resolve Completely.")
			continue
		}
		tmp, err := CreateTemp(string(ScriptText))
		if err != nil {
			fmt.Println("[e]Failed to write contents of", url, "to TempFile")
		}
		TmpFiles = append(TmpFiles, tmp)

	}

	return TmpFiles, nil
}

func CreateTemp(Content string) (string, error) {
	file, err := os.CreateTemp("/tmp", "afv-*.js")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(Content)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}
