package contracts

type RouteMethod string

const (
	RouteMethodGet    RouteMethod = "GET"
	RouteMethodPost   RouteMethod = "POST"
	RouteMethodPut    RouteMethod = "PUT"
	RouteMethodPatch  RouteMethod = "PATCH"
	RouteMethodDelete RouteMethod = "DELETE"
)

type RouteContract struct {
	Method RouteMethod
	Path   string
}
