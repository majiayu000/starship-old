'use client';

import React from 'react';
import { ContentMetadata, ContentRendererOptions } from '../../types';

interface LineBreakComponentProps {
  metadata?: ContentMetadata;
  options?: ContentRendererOptions;
}

/**
 * Component for rendering line breaks
 */
const LineBreakComponent: React.FC<LineBreakComponentProps> = () => {
  return <br />;
};

export default LineBreakComponent; 