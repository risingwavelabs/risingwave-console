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
     * Database host address
     */
    host: string;
    /**
     * SQL connection port
     */
    sqlPort: number;
    /**
     * Metadata service port
     */
    metaPort: number;
    /**
     * Database username (optional)
     */
    username?: string;
    /**
     * Database password (optional)
     */
    password?: string;
};

