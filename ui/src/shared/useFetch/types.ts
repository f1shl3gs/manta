import {ReactNode} from 'react'

type NonUndefined<A> = A extends undefined ? never : A

export type FunctionKeys<T extends object> = {
  [K in keyof T]-?: NonUndefined<T[K]> extends Function ? K : never
}[keyof T]

export declare type NonFunctionKeys<T extends object> = {
  [K in keyof T]-?: NonUndefined<T[K]> extends Function ? never : K
}[keyof T]

export enum HTTPMethod {
  DELETE = 'DELETE',
  GET = 'GET',
  HEAD = 'HEAD',
  OPTIONS = 'OPTIONS',
  PATCH = 'PATCH',
  POST = 'POST',
  PUT = 'PUT',
}

export interface DoFetchArgs {
  url: string
  options: RequestInit
  response: {
    id: string
  }
}

export interface FetchContextTypes {
  url: string
  options: IncomingOptions
  graphql?: boolean
}

export interface FetchProviderProps {
  url?: string
  options?: IncomingOptions
  graphql?: boolean
  children: ReactNode
}

export type BodyOnly = (body: BodyInit | object) => Promise<any>

export type RouteOnly = (route: string) => Promise<any>

export type RouteAndBodyOnly = (
  route: string,
  body: BodyInit | object
) => Promise<any>

export type RouteOrBody = string | BodyInit | object
export type UFBody = BodyInit | object
export type RetryOpts = {attempt: number; error?: Error; response?: Response}

export type NoArgs = () => Promise<any>

export type FetchData = (
  routeOrBody?: string | BodyInit | object,
  body?: BodyInit | object
) => Promise<any>

export type RequestInitJSON = RequestInit & {
  headers: {
    'Content-Type': string
  }
}

export interface ReqMethods {
  get: (route?: string) => Promise<any>
  post: FetchData
  patch: FetchData
  put: FetchData
  del: FetchData
  delete: FetchData
  query: (query: string, variables?: BodyInit | object) => Promise<any>
  mutate: (mutation: string, variables?: BodyInit | object) => Promise<any>
  abort: () => void
}

export interface Data<TData> {
  data: TData | undefined
}

export interface ReqBase<TData> {
  data: TData | undefined
  loading: boolean
  error: Error
}

export interface Res<TData> extends Response {
  data?: TData | undefined
}

export type Req<TData = any> = ReqMethods & ReqBase<TData>

export type UseFetchArgs = [
  (string | IncomingOptions | OverwriteGlobalOptions)?,
  (IncomingOptions | OverwriteGlobalOptions | any[])?,
  any[]?
]

export type UseFetchArrayReturn<TData> = [
  Req<TData>,
  Res<TData>,
  boolean,
  Error
]

export type UseFetchObjectReturn<TData> = ReqBase<TData> &
  ReqMethods & {
    request: Req<TData>
    response: Res<TData>
  }

export type UseFetch<TData> = UseFetchArrayReturn<TData> &
  UseFetchObjectReturn<TData>

export type Interceptors<TData = any> = {
  request?: ({
    options,
    url,
    path,
    route,
  }: {
    options: RequestInit
    url?: string
    path?: string
    route?: string
  }) => Promise<RequestInit> | RequestInit
  response?: ({response}: {response: Res<TData>}) => Promise<Res<TData>>
}

export interface CustomOptions {
  data: any
  interceptors: Interceptors
  loading: boolean
  onAbort: () => void
  onError: OnError
  onNewData: (currData: any, newData: any) => any
  onTimeout: () => void
  persist: boolean
  perPage: number
  responseType: ResponseType
  retries: number
  retryOn: RetryOn
  retryDelay: RetryDelay
  suspense: boolean
  timeout: number
}

// these are the possible options that can be passed
export type IncomingOptions = Partial<CustomOptions> &
  Omit<RequestInit, 'body'> & {body?: BodyInit | object | null}
// these options have `context` and `defaults` applied so
// the values should all be filled
export type Options = CustomOptions &
  Omit<RequestInit, 'body'> & {body?: BodyInit | object | null}

export type OverwriteGlobalOptions = (options: Options) => Options

export type RetryOn =
  | (<TData = any>({attempt, error, response}: RetryOpts) => Promise<boolean>)
  | number[]
export type RetryDelay =
  | (<TData = any>({attempt, error, response}: RetryOpts) => number)
  | number

export type BodyInterfaceMethods = Exclude<
  FunctionKeys<Body>,
  'body' | 'bodyUsed' | 'formData'
>
export type ResponseType = BodyInterfaceMethods | BodyInterfaceMethods[]

export type OnError = ({error}: {error: Error}) => void

export type UseFetchArgsReturn = {
  host: string
  path?: string
  customOptions: {
    interceptors: Interceptors
    onAbort: () => void
    onError: OnError
    onNewData: (currData: any, newData: any) => any
    onTimeout: () => void
    perPage: number
    persist: boolean
    responseType: ResponseType
    retries: number
    retryDelay: RetryDelay
    retryOn: RetryOn | undefined
    suspense: boolean
    timeout: number
    // defaults
    loading: boolean
    data?: any
  }
  requestInit: RequestInit
  dependencies?: any[]
}

/**
 * Helpers
 */
export type ValueOf<T> = T[keyof T]

export type NonObjectKeysOf<T> = {
  [K in keyof T]: T[K] extends Array<any> ? K : T[K] extends object ? never : K
}[keyof T]

export type ObjectValuesOf<T extends Record<string, any>> = Exclude<
  Exclude<Extract<ValueOf<T>, object>, never>,
  Array<any>
>

export type UnionToIntersection<U> = (
  U extends any ? (k: U) => void : never
) extends (k: infer I) => void
  ? I
  : never

export type Flatten<T> = Pick<T, NonObjectKeysOf<T>> &
  UnionToIntersection<ObjectValuesOf<T>>
