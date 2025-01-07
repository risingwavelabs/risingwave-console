/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type DiagnosticConfig = {
    /**
     * Interval between diagnostic data collections (e.g., '15m', '30m', '1h', '6h', '12h', '24h')
     */
    interval: string;
    /**
     * How long to retain diagnostic data (e.g., '1d', '7d', '14d', '30d', '90d')
     */
    expiration: string;
    /**
     * Whether to keep diagnostic data indefinitely
     */
    noExpiration: boolean;
};

