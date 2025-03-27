import { ContentTypeEnum } from '../types';

// Import all component types
import TextComponent from './elements/TextComponent';
import FormulaComponent from './elements/FormulaComponent';
import ImageComponent from './elements/ImageComponent';
import HtmlComponent from './elements/HtmlComponent';
import BlankComponent from './elements/BlankComponent';
import FigureComponent from './elements/FigureComponent';
import TableComponent from './elements/TableComponent';
import LineBreakComponent from './elements/LineBreakComponent';

/**
 * Registry of all available content components
 */
const ComponentRegistry = {
  // Standard components
  [ContentTypeEnum.STRING]: TextComponent,
  [ContentTypeEnum.FORMULA]: FormulaComponent,
  [ContentTypeEnum.IMAGE]: ImageComponent,
  [ContentTypeEnum.HTML]: HtmlComponent,
  [ContentTypeEnum.BLANK]: BlankComponent,
  [ContentTypeEnum.FIGURE]: FigureComponent,
  [ContentTypeEnum.TABLE]: TableComponent,
  [ContentTypeEnum.LINE_BREAK]: LineBreakComponent,
};

/**
 * Get a component by its type
 * @param type The content type to get a component for
 * @returns The component or undefined if not found
 */
export const getComponentByType = (type: ContentTypeEnum | string) => {
  return ComponentRegistry[type as ContentTypeEnum];
};

/**
 * Register a new component for a specific content type
 * @param type The content type to register the component for
 * @param component The component to register
 */
export const registerComponent = (type: ContentTypeEnum | string, component: React.ComponentType<any>) => {
  ComponentRegistry[type as ContentTypeEnum] = component;
};

export default ComponentRegistry; 