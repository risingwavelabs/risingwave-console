/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { EventTaskError } from './EventTaskError';
export type EventSpec = {
    type: EventSpec.type;
    taskError?: EventTaskError;
};
export namespace EventSpec {
    export enum type {
        TASK_ERROR = 'TaskError',
    }
}

