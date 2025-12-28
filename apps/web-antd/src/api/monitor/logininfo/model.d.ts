import type { BaseCollectionModel } from 'pocketbase';

export interface LoginLog extends BaseCollectionModel {
  id: string;
  tenant_id: string;
  user_name: string;
  status: string;
  ipaddr: string;
  login_location: string;
  browser: string;
  os: string;
  msg: string;
  login_time: string;
  client_key: string;
}
