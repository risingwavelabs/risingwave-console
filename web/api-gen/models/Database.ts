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
    clusterID: number;
    /**
     * ID of the organization this database belongs to
     */
    organizationID: number;
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
    createdAt: string;
    /**
     * Last update timestamp
     */
    updatedAt: string;
    /**
     * List of tables in the database
     */
    tables?: Array<Table>;
};

