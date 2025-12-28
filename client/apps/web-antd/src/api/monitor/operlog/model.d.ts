import type { BaseCollectionModel } from 'pocketbase';

export interface OperationLog extends BaseCollectionModel {
  id: string;
  tenant_id: string;
  title: string;
  business_type: string;
  business_types?: any;
  method: string;
  request_method: string;
  operator_type: number;
  oper_name: string;
  dept_name: string;
  oper_url: string;
  oper_ip: string;
  oper_location: string;
  oper_param: string;
  json_result: string;
  status: string;
  error_msg: string;
  oper_time: string;
  cost_time: number;
}
