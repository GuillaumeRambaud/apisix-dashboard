/**
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import {
  Alert,
  Button,
  Card,
  Divider,
  Group,
  Loader,
  Space,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { createFileRoute } from '@tanstack/react-router';
import { getDefaultStore } from 'jotai';
import React, { useEffect,useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { getRouteListQueryOptions } from '@/apis/hooks';
import PageHeader from '@/components/page/PageHeader';
import { API_HEADER_KEY } from '@/config/constant';
import { queryClient } from '@/config/global';
import { adminKeyAtom } from '@/stores/global';
import { pageSearchSchema } from '@/types/schema/pageSearch';
import IconDownload from '~icons/material-symbols/download';
import IconInfo from '~icons/material-symbols/info-outline';
import IconUpload from '~icons/material-symbols/upload';

function RouteComponent() {
  const { t } = useTranslation();
  const adminKey = getDefaultStore().get(adminKeyAtom);
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const [importResult, setImportResult] = useState<string | null>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [lastImport, setLastImport] = useState<string | null>(null);

  // Fetch last import date from server
  const fetchLastImport = async () => {
    try {
      const res = await fetch('/apisix/admin/import/last', {
        headers: { [API_HEADER_KEY]: adminKey },
      });
      if (!res.ok) return;
      const data = await res.json();
      setLastImport(data.last_import);
    } catch (err) {
      if (process.env.NODE_ENV === 'development') {
        console.error('Failed to fetch last import:', err);
      }
    }
  };

  useEffect(() => {
    fetchLastImport();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [adminKey]);

  const handleImport = () => {
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
      fileInputRef.current.click();
    }
  };

  const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    try {
      setIsLoading(true);
      setErrorMsg(null);
      setImportResult(null);

      const content = await file.text();

      const res = await fetch('/apisix/admin/import', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-yaml',
          [API_HEADER_KEY]: adminKey,
        },
        body: content,
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `Import failed: ${res.status}`);
      }

      const result: { message?: string; last_import?: string } = await res.json();
      const message =
        result.message ||
        t('importExport.messages.importSuccess', 'Configuration imported successfully.');

      setImportResult(message);
      if (result.last_import) setLastImport(result.last_import);

      notifications.show({
        title: t('importExport.titles.importSuccess', 'Import successful'),
        message,
        color: 'green',
        icon: <IconUpload />,
      });
    } catch (err) {
      const msg =
        err instanceof Error
          ? err.message
          : t('importExport.errors.importFailed', 'Failed to import configuration');
      setErrorMsg(msg);
      notifications.show({
        title: t('importExport.titles.importFailed', 'Import failed'),
        message: msg,
        color: 'red',
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleExport = async () => {
    try {
      setIsLoading(true);
      const res = await fetch('/apisix/admin/export', {
        method: 'GET',
        headers: { [API_HEADER_KEY]: adminKey },
      });

      if (!res.ok) {
        throw new Error(`Export failed: ${res.status}`);
      }

      const blob = await res.blob();
      const timestamp = Date.now();
      const filename = `apisix-config_${timestamp}.yaml`;

      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);

      notifications.show({
        title: t('importExport.titles.exportSuccess', 'Export successful'),
        message: t('importExport.messages.exportFile', {
          defaultValue: 'Configuration exported as {{filename}}',
          filename,
        }),
        color: 'blue',
        icon: <IconDownload />,
      });
    } catch (err) {
      const msg =
        err instanceof Error
          ? err.message
          : t('importExport.errors.exportFailed', 'Failed to export configuration');
      notifications.show({
        title: t('importExport.titles.exportFailed', 'Export failed'),
        message: msg,
        color: 'red',
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <>
      <PageHeader title={t('importExport.sources.importExport', 'Import / Export')} />

      {/* Hidden file input */}
      <input
        type="file"
        accept=".yaml,.yml"
        ref={fileInputRef}
        style={{ display: 'none' }}
        onChange={handleFileChange}
      />

      <Space h="md" />

      <Card shadow="md" radius="md" withBorder>
        <Stack gap="md">
          <Title order={4}>{t('importExport.titles.configManagement', 'Configuration Management')}</Title>
          <Text c="dimmed">
            {t(
              'importExport.descriptions.importExport',
              'Use these tools to import or export your APISIX configuration.'
            )}
          </Text>

          {lastImport && (
            <Text style={{ color: 'blue', fontWeight: 500 }}>
              {t('importExport.labels.lastImport', { defaultValue: 'Last import: {{date}}', date: lastImport })}
            </Text>
          )}

          <Divider my="sm" />

          <Group>
            <Button
              onClick={handleImport}
              variant="light"
              disabled={isLoading}
              leftSection={
                isLoading ? (
                  <Loader size="xs" color="blue" />
                ) : (
                  <IconUpload className="h-4 w-4" />
                )
              }
            >
              {t('importExport.actions.import', 'Import Configuration')}
            </Button>

            <Button
              onClick={handleExport}
              variant="filled"
              disabled={isLoading}
              leftSection={<IconDownload className="h-4 w-4" />}
            >
              {t('importExport.actions.export', 'Export Configuration')}
            </Button>
          </Group>

          {importResult && (
            <Alert color="green" title={t('importExport.titles.importCompleted', 'Import completed')} icon={<IconInfo />}>
              {importResult}
            </Alert>
          )}

          {errorMsg && (
            <Alert color="red" title={t('importExport.titles.importFailed', 'Import failed')} icon={<IconInfo />}>
              {errorMsg}
            </Alert>
          )}
        </Stack>
      </Card>
    </>
  );
}

export const Route = createFileRoute('/import_export/')({
  component: RouteComponent,
  validateSearch: pageSearchSchema,
  loaderDeps: ({ search }) => search,
  loader: ({ deps }) =>
    queryClient.ensureQueryData(getRouteListQueryOptions(deps)),
});
