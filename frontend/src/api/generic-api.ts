import { GenericItem, PaginatedResult, QueryParams, StatusUpdate } from '@/models/generic-item';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const genericApi = {
  // Get all items with pagination and filtering
  async getAll(params: QueryParams): Promise<PaginatedResult<GenericItem>> {
    // Build query string
    const queryParams = new URLSearchParams();
    queryParams.append('page', params.page.toString());
    queryParams.append('pageSize', params.pageSize.toString());
    
    if (params.status) {
      queryParams.append('status', params.status);
    }
    
    if (params.searchTerm) {
      queryParams.append('search', params.searchTerm);
    }
    
    if (params.sortBy) {
      queryParams.append('sortBy', params.sortBy);
    }
    
    if (params.sortOrder) {
      queryParams.append('sortOrder', params.sortOrder);
    }
    
    // Make the request
    const response = await fetch(`${API_BASE_URL}/items?${queryParams.toString()}`);
    
    if (!response.ok) {
      throw new Error(`Error fetching items: ${response.statusText}`);
    }
    
    const data = await response.json();
    return data.data;
  },
  
  // Get an item by ID
  async getById(id: string): Promise<GenericItem> {
    const response = await fetch(`${API_BASE_URL}/items/${id}`);
    
    if (!response.ok) {
      throw new Error(`Error fetching item: ${response.statusText}`);
    }
    
    const data = await response.json();
    return data.data;
  },
  
  // Create a new item
  async create(item: Omit<GenericItem, 'id' | 'createdAt' | 'updatedAt'>): Promise<GenericItem> {
    const response = await fetch(`${API_BASE_URL}/items`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(item),
    });
    
    if (!response.ok) {
      throw new Error(`Error creating item: ${response.statusText}`);
    }
    
    const data = await response.json();
    return data.data;
  },
  
  // Update an existing item
  async update(id: string, item: Omit<GenericItem, 'id' | 'createdAt' | 'updatedAt'>): Promise<GenericItem> {
    const response = await fetch(`${API_BASE_URL}/items/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(item),
    });
    
    if (!response.ok) {
      throw new Error(`Error updating item: ${response.statusText}`);
    }
    
    const data = await response.json();
    return data.data;
  },
  
  // Update an item's status
  async updateStatus(id: string, update: StatusUpdate): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/items/${id}/status`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(update),
    });
    
    if (!response.ok) {
      throw new Error(`Error updating item status: ${response.statusText}`);
    }
  },
  
  // Delete an item
  async delete(id: string): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/items/${id}`, {
      method: 'DELETE',
    });
    
    if (!response.ok) {
      throw new Error(`Error deleting item: ${response.statusText}`);
    }
  },
};
