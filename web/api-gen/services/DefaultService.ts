/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Database } from '../models/Database';
import type { Table } from '../models/Table';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class DefaultService {
    /**
     * List all databases
     * Retrieve a list of all databases and their tables
     * @returns Database Successfully retrieved database list
     * @throws ApiError
     */
    public static listDatabases(): CancelablePromise<Array<Database>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/databases',
        });
    }
    /**
     * Create a new database
     * Create a new database
     * @param requestBody
     * @returns Database Database created successfully
     * @throws ApiError
     */
    public static createDatabase(
        requestBody: {
            /**
             * Name of the database
             */
            name: string;
        },
    ): CancelablePromise<Database> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/databases',
            body: requestBody,
            mediaType: 'application/json',
        });
    }
    /**
     * Get database details
     * Retrieve details of a specific database
     * @param databaseId
     * @returns Database Successfully retrieved database
     * @throws ApiError
     */
    public static getDatabase(
        databaseId: string,
    ): CancelablePromise<Database> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/databases/{databaseId}',
            path: {
                'databaseId': databaseId,
            },
        });
    }
    /**
     * Update database
     * Update a specific database
     * @param databaseId
     * @param requestBody
     * @returns Database Database updated successfully
     * @throws ApiError
     */
    public static updateDatabase(
        databaseId: string,
        requestBody: {
            /**
             * New name of the database
             */
            name?: string;
        },
    ): CancelablePromise<Database> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/databases/{databaseId}',
            path: {
                'databaseId': databaseId,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }
    /**
     * Delete database
     * Delete a specific database
     * @param databaseId
     * @returns void
     * @throws ApiError
     */
    public static deleteDatabase(
        databaseId: string,
    ): CancelablePromise<void> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/databases/{databaseId}',
            path: {
                'databaseId': databaseId,
            },
        });
    }
    /**
     * List database tables
     * Retrieve all tables in a specific database
     * @param databaseId
     * @returns Table Successfully retrieved tables
     * @throws ApiError
     */
    public static listDatabaseTables(
        databaseId: string,
    ): CancelablePromise<Array<Table>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/databases/{databaseId}/tables',
            path: {
                'databaseId': databaseId,
            },
        });
    }
    /**
     * Create table
     * Create a new table in the specified database
     * @param databaseId
     * @param requestBody
     * @returns Table Table created successfully
     * @throws ApiError
     */
    public static createTable(
        databaseId: string,
        requestBody: {
            /**
             * Name of the table
             */
            name: string;
        },
    ): CancelablePromise<Table> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/databases/{databaseId}/tables',
            path: {
                'databaseId': databaseId,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }
    /**
     * Get table details
     * Retrieve details of a specific table
     * @param databaseId
     * @param tableId
     * @returns Table Successfully retrieved table
     * @throws ApiError
     */
    public static getTable(
        databaseId: string,
        tableId: string,
    ): CancelablePromise<Table> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/databases/{databaseId}/tables/{tableId}',
            path: {
                'databaseId': databaseId,
                'tableId': tableId,
            },
        });
    }
    /**
     * Update table
     * Update a specific table
     * @param databaseId
     * @param tableId
     * @param requestBody
     * @returns Table Table updated successfully
     * @throws ApiError
     */
    public static updateTable(
        databaseId: string,
        tableId: string,
        requestBody: {
            /**
             * New name of the table
             */
            name?: string;
        },
    ): CancelablePromise<Table> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/databases/{databaseId}/tables/{tableId}',
            path: {
                'databaseId': databaseId,
                'tableId': tableId,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }
    /**
     * Delete table
     * Delete a specific table
     * @param databaseId
     * @param tableId
     * @returns void
     * @throws ApiError
     */
    public static deleteTable(
        databaseId: string,
        tableId: string,
    ): CancelablePromise<void> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/databases/{databaseId}/tables/{tableId}',
            path: {
                'databaseId': databaseId,
                'tableId': tableId,
            },
        });
    }
}
