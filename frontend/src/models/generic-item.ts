export interface BaseEntity {
  id: string;
  createdAt: string;
  updatedAt: string;
}

export interface GenericItem extends BaseEntity {
  name: string;
  description: string;
  status: string; // 'draft', 'active', 'inactive', 'archived'
}

export interface StatusUpdate {
  status: string;
  comment?: string;
  updatedBy?: string;
}

export interface PaginatedResult<T> {
  items: T[];
  totalItems: number;
  totalPages: number;
  page: number;
  pageSize: number;
}

export interface QueryParams {
  page: number;
  pageSize: number;
  status?: string;
  searchTerm?: string;
  sortBy?: string;
  sortOrder?: string;
}
