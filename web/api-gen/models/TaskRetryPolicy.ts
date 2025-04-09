/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type TaskRetryPolicy = {
    /**
     * Interval of the retry policy, e.g. 1h, 1d, 1w, 1m
     */
    interval: string;
    /**
     * Whether to always retry the task on failure
     */
    always_retry_on_failure?: boolean;
};

