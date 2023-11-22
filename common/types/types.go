package types

// Zeroer, as implemented by time.Time, will be used by the truth template
// funcs in Hugo (if, with, not, and or).
type Zeroer interface {
  IsZero() bool
}
