'use client';

import React from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { TableColumn } from '@/strategies/ReviewStrategy';

interface PaginationProps {
  current: number;
  pageSize: number;
  total: number;
  onChange: (page: number, pageSize: number) => void;
}

interface DataTableProps<T> {
  columns: TableColumn[];
  data: T[];
  loading?: boolean;
  pagination?: PaginationProps;
  onRowClick?: (record: T) => void;
}

export function DataTable<T extends { originalId: string }>({
  columns,
  data,
  loading = false,
  pagination,
  onRowClick,
}: DataTableProps<T>) {
  const handlePageChange = (page: number) => {
    if (pagination) {
      pagination.onChange(page, pagination.pageSize);
    }
  };

  const renderPagination = () => {
    if (!pagination) return null;

    const { current, pageSize, total } = pagination;
    const totalPages = Math.ceil(total / pageSize);
    
    return (
      <div className="flex items-center justify-between py-4">
        <div className="text-sm text-gray-500">
          总计 {total} 条记录，共 {totalPages} 页
        </div>
        <div className="flex space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => handlePageChange(current - 1)}
            disabled={current <= 1}
          >
            上一页
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => handlePageChange(current + 1)}
            disabled={current >= totalPages}
          >
            下一页
          </Button>
        </div>
      </div>
    );
  };

  return (
    <div>
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              {columns.map((column) => (
                <TableHead key={column.key} style={{ width: column.width }}>
                  {column.title}
                </TableHead>
              ))}
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              Array.from({ length: pagination?.pageSize || 10 }).map((_, index) => (
                <TableRow key={index}>
                  {columns.map((column) => (
                    <TableCell key={`${index}-${column.key}`}>
                      <Skeleton className="h-5 w-full" />
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : data.length === 0 ? (
              <TableRow>
                <TableCell colSpan={columns.length} className="h-24 text-center">
                  没有数据
                </TableCell>
              </TableRow>
            ) : (
              data.map((record: any) => (
                <TableRow
                  key={record.originalId}
                  className={onRowClick ? 'cursor-pointer hover:bg-gray-50' : ''}
                  onClick={() => onRowClick?.(record)}
                >
                  {columns.map((column) => (
                    <TableCell key={`${record.originalId}-${column.key}`}>
                      {column.render
                        ? column.render(record[column.dataIndex], record)
                        : record[column.dataIndex]}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
      {renderPagination()}
    </div>
  );
} 