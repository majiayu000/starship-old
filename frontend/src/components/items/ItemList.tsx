'use client';

import React, { useState, useEffect } from 'react';
import { GenericItem, PaginatedResult, QueryParams } from '@/models/generic-item';
import { genericApi } from '@/api/generic-api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';

import { toast } from 'sonner';

interface ItemListProps {
  onItemSelect?: (item: GenericItem) => void;
}

export default function ItemList({ onItemSelect }: ItemListProps) {
  // State
  const [items, setItems] = useState<PaginatedResult<GenericItem> | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [queryParams, setQueryParams] = useState<QueryParams>({
    page: 1,
    pageSize: 10,
    sortBy: 'createdAt',
    sortOrder: 'desc',
  });
  const [searchTerm, setSearchTerm] = useState<string>('');
  const [statusFilter, setStatusFilter] = useState<string>('all');

  // Load items
  const loadItems = async () => {
    setLoading(true);
    try {
      const result = await genericApi.getAll(queryParams);
      setItems(result);
    } catch (error) {
      console.error('Error loading items:', error);
      toast.error('Failed to load items');
    } finally {
      setLoading(false);
    }
  };

  // Effect to load items when query params change
  useEffect(() => {
    loadItems();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [queryParams]);

  // Handle search
  const handleSearch = () => {
    setQueryParams({
      ...queryParams,
      page: 1, // Reset to first page
      searchTerm,
    });
  };

  // Handle status filter change
  const handleStatusChange = (value: string) => {
    setStatusFilter(value);
    setQueryParams({
      ...queryParams,
      page: 1, // Reset to first page
      status: value === 'all' ? '' : value, // Use empty string for 'all'
    });
  };

  // Handle page change
  const handlePageChange = (page: number) => {
    setQueryParams({
      ...queryParams,
      page,
    });
  };

  // Handle sort change
  const handleSortChange = (field: string) => {
    const newSortOrder = queryParams.sortBy === field && queryParams.sortOrder === 'asc' ? 'desc' : 'asc';
    setQueryParams({
      ...queryParams,
      sortBy: field,
      sortOrder: newSortOrder,
    });
  };

  // Render status badge
  const renderStatusBadge = (status: string) => {
    let variant: 'default' | 'secondary' | 'destructive' | 'outline' = 'default';

    switch (status) {
      case 'active':
        variant = 'default';
        break;
      case 'draft':
        variant = 'secondary';
        break;
      case 'inactive':
        variant = 'outline';
        break;
      case 'archived':
        variant = 'destructive';
        break;
    }

    return <Badge variant={variant}>{status}</Badge>;
  };

  // Render pagination
  const renderPagination = () => {
    if (!items) return null;

    const { page, totalPages } = items;

    // Calculate page range
    const pageNumbers: number[] = [];
    const maxPagesToShow = 5;

    let startPage = Math.max(1, page - Math.floor(maxPagesToShow / 2));
    const endPage = Math.min(totalPages, startPage + maxPagesToShow - 1);

    if (endPage - startPage + 1 < maxPagesToShow) {
      startPage = Math.max(1, endPage - maxPagesToShow + 1);
    }

    for (let i = startPage; i <= endPage; i++) {
      pageNumbers.push(i);
    }

    return (
      <div className="flex items-center justify-center space-x-2">
        <Button
          variant="outline"
          size="sm"
          onClick={() => handlePageChange(Math.max(1, page - 1))}
          disabled={page <= 1}
        >
          Previous
        </Button>

        {pageNumbers.map((pageNum) => (
          <Button
            key={pageNum}
            variant={pageNum === page ? "default" : "outline"}
            size="sm"
            onClick={() => handlePageChange(pageNum)}
          >
            {pageNum}
          </Button>
        ))}

        <Button
          variant="outline"
          size="sm"
          onClick={() => handlePageChange(Math.min(totalPages, page + 1))}
          disabled={page >= totalPages}
        >
          Next
        </Button>
      </div>
    );
  };

  return (
    <div className="space-y-4">
      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-4">
        <div className="flex-1">
          <Input
            placeholder="Search items..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
          />
        </div>

        <div className="w-full sm:w-48">
          <Select value={statusFilter} onValueChange={handleStatusChange}>
            <SelectTrigger>
              <SelectValue placeholder="Filter by status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All statuses</SelectItem>
              <SelectItem value="draft">Draft</SelectItem>
              <SelectItem value="active">Active</SelectItem>
              <SelectItem value="inactive">Inactive</SelectItem>
              <SelectItem value="archived">Archived</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <Button onClick={handleSearch}>Search</Button>
      </div>

      {/* Table */}
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead
                className="cursor-pointer"
                onClick={() => handleSortChange('name')}
              >
                Name
                {queryParams.sortBy === 'name' && (
                  <span className="ml-1">{queryParams.sortOrder === 'asc' ? '↑' : '↓'}</span>
                )}
              </TableHead>
              <TableHead>Description</TableHead>
              <TableHead
                className="cursor-pointer"
                onClick={() => handleSortChange('status')}
              >
                Status
                {queryParams.sortBy === 'status' && (
                  <span className="ml-1">{queryParams.sortOrder === 'asc' ? '↑' : '↓'}</span>
                )}
              </TableHead>
              <TableHead
                className="cursor-pointer"
                onClick={() => handleSortChange('createdAt')}
              >
                Created
                {queryParams.sortBy === 'createdAt' && (
                  <span className="ml-1">{queryParams.sortOrder === 'asc' ? '↑' : '↓'}</span>
                )}
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={4} className="text-center py-8">
                  Loading...
                </TableCell>
              </TableRow>
            ) : items && items.items.length > 0 ? (
              items.items.map((item) => (
                <TableRow
                  key={item.id}
                  className="cursor-pointer hover:bg-muted/50"
                  onClick={() => onItemSelect && onItemSelect(item)}
                >
                  <TableCell className="font-medium">{item.name}</TableCell>
                  <TableCell className="max-w-xs truncate">{item.description}</TableCell>
                  <TableCell>{renderStatusBadge(item.status)}</TableCell>
                  <TableCell>{new Date(item.createdAt).toLocaleDateString()}</TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={4} className="text-center py-8">
                  No items found
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      {/* Pagination */}
      {items && items.totalPages > 1 && (
        <div className="flex justify-center mt-4">
          {renderPagination()}
        </div>
      )}

      {/* Summary */}
      {items && (
        <div className="text-sm text-muted-foreground">
          Showing {items.items.length} of {items.totalItems} items
        </div>
      )}
    </div>
  );
}
