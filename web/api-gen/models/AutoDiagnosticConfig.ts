/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type AutoDiagnosticConfig = {
    /**
     * Whether to enable automatic diagnostics
     */
    enabled: boolean;
    /**
     * Cron expression for diagnostic data collection (e.g., '0 0 * * *')
     */
    cronExpression: string;
    /**
     * How long to retain diagnostic data (e.g., '1d', '7d', '14d', '30d', '90d')
     */
    retentionDuration?: string;
};

