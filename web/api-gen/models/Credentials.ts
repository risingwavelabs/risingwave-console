/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type Credentials = {
    /**
     * JWT access token
     */
    access_token: string;
    /**
     * JWT refresh token for obtaining new access tokens
     */
    refresh_token: string;
    /**
     * Token type
     */
    token_type: Credentials.token_type;
};
export namespace Credentials {
    /**
     * Token type
     */
    export enum token_type {
        BEARER = 'Bearer',
    }
}

