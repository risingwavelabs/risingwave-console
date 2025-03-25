/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type MetricsStoreLabelMatcher = {
    op: MetricsStoreLabelMatcher.op;
    /**
     * Label key
     */
    key: string;
    /**
     * Label value
     */
    value: string;
};
export namespace MetricsStoreLabelMatcher {
    export enum op {
        EQ = 'EQ',
        NEQ = 'NEQ',
        RE = 'RE',
        NRE = 'NRE',
    }
}

