'use client';

import React from 'react';
import BaseComponent from '../base/BaseComponent';
import { ContentMetadata, ContentRendererOptions } from '../../types';
import { sanitizeHtml } from '@/lib/sanitizeHtml';

interface HtmlComponentProps {
  content: string;
  metadata?: ContentMetadata;
  options?: ContentRendererOptions;
}

/**
 * Component for safely rendering HTML content
 */
const HtmlComponent: React.FC<HtmlComponentProps> = ({
  content,
  metadata = {},
  options = {}
}) => {
  if (!content) return null;

  const sanitizedContent = sanitizeHtml(content);

  return (
    <BaseComponent metadata={metadata}>
      <div
        className="html-content"
        dangerouslySetInnerHTML={{ __html: sanitizedContent }}
      />
    </BaseComponent>
  );
};

export default HtmlComponent;
