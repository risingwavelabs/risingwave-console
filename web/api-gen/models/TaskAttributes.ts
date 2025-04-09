/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { TaskCronjob } from './TaskCronjob';
import type { TaskRetryPolicy } from './TaskRetryPolicy';
export type TaskAttributes = {
    /**
     * If the task is created by a user, this field will be the organization ID of the user
     */
    orgID?: number;
    /**
     * Timeout of the task, e.g. 1h, 1d, 1w, 1m
     */
    timeout?: string;
    cronjob?: TaskCronjob;
    retryPolicy?: TaskRetryPolicy;
};

