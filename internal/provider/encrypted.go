package provider

import (
	"encoding/json"
	"fmt"
)

type HashingType string

const (
	HashingSoloType   HashingType = "solo"
	HashingCustomType HashingType = "custom"
)

type EncryptedValue struct {
	HashingType HashingType `json:"hashing_type"`
	Encrypted   *string     `json:"encrypted,omitempty"`
	Plaintext   *string     `json:"plaintext,omitempty"`
}

type EncryptedValuesManager struct {
	preferServerState bool                       `json:"-"`
	EncryptedValues   map[string]*EncryptedValue `json:"values"`
}

type EncryptedHashingTypeMismatchError struct {
	key      string
	expected HashingType
	actual   HashingType
}

func (o EncryptedHashingTypeMismatchError) Error() string {
	return fmt.Sprintf("unexpected hashing type for key '%s': %s != %s", o.key, o.expected, o.actual)
}

func NewEncryptedValuesManager(payload []byte, preferServerState bool) (*EncryptedValuesManager, error) {
	manager := &EncryptedValuesManager{
		preferServerState: preferServerState,
	}

	if payload != nil {
		err := json.Unmarshal(payload, manager)
		if err != nil {
			return nil, err
		}
	} else {
		manager.EncryptedValues = make(map[string]*EncryptedValue)
	}

	return manager, nil
}

func (o EncryptedValuesManager) PreferServerState() bool {
	return o.preferServerState
}

func (o *EncryptedValuesManager) StorePlaintextValue(key string, hashing_type HashingType, value string) error {
	values, found := o.EncryptedValues[key]
	if !found {
		values = &EncryptedValue{
			Plaintext:   &value,
			HashingType: hashing_type,
		}
		o.EncryptedValues[key] = values
	} else if values.HashingType != hashing_type {
		return EncryptedHashingTypeMismatchError{
			key:      key,
			expected: hashing_type,
			actual:   values.HashingType,
		}
	} else {
		values.Plaintext = &value
	}

	return nil
}

func (o *EncryptedValuesManager) StoreEncryptedValue(key string, hashing_type HashingType, value string) error {
	values, found := o.EncryptedValues[key]
	if !found {
		values = &EncryptedValue{
			Encrypted:   &value,
			HashingType: hashing_type,
		}
		o.EncryptedValues[key] = values
	} else if values.HashingType != hashing_type {
		return EncryptedHashingTypeMismatchError{
			key:      key,
			expected: hashing_type,
			actual:   values.HashingType,
		}
	} else {
		values.Encrypted = &value
	}

	return nil
}

func (o EncryptedValuesManager) GetPlaintextValue(key string) (string, bool) {
	if values, found := o.EncryptedValues[key]; !found {
		return "", false
	} else if values.Plaintext == nil {
		return "", false
	} else {
		return *values.Plaintext, true
	}
}

func (o EncryptedValuesManager) GetEncryptedValue(key string) (string, bool) {
	if values, found := o.EncryptedValues[key]; !found {
		return "", false
	} else if values.Encrypted == nil {
		return "", false
	} else {
		return *values.Encrypted, true
	}
}
