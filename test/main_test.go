package test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	handler "github.com/MdSadiqMd/Sitemap-Builder/internal"
)

func TestGenerateSitemap(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		maxDepth    int
		expectError bool
	}{
		{
			name:        "valid URL with depth 1",
			url:         "https://example.com",
			maxDepth:    1,
			expectError: false,
		},
		{
			name:        "valid URL with depth 2",
			url:         "https://google.com",
			maxDepth:    2,
			expectError: false,
		},
		{
			name:        "zero depth",
			url:         "https://example.com",
			maxDepth:    0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := handler.GenerateSitemap(tt.url, tt.maxDepth, &buf)
			if (err != nil) != tt.expectError {
				t.Errorf("GenerateSitemap() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError {
				output := buf.String()
				if !strings.Contains(output, xml.Header) {
					t.Error("Output should contain XML header")
				}
				if !strings.Contains(output, handler.XMLNS_CONST) {
					t.Errorf("Output should contain handler.XMLNS_CONST: %s", handler.XMLNS_CONST)
				}

				var result handler.UrlSet
				xmlContent := strings.TrimPrefix(output, xml.Header)
				xmlContent = strings.TrimSpace(xmlContent)
				if err := xml.Unmarshal([]byte(xmlContent), &result); err != nil {
					t.Errorf("Generated XML is invalid: %v", err)
				}
				if result.Xmlns != handler.XMLNS_CONST {
					t.Errorf("Expected handler.XMLNS_CONST %s, got %s", handler.XMLNS_CONST, result.Xmlns)
				}
				t.Logf("Generated %d URLs for %s with depth %d", len(result.Urls), tt.url, tt.maxDepth)
			}
		})
	}
}

func TestXMLStructures(t *testing.T) {
	t.Run("loc struct XML marshaling", func(t *testing.T) {
		location := handler.Loc{Value: "https://example.com"}
		data, err := xml.Marshal(location)
		if err != nil {
			t.Fatalf("Failed to marshal loc: %v", err)
		}

		expected := "<loc><loc>https://example.com</loc></loc>"
		if string(data) != expected {
			t.Errorf("Expected %s, got %s", expected, string(data))
		}
	})

	t.Run("urlset struct XML marshaling", func(t *testing.T) {
		urls := handler.UrlSet{
			Xmlns: handler.XMLNS_CONST,
			Urls: []handler.Loc{
				{Value: "https://example.com"},
				{Value: "https://example.com/page1"},
			},
		}

		data, err := xml.MarshalIndent(urls, "", " ")
		if err != nil {
			t.Fatalf("Failed to marshal urlset: %v", err)
		}

		output := string(data)
		if !strings.Contains(output, `handler.XMLNS_CONST="`+handler.XMLNS_CONST+`"`) {
			t.Error("XML should contain correct handler.XMLNS_CONST attribute")
		}
		if !strings.Contains(output, "https://example.com") {
			t.Error("XML should contain the URLs")
		}

		var result handler.UrlSet
		if err := xml.Unmarshal(data, &result); err != nil {
			t.Errorf("Generated XML should be valid: %v", err)
		}
		if len(result.Urls) != 2 {
			t.Errorf("Expected 2 URLs, got %d", len(result.Urls))
		}
	})
}

func TestConstants(t *testing.T) {
	t.Run("handler.XMLNS_CONST constant", func(t *testing.T) {
		expected := "https://www.sitemaps.org/schemas/sitemap/0.9"
		if handler.XMLNS_CONST != expected {
			t.Errorf("Expected handler.XMLNS_CONST %s, got %s", expected, handler.XMLNS_CONST)
		}
	})
}

func TestMainFunction(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "default arguments",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "custom URL and depth",
			args:    []string{"-url", "https://example.com", "-depth", "1"},
			wantErr: false,
		},
		{
			name:    "zero depth",
			args:    []string{"-depth", "0"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("go", "build", "-o", "sitemap-test", "../cmd")
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to build binary: %v", err)
			}
			defer os.Remove("sitemap-test")

			cmd = exec.Command("./sitemap-test", tt.args...)
			output, err := cmd.Output()
			if (err != nil) != tt.wantErr {
				t.Errorf("main() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				outputStr := string(output)
				if !strings.Contains(outputStr, xml.Header) {
					t.Error("Output should contain XML header")
				}
				if !strings.Contains(outputStr, handler.XMLNS_CONST) {
					t.Error("Output should contain correct handler.XMLNS_CONST")
				}

				xmlContent := strings.TrimPrefix(outputStr, xml.Header)
				xmlContent = strings.TrimSpace(xmlContent)

				var result handler.UrlSet
				if err := xml.Unmarshal([]byte(xmlContent), &result); err != nil {
					t.Errorf("Output should be valid XML: %v", err)
				}
				t.Logf("Integration test generated %d URLs", len(result.Urls))
			}
		})
	}
}

func BenchmarkGenerateSitemap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		handler.GenerateSitemap("https://example.com", 1, &buf)
	}
}

func TestErrorHandling(t *testing.T) {
	t.Run("invalid writer", func(t *testing.T) {
		failingWriter := &failingWriter{}
		err := handler.GenerateSitemap("https://example.com", 1, failingWriter)
		if err == nil {
			t.Error("Should return error when writer fails")
		}
	})
}

type failingWriter struct{}

func (fw *failingWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("write failed")
}
