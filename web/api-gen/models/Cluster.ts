/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type Cluster = {
    ID: number;
    OrgID: number;
    name: string;
    host: string;
    sqlPort: number;
    metaPort: number;
    httpPort: number;
    version: string;
    /**
     * ID of the metrics store this cluster belongs to
     */
    metricsStoreID?: number;
    createdAt: string;
    updatedAt: string;
};

