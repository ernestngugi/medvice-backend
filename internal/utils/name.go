package utils

import "fmt"

func ValidateSingleName(name string) error {

	err := emptyName(name)
	if err != nil {
		return err
	}

	err = validateNameLength(name)
	if err != nil {
		return err
	}

	return nil
}

func emptyName(name string) error {

	if name == "" {
		return fmt.Errorf("name caanot be empty")
	}

	return nil
}

func validateNameLength(name string) error {

	if len(name) > 50 {
		return fmt.Errorf("name length exceed 50 characters")
	}
	return nil
}
