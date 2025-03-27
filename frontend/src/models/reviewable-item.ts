export interface ReviewableItem {
  originalId: string;
  reviewStatus: string; // 'pending', 'approved', 'rejected'
  reviewComment?: string;
  reviewedAt?: string | null;
  reviewedBy?: string | null;
}

export interface ReviewStatusUpdate {
  status: string;
  comment: string;
  reviewedBy?: string;
}

export interface PaginatedResult<T> {
  items: T[];
  totalItems: number;
  totalPages: number;
  page: number;
  pageSize: number;
  dataSource: string;
}

export interface QueryParams {
  page: number;
  pageSize: number;
  reviewStatus?: string;
  domain?: string;
  skill?: string;
  searchTerm?: string;
  sortBy?: string;
  sortOrder?: string;
} 