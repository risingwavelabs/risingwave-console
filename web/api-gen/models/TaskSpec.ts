/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { TaskSpecAutoBackup } from './TaskSpecAutoBackup';
import type { TaskSpecAutoDiagnostic } from './TaskSpecAutoDiagnostic';
import type { TaskSpecDeleteClusterDiagnostic } from './TaskSpecDeleteClusterDiagnostic';
import type { TaskSpecDeleteOpaqueKey } from './TaskSpecDeleteOpaqueKey';
import type { TaskSpecDeleteSnapshot } from './TaskSpecDeleteSnapshot';
export type TaskSpec = {
    type: TaskSpec.type;
    autoBackup?: TaskSpecAutoBackup;
    autoDiagnostic?: TaskSpecAutoDiagnostic;
    deleteSnapshot?: TaskSpecDeleteSnapshot;
    deleteClusterDiagnostic?: TaskSpecDeleteClusterDiagnostic;
    deleteOpaqueKey?: TaskSpecDeleteOpaqueKey;
};
export namespace TaskSpec {
    export enum type {
        AUTO_BACKUP = 'auto-backup',
        AUTO_DIAGNOSTIC = 'auto-diagnostic',
        DELETE_SNAPSHOT = 'delete-snapshot',
        DELETE_CLUSTER_DIAGNOSTIC = 'delete-cluster-diagnostic',
        DELETE_OPAQUE_KEY = 'delete-opaque-key',
    }
}

