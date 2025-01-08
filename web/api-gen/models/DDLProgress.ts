/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type DDLProgress = {
    ID: number;
    statement: string;
    /**
     * Progress of the materialized view creation
     */
    progress: string;
    /**
     * When the DDL operation was initialized
     */
    initializedAt: string;
};

