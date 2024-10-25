package errors

const prefix = "orders_"

const UnexpectedErrorMessage = "unexpected error occurred"

const (
	OrderGetInvalidParams     = prefix + "get_invalid_params"
	OrderGetUnAuthorized      = prefix + "get_unauthorized"
	OrderGetNotFound          = prefix + "get_not_found"
	OrderGetRateLimitExceeded = prefix + "get_rate_limit_exceeded"
	OrderGetServerError       = prefix + "get_server_error"

	OrderCreateInvalidInput      = prefix + "create_invalid_input"
	OrderCreateUnauthorized      = prefix + "create_unauthorized"
	OrderCreateServerError       = prefix + "create_server_error"
	OrderCreateRateLimitExceeded = prefix + "create_rate_limit_exceeded"

	OrderUpdateInvalidInput      = prefix + "update_invalid_input"
	OrderUpdateUnauthorized      = prefix + "update_unauthorized"
	OrderUpdateNotFound          = prefix + "update_not_found"
	OrderUpdateRateLimitExceeded = prefix + "update_rate_limit_exceeded"
	OrderUpdateServerError       = prefix + "update_server_error"

	OrderDeleteInvalidID         = prefix + "delete_invalid_order_id"
	OrderDeleteUnauthorized      = prefix + "delete_unauthorized"
	OrderDeleteNotFound          = prefix + "delete_not_found"
	OrderDeleteRateLimitExceeded = prefix + "delete_rate_limit_exceeded"
	OrderDeleteServerError       = prefix + "delete_server_error"
)
