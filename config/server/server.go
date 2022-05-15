package server

import (
	"context"
	"log"
	"panagiotisptr/config/protos"
	"panagiotisptr/config/server/config"
)

type ConfigServer struct {
	logger  *log.Logger
	service *config.Service
	protos.UnimplementedConfigServer
}

func NewConfigServer(logger *log.Logger) (*ConfigServer, error) {
	return &ConfigServer{
		logger:  logger,
		service: config.NewService(logger),
	}, nil
}

func (as *ConfigServer) SetParam(
	ctx context.Context,
	request *protos.SetParamRequest,
) (*protos.SetParamResponse, error) {
	err := as.service.SetParam(request.ParamName, request.ParamValue)

	return &protos.SetParamResponse{
		Success: err == nil,
	}, err
}

func (as *ConfigServer) UnsetParam(
	ctx context.Context,
	request *protos.UnsetParamRequest,
) (*protos.UnsetParamResponse, error) {
	err := as.service.UnsetParam(request.ParamName)

	return &protos.UnsetParamResponse{
		Success: err == nil,
	}, err
}

func (as *ConfigServer) GetParam(
	ctx context.Context,
	request *protos.GetParamRequest,
) (*protos.GetParamResponse, error) {
	response := &protos.GetParamResponse{
		ParamName:  "",
		ParamValue: "",
	}
	param, err := as.service.GetParam(request.ParamName)
	if err != nil {
		response.ParamName = param.Name
		response.ParamValue = param.Value
	}

	return response, err
}
