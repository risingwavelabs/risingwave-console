/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type AutoBackupConfig = {
    /**
     * Whether automatic snapshots are enabled
     */
    enabled: boolean;
    /**
     * Cron expression for automatic snapshots (e.g., '0 0 * * *')
     */
    cronExpression: string;
    /**
     * Number of automatic snapshots to retain
     */
    keepLast: number;
};

