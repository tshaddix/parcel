package decoding

type (
	RequestDecodeError struct {
		FromType string
		ToType   string
		Err      error
	}
)

func (e *RequestDecodeError) Error() string {
	return "Decode Error: Bad type conversion (" + e.FromType + " to " + e.ToType + ")"
}
