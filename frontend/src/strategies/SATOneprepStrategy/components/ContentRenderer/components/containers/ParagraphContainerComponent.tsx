'use client';

import React from 'react';
import { ContentItemType, ContentMetadata, ContentRendererOptions, ContentTypeEnum } from '../../types';
import ContentItemComponent from '../ContentItemComponent';

interface ParagraphContainerComponentProps {
  children?: React.ReactNode;
  items?: ContentItemType[];
  metadata?: ContentMetadata;
  options?: ContentRendererOptions;
}

/**
 * Component for grouping content items into paragraphs
 */
const ParagraphContainerComponent: React.FC<ParagraphContainerComponentProps> = ({
  children,
  items = [],
  metadata = {},
  options = {}
}) => {
  const styles = metadata?.style || [];
  let containerClasses = 'paragraph-container';
  
  // Apply paragraph-specific styles
  if (styles.includes('indent')) containerClasses += ' pl-8';
  if (styles.includes('center')) containerClasses += ' text-center';
  if (styles.includes('right')) containerClasses += ' text-right';
  
  // Use either children or render items
  return (
    <div className={containerClasses}>
      {children || items.map((item, index) => (
        <ContentItemComponent 
          key={`p-item-${index}`}
          item={item}
          index={index}
          isLast={index === items.length - 1}
          totalItems={items.length}
          options={options}
        />
      ))}
    </div>
  );
};

export default ParagraphContainerComponent; 