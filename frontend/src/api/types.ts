import { ReviewableItem, ReviewStatusUpdate, PaginatedResult, QueryParams } from '../models/reviewable-item';
import { SATOneprepItem } from '../models/sat-oneprep-item';
import { SATIXLItem } from '../models/sat-ixl-item';

export type DataSourceItemType = {
  'sat_oneprep': SATOneprepItem;
  'sat_ixl': SATIXLItem;
};

export type FilterOptions = {
  reviewStatus: string[];
  domain?: string[];
  skill?: string[];
  [key: string]: string[] | undefined;
}; 