package service

import (
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
)

func clusterToApi(cluster *querier.Cluster) *apigen.Cluster {
	return &apigen.Cluster{
		ID:                 cluster.ID,
		OrganizationID:     cluster.OrganizationID,
		Name:               cluster.Name,
		Host:               cluster.Host,
		Version:            cluster.Version,
		SqlPort:            cluster.SqlPort,
		MetaPort:           cluster.MetaPort,
		HttpPort:           cluster.HttpPort,
		CreatedAt:          cluster.CreatedAt,
		UpdatedAt:          cluster.UpdatedAt,
		PrometheusEndpoint: cluster.PrometheusEndpoint,
	}
}
