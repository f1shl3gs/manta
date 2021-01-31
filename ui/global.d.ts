declare module '*.md' {
  const value: string
  export default value
}

declare module '*.svg' {
  export const ReactComponent: SFC<SVGProps<SVGSVGElement>>
  const src: string
  export default src
}

declare module '*.png' {
  const value: any
  export = value
}
