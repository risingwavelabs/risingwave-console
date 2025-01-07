/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Table } from './Table';
export type Database = {
    /**
     * Unique identifier of the database
     */
    id: number;
    /**
     * Name of the database
     */
    name: string;
    /**
     * ID of the cluster this database belongs to
     */
    cluster_id: number;
    /**
     * ID of the organization this database belongs to
     */
    organization_id: number;
    /**
     * Database username (optional)
     */
    username?: string;
    /**
     * Database password (optional)
     */
    password?: string;
    /**
     * Creation timestamp
     */
    created_at: string;
    /**
     * Last update timestamp
     */
    updated_at: string;
    /**
     * List of tables in the database
     */
    tables?: Array<Table>;
};

