// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"api-gateway/internal/config"
	"api-gateway/internal/load_balancers"
	"api-gateway/internal/services"
	"api-gateway/internal/services/interfaces"
)

// Injectors from wire.go:

func InitBackEndService() (interfaces.BackEndService, error) {
	backEndService := services.NewBackEndService()
	return backEndService, nil
}

func InitLoadBalancer() (interfaces.LoadBalancer, error) {
	loadBalancer := ProviderLoadBalancer()
	return loadBalancer, nil
}

func InitProxyService() (interfaces.ProxyService, error) {
	loadBalancer, err := InitLoadBalancer()
	if err != nil {
		return nil, err
	}
	backEndService, err := InitBackEndService()
	if err != nil {
		return nil, err
	}
	proxyService := services.NewProxyService(loadBalancer, backEndService)
	return proxyService, nil
}

// wire.go:

func ProviderLoadBalancer() interfaces.LoadBalancer {
	switch config.CFG.BalancerType {
	case "round_robin":
		return loadbalancers.NewRoundRobinBalancer()
	case "least_connections":
		return loadbalancers.NewLeastConnectionsBalancer()
	default:
		return loadbalancers.NewRoundRobinBalancer()
	}
}
