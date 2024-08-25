package structmerge

import (
	"reflect"
	"strings"
	"time"
)

var (
	ErrInvalidDestination = newMergeError("destination must be a pointer to a struct")
	ErrInvalidSource      = newMergeError("source must be a struct")
	ErrTypeMismatch       = newMergeError("source and destination types do not match")
)

type MergeError struct {
	message string
}

func (e *MergeError) Error() string {
	return e.message
}

func newMergeError(message string) *MergeError {
	return &MergeError{message: message}
}

// Merger interface allows custom structs that don't export fields to be merged.
// The pointer receiver on type should implement Merge.
type Merger interface {

	// Marge pointer receiver with the source value.
	Merge(src reflect.Value) error
}

// MergeOption defined the behavior for merging fields.
type MergeOption int

const (
	// IncludeAll includes all fields in the merge.
	IncludeAll MergeOption = iota

	// ExcludeEmpty excludes empty fields from source.
	ExcludeEmpty

	// OverwriteEmpty overwrites empty fields in destination
	OverwriteEmpty
)

// Config holds configuration for the merge operation.
type Config struct {
	Option  MergeOption
	Include []string // Fields to include in the destination
	Exclude []string // Fields to exclude from destination struct
}

// Merge combines two structs of the same type based on the provided configuration
// The default configuration is to include all fields.
func Merge(dst, src interface{}, cfg ...Config) error {
	defaultConfig := Config{Option: IncludeAll}
	if len(cfg) > 0 {
		defaultConfig = cfg[0]
	}
	return mergeValues(reflect.ValueOf(dst), reflect.ValueOf(src), defaultConfig, "")
}

func mergeValues(dst, src reflect.Value, cfg Config, prefix string) error {
	if dst.Kind() != reflect.Ptr || dst.Elem().Kind() != reflect.Struct {
		return ErrInvalidDestination
	}

	if dst.IsNil() {
		dst.Set(reflect.New(dst.Type().Elem()))
	}
	dst = dst.Elem()

	if src.Kind() != reflect.Struct {
		return ErrInvalidSource
	}

	if dst.Type() != src.Type() {
		return ErrTypeMismatch
	}

	// Check if it's time.Time and copy it directly
	if dst.CanInterface() {
		if _, ok := dst.Interface().(time.Time); ok {
			dst.Set(src)
			return nil
		}
	}

	// Check if a struct implements the Merger interface
	ifacetype := reflect.TypeOf((*Merger)(nil)).Elem()
	if dst.CanAddr() && dst.Addr().Type().Implements(ifacetype) {
		merger := dst.Addr().Interface().(Merger)
		return merger.Merge(src)
	}

	includeMap := make(map[string]bool)
	for _, f := range cfg.Include {
		includeMap[f] = true
	}

	excludeMap := make(map[string]bool)
	for _, f := range cfg.Exclude {
		excludeMap[f] = true
	}

	for i := 0; i < dst.NumField(); i++ {
		field := dst.Type().Field(i)
		fieldName := field.Name
		fullFieldName := prefix + fieldName

		// Check if field should be included or excluded
		if len(cfg.Include) > 0 && !shouldInclude(fullFieldName, includeMap) {
			continue // Skip if not included
		}

		if excludeMap[fullFieldName] {
			continue // Skip if excluded
		}

		dstField := dst.Field(i)
		srcField := src.Field(i)

		// Check if a specific field implements merger
		if dstField.CanAddr() && dstField.Addr().Type().Implements(ifacetype) {
			merger := dstField.Addr().Interface().(Merger)
			merger.Merge(srcField)
			continue
		}

		// Only set if the field is settable
		if !dstField.CanSet() {
			continue
		}

		// Handle nested struct merging
		if dstField.Kind() == reflect.Struct {
			// Recursively merge nested structs
			err := mergeValues(dstField.Addr(), srcField, cfg, fullFieldName+".")
			if err != nil {
				return err
			}
		} else {
			shouldSet := true
			switch cfg.Option {
			case ExcludeEmpty:
				shouldSet = !isZero(srcField)
			case OverwriteEmpty:
				shouldSet = isZero(dstField)
			}

			if shouldSet || cfg.Option == IncludeAll {
				dstField.Set(srcField)
			}
		}
	}

	return nil
}

func shouldInclude(fullFieldName string, includeMap map[string]bool) bool {
	// Check if the exact full field name is in the include map
	if includeMap[fullFieldName] {
		return true
	}

	if !strings.Contains(fullFieldName, ".") {
		for key := range includeMap {
			if strings.HasPrefix(key, fullFieldName) {
				return true
			}
		}
	}

	return false
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return false
}
