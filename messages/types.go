package messages

type (
	Handler interface {
		IsApplicable(message []interface{}) bool
		Process(message []interface{}) error
	}
)
