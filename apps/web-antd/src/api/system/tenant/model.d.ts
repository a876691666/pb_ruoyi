import type { BaseCollectionModel } from 'pocketbase';

export interface Tenant extends BaseCollectionModel {
  account_count: number;
  address?: string;
  company_name: string;
  contact_phone: string;
  contact_user_name: string;
  domain?: string;
  expire_time?: string;
  intro: string;
  license_number?: any;
  id: string;
  remark?: string;
  status: string;
  package_id: string;
}
