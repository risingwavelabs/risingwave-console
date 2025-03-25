/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { MetricsStoreLabelMatcherList } from './MetricsStoreLabelMatcherList';
import type { MetricsStoreSpec } from './MetricsStoreSpec';
export type MetricsStoreCreate = {
    name: string;
    spec: MetricsStoreSpec;
    defaultLabels?: MetricsStoreLabelMatcherList;
};

