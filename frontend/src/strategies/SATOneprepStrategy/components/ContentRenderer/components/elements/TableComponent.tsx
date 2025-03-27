'use client';

import React from 'react';
import BaseComponent from '../base/BaseComponent';
import { ContentMetadata, ContentRendererOptions } from '../../types';
import HtmlComponent from './HtmlComponent';

interface TableContent {
  content: string;
  caption?: string;
}

interface TableComponentProps {
  content: TableContent;
  metadata?: ContentMetadata;
  options?: ContentRendererOptions;
}

/**
 * Component for rendering tables
 */
const TableComponent: React.FC<TableComponentProps> = ({
  content,
  metadata = {},
  options = {}
}) => {
  if (!content || !content.content) return null;
  
  return (
    <BaseComponent metadata={metadata} className="table-component my-4">
      <div className="table-wrapper overflow-x-auto">
        <HtmlComponent content={content.content} />
      </div>
      
      {content.caption && (
        <div className="table-caption text-sm text-gray-500 mt-1 text-center">
          {content.caption}
        </div>
      )}
    </BaseComponent>
  );
};

export default TableComponent; 