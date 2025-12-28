import type { BaseCollectionModel } from 'pocketbase';

export interface Client extends BaseCollectionModel {
  id: number;
  clientId: string;
  clientKey: string;
  clientSecret: string;
  grantTypeList: string[];
  grantType: string;
  deviceType: string;
  activeTimeout: number;
  timeout: number;
  status: string;
}
