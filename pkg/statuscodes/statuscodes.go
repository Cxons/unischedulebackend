package statuscodes



type Status struct {
	Code int
	Message string
}



var (
    Continue           = Status{Code: 100, Message: "Continue"}
    SwitchingProtocols = Status{Code: 101, Message: "Switching Protocols"}
    Processing         = Status{Code: 102, Message: "Processing"}

    OK                   = Status{Code: 200, Message: "OK"}
    Created              = Status{Code: 201, Message: "Created"}
    Accepted             = Status{Code: 202, Message: "Accepted"}
    NonAuthoritativeInfo = Status{Code: 203, Message: "Non-Authoritative Information"}
    NoContent            = Status{Code: 204, Message: "No Content"}
    ResetContent         = Status{Code: 205, Message: "Reset Content"}
    PartialContent       = Status{Code: 206, Message: "Partial Content"}

    MultipleChoices   = Status{Code: 300, Message: "Multiple Choices"}
    MovedPermanently  = Status{Code: 301, Message: "Moved Permanently"}
    Found             = Status{Code: 302, Message: "Found"}
    SeeOther          = Status{Code: 303, Message: "See Other"}
    NotModified       = Status{Code: 304, Message: "Not Modified"}
    TemporaryRedirect = Status{Code: 307, Message: "Temporary Redirect"}
    PermanentRedirect = Status{Code: 308, Message: "Permanent Redirect"}

    BadRequest            = Status{Code: 400, Message: "Bad Request"}
    Unauthorized          = Status{Code: 401, Message: "Unauthorized"}
    PaymentRequired       = Status{Code: 402, Message: "Payment Required"}
    Forbidden             = Status{Code: 403, Message: "Forbidden"}
    NotFound              = Status{Code: 404, Message: "Not Found"}
    MethodNotAllowed      = Status{Code: 405, Message: "Method Not Allowed"}
    NotAcceptable         = Status{Code: 406, Message: "Not Acceptable"}
    ProxyAuthRequired     = Status{Code: 407, Message: "Proxy Authentication Required"}
    RequestTimeout        = Status{Code: 408, Message: "Request Timeout"}
    Conflict              = Status{Code: 409, Message: "Conflict"}
    Gone                  = Status{Code: 410, Message: "Gone"}
    LengthRequired        = Status{Code: 411, Message: "Length Required"}
    PreconditionFailed    = Status{Code: 412, Message: "Precondition Failed"}
    PayloadTooLarge       = Status{Code: 413, Message: "Payload Too Large"}
    URITooLong            = Status{Code: 414, Message: "URI Too Long"}
    UnsupportedMediaType  = Status{Code: 415, Message: "Unsupported Media Type"}
    RangeNotSatisfiable   = Status{Code: 416, Message: "Range Not Satisfiable"}
    ExpectationFailed     = Status{Code: 417, Message: "Expectation Failed"}
    ImATeapot             = Status{Code: 418, Message: "I'm a teapot"}
    MisdirectedRequest    = Status{Code: 421, Message: "Misdirected Request"}
    UnprocessableEntity   = Status{Code: 422, Message: "Unprocessable Entity"}
    Locked                = Status{Code: 423, Message: "Locked"}
    FailedDependency      = Status{Code: 424, Message: "Failed Dependency"}
    TooManyRequests       = Status{Code: 429, Message: "Too Many Requests"}

    InternalServerError   = Status{Code: 500, Message: "Internal Server Error"}
    NotImplemented        = Status{Code: 501, Message: "Not Implemented"}
    BadGateway            = Status{Code: 502, Message: "Bad Gateway"}
    ServiceUnavailable    = Status{Code: 503, Message: "Service Unavailable"}
    GatewayTimeout        = Status{Code: 504, Message: "Gateway Timeout"}
    HTTPVersionNotSupported = Status{Code: 505, Message: "HTTP Version Not Supported"}
)

var statusMap = map[string]Status{
    "Continue":                     Continue,
    "Switching Protocols":          SwitchingProtocols,
    "Processing":                   Processing,
    "OK":                           OK,
    "Created":                      Created,
    "Accepted":                     Accepted,
    "Non-Authoritative Information": NonAuthoritativeInfo,
    "No Content":                   NoContent,
    "Reset Content":                ResetContent,
    "Partial Content":              PartialContent,
    "Multiple Choices":             MultipleChoices,
    "Moved Permanently":            MovedPermanently,
    "Found":                        Found,
    "See Other":                    SeeOther,
    "Not Modified":                 NotModified,
    "Temporary Redirect":           TemporaryRedirect,
    "Permanent Redirect":           PermanentRedirect,
    "Bad Request":                  BadRequest,
    "Unauthorized":                 Unauthorized,
    "Payment Required":             PaymentRequired,
    "Forbidden":                    Forbidden,
    "Not Found":                    NotFound,
    "Method Not Allowed":           MethodNotAllowed,
    "Not Acceptable":               NotAcceptable,
    "Proxy Authentication Required": ProxyAuthRequired,
    "Request Timeout":              RequestTimeout,
    "Conflict":                     Conflict,
    "Gone":                         Gone,
    "Length Required":              LengthRequired,
    "Precondition Failed":          PreconditionFailed,
    "Payload Too Large":            PayloadTooLarge,
    "URI Too Long":                 URITooLong,
    "Unsupported Media Type":       UnsupportedMediaType,
    "Range Not Satisfiable":        RangeNotSatisfiable,
    "Expectation Failed":           ExpectationFailed,
    "I'm a teapot":                 ImATeapot,
    "Misdirected Request":          MisdirectedRequest,
    "Unprocessable Entity":         UnprocessableEntity,
    "Locked":                       Locked,
    "Failed Dependency":            FailedDependency,
    "Too Many Requests":            TooManyRequests,
    "Internal Server Error":        InternalServerError,
    "Not Implemented":              NotImplemented,
    "Bad Gateway":                  BadGateway,
    "Service Unavailable":          ServiceUnavailable,
    "Gateway Timeout":              GatewayTimeout,
    "HTTP Version Not Supported":   HTTPVersionNotSupported,
}

func RetrieveCodeFromStatusMessage(message string) int {
    status, ok := statusMap[message]
    if !ok {
        return 0 // or return an error
    }
    return status.Code
}