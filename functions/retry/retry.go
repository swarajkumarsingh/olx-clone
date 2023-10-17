package retry

import (
	"errors"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Stop indicates manually stopping of custom retry
type Stop struct {
	error
}

// NewStop function can be used to implement a custom stop message
func NewStop(message string) Stop {
	return Stop{errors.New(message)}
}

// CustomRetry retries a function f for attempts time with given duration
func CustomRetry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(Stop); ok {
			// Return the original error for later checking
			return s.error
		}
		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return CustomRetry(attempts, 2*sleep, f)
		}
		return err
	}

	return nil
}

// func try() error {
// 	req, err := http.NewRequest(
// 		"Get",
// 		"https://google.com",
// 		nil,
// 	)
// 	if err != nil {
// 		log.Println("unable to make request: %s", err)
// 	}
// 	// Execute the request
// 	return CustomRetry(3, time.Second, func() error {
// 		resp, err := http.DefaultClient.Do(req)
// 		if err != nil {
// 			// This error will result in a retry
// 			return err
// 		}
// 		defer resp.Body.Close()

// 		s := resp.StatusCode
// 		switch {
// 		case s >= 500:
// 			// Retry
// 			return fmt.Errorf("server error: %v", s)
// 		case s >= 400:
// 			// Don't retry, it was client's fault
// 			// return stop{fmt.Errorf("client error: %v", s)}
// 			return fmt.Errorf("client error: %v", s)
// 		default:
// 			// Happy
// 			return nil
// 		}
// 	})
// }
