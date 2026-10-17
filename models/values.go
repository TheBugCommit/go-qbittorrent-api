package models

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// setStringParam sets key to *value, skipping nil pointers.
func setStringParam(values url.Values, key string, value *string) {
	if value != nil {
		values.Set(key, *value)
	}
}

// setBoolParam sets key to *value, skipping nil pointers.
func setBoolParam(values url.Values, key string, value *bool) {
	if value != nil {
		values.Set(key, strconv.FormatBool(*value))
	}
}

// setIntParam sets key to *value, skipping nil pointers.
func setIntParam(values url.Values, key string, value *int) {
	if value != nil {
		values.Set(key, strconv.Itoa(*value))
	}
}

// setInt64Param sets key to *value, skipping nil pointers.
func setInt64Param(values url.Values, key string, value *int64) {
	if value != nil {
		values.Set(key, strconv.FormatInt(*value, 10))
	}
}

// setFloat64Param sets key to *value, skipping nil pointers.
func setFloat64Param(values url.Values, key string, value *float64) {
	if value != nil {
		values.Set(key, strconv.FormatFloat(*value, 'f', -1, 64))
	}
}

// setListParam sets key to the elements of value joined by sep. An empty slice
// leaves key unset.
func setListParam(values url.Values, key string, value []string, sep string) {
	if len(value) > 0 {
		values.Set(key, strings.Join(value, sep))
	}
}

// stringParam returns a pointer to the value of key, or nil when key is absent.
func stringParam(values url.Values, key string) *string {
	if !values.Has(key) {
		return nil
	}

	value := values.Get(key)

	return &value
}

// boolParam parses the value of key, or returns nil when key is absent.
func boolParam(values url.Values, key string) (*bool, error) {
	if !values.Has(key) {
		return nil, nil
	}

	value, err := strconv.ParseBool(values.Get(key))
	if err != nil {
		return nil, fmt.Errorf("parse bool %q: %w", key, err)
	}

	return &value, nil
}

// intParam parses the value of key, or returns nil when key is absent.
func intParam(values url.Values, key string) (*int, error) {
	if !values.Has(key) {
		return nil, nil
	}

	value, err := strconv.Atoi(values.Get(key))
	if err != nil {
		return nil, fmt.Errorf("parse int %q: %w", key, err)
	}

	return &value, nil
}

// int64Param parses the value of key, or returns nil when key is absent.
func int64Param(values url.Values, key string) (*int64, error) {
	if !values.Has(key) {
		return nil, nil
	}

	value, err := strconv.ParseInt(values.Get(key), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse int64 %q: %w", key, err)
	}

	return &value, nil
}

// int64Param parses the value of key, or returns nil when key is absent.
func float64Param(values url.Values, key string) (*float64, error) {
	if !values.Has(key) {
		return nil, nil
	}

	value, err := strconv.ParseFloat(values.Get(key), 64)
	if err != nil {
		return nil, fmt.Errorf("parse float64 %q: %w", key, err)
	}

	return &value, nil
}

// listParam splits the value of key by sep, or returns nil when key is absent.
func listParam(values url.Values, key string, sep string) []string {
	if !values.Has(key) {
		return nil
	}

	return strings.Split(values.Get(key), sep)
}
