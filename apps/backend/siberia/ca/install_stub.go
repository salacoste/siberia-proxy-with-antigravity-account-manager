//go:build !darwin

package ca

import "fmt"

func (s *Service) InstallCert() error {
	return fmt.Errorf("certificate installation not yet supported on this OS")
}

func (s *Service) CheckTrust() bool {
	return false
}
