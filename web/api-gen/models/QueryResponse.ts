/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Column } from './Column';
export type QueryResponse = {
    columns: Array<Column>;
    rows: Array<Record<string, any>>;
    /**
     * Number of rows affected by the query
     */
    rowsAffected: number;
    /**
     * Error message if the query failed
     */
    error?: string;
};

