/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { TaskCronjob } from './TaskCronjob';
import type { TaskRetryPolicy } from './TaskRetryPolicy';
export type TaskAttributes = {
    /**
     * Timeout of the task, e.g. 1h, 1d, 1w, 1m
     */
    timeout?: string;
    cronjob?: TaskCronjob;
    retryPolicy?: TaskRetryPolicy;
};

