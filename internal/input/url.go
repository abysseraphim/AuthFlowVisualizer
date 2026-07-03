package input

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// chromePaths are candidate locations for a headless Chrome/Chromium binary.
var chromePaths = []string{
	"/opt/google/chrome/chrome",
	"/usr/bin/chromium-browser",
	"/usr/bin/chromium",
	"/usr/bin/google-chrome",
	"/snap/bin/chromium",
}

// findChrome returns the first available Chrome binary, or "" if none found.
func findChrome() string {
	// prefer PATH
	if p, err := exec.LookPath("google-chrome"); err == nil {
		return p
	}
	if p, err := exec.LookPath("chromium-browser"); err == nil {
		return p
	}
	if p, err := exec.LookPath("chromium"); err == nil {
		return p
	}
	for _, p := range chromePaths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// fetchRendered fetches the fully-rendered DOM of a URL using headless Chrome.
func fetchRendered(targetURL string) ([]byte, error) {
	chromeBin := findChrome()
	if chromeBin == "" {
		fmt.Println("[!]Headless browser not found, falling back to static fetch")
		return staticFetch(targetURL)
	}

	fmt.Println("[*]Headless browser found:", chromeBin)

	tmpHTML, err := os.CreateTemp("", "afv-dom-*.html")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpHTML.Name())
	tmpHTML.Close()

	cmd := exec.Command(
		chromeBin,
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"--dump-dom",
		"--timeout=15000",
		targetURL,
	)

	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf

	done := make(chan error, 1)
	if err := cmd.Start(); err != nil {
		fmt.Println("[!]Failed to launch headless browser, falling back to static fetch:", err)
		return staticFetch(targetURL)
	}
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			fmt.Println("[!]Headless browser exited with error, falling back to static fetch:", err)
			return staticFetch(targetURL)
		}
	case <-time.After(20 * time.Second):
		cmd.Process.Kill()
		fmt.Println("[!]Headless browser timed out, falling back to static fetch")
		return staticFetch(targetURL)
	}

	body := out.Bytes()
	if len(body) == 0 {
		fmt.Println("[!]Headless browser returned empty body, falling back to static fetch")
		return staticFetch(targetURL)
	}
	return body, nil
}

func staticFetch(targetURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0 Safari/537.36")
	req.Header.Set("Referer", targetURL)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New("[e]status not okay")
	}
	return io.ReadAll(resp.Body)
}

func URLHandler(path string) ([]string, error) {
	TmpFiles := []string{}
	SRCs := []string{}

	if !strings.Contains(path, "http://") && !strings.Contains(path, "https://") {
		path = "https://" + path
	}

	body, err := fetchRendered(path)
	if err != nil {
		return TmpFiles, err
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return TmpFiles, err
	}

	var getScripts func(*html.Node)
	getScripts = func(node *html.Node) {
		if node == nil {
			return
		}

		if node.Type == html.ElementNode && node.Data == "script" {
			for _, attr := range node.Attr {
				if strings.ToLower(attr.Key) == "src" {
					SRCs = append(SRCs, attr.Val)
					break
				}
			}

			if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
				tmp, err := CreateTemp(string(node.FirstChild.Data))
				if err != nil {
					fmt.Println("[e]Failed to write Contents of Inline Script to TempFile")
				}
				TmpFiles = append(TmpFiles, tmp)
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			getScripts(child)
		}
	}

	getScripts(doc)

	baseURL, err := url.Parse(path)
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
		fullURL := baseURL.ResolveReference(scriptURL)
		resolvedSRCs = append(resolvedSRCs, fullURL.String())
	}

	fmt.Println("[*]Resolved URLs:")
	for _, u := range resolvedSRCs {
		fmt.Println("   ", u)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	for _, rawURL := range resolvedSRCs {
		req, err := http.NewRequest("GET", rawURL, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0 Safari/537.36")
		req.Header.Set("Referer", path)

		response, err := client.Do(req)
		if err != nil {
			fmt.Println("[!]URL:", rawURL, "didn't resolve")
			continue
		}
		ScriptText, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			fmt.Println("[!]URL:", rawURL, "didn't resolve completely")
			continue
		}
		tmp, err := CreateTemp(string(ScriptText))
		if err != nil {
			fmt.Println("[e]Failed to write contents of", rawURL, "to TempFile")
			continue
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
