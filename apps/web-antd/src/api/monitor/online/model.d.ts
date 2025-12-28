import type { BaseCollectionModel } from 'pocketbase';

export interface OnlineUser extends BaseCollectionModel {
  id: string;
  deptName: string;
  userName: string;
  ipaddr: string;
  loginLocation: string;
  browser: string;
  os: string;
  loginTime: number;
  deviceType: string;
}
