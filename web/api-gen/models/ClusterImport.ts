/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type ClusterImport = {
    /**
     * Name of the cluster
     */
    name: string;
    /**
     * Cluster host address
     */
    host: string;
    /**
     * SQL connection port
     */
    sqlPort: number;
    /**
     * Metadata node port
     */
    metaPort: number;
    /**
     * HTTP port
     */
    httpPort: number;
    /**
     * Version of the cluster
     */
    version: string;
    /**
     * ID of the metrics store this cluster belongs to
     */
    metricsStoreID?: number;
};

