/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type DatabaseConnectInfo = {
    /**
     * Name of the database
     */
    name: string;
    /**
     * ID of the cluster this database belongs to
     */
    clusterID: number;
    /**
     * Database username
     */
    username: string;
    /**
     * Database password (optional)
     */
    password?: string;
    /**
     * Database name
     */
    database: string;
};

