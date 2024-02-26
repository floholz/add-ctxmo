package src

import (
	"errors"
	"os"
)

func NoEmptyStringValidator(s string) error {
	if s == "" {
		return errors.New("no empty strings allowed")
	}
	return nil
}

func ValidPathValidator(s string) error {
	_, err := os.Stat(s)
	if err == nil {
		return nil
	}
	return err
}
