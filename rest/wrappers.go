
package rest

// great, these are needed to do the structural typing for the tests...
type HttpReq interface {
	FormValue(key string) string
}
type HttpHeader interface {
	Get(key string) string
}
