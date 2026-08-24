package icx

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestICXClientLocalFallback(t *testing.T) {
	client := NewClient(Config{
		SpaceID:       "test_space",
		LocalFallback: true,
		Timeout:       200 * time.Millisecond,
	})

	ctx := context.Background()

	ingestResp, err := client.IngestTextWithContext(ctx, IngestTextRequest{
		Text:     "The SEC EDGAR system requires Form 10-K and 10-Q filings for public companies.",
		Filename: "sec_guide.md",
		SpaceID:  "test_space",
		Family:   "skill.definition",
	})
	if err != nil {
		t.Fatalf("unexpected error during ingest: %v", err)
	}

	if !ingestResp.Success || !ingestResp.LocalFallback {
		t.Errorf("expected successful local fallback ingest")
	}
	if ingestResp.Status != statusLocalFallback {
		t.Errorf("expected status %q, got %q", statusLocalFallback, ingestResp.Status)
	}
	if ingestResp.ContentHash == "" {
		t.Errorf("expected non-empty ContentHash")
	}

	batchResp, err := client.IngestBatchWithContext(ctx, IngestBatchRequest{
		Items: []IngestTextRequest{
			{Text: "Item one data.", Filename: "item1.md"},
			{Text: "Item two data.", Filename: "item2.md"},
		},
		SpaceID: "test_space",
	})
	if err != nil {
		t.Fatalf("unexpected error during batch ingest: %v", err)
	}
	if !batchResp.Success || batchResp.TotalItems != 2 {
		t.Errorf("expected 2 items batch ingested, got %d", batchResp.TotalItems)
	}
	if !batchResp.LocalFallback {
		t.Errorf("expected batch ingest to be labeled local fallback")
	}

	recallResp, err := client.RecallWithContext(ctx, RecallRequest{
		Query:   "10-K filings SEC",
		SpaceID: "test_space",
		TopK:    5,
	})
	if err != nil {
		t.Fatalf("unexpected error during recall: %v", err)
	}

	if !recallResp.Success || !recallResp.LocalFallback {
		t.Errorf("expected local fallback recall")
	}
	if len(recallResp.Matches) == 0 {
		t.Errorf("expected at least 1 matching register in local fallback recall")
	}
	if recallResp.Grounding != 0 {
		t.Errorf("local fallback must not invent a grounding score, got %v", recallResp.Grounding)
	}
}

func TestICXClientFailClosedWhenKeyed(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "http://127.0.0.1:59999/unreachable",
		APIKey:  "clabs_test",
		SpaceID: "test_space",
		Timeout: 200 * time.Millisecond,
	})

	_, err := client.IngestTextWithContext(context.Background(), IngestTextRequest{
		Text:     "should not be stored as a fake lattice write",
		Filename: "nope.md",
		SpaceID:  "test_space",
	})
	if err == nil {
		t.Fatal("expected ingest error when API key is set and host is unreachable")
	}
	if !strings.Contains(err.Error(), "icx ingest") {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = client.RecallWithContext(context.Background(), RecallRequest{
		Query:   "anything",
		SpaceID: "test_space",
	})
	if err == nil {
		t.Fatal("expected recall error when API key is set and host is unreachable")
	}
}

func TestICXClientHeaders(t *testing.T) {
	client := NewClient(Config{
		BaseURL:      "https://icx.api.caleralabs.com/v1",
		APIKey:       "clabs_key_123",
		SpaceID:      "space_alpha",
		BYOKProvider: "gemini",
		BYOKKey:      "AIzaSy_key_456",
		BYOKModel:    "gemini-3.5-flash-lite",
	})

	headers := client.buildHeaders("space_custom")
	if headers.Get("Authorization") != "Bearer clabs_key_123" {
		t.Errorf("expected Authorization header with Bearer token")
	}
	if headers.Get("X-Space-ID") != "space_custom" {
		t.Errorf("expected custom space ID in header")
	}
	if headers.Get("X-LLM-Provider") != "gemini" {
		t.Errorf("expected X-LLM-Provider header")
	}
	if headers.Get("X-LLM-API-Key") != "" {
		t.Errorf("OSS client must not forward BYOK keys, got %q", headers.Get("X-LLM-API-Key"))
	}
}
