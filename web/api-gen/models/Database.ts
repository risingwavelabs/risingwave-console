/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Table } from './Table';
export type Database = {
    /**
     * Unique identifier of the database
     */
    id: string;
    /**
     * Name of the database
     */
    name: string;
    /**
     * List of tables in the database
     */
    tables: Array<Table>;
};

