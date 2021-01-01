import React, {ComponentType} from 'react'

type ErrorMessageComponent = ComponentType<{error: Error}>

const DefaultErrorMessage: ErrorMessageComponent = () => {
  return (
    <p
      className="default-error-message"
      style={{display: 'flex', placeContent: 'center'}}
    >
      An error has occurred. Please report the issue&nbsp;
      <a href="/">here</a>.
    </p>
  )
}

export default DefaultErrorMessage
