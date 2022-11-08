import React, {Component, ComponentType, ErrorInfo} from 'react'
import DefaultErrorMessage from 'src/shared/components/DefaultErrorMessage'

export type ErrorMessageComponent = ComponentType<{error: Error}>

// See https://docs.honeybadger.io/lib/javascript/guides/reporting-errors.html#additional-options
interface HoneyBadgerAdditionalOptions {
  component?: string
  context?: {[key: string]: any}
  cookies?: {[key: string]: any}
  name?: string
  params?: {[key: string]: any}
}

export const reportError = (
  error: Error,
  additionalOptions?: HoneyBadgerAdditionalOptions
): void => {
  let additionalContext = {}
  if (additionalOptions && additionalOptions.context) {
    additionalContext = {...additionalOptions.context}
  }

  const _context = {
    ...additionalContext,
    // ...getUserFlags(),
  }

  let options: HoneyBadgerAdditionalOptions = {}
  if (additionalOptions) {
    options = {...additionalOptions}

    delete options.context // already included in the above context object
  }

  console.error(error)
}

export const parseComponentName = (errorInfo: ErrorInfo): string => {
  return errorInfo.componentStack
    .trim()
    .split('\n')
    .map(s => s.split(' ')[1])[0]
}

export function ErrorHandlingWith(Error: ErrorMessageComponent) {
  return <P, S, T extends {new (...args: any[]): Component<P, S>}>(
    constructor: T
  ) => {
    class Wrapped extends constructor {
      public static get displayName(): string {
        return constructor.name
      }

      // @ts-ignore
      private error: Error = null

      public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        this.error = error
        this.forceUpdate()

        reportError(error, {component: parseComponentName(errorInfo)})
      }

      public render() {
        if (this.error) {
          return <Error error={this.error} />
        }

        return super.render()
      }
    }

    return Wrapped
  }
}

export const ErrorHandling = ErrorHandlingWith(DefaultErrorMessage)
