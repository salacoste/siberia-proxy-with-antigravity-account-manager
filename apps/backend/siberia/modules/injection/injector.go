package injection

import (
	"fmt"
	"time"
)

// Injector handles token injection into external app DBs
type Injector interface {
	Inject(dbPath string, accessToken, refreshToken string, expiry time.Time) error
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (i *Service) Inject(dbPath string, accessToken, refreshToken string, expiry time.Time) error {
	fmt.Printf("[Injector] Injecting into DB: %s\n", dbPath)
	fmt.Printf("           Access: %s...\n", accessToken[:5]) // Log partial for visual verify
	fmt.Printf("           Refresh: %s...\n", refreshToken[:5])

	// Simulate Protobuf decoding/encoding delay
	time.Sleep(1 * time.Second)

	fmt.Printf("[Injector] Injection Complete.\n")
	return nil
}
