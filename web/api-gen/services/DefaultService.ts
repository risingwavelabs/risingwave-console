/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Cluster } from '../models/Cluster';
import type { ClusterCreate } from '../models/ClusterCreate';
import type { Credentials } from '../models/Credentials';
import type { Database } from '../models/Database';
import type { DatabaseConnectInfo } from '../models/DatabaseConnectInfo';
import type { DiagnosticConfig } from '../models/DiagnosticConfig';
import type { DiagnosticData } from '../models/DiagnosticData';
import type { RefreshTokenRequest } from '../models/RefreshTokenRequest';
import type { SignInRequest } from '../models/SignInRequest';
import type { Snapshot } from '../models/Snapshot';
import type { SnapshotConfig } from '../models/SnapshotConfig';
import type { SnapshotCreate } from '../models/SnapshotCreate';
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
        requestBody: DatabaseConnectInfo,
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
     * @param id
     * @returns Database Successfully retrieved database
     * @throws ApiError
     */
    public static getDatabase(
        id: number,
    ): CancelablePromise<Database> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/databases/{id}',
            path: {
                'id': id,
            },
        });
    }
    /**
     * Update database
     * Update a specific database
     * @param id
     * @param requestBody
     * @returns Database Database updated successfully
     * @throws ApiError
     */
    public static updateDatabase(
        id: number,
        requestBody: DatabaseConnectInfo,
    ): CancelablePromise<Database> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/databases/{id}',
            path: {
                'id': id,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }
    /**
     * Delete database
     * Delete a specific database
     * @param id
     * @returns void
     * @throws ApiError
     */
    public static deleteDatabase(
        id: number,
    ): CancelablePromise<void> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/databases/{id}',
            path: {
                'id': id,
            },
        });
    }
    /**
     * List all clusters
     * Retrieve a list of all database clusters
     * @returns Cluster Successfully retrieved cluster list
     * @throws ApiError
     */
    public static listClusters(): CancelablePromise<Array<Cluster>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/clusters',
        });
    }
    /**
     * Create a new cluster
     * Create a new database cluster
     * @param requestBody
     * @returns Cluster Cluster created successfully
     * @throws ApiError
     */
    public static createCluster(
        requestBody: ClusterCreate,
    ): CancelablePromise<Cluster> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/clusters',
            body: requestBody,
            mediaType: 'application/json',
        });
    }
    /**
     * Get cluster details
     * Retrieve details of a specific cluster
     * @param id
     * @returns Cluster Successfully retrieved cluster
     * @throws ApiError
     */
    public static getCluster(
        id: string,
    ): CancelablePromise<Cluster> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/clusters/{id}',
            path: {
                'id': id,
            },
        });
    }
    /**
     * Update cluster
     * Update a specific cluster
     * @param id
     * @param requestBody
     * @returns Cluster Cluster updated successfully
     * @throws ApiError
     */
    public static updateCluster(
        id: string,
        requestBody: ClusterCreate,
    ): CancelablePromise<Cluster> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/clusters/{id}',
            path: {
                'id': id,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }
    /**
     * Delete cluster
     * Delete a specific cluster
     * @param id
     * @returns void
     * @throws ApiError
     */
    public static deleteCluster(
        id: string,
    ): CancelablePromise<void> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/clusters/{id}',
            path: {
                'id': id,
            },
        });
    }
    /**
     * List cluster snapshots
     * Retrieve a list of all snapshots for a specific cluster
     * @param id
     * @returns Snapshot Successfully retrieved snapshot list
     * @throws ApiError
     */
    public static listClusterSnapshots(
        id: string,
    ): CancelablePromise<Array<Snapshot>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/clusters/{id}/snapshots',
            path: {
                'id': id,
            },
        });
    }
    /**
     * Create a new snapshot
     * Create a new metadata snapshot for a specific cluster
     * @param id
     * @param requestBody
     * @returns Snapshot Snapshot created successfully
     * @throws ApiError
     */
    public static createClusterSnapshot(
        id: string,
        requestBody: SnapshotCreate,
    ): CancelablePromise<Snapshot> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/clusters/{id}/snapshots',
            path: {
                'id': id,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }
    /**
     * Delete snapshot
     * Delete a specific snapshot
     * @param id
     * @param snapshotId
     * @returns void
     * @throws ApiError
     */
    public static deleteClusterSnapshot(
        id: string,
        snapshotId: string,
    ): CancelablePromise<void> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/clusters/{id}/snapshots/{snapshotId}',
            path: {
                'id': id,
                'snapshotId': snapshotId,
            },
        });
    }
    /**
     * Restore snapshot
     * Restore cluster metadata from a specific snapshot
     * @param id
     * @param snapshotId
     * @returns any Snapshot restored successfully
     * @throws ApiError
     */
    public static restoreClusterSnapshot(
        id: string,
        snapshotId: string,
    ): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/clusters/{id}/snapshots/{snapshotId}',
            path: {
                'id': id,
                'snapshotId': snapshotId,
            },
        });
    }
    /**
     * Get snapshot configuration
     * Get automatic snapshot configuration for a cluster
     * @param id
     * @returns SnapshotConfig Successfully retrieved snapshot configuration
     * @throws ApiError
     */
    public static getClusterSnapshotConfig(
        id: string,
    ): CancelablePromise<SnapshotConfig> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/clusters/{id}/snapshot-config',
            path: {
                'id': id,
            },
        });
    }
    /**
     * Update snapshot configuration
     * Update automatic snapshot configuration for a cluster
     * @param id
     * @param requestBody
     * @returns SnapshotConfig Snapshot configuration updated successfully
     * @throws ApiError
     */
    public static updateClusterSnapshotConfig(
        id: string,
        requestBody: SnapshotConfig,
    ): CancelablePromise<SnapshotConfig> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/clusters/{id}/snapshot-config',
            path: {
                'id': id,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }
    /**
     * List diagnostic data
     * Retrieve diagnostic data for a specific cluster with optional date range filtering
     * @param id
     * @param from Start date for filtering diagnostic data
     * @param to End date for filtering diagnostic data
     * @param page Page number for pagination
     * @param perPage Number of items per page
     * @returns any Successfully retrieved diagnostic data
     * @throws ApiError
     */
    public static listClusterDiagnostics(
        id: string,
        from?: string,
        to?: string,
        page: number = 1,
        perPage: number = 20,
    ): CancelablePromise<{
        items: Array<DiagnosticData>;
        /**
         * Total number of diagnostic entries
         */
        total: number;
        /**
         * Current page number
         */
        page: number;
        /**
         * Number of items per page
         */
        per_page: number;
    }> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/clusters/{id}/diagnostics',
            path: {
                'id': id,
            },
            query: {
                'from': from,
                'to': to,
                'page': page,
                'per_page': perPage,
            },
        });
    }
    /**
     * Get diagnostic configuration
     * Get diagnostic data collection configuration for a cluster
     * @param id
     * @returns DiagnosticConfig Successfully retrieved diagnostic configuration
     * @throws ApiError
     */
    public static getClusterDiagnosticConfig(
        id: string,
    ): CancelablePromise<DiagnosticConfig> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/clusters/{id}/diagnostics/config',
            path: {
                'id': id,
            },
        });
    }
    /**
     * Update diagnostic configuration
     * Update diagnostic data collection configuration for a cluster
     * @param id
     * @param requestBody
     * @returns DiagnosticConfig Diagnostic configuration updated successfully
     * @throws ApiError
     */
    public static updateClusterDiagnosticConfig(
        id: string,
        requestBody: DiagnosticConfig,
    ): CancelablePromise<DiagnosticConfig> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/clusters/{id}/diagnostics/config',
            path: {
                'id': id,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }
    /**
     * Sign in user
     * Authenticate user and return access token
     * @param requestBody
     * @returns Credentials Successfully authenticated
     * @throws ApiError
     */
    public static signIn(
        requestBody: SignInRequest,
    ): CancelablePromise<Credentials> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/auth/sign-in',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `Invalid credentials`,
            },
        });
    }
    /**
     * Refresh access token
     * Get a new access token using a refresh token
     * @param requestBody
     * @returns Credentials Successfully refreshed token
     * @throws ApiError
     */
    public static refreshToken(
        requestBody: RefreshTokenRequest,
    ): CancelablePromise<Credentials> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/auth/refresh',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `Invalid or expired refresh token`,
            },
        });
    }
}
