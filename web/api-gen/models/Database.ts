/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Schema } from './Schema';
export type Database = {
    /**
     * Unique identifier of the database
     */
    ID: number;
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
    OrgID: number;
    /**
     * Database username
     */
    username: string;
    /**
     * Database name
     */
    database: string;
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
     * List of schemas in the database
     */
    schemas?: Array<Schema>;
};

