/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { TaskSpecAutoBackup } from './TaskSpecAutoBackup';
import type { TaskSpecAutoDiagnostic } from './TaskSpecAutoDiagnostic';
export type TaskSpec = {
    type: TaskSpec.type;
    autoBackup?: TaskSpecAutoBackup;
    autoDiagnostic?: TaskSpecAutoDiagnostic;
};
export namespace TaskSpec {
    export enum type {
        AUTO_BACKUP = 'auto-backup',
        AUTO_DIAGNOSTIC = 'auto-diagnostic',
    }
}

