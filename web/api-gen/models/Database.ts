/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Relation } from './Relation';
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
    organizationID: number;
    /**
     * Database username
     */
    username: string;
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
     * List of relations in the database
     */
    relations?: Array<Relation>;
};

