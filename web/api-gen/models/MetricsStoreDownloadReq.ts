/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type MetricsStoreDownloadReq = {
    /**
     * Step of the metrics store, e.g. 1h, 1d, 1w, 1m, 1s
     */
    step?: string;
    /**
     * Start time of the metrics store
     */
    start?: string;
    /**
     * End time of the metrics store
     */
    end?: string;
    /**
     * (0, 1], if OOM, reduce the memory usage in Prometheus instance by this ratio (default: 1)
     *
     */
    queryRatio?: number;
    /**
     * query to get the metrics, e.g. `{namespace="risingwave-console"}`
     *
     */
    query?: string;
};

