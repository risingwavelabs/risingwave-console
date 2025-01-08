/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Column } from './Column';
export type Table = {
    /**
     * Unique identifier of the table
     */
    id: string;
    /**
     * Name of the table
     */
    name: string;
    /**
     * Type of the relation
     */
    type: Table.type;
    /**
     * List of columns in the table
     */
    columns: Array<Column>;
};
export namespace Table {
    /**
     * Type of the relation
     */
    export enum type {
        TABLE = 'table',
        SOURCE = 'source',
        SINK = 'sink',
        MATERIALIZED_VIEW = 'materializedView',
    }
}

