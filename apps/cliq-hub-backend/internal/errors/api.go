package errors

type APIError struct {
    Error   string                 `json:"error"`
    Message string                 `json:"message"`
    Meta    map[string]interface{} `json:"meta,omitempty"`
}

func New(code, msg string) APIError {
    return APIError{Error: code, Message: msg}
}

func (e APIError) WithMeta(key string, value interface{}) APIError {
    if e.Meta == nil {
        e.Meta = make(map[string]interface{})
    }
    e.Meta[key] = value
    return e
}