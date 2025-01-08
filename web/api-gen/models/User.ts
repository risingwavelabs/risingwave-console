/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type User = {
    /**
     * Unique identifier of the user
     */
    ID: number;
    /**
     * User's name
     */
    name: string;
    /**
     * ID of the organization this user belongs to
     */
    organizationID: number;
    /**
     * Creation timestamp
     */
    createdAt: string;
    /**
     * Last update timestamp
     */
    updatedAt: string;
};

