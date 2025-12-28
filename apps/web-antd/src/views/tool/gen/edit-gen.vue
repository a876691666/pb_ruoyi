<script setup lang="ts">
import type { GenInfo } from '#/api/tool/gen/model';

import { onMounted, provide, ref, unref, useTemplateRef } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';
import { useTabs } from '@vben/hooks';

import { Card, Skeleton, TabPane, Tabs } from 'ant-design-vue';

import { editSave, genInfo } from '#/api/tool/gen';

import { BasicSetting, GenConfig } from './edit-steps';
import { serializeGenData } from './gen-data-serializer';

const { setTabTitle, closeCurrentTab } = useTabs();
const routes = useRoute();
// 获取路由参数
const tableId = routes.params.tableId as string;

const genInfoData = ref<GenInfo>();

provide('genInfoData', genInfoData);

onMounted(async () => {
  const resp = await genInfo(tableId);
  // 需要做菜单转换 严格相等 才能选中回显
  // resp.info.parentMenuId = safeParseNumber(resp.info.parentMenuId);
  genInfoData.value = resp;
  setTabTitle(`生成配置: ${resp.name}`);
});

const currentTab = ref<'fields' | 'setting'>('setting');
const basicSettingRef = useTemplateRef('basicSettingRef');
const genConfigRef = useTemplateRef('genConfigRef');

const router = useRouter();
async function handleSave() {
  try {
    // 校验tab1
    const settingValidate = await basicSettingRef.value?.validateForm();
    if (!settingValidate) {
      currentTab.value = 'setting';
      return;
    }
    // 校验tab2
    const genConfigValidate = await genConfigRef.value?.validateTable();
    if (!genConfigValidate) {
      currentTab.value = 'fields';
      return;
    }
    // 获取表单数据
    const formValues = await basicSettingRef.value?.getFormValues();
    // 获取表格数据
    const tableRecords = genConfigRef.value?.getTableRecords() ?? [];
    // 序列化数据
    const requestData = serializeGenData(
      unref(genInfoData)!,
      formValues,
      tableRecords,
    );
    // 保存
    await editSave(requestData);
    // 关闭 & 跳转
    await closeCurrentTab();
    router.push({ path: '/tool/gen', replace: true });
  } catch (error) {
    console.error(error);
  }
}
</script>

<template>
  <Page :auto-content-height="true">
    <Card
      class="h-full"
      v-if="genInfoData"
      :body-style="{ padding: '0 16px 16px' }"
    >
      <Tabs v-model:active-key="currentTab" size="middle">
        <template #rightExtra>
          <!-- 因为编辑表格判断点击单元格之外的元素会取消编辑状态，此时需要事件拦截 -->
          <a-button
            class="vxe-table--ignore-clear"
            type="primary"
            @click="handleSave"
          >
            保存配置
          </a-button>
        </template>
        <TabPane key="setting" tab="生成信息" :force-render="true">
          <BasicSetting ref="basicSettingRef" />
        </TabPane>
        <TabPane key="fields" tab="字段信息" :force-render="true">
          <GenConfig ref="genConfigRef" />
        </TabPane>
      </Tabs>
    </Card>
    <Skeleton v-else :active="true" />
  </Page>
</template>
