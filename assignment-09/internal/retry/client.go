package retry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PaymentClient struct {
	Client *http.Client
	MaxRetries int
	BaseDelay time.Duration
	MaxDelay time.Duration
}

func (p *PaymentClient) ExecutePayment(ctx context.Context, req *http.Request) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt < p.MaxRetries; attempt++ {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		request := req.Clone(ctx)
		resp, err := p.Client.Do(request)
		if err == nil && resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			defer resp.Body.Close()
			body, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				return nil, readErr
			}
			fmt.Printf("Attempt %d: Success!\n", attempt+1)
			return body, nil
		}
		if resp != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
		lastErr = err
		if !IsRetryable(resp, err) {
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
		if attempt == p.MaxRetries-1 {
			break
		}
		delay := CalculateBackoff(attempt, p.BaseDelay, p.MaxDelay)
		if err != nil {
			fmt.Printf("Attempt %d failed: %v, waiting %v...\n", attempt+1, err, delay)
		} else {
			fmt.Printf("Attempt %d failed: status %d, waiting %v...\n", attempt+1, resp.StatusCode, delay)
		}
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("payment failed after %d attempts", p.MaxRetries)
}
