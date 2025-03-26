/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { TaskSpec } from './TaskSpec';
export type Task = {
    ID: number;
    workerName?: string;
    spec: TaskSpec;
    status: Task.status;
    remaining: number;
    createdAt?: string;
    updatedAt?: string;
};
export namespace Task {
    export enum status {
        PENDING = 'pending',
        RUNNING = 'running',
        COMPLETED = 'completed',
        FAILED = 'failed',
        PAUSED = 'paused',
        CANCELLED = 'cancelled',
    }
}

