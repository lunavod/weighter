package scales

import "reflect"

type GetMassaResponse struct {
	Response
	Weight   uint32
	Division byte
	Stable   bool
	Net      bool
	Zero     bool
}

func ReadGetMassaResponse(s Connection) (GetMassaResponse, error) {
	resp := GetMassaResponse{}

	err := resp.Response.Read(s)
	if err != nil {
		return resp, err
	}

	FillResponseStruct(resp.Raw(), reflect.ValueOf(&resp).Elem())
	return resp, nil
}

type GetNameResponse struct {
	Response
	ScalesID uint32
	Name     string
}

func ReadGetNameResponse(s Connection) (GetNameResponse, error) {
	resp := GetNameResponse{}

	err := resp.Response.Read(s)
	if err != nil {
		return resp, err
	}

	FillResponseStruct(resp.Raw(), reflect.ValueOf(&resp).Elem())
	return resp, nil
}

type GetWifiIpResponse struct {
	Response
	Ip      [4]byte
	Mask    [4]byte
	Gateway [4]byte
	IpAP    [4]byte
	Port    uint16
}

func ReadGetWifiIpResponse(s Connection) (GetWifiIpResponse, error) {
	resp := GetWifiIpResponse{}

	err := resp.Response.Read(s)
	if err != nil {
		return resp, err
	}

	FillResponseStruct(resp.Raw(), reflect.ValueOf(&resp).Elem())
	return resp, nil
}