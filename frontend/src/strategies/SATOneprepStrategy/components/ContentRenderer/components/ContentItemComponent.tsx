'use client';

import React from 'react';
import { ContentItemProps, ContentTypeEnum } from '../types';
import { getComponentByType } from './ComponentRegistry';

// Import CSS
import './ContentRenderer.css';

/**
 * Factory component that renders the appropriate component based on content type
 */
const ContentItemComponent: React.FC<ContentItemProps> = ({
  item,
  index,
  isLast,
  totalItems,
  options
}) => {
  if (!item) return null;
  
  const { type } = item;
  
  // Get the appropriate component from the registry
  const Component = getComponentByType(type);
  
  if (!Component) {
    console.warn(`Unknown content type: ${type}`);
    return null;
  }
  
  // Use a type assertion since TypeScript can't infer the correct content structure
  // from the union type, but at runtime this will work correctly
  const componentProps = {
    // We need to use type assertion to handle the union type
    content: (item as any).content,
    metadata: item.metadata,
    options
  };
  
  // Render the component with the appropriate props
  return <Component {...componentProps} />;
};

export default ContentItemComponent; 