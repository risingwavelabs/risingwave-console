/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { TaskCronjob } from './TaskCronjob';
import type { TaskScheduled } from './TaskScheduled';
export type TaskAttributes = {
    /**
     * If the task is created by a user, this field will be the organization ID of the user
     */
    orgID?: number;
    cronjob?: TaskCronjob;
    scheduled?: TaskScheduled;
};

