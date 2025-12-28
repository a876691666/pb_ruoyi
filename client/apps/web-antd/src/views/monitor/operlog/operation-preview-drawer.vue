<script setup lang="ts">
import type { OperationLog } from '#/api/monitor/operlog/model';

import { computed, shallowRef } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { DictEnum } from '@vben/constants';

import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';

import {
  renderDict,
  renderHttpMethodTag,
  renderJsonPreview,
} from '#/utils/render';

const [BasicDrawer, drawerApi] = useVbenDrawer({
  onOpenChange: handleOpenChange,
  onClosed() {
    currentLog.value = null;
  },
});

const currentLog = shallowRef<null | OperationLog>(null);
function handleOpenChange(open: boolean) {
  if (!open) {
    return null;
  }
  const { record } = drawerApi.getData() as { record: any };
  currentLog.value = record;
}

const actionInfo = computed(() => {
  if (!currentLog.value) {
    return '-';
  }
  const data = currentLog.value;
  return `账号: ${data.oper_name || '-'} / ${data.dept_name || '-'} / ${data.oper_ip || '-'} / ${data.oper_location || '-'}`;
});
</script>

<template>
  <BasicDrawer :footer="false" class="w-[600px]" title="查看日志">
    <Descriptions v-if="currentLog" size="small" bordered :column="1">
      <DescriptionsItem label="日志编号" :label-style="{ minWidth: '120px' }">
        {{ currentLog.id }}
      </DescriptionsItem>
      <DescriptionsItem label="操作结果">
        <component
          :is="renderDict(currentLog.status, DictEnum.SYS_COMMON_STATUS)"
        />
      </DescriptionsItem>
      <DescriptionsItem label="操作模块">
        <div class="flex items-center">
          <Tag>{{ currentLog.title }}</Tag>
          <component
            :is="renderDict(currentLog.business_type, DictEnum.SYS_OPER_TYPE)"
          />
        </div>
      </DescriptionsItem>
      <DescriptionsItem label="操作信息">
        {{ actionInfo }}
      </DescriptionsItem>
      <DescriptionsItem label="请求信息">
        <component :is="renderHttpMethodTag(currentLog.request_method)" />
        {{ currentLog.oper_url }}
      </DescriptionsItem>
      <DescriptionsItem v-if="currentLog.error_msg" label="异常信息">
        <span class="font-semibold text-red-600">
          {{ currentLog.error_msg }}
        </span>
      </DescriptionsItem>
      <DescriptionsItem label="方法">
        {{ currentLog.method }}
      </DescriptionsItem>
      <DescriptionsItem label="请求参数">
        <div class="max-h-[300px] overflow-y-auto">
          <component :is="renderJsonPreview(currentLog.oper_param)" />
        </div>
      </DescriptionsItem>
      <DescriptionsItem v-if="currentLog.json_result" label="响应参数">
        <div class="max-h-[300px] overflow-y-auto">
          <component :is="renderJsonPreview(currentLog.json_result)" />
        </div>
      </DescriptionsItem>
      <DescriptionsItem label="请求耗时">
        {{ `${currentLog.cost_time} ms` }}
      </DescriptionsItem>
      <DescriptionsItem label="操作时间">
        {{ `${currentLog.oper_time}` }}
      </DescriptionsItem>
    </Descriptions>
  </BasicDrawer>
</template>
