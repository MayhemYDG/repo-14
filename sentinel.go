package cerrors

// Sentinel allows for the creation of constant error types.
// See https://dave.cheney.net/2016/04/07/constant-errors#constant-errors.
type Sentinel string

func (e Sentinel) Error() string { return string(e) }
