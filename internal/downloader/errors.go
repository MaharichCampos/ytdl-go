package downloader

import "errors"

// Error categories are used to map failures to stable exit codes and messages.
type ErrorCategory string

const (
	CategoryUnknown     ErrorCategory = "unknown"
	CategoryInvalidURL  ErrorCategory = "invalid_url"
	CategoryUnsupported ErrorCategory = "unsupported"
	CategoryRestricted  ErrorCategory = "restricted"
	CategoryNetwork     ErrorCategory = "network"
	CategoryFilesystem  ErrorCategory = "filesystem"
)

// CategorizedError wraps an error with a semantic category.
type CategorizedError struct {
	Category ErrorCategory
	Err      error
}

func (e CategorizedError) Error() string {
	return e.Err.Error()
}

func (e CategorizedError) Unwrap() error {
	return e.Err
}

// wrapCategory ensures an error carries the desired category.
func wrapCategory(category ErrorCategory, err error) error {
	if err == nil {
		return nil
	}
	var ce CategorizedError
	if errors.As(err, &ce) {
		return err
	}
	return CategorizedError{Category: category, Err: err}
}

// errorCategory extracts the category from a wrapped error.
func errorCategory(err error) ErrorCategory {
	var ce CategorizedError
	if errors.As(err, &ce) {
		return ce.Category
	}
	return CategoryUnknown
}

// CategoryOf exposes the detected error category for callers outside this package.
func CategoryOf(err error) ErrorCategory {
	return errorCategory(err)
}

// ExitCode maps a categorized error to a stable non-zero exit code.
func ExitCode(err error) int {
	switch errorCategory(err) {
	case CategoryInvalidURL:
		return 2
	case CategoryUnsupported:
		return 3
	case CategoryRestricted:
		return 4
	case CategoryNetwork:
		return 5
	case CategoryFilesystem:
		return 6
	case CategoryUnknown:
		if err != nil {
			return 1
		}
		return 0
	default:
		return 1
	}
}
