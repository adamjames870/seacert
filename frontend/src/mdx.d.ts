declare module '*.mdx' {
  import { ReactNode } from 'react'
  export const meta: {
    title: string
    description: string
    slug: string
  }
  const MDXComponent: (props: any) => ReactNode
  export default MDXComponent
}
