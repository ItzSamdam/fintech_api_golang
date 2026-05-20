package response

type ErrorResponse struct {
    Error   string            `json:"error"`
    Message string            `json:"message"`
    Code    string            `json:"code,omitempty"`
    Details map[string]string `json:"details,omitempty"`
}

type ValidationErrorResponse struct {
    Error   string                `json:"error"`
    Message string                `json:"message"`
    Errors  map[string]string     `json:"errors"`
}

type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Total      int64       `json:"total"`
    Page       int         `json:"page"`
    Limit      int         `json:"limit"`
    TotalPages int         `json:"total_pages"`
    HasNext    bool        `json:"has_next"`
    HasPrev    bool        `json:"has_prev"`
}