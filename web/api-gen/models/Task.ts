/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { TaskSpec } from './TaskSpec';
export type Task = {
    ID: number;
    spec: TaskSpec;
    status: Task.status;
    startedAt?: string;
    createdAt: string;
    updatedAt: string;
};
export namespace Task {
    export enum status {
        PENDING = 'pending',
        COMPLETED = 'completed',
        FAILED = 'failed',
        PAUSED = 'paused',
    }
}

