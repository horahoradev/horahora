// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: party.proto

package protocol

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on NewPartyRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *NewPartyRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on NewPartyRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// NewPartyRequestMultiError, or nil if none found.
func (m *NewPartyRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *NewPartyRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserID

	if len(errors) > 0 {
		return NewPartyRequestMultiError(errors)
	}

	return nil
}

// NewPartyRequestMultiError is an error wrapping multiple validation errors
// returned by NewPartyRequest.ValidateAll() if the designated constraints
// aren't met.
type NewPartyRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m NewPartyRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m NewPartyRequestMultiError) AllErrors() []error { return m }

// NewPartyRequestValidationError is the validation error returned by
// NewPartyRequest.Validate if the designated constraints aren't met.
type NewPartyRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e NewPartyRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e NewPartyRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e NewPartyRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e NewPartyRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e NewPartyRequestValidationError) ErrorName() string { return "NewPartyRequestValidationError" }

// Error satisfies the builtin error interface
func (e NewPartyRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sNewPartyRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = NewPartyRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = NewPartyRequestValidationError{}

// Validate checks the field values on NewPartyResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *NewPartyResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on NewPartyResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// NewPartyResponseMultiError, or nil if none found.
func (m *NewPartyResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *NewPartyResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for PartyID

	if len(errors) > 0 {
		return NewPartyResponseMultiError(errors)
	}

	return nil
}

// NewPartyResponseMultiError is an error wrapping multiple validation errors
// returned by NewPartyResponse.ValidateAll() if the designated constraints
// aren't met.
type NewPartyResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m NewPartyResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m NewPartyResponseMultiError) AllErrors() []error { return m }

// NewPartyResponseValidationError is the validation error returned by
// NewPartyResponse.Validate if the designated constraints aren't met.
type NewPartyResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e NewPartyResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e NewPartyResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e NewPartyResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e NewPartyResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e NewPartyResponseValidationError) ErrorName() string { return "NewPartyResponseValidationError" }

// Error satisfies the builtin error interface
func (e NewPartyResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sNewPartyResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = NewPartyResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = NewPartyResponseValidationError{}

// Validate checks the field values on PartyRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PartyRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PartyRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PartyRequestMultiError, or
// nil if none found.
func (m *PartyRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *PartyRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserID

	// no validation rules for PartyID

	if len(errors) > 0 {
		return PartyRequestMultiError(errors)
	}

	return nil
}

// PartyRequestMultiError is an error wrapping multiple validation errors
// returned by PartyRequest.ValidateAll() if the designated constraints aren't met.
type PartyRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PartyRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PartyRequestMultiError) AllErrors() []error { return m }

// PartyRequestValidationError is the validation error returned by
// PartyRequest.Validate if the designated constraints aren't met.
type PartyRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PartyRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PartyRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PartyRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PartyRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PartyRequestValidationError) ErrorName() string { return "PartyRequestValidationError" }

// Error satisfies the builtin error interface
func (e PartyRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPartyRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PartyRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PartyRequestValidationError{}

// Validate checks the field values on Empty with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Empty) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Empty with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in EmptyMultiError, or nil if none found.
func (m *Empty) ValidateAll() error {
	return m.validate(true)
}

func (m *Empty) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return EmptyMultiError(errors)
	}

	return nil
}

// EmptyMultiError is an error wrapping multiple validation errors returned by
// Empty.ValidateAll() if the designated constraints aren't met.
type EmptyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m EmptyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m EmptyMultiError) AllErrors() []error { return m }

// EmptyValidationError is the validation error returned by Empty.Validate if
// the designated constraints aren't met.
type EmptyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EmptyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EmptyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EmptyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EmptyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EmptyValidationError) ErrorName() string { return "EmptyValidationError" }

// Error satisfies the builtin error interface
func (e EmptyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEmpty.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EmptyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EmptyValidationError{}

// Validate checks the field values on LeaderResponse with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LeaderResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LeaderResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in LeaderResponseMultiError,
// or nil if none found.
func (m *LeaderResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *LeaderResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for IsLeader

	if len(errors) > 0 {
		return LeaderResponseMultiError(errors)
	}

	return nil
}

// LeaderResponseMultiError is an error wrapping multiple validation errors
// returned by LeaderResponse.ValidateAll() if the designated constraints
// aren't met.
type LeaderResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LeaderResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LeaderResponseMultiError) AllErrors() []error { return m }

// LeaderResponseValidationError is the validation error returned by
// LeaderResponse.Validate if the designated constraints aren't met.
type LeaderResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LeaderResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LeaderResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LeaderResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LeaderResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LeaderResponseValidationError) ErrorName() string { return "LeaderResponseValidationError" }

// Error satisfies the builtin error interface
func (e LeaderResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLeaderResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LeaderResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LeaderResponseValidationError{}

// Validate checks the field values on PartyState with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PartyState) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PartyState with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PartyStateMultiError, or
// nil if none found.
func (m *PartyState) ValidateAll() error {
	return m.validate(true)
}

func (m *PartyState) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetUsers() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PartyStateValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PartyStateValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PartyStateValidationError{
					field:  fmt.Sprintf("Users[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	for idx, item := range m.GetVideos() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PartyStateValidationError{
						field:  fmt.Sprintf("Videos[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PartyStateValidationError{
						field:  fmt.Sprintf("Videos[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PartyStateValidationError{
					field:  fmt.Sprintf("Videos[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return PartyStateMultiError(errors)
	}

	return nil
}

// PartyStateMultiError is an error wrapping multiple validation errors
// returned by PartyState.ValidateAll() if the designated constraints aren't met.
type PartyStateMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PartyStateMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PartyStateMultiError) AllErrors() []error { return m }

// PartyStateValidationError is the validation error returned by
// PartyState.Validate if the designated constraints aren't met.
type PartyStateValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PartyStateValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PartyStateValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PartyStateValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PartyStateValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PartyStateValidationError) ErrorName() string { return "PartyStateValidationError" }

// Error satisfies the builtin error interface
func (e PartyStateValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPartyState.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PartyStateValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PartyStateValidationError{}

// Validate checks the field values on User with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *User) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on User with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in UserMultiError, or nil if none found.
func (m *User) ValidateAll() error {
	return m.validate(true)
}

func (m *User) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserID

	// no validation rules for IsLeader

	// no validation rules for Username

	if len(errors) > 0 {
		return UserMultiError(errors)
	}

	return nil
}

// UserMultiError is an error wrapping multiple validation errors returned by
// User.ValidateAll() if the designated constraints aren't met.
type UserMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserMultiError) AllErrors() []error { return m }

// UserValidationError is the validation error returned by User.Validate if the
// designated constraints aren't met.
type UserValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserValidationError) ErrorName() string { return "UserValidationError" }

// Error satisfies the builtin error interface
func (e UserValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUser.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserValidationError{}

// Validate checks the field values on Video with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Video) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Video with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in VideoMultiError, or nil if none found.
func (m *Video) ValidateAll() error {
	return m.validate(true)
}

func (m *Video) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Title

	// no validation rules for ID

	// no validation rules for Location

	if len(errors) > 0 {
		return VideoMultiError(errors)
	}

	return nil
}

// VideoMultiError is an error wrapping multiple validation errors returned by
// Video.ValidateAll() if the designated constraints aren't met.
type VideoMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m VideoMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m VideoMultiError) AllErrors() []error { return m }

// VideoValidationError is the validation error returned by Video.Validate if
// the designated constraints aren't met.
type VideoValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e VideoValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e VideoValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e VideoValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e VideoValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e VideoValidationError) ErrorName() string { return "VideoValidationError" }

// Error satisfies the builtin error interface
func (e VideoValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sVideo.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = VideoValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = VideoValidationError{}
