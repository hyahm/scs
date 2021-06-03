package client

import "errors"

var ErrPnameIsEmpty = errors.New("pname is empty")
var ErrNameIsEmpty = errors.New("name is empty")
var ErrToken = errors.New("token error")
var ErrStatusNetworkAuthenticationRequired = errors.New("StatusNetworkAuthenticationRequired")
var ErrResponseData = errors.New("error response data")
var ErrFoundPname = errors.New("not found pname")
var ErrFoundName = errors.New("not found name")
var ErrFoundPnameOrName = errors.New("not found pname or name")
var ErrWaitReload = errors.New("waiting for last reload complete")
var ErrHttps = errors.New("Client sent an HTTP request to an HTTPS server.")
