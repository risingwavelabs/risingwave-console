/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type ClusterCreate = {
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
     * Version of the cluster
     */
    version: string;
};

