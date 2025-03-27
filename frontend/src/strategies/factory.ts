import SATOneprepStrategy from './SATOneprepStrategy';
import SATIXLStrategy from './SATIXLStrategy';
import { ReviewStrategy } from './ReviewStrategy';
import { ReviewableItem } from '@/models/reviewable-item';

export class StrategyFactory {
  private static strategies: ReviewStrategy<any>[] = [
    new SATOneprepStrategy(),
    new SATIXLStrategy(),
  ];

  static getStrategy(dataSource: string): ReviewStrategy<any> | null {
    return this.strategies.find(strategy => strategy.canHandle(dataSource)) || null;
  }

  static isDataSourceSupported(dataSource: string): boolean {
    return this.strategies.some(strategy => strategy.canHandle(dataSource));
  }

  static getAvailableDataSources(): string[] {
    return ['sat_oneprep', 'sat_ixl'];
  }
} 