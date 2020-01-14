package sfoxapi

func (api *SFOXAPI) reportError(source ErrorSourceKey, err error) {
	if err != nil && api.ErrorMonitor != nil {
		api.ErrorMonitor.RecordError(CreateOrderKey, err)
	}
}

type ResponseBodyError struct {
	Underlying   error
	ResponseBody string
}

func (e *ResponseBodyError) Error() string {
	return e.Underlying.Error() + " " + e.ResponseBody
}
