/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Column } from './Column';
export type Relation = {
    /**
     * Unique identifier of the table
     */
    ID: number;
    /**
     * Name of the table
     */
    name: string;
    /**
     * Type of the relation
     */
    type: Relation.type;
    /**
     * List of columns in the table
     */
    columns: Array<Column>;
    dependencies: Array<number>;
};
export namespace Relation {
    /**
     * Type of the relation
     */
    export enum type {
        TABLE = 'table',
        SOURCE = 'source',
        SINK = 'sink',
        MATERIALIZED_VIEW = 'materializedView',
        SYSTEM_TABLE = 'system table',
    }
}

