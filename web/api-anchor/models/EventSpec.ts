/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { EventTaskCompleted } from './EventTaskCompleted';
import type { EventTaskError } from './EventTaskError';
export type EventSpec = {
    type: EventSpec.type;
    taskError?: EventTaskError;
    taskCompleted?: EventTaskCompleted;
};
export namespace EventSpec {
    export enum type {
        TASK_ERROR = 'TaskError',
        TASK_COMPLETED = 'TaskCompleted',
    }
}

