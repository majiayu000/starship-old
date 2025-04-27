// Define filter options type
// This is used for filtering items by status or other properties

export type FilterOptions = {
  status: string[];
  [key: string]: string[] | undefined;
};