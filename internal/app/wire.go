//go:build wireinject
// +build wireinject

package app

import (
	"api-gateway/internal/config"
	loadbalancers "api-gateway/internal/load_balancers"
	"api-gateway/internal/services"
	"api-gateway/internal/services/interfaces"

	"github.com/google/wire"
)

func InitBackEndService() (interfaces.BackEndService, error) {
	wire.Build(
		services.NewBackEndService,
	)
	return nil, nil
}

func ProviderLoadBalancer() interfaces.LoadBalancer {
	switch config.CFG.BalancerType {
	case "round_robin":
		return loadbalancers.NewRoundRobinBalancer()
	case "least_connections":
		return loadbalancers.NewLeastConnectionsBalancer()
	default:
		return loadbalancers.NewRoundRobinBalancer() // default load balancer
	}
}

func InitLoadBalancer() (interfaces.LoadBalancer, error) {
	wire.Build(
		ProviderLoadBalancer,
	)
	return nil, nil
}

func InitProxyService() (interfaces.ProxyService, error) {
	wire.Build(
		InitBackEndService,
		InitLoadBalancer,
		services.NewProxyService,
	)
	return nil, nil
}
