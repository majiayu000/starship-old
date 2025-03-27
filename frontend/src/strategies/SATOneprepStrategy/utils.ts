import { QueryParams } from '@/models/reviewable-item';

export const formatSATOneprepFilters = (filters: any): Partial<QueryParams> => {
  return {
    reviewStatus: filters.status || undefined,
    domain: filters.domain || undefined,
    skill: filters.skill || undefined,
    searchTerm: filters.searchTerm || undefined,
    sortBy: filters.sortBy || undefined,
    sortOrder: filters.sortOrder || undefined
  };
}; 