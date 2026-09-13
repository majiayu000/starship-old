import DOMPurify from 'isomorphic-dompurify';

/**
 * Allowlist sanitizer for untrusted HTML before dangerouslySetInnerHTML.
 * Strips scripts, event handlers, javascript: URLs, and other XSS vectors.
 */
export function sanitizeHtml(dirty: string): string {
  if (!dirty) return '';
  return DOMPurify.sanitize(dirty, {
    USE_PROFILES: { html: true },
  });
}

/**
 * Allowlist sanitizer for untrusted SVG markup before dangerouslySetInnerHTML.
 * Permits SVG structure while removing scripts, foreignObject abuse, and handlers.
 */
export function sanitizeSvg(dirty: string): string {
  if (!dirty) return '';
  return DOMPurify.sanitize(dirty, {
    USE_PROFILES: { svg: true, svgFilters: true },
  });
}

/**
 * Sanitizer for mixed HTML/SVG content (figures, tables with inline SVG, etc.).
 */
export function sanitizeRichHtml(dirty: string): string {
  if (!dirty) return '';
  return DOMPurify.sanitize(dirty, {
    USE_PROFILES: { html: true, svg: true, svgFilters: true },
  });
}
