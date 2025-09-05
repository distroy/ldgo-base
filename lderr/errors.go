/*
 * Copyright (C) distroy
 */

package lderr

import (
	"context"
	"net/http"
)

const (
	errCodeSuccess   = 0
	errStatusSuccess = http.StatusOK
	errMessageSucess = "success"

	errCodeUnkown    = -1
	errStatusUnkonw  = http.StatusOK
	errMessageUnkonw = "unknown error"
)

// error definitions: [-1, -999]
var (
	ErrSuccess          = newInt(http.StatusOK, 0, "success")
	ErrUnkown           = newInt(http.StatusOK, -1, "unknown error")
	ErrServicePanic     = newInt(http.StatusServiceUnavailable, -2, "service panic")
	ErrClientPanic      = newInt(http.StatusServiceUnavailable, -3, "client panic")
	ErrInvalidParameter = newInt(http.StatusOK, -4, "invalid parameter")

	ErrCtxCanceled          = newIntByErr(http.StatusOK, -11, context.Canceled)
	ErrCtxDeadlineExceeded  = newIntByErr(http.StatusOK, -12, context.DeadlineExceeded)
	ErrCtxDeadlineNotEnough = newInt(http.StatusOK, -13, "context deadline not enough")

	ErrReflectError        = newInt(http.StatusOK, -21, "reflect error")
	ErrReflectTargetNotPtr = newInt(http.StatusOK, -22, "reflect target is not pointer")
	ErrReflectTargetNilPtr = newInt(http.StatusOK, -23, "reflect target is nil pointer")
	ErrReflectTypeNotEqual = newInt(http.StatusOK, -24, "reflect types of target and source are not equal")

	ErrNumberOverflow        = newInt(http.StatusOK, -31, "number overflow")
	ErrInvalidNumberSyntax   = newInt(http.StatusOK, -32, "invalid number syntax")
	ErrInvalidConvertType    = newInt(http.StatusOK, -33, "invalid convert type")
	ErrInvalidTemplateSyntax = newInt(http.StatusOK, -34, "invalid template syntax")

	ErrReadFile = newInt(http.StatusOK, -41, "read file fail")

	ErrMarshal   = newInt(http.StatusOK, -51, "marshal error")
	ErrUnmarshal = newInt(http.StatusOK, -52, "unmarshal error")

	ErrNonAuthoritativeInfo = newInt(http.StatusNonAuthoritativeInfo, -101, "http non authoritative info")
	ErrUnauthorized         = newInt(http.StatusUnauthorized, -102, "http unauthorized")
	ErrInternalServerError  = newInt(http.StatusInternalServerError, -103, "http internal server error")
	ErrTooManyRequests      = newInt(http.StatusTooManyRequests, -104, "too many requests")

	ErrInvalidRequestType = newInt(http.StatusOK, -201, "input request type must be pointer to struct")
	ErrParseRequest       = newInt(http.StatusOK, -202, "parse request error")
	ErrRequestValidation  = newInt(http.StatusOK, -203, "request validation error")
	ErrDuplicateRequest   = newInt(http.StatusOK, -204, "duplicate request")

	ErrHttpCall          = newInt(http.StatusOK, -301, "http call request error")
	ErrHttpReadBody      = newInt(http.StatusOK, -302, "http read body error")
	ErrHttpNewRequest    = newInt(http.StatusOK, -303, "http new request error")
	ErrHttpInvalidStatus = newInt(http.StatusOK, -304, "http invalid status")
	ErrHttpRenderBody    = newInt(http.StatusOK, -305, "http render body fail")

	ErrRpcCall      = newInt(http.StatusOK, -311, "rpc call error")
	ErrRpcTimeout   = newInt(http.StatusOK, -312, "rpc timeout")
	ErrRpcBokenPipe = newInt(http.StatusOK, -313, "rpc boken pip")
	ErrRpcBizError  = newInt(http.StatusOK, -314, "rpc call with biz error")

	ErrDataNotFound  = newInt(http.StatusOK, -801, "data not found")
	ErrDuplicateData = newInt(http.StatusOK, -802, "duplicate data")

	ErrDbQuery    = newInt(http.StatusOK, -811, "database query error")
	ErrDbExcute   = newInt(http.StatusOK, -812, "database execute error")
	ErrDbTxBegin  = newInt(http.StatusOK, -813, "database tx begin error")
	ErrDbTxCommit = newInt(http.StatusOK, -814, "database tx commit error")

	ErrCacheRead      = newInt(http.StatusOK, -821, "cache read error")
	ErrCacheWrite     = newInt(http.StatusOK, -822, "cache write error")
	ErrCacheTimeout   = newInt(http.StatusOK, -823, "cache timeout")
	ErrCacheMarshal   = newInt(http.StatusOK, -824, "cache codec marshal error")
	ErrCacheUnmarshal = newInt(http.StatusOK, -825, "cache codec unmarshal error")

	ErrCacheMutexLocked    = newInt(http.StatusOK, -831, "cache mutex had been locked")
	ErrCacheMutexNotExists = newInt(http.StatusOK, -832, "cache mutex is not exists")
	ErrCacheMutexNotMatch  = newInt(http.StatusOK, -833, "cache mutex is not match")
)
