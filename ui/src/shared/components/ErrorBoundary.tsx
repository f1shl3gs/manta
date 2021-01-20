// Libraries
import React, {Component, ComponentType, ErrorInfo} from 'react'

// Components
import DefaultErrorMessage from './DefaultErrorMessage'

export type ErrorMessageComponent = ComponentType<{error: Error}>

interface ErrorBoundaryProps {
  errorComponent: ErrorMessageComponent
}

interface ErrorBoundaryState {
  error: Error | null
}

class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  public static defaultProps = {errorComponent: DefaultErrorMessage}

  public state: ErrorBoundaryState = {error: null}

  public static getDerivedStateFromError(error: Error) {
    return {error}
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // reportError(error, { component: parseComponentName(errorInfo) });
    console.log('error:', error)
    console.log('errorInfo:', errorInfo)
  }

  public render() {
    const {error} = this.state

    if (error) {
      return <this.props.errorComponent error={error} />
    }

    return this.props.children
  }
}

export default ErrorBoundary
