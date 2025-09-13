package grpc

import (
	"google.golang.org/grpc/codes"
)

// IntToCode converts an integer to gRPC codes.Code
// If the int value is not a valid gRPC code, it returns codes.Unknown as default
func IntToCode(codeInt int) codes.Code {
	// Validate if the int value is within the valid range of gRPC codes
	if codeInt < 0 || codeInt > int(codes.Unauthenticated) {
		return codes.Unknown
	}

	// Convert int to codes.Code
	code := codes.Code(codeInt)

	// Additional validation for known invalid codes
	switch code {
	case codes.OK, codes.Canceled, codes.Unknown, codes.InvalidArgument,
		codes.DeadlineExceeded, codes.NotFound, codes.AlreadyExists,
		codes.PermissionDenied, codes.ResourceExhausted, codes.FailedPrecondition,
		codes.Aborted, codes.OutOfRange, codes.Unimplemented, codes.Internal,
		codes.Unavailable, codes.DataLoss, codes.Unauthenticated:
		return code
	default:
		return codes.Unknown
	}
}

// IntToCodeWithDefault converts an integer to gRPC codes.Code with a custom default
// If the int value is not a valid gRPC code, it returns the provided default code
func IntToCodeWithDefault(codeInt int, defaultCode codes.Code) codes.Code {
	// Validate if the int value is within the valid range of gRPC codes
	if codeInt < 0 || codeInt > int(codes.Unauthenticated) {
		return defaultCode
	}

	// Convert int to codes.Code
	code := codes.Code(codeInt)

	// Additional validation for known invalid codes
	switch code {
	case codes.OK, codes.Canceled, codes.Unknown, codes.InvalidArgument,
		codes.DeadlineExceeded, codes.NotFound, codes.AlreadyExists,
		codes.PermissionDenied, codes.ResourceExhausted, codes.FailedPrecondition,
		codes.Aborted, codes.OutOfRange, codes.Unimplemented, codes.Internal,
		codes.Unavailable, codes.DataLoss, codes.Unauthenticated:
		return code
	default:
		return defaultCode
	}
}

// GetCodeFromHTTPStatus converts HTTP status code to appropriate gRPC codes.Code
// This is useful when converting HTTP responses to gRPC responses
func GetCodeFromHTTPStatus(httpStatus int) codes.Code {
	switch {
	case httpStatus >= 200 && httpStatus < 300:
		return codes.OK
	case httpStatus == 400:
		return codes.InvalidArgument
	case httpStatus == 401:
		return codes.Unauthenticated
	case httpStatus == 403:
		return codes.PermissionDenied
	case httpStatus == 404:
		return codes.NotFound
	case httpStatus == 409:
		return codes.AlreadyExists
	case httpStatus == 412:
		return codes.FailedPrecondition
	case httpStatus == 429:
		return codes.ResourceExhausted
	case httpStatus >= 500 && httpStatus < 600:
		return codes.Internal
	default:
		return codes.Unknown
	}
}
