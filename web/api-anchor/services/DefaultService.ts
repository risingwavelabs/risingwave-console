/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Credentials } from '../models/Credentials';
import type { Event } from '../models/Event';
import type { Org } from '../models/Org';
import type { RefreshTokenRequest } from '../models/RefreshTokenRequest';
import type { SignInRequest } from '../models/SignInRequest';
import type { Task } from '../models/Task';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class DefaultService {
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
     * Sign out user
     * Sign out user and invalidate all tokens
     * @returns any Successfully signed out
     * @throws ApiError
     */
    public static signOut(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/auth/sign-out',
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
    /**
     * Get all tasks
     * Get all tasks
     * @returns Task Successfully retrieved tasks
     * @throws ApiError
     */
    public static listTasks(): CancelablePromise<Array<Task>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/tasks',
        });
    }
    /**
     * Get all events
     * Get all events
     * @returns Event Successfully retrieved events
     * @throws ApiError
     */
    public static listEvents(): CancelablePromise<Array<Event>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/events',
        });
    }
    /**
     * Get all organizations of which the user is a member
     * Get all organizations of which the user is a member
     * @returns Org Successfully retrieved organizations
     * @throws ApiError
     */
    public static listOrgs(): CancelablePromise<Array<Org>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/orgs',
        });
    }
}
