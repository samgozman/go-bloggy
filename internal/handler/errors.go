package handler

const (
	errRequestBodyBinding = "ERR_REQUEST_BODY_BINDING"
	errBodyValidation     = "ERR_BODY_VALIDATION"
	errForbidden          = "ERR_FORBIDDEN"
	errExchangeCode       = "ERR_EXCHANGE_CODE"
	errGetUserInfo        = "ERR_GET_USER_INFO"
	errCreateToken        = "ERR_CREATE_TOKEN" //nolint:gosec
	errCreateUser         = "ERR_CREATE_USER"
	errUnauthorized       = "ERR_UNAUTHORIZED"
	errGetUser            = "ERR_GET_USER"
	errDuplicatePost      = "ERR_DUPLICATE_POST"
	errCreatePost         = "ERR_CREATE_POST"
	errValidationFailed   = "ERR_VALIDATION_FAILED"
	errPostNotFound       = "ERR_POST_NOT_FOUND"
	errParamValidation    = "ERR_PARAM_VALIDATION"
	errGetPosts           = "ERR_GET_POSTS"
	errGetPostsCount      = "ERR_GET_POSTS_COUNT"
	errUpdatePost         = "ERR_UPDATE_POST"
	errCreateSubscription = "ERR_CREATE_SUBSCRIPTION"
	errGetSubscription    = "ERR_GET_SUBSCRIPTION"
	errDeleteSubscription = "ERR_DELETE_SUBSCRIPTION"
	errValidationEmail    = "ERR_VALIDATION_EMAIL"
)
