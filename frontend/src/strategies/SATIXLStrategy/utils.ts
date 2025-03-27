import { QueryParams } from '@/models/reviewable-item';

export function formatSATIXLFilters(filters: any): Partial<QueryParams> {
  return {
    reviewStatus: filters?.reviewStatus,
    skill: filters?.skill,
    searchTerm: filters?.searchTerm,
    sortBy: filters?.sortBy,
    sortOrder: filters?.sortOrder
  };
} 