/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type Credentials = {
    /**
     * JWT access token
     */
    accessToken: string;
    /**
     * JWT refresh token for obtaining new access tokens
     */
    refreshToken: string;
    /**
     * Token type
     */
    tokenType: Credentials.tokenType;
};
export namespace Credentials {
    /**
     * Token type
     */
    export enum tokenType {
        BEARER = 'Bearer',
    }
}

