interface Props {
  html: string
  className?: string
}

export function MarkdownContent({ html, className = '' }: Props) {
  return (
    <div
      className={`markdown-content ${className}`}
      dangerouslySetInnerHTML={{ __html: html }}
    />
  )
}
